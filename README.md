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

## Usage

To add an entry as a log call `caplog`, which will
open `vi` text editor. The message should follow git commit
message conventions to provide a more clear log entry as content.

```
$ caplog
```

To add a ''quick'' entry log. Call `caplog` with one argument.
The argument will be used as the log entry.

```
$ caplog "Some entry text"
```

## Log history

Logs are created under `$HOME/.caplog/logbook`, which is initialised as a git repository.
The logs are written directly to main branch following `<hash>_<timestamp>.log` pattern. 
After file is created it will automatically be committed to the `.caplog/logbook` repository.

### Log entries

Created logs are human readable—well it depends of the writer.

```log
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
grep -r <keyword> ~/.caplog/logbook/
```
