# changelogger

Changelogger is a CLI tool for maintaining changelog files and version releases.

## Installation

To install it, you can either use the following one-liner or download the [binary](https://github.com/MarkusFreitag/changelogger/releases/latest/download/changelogger_linux_amd64.deb), [deb](https://github.com/MarkusFreitag/changelogger/releases/latest/download/changelogger_linux_amd64.rpm) or [rpm](https://github.com/MarkusFreitag/changelogger/releases/latest/download/changelogger_linux_amd64.tar.gz) package from the releases.
```bash
curl -Ls https://git.io/changelogger.sh | sh
```

When you using it, every 24h, it checks whether new updates are available. If so, it can be updated using its `update` command.
```bash
changelogger update
```

## Usage

### Add a new entry

To initialize the file or create a new entry, simply run the tool without additional commands. It will open an editor where you can enter your changes. In the editor you will see the recent changes for your user. Just add leave these lines and add your new ones below it.

### Make a new release

You will be asked, what kind of version bump you would like to do. Then the last version will be bumped and set for the last unreleased entry. This will not only update the version but also the author and the timestamp.
```bash
changelogger release new
```
