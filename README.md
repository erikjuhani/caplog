# caplog

Simple logging system for code maneuvering captains.

!!! Caplog is in a prototype stage !!!

## Prerequisites

* [Git](https://git-scm.com/)
* [The Go Programming Language](https://go.dev/dl/)

## Installation

Currently caplog is in a stage, where it needs to be compiled/installed from source.

To use caplog locally you must do the following commands:

```
$ git clone git@github.com:erikjuhani/caplog.git
$ cd caplog
$ go install
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

### Finding entries

The logs are human readable and can be looked or parsed with tooling designed for text files. For example with grep.

```
grep -r <keyword> ~/.caplog/logbook/
```
