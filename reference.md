# Caplog command and functionality referecence

## Config

- editor=vi|vim|code|etc
- repository=<gitrepopath>
- remote_repository=<gitremoterepopath>

## Commands

```
caplog
```

- opens text-editor
- after save writes to file

```
caplog "New log entry"
```

- directly writes the preceeding text to file

```
caplog -t tag0 tag1
```

- opens text-editor if text not given
- adds tags to respective log entry

## Logical structure for packages

cmd/root.go?
fs/ - filesystem operations

/core.go - main functionality of the program / break up later

## Log entry structure

Three options (last one is flat):
1. Saved to filesystem as .caplog/<gitrepo>/<year>/<month>/<day>.log
2. Saved to filesystem as .caplog/<gitrepo>/<year>/<month>/<date>.log
3. Saved to filesystem as .caplog/<gitrepo>/<date>.log

## Log entry structure

`<day>/<month>/<year>

<hours>:<minutes>   <summary>

                    <body>

                    tags: [<tag>]
`

### Example log entries

`14/05/2022

12:27   Initialised caplog repository with dependencies

        tags: caplog, go 

12:44   Wrote reference documentation for caplog

        The reference document is only for development purposes,
        so it does not reflect the goals or future of the
        application itself.

        Part of figuring out what would be a nice system to use
        for logging different things with minimal effort on making
        those particular logs.

        tags: caplog, documentation, thoughts
`
