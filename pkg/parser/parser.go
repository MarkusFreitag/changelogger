package parser // import "github.com/MarkusFreitag/changelogger/pkg/parser"

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/MarkusFreitag/changelogger/pkg/stringutil"
	"github.com/Masterminds/semver"
)

var (
	RgxFirstLine = regexp.MustCompile(`#\s[A-Za-z]+\s[A-Za-z]+\s(?P<version>v?\d+.\d+.\d+(?:-\d+)?)\s\((?P<date>\d{4}-\d{2}-\d{2})\)`)
	RgxLastLine  = regexp.MustCompile(`\*Released\sby\s(?P<name>.*?)\s<(?P<mail>.*?)>`)
)

func reSubMatchMap(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	if match == nil {
		return nil
	}
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 && name != "" {
			subMatchMap[name] = match[i]
		}
	}
	return subMatchMap
}

type Author struct {
	Name  string
	Email string
}

type Changes map[string]string

func (c Changes) FormatByAuthor(author string, indent int) string {
	output := stringutil.IncrIndent(fmt.Sprintf("* **%s**", author), indent)
	if indent == 0 {
		indent = 2
	} else {
		indent += indent
	}
	if changes, ok := c[author]; ok {
		output += "\n" + stringutil.IncrIndent(changes, indent)
	}
	return output + "\n"
}

func (c Changes) Format(indent int) string {
	var index int
	parts := make([]string, len(c))
	for author, changes := range c {
		parts[index] = stringutil.IncrIndent(fmt.Sprintf("* **%s**", author), indent)
		parts[index] += "\n" + stringutil.IncrIndent(changes, indent*indent)
		index++
	}
	return strings.Join(parts, "\n")
}

type Releases []*Release

func PrependRelease(releases Releases, release *Release) Releases {
	buf := make(Releases, len(releases)+1)
	buf[0] = release
	for index, r := range releases {
		buf[index+1] = r
	}
	return buf
}

type Release struct {
	Header   string
	Footer   string
	Version  *semver.Version
	Date     time.Time
	Changes  Changes
	By       Author
	Released bool
}

func NewRelease() *Release {
	return &Release{
		Header:  "# For next release",
		Footer:  "*Not released yet*",
		Changes: make(Changes),
	}
}

func (r *Release) SortedAuthors() []string {
	authors := make([]string, 0)
	for author := range r.Changes {
		authors = append(authors, author)
	}
	sort.Strings(authors)
	return authors
}

func (r *Release) MarshalJSON() ([]byte, error) {
	type version struct {
		Major int64  `json:"major"`
		Minor int64  `json:"minor"`
		Patch int64  `json:"patch"`
		Ext   string `json:"ext"`
		Full  string `json:"full"`
	}
	rel := &struct {
		Header   string  `json:"header"`
		Footer   string  `json:"footer"`
		Version  version `json:"version"`
		Date     string  `json:"date"`
		Changes  Changes `json:"changes"`
		By       Author  `json:"by"`
		Released bool    `json:"released"`
	}{
		Header:   r.Header,
		Footer:   r.Footer,
		Changes:  r.Changes,
		By:       r.By,
		Released: r.Released,
	}

	if r.Released {
		rel.Version = version{
			Major: r.Version.Major(),
			Minor: r.Version.Minor(),
			Patch: r.Version.Patch(),
			Ext:   r.Version.Prerelease(),
			Full:  r.Version.String(),
		}
		rel.Date = r.Date.Format("2006-01-02")
	}
	return json.Marshal(rel)
}

func (r *Release) GenerateHeader(bump string) {
	if r.Version == nil {
		r.Header = "# For next Release"
	}
	r.Header = fmt.Sprintf("# %s Release %s (%s)", strings.Title(bump), r.Version.Original(), r.Date.Format("2006-01-02"))
}

func (r *Release) GenerateFooter() {
	if r.Version == nil {
		r.Footer = "*Not released yet*"
	}
	r.Footer = fmt.Sprintf("*Released by %s <%s>*", r.By.Name, r.By.Email)
}

func (r *Release) Show() string {
	parts := make([]string, 3)
	parts[0] = r.Header
	for _, author := range r.SortedAuthors() {
		parts[1] += r.Changes.FormatByAuthor(author, 2)
	}
	parts[2] = r.Footer
	return strings.Join(parts, "\n") + "\n"
}

func (r *Release) ParseHeader() {
	matches := reSubMatchMap(RgxFirstLine, r.Header)
	if matches != nil {
		v, err := semver.NewVersion(matches["version"])
		if err == nil {
			r.Released = true
			r.Version = v
		}
		d, err := time.Parse("2006-01-02", matches["date"])
		if err == nil {
			r.Date = d
		}
	}
}

func (r *Release) ParseFooter() {
	matches := reSubMatchMap(RgxLastLine, r.Footer)
	if matches != nil {
		r.By = Author{Name: matches["name"], Email: matches["mail"]}
	}
}

func (r *Release) ParseBody(body string) {
	body = stringutil.DecrIndent(body, stringutil.IndentLvl(body))
	var part string
	for {
		part, body = stringutil.SplitBlock(body, "\n*")
		authorStr, changesStr := stringutil.SplitBlock(part, "\n  *")
		authorStr = strings.TrimPrefix(authorStr, "* **")
		authorStr = strings.TrimSuffix(authorStr, "**")
		changesStr = stringutil.DecrIndent(changesStr, stringutil.IndentLvl(changesStr))
		r.Changes[authorStr] = changesStr
		if body == "" {
			break
		}
	}
}

func parseRelease(block string) (*Release, error) {
	header, body, footer := splitReleaseBlock(block)
	r := &Release{
		Header:   header,
		Footer:   footer,
		Changes:  make(Changes),
		Released: false,
	}
	r.ParseHeader()
	r.ParseFooter()
	r.ParseBody(body)
	return r, nil
}

func splitReleaseBlock(block string) (string, string, string) {
	parts := strings.SplitN(block, "\n", 2)
	header := parts[0]
	block = strings.TrimPrefix(block, header+"\n")

	block = strings.TrimRight(block, "\n")
	parts = strings.Split(block, "\n")
	footer := parts[len(parts)-1]
	block = strings.TrimSuffix(block, footer)
	return header, strings.TrimSuffix(block, "\n\n"), footer
}

func blockScanner(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := strings.Index(string(data), "*\n\n#"); i >= 0 {
		return i + 3, data[0 : i+1], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return
}

func ReadFile(filename string) (Releases, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	releases := make(Releases, 0)
	scanner := bufio.NewScanner(file)
	scanner.Split(blockScanner)
	for scanner.Scan() {
		rel, err := parseRelease(scanner.Text())
		if err != nil {
			return nil, err
		}
		releases = append(releases, rel)
	}

	if scanner.Err() != nil {
		return nil, err
	}
	return releases, nil
}
