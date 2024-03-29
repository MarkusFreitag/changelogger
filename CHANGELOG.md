# Patch Release v0.6.1 (2023-11-23)
  * **Tom Siewert**
    * debian: lowercase unstable

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Minor Release v0.6.0 (2023-01-18)
  * **Markus Freitag**
    * pkg/gitconfig: replace character-based with ini config parser

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Patch Release v0.5.3 (2022-08-23)
  * **Tom Siewert**
    * pkg/gitconfig: Add ~/.config/git/config to possible files
    * ci: Use official golangci-lint action

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Patch Release v0.5.2 (2022-05-10)
  * **Tom Siewert**
    * ci: Migrate to GitHub Workflows
    * install.sh: Update to latest to support arm64 properly

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Patch Release v0.5.1 (2022-03-23)
  * **Markus Freitag**
    * add arm64 builds
    * update module dependencies

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Minor Release v0.5.0 (2020-10-19)
  * **Markus Freitag**
    * - cmd/debian: Add version and date flag for experimental releases

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Patch Release v0.4.3 (2020-04-18)
  * **Markus Freitag**
    * Move darwin from goarch to goos

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Patch Release v0.4.2 (2020-04-18)
  * **Markus Freitag**
    * Add binary builds for darwin

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Patch Release v0.4.1 (2020-03-26)
  * **Markus Freitag**
    * Add missing whitespaces to debian changelog template

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Minor Release v0.4.0 (2020-03-05)
  * **Markus Freitag**
    * Drop auto update checker
    * Prefix changelog entry with an asterisk if missing
    * Remove empty trailing lines
    * Unify error handling

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Patch Release v0.3.3 (2020-02-12)
  * **Markus Freitag**
    * Fix goreleaser config

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Patch Release v0.3.2 (2020-02-12)
  * **Markus Freitag**
    * **bugfix** Fix early return when generating debian changelog

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Patch Release v0.3.1 (2019-07-30)
  * **Markus Freitag**
    * Remove 'v' from version string when generating debian changelog

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Minor Release v0.3.0 (2019-07-26)
  * **Markus Freitag**
    * Add CLI command `debian`
      * `debian dummy` creates a changelog file containing only the latest
        release with the hint to check CHANGELOG.md
      * `debian full` generates a debian formated changelog out of CHANGELOG.md

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Patch Release v0.2.5 (2019-07-24)
  * **Markus Freitag**
    * Sort authors within a release

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Patch Release v0.2.4 (2019-05-06)
  * **Markus Freitag**
    * Set default versioning to semver

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Patch Release v0.2.3 (2019-02-12)
  * **Markus Freitag**
    * Enable syntax highlighting in editor
    * Run update check just every 24h

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Patch Release v0.2.2 (2019-02-01)
  * **Markus Freitag**
    * Do not exit program when update check fails

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Patch Release v0.2.1 (2019-02-01)
  * **Markus Freitag**
    * Bugfix for the update checking routine

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Minor Release v0.2.0 (2019-02-01)
  * **Markus Freitag**
    * Implement update routine
      * automatic update checker in the init phase
      * `update` command to selfupdate the binary
    * Implement `json` command
    * Add prerun check for minimal version

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Minor Release v0.1.0 (2019-01-31)
  * **Markus Freitag**
    * Initial release, implemented features
      * create or update CHANGELOG files
      * bump version for a new release
      * show information about the last release

*Released by Markus Freitag <fmarkus@mailbox.org>*
