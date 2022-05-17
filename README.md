# caplog

Simple logging system for code maneuvering captains.

!!! Caplog is in a prototype stage !!!

## Prerequisites

* [Git](https://git-scm.com/)
* [The Go Programming Language](https://go.dev/dl/)

## Installation

Installation can be done by simply using `go install`.

```
$ go install github.com/erikjuhani/caplog@latest
```

## Configuration

## Usage

To add an entry as a log call `caplog`, which will open by default `vi` text editor.
The message should follow git commit message conventions to provide a more clear log entry as content.

The default editor can be changed to any preferred editor by providing a configuration file `.caplog/config`.

```toml
editor="nvim"
```

```
$ caplog
```

To add a ''quick'' entry log. Call `caplog` with one argument.
The argument will be used as the log entry.

```
$ caplog "Some entry text"
```

## Log history

Logs are created by default under `$HOME/.caplog/capbook`, which is initialized as a git repository.
The logs are written directly to main branch following `<timestamp>_<hash>.log` pattern.
After file is created it will automatically be committed to the `.caplog/capbook` repository.

The default location can be changed to any preferred location. It can also be an existing git repository.

NOTE: if using existing git repository, currently the logs will always be added to root.

```toml
[git]
  local_repository="~/mybook"
```

### Log entries

Created logs are in human readable text format.
The writer has all the freedom of composing the message.

```log
// 2022-05-16T19:20:17_49b13c5.log
19:20	Hello this is my first log entry!

Used as an example to provide some idea of the log entry.

You can write anything here and even use keywords or tags
to provide easier content seeking capabilities.

tags: example, caplog
```

### Log storage

Logs are stored as files in the filesystem and ultimately the changes
are stored inside a git repository. Each written log entry is completely reflected
in the git commit message, which enables users to traverse the log history using
familiar tools like `git log`.

### Finding log entries

The logs are human readable and can be looked or parsed with tooling designed for text files. For example with grep.

```
grep -r <keyword> ~/.caplog/caplog/
```

## TODO (in priority order):

- [ ] If repository is remote tracked try `push --force-with-lease`

- [ ] Implement `-t|--tags` for adding tags to a log entry

- [ ] Implement `-p|--page` to separate log entries into pages by default uses root `capbook/<page>`

- [ ] Implement `-c|--config` to change configuration values `git.local_repository=~/mylogs`

- [ ] Implement `-g|--get-dir` return capbook directory for easier use of find operations

- [ ] Implement `-d|--dry-run` run the command dry, no filesystem or git changes

- [ ] Use github actions to create separate binaries for different system architectures

- [ ] Add brew formula for easier install with homebrew

- [ ] Add more installation methods to README `wget` and `brew`

- [ ] Add charmbracelet TUI library to create better loading experience

- [ ] Maybe using git as a detached process if experience is really slow?

- [ ] Maybe use libgit2 instead of calling `git` executable?

- [ ] Log file encryption -> commit message obfuscated

