---
theme: dark
---

# Whoami

Erik Kinnunen

Software man by day, byte pusher by night 

[github](https://github.com/erikjuhani)

---

# Caplog

A journaling utility tool made for terminal users

Made with _Go_ (the best programming language ❤️ )

<br/>

_Special thanks to_

[Andy Davies](https://github.com/Pondidum)

[Sampo Siltanen](https://github.com/ssiltanen)

[Raine Virta](https://github.com/raine)

[Taylor Thompson](https://github.com/jamestthompson3)

---

# I wanted

A simple tool to store ideas, conversations and plans

A personal shareable journal tool, which has the least amount of 
friction to my own workflow

To traverse the log history by using keywords or tags

Use git for easy sharing

---

# Workflow

Write quick logs by giving a string as the only argument

```bash
caplog "Quick log"
```

Or compose longer logs using chosen text editor

```bash
caplog
```

Uses git as the backend to store the logs

---

# Configuration

Can be configured as _good_ software should 😉

```toml
editor="nvim"
[git]
  local_repository="~/mybook"
```
---

# Logs history

Logs are stored as human readable text

Logs are written directly to the active branch 
following `<timestamp>_<hash>.log` pattern

Logs are automatically committed to repository

If git has remote set will do automatic rebasing and push

---

# Finding logs

Using grep
```bash
grep -r test ~/.caplog/capbook/
```

Using git
```bash
cd ~/.caplog/capbook/
git log
```

---

# Planned features

`-t|--tags` for adding tags to a log entry

`-g|--get-dir` return capbook directory for easier use of find operations

`-p|--page` to separate log entries into pages by default uses root `capbook/<page>`

`-c|--config` to change configuration values `git.local_repository=~/mylogs`

Log file encryption + commit message obfuscated

---

# Demo
