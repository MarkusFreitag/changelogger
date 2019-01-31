package editor

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"

	shellquote "github.com/kballard/go-shellquote"
)

func Open(initialValue *string) error {
	var editor string
	if val, ok := os.LookupEnv("EDITOR"); ok {
		editor = val
	}
	if val, ok := os.LookupEnv("VISUAL"); ok {
		editor = val
	}
	if editor == "" {
		return errors.New("set EDITOR or VISUAL variable")
	}

	file, err := ioutil.TempFile("", "changelogger")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	if _, err := file.WriteString(*initialValue); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	args, err := shellquote.Split(editor)
	if err != nil {
		return err
	}
	args = append(args, file.Name())

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	raw, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return err
	}

	*initialValue = string(raw)
	return nil
}
