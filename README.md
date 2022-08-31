# caplog

Simple logging system for code maneuvering captains.

**⚠️Caplog is in a prototype stage⚠️**

## Prerequisites

* [Git](https://git-scm.com/)
* [The Go Programming Language](https://go.dev/dl/) (version `1.18` or higher is required)

## Installation

Installation can be done simply by using `go install`.

```bash
go install github.com/erikjuhani/caplog@latest
```

## Usage

To add an entry as a log call `caplog`, which will open by default `vi` text editor.
The message should follow git commit message conventions to provide a more clear log entry as content.

```bash
caplog
```

The default editor can be changed to any preferred editor by providing a configuration file `.caplog/config`.

```toml
editor="nvim"
```

To add a ''quick'' entry log. Call `caplog` with one argument.
The argument will be used as the log entry.

```bash
caplog "Some entry text"
```

## Log history

Logs are created by default under `$HOME/.caplog/capbook`, which is initialized as a git repository.
The logs are written directly to main branch following `<day>-<month>-<year>.log.md` pattern.
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

The logs can be easily read with markdown viewers like `glow`.

```log
// 16-05-2022.log.md
---
date: Monday, May 16, 2022
---

19:20	Hello this is my first log entry!

	Used as an example to provide some idea of the log entry.

	You can write anything here and even use keywords or tags
	to provide easier content seeking capabilities.

	tags: example, caplog
```

#### Tagging

Logs can be tagged by either writing it in the log entry or using caplog `-t` flag.

Logs can be tagged with one or multiple tags at once, these tags will be added to the end of the log entry.

```bash
caplog "New entry with tags" -t tag0 -t tag1

caplog "New entry with tags - comma separation" -t tag0,tag1
```

### Pages

Logs can be grouped under sub-directories or what I like to call _pages_.

To write log entry under a page you need to provide a `-p` flag with the name of the sub-directory.

```bash
caplog "New entry in to different page" -p subpage
```

### Configuration

Configuration can either be adjusted by manually writing to the caplog config file or by
using config flag option to provide configuration changes through cli.

```bash
caplog -c git.local_repository=~/mybook
```

User can also provide multiple configuration values at once.

```bash
caplog -c git.local_repository=~/mybook -c editor=vim
```

### Log storage

Logs are stored as files in the filesystem and ultimately the changes
are stored inside a git repository. Each written log entry is completely reflected
in the git commit message, which enables users to traverse the log history using
familiar tools like `git log`.

### Finding log entries

The logs are human readable and can be looked or parsed with tooling designed for text files. For example with grep.

Using `grep` command to find certain logs with `<keyword>`. Use for example `cat` to view the actual found logs.

```bash
grep --exclude-dir=.git -lrF <keyword> ~/.caplog/capbook | xargs cat
```

With directory flag using bash

```bash
grep --exclude-dir=.git -lrF <keyword> $(caplog -g) | xargs cat
```
