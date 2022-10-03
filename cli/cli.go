package cli

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/erikjuhani/caplog/config"
	"github.com/erikjuhani/caplog/core"
	"github.com/erikjuhani/miniflag"
)

var (
	dir       = miniflag.Flag("getdir", "g", false, "Returns the local repository directory")
	page      = miniflag.Flag("page", "p", "", "Saves log entry to <sub-directory>/<page>")
	workspace = miniflag.Flag("workspace", "w", "", "Changes workspace to given <workspace> if it exists")
	tags      = miniflag.Flag("tag", "t", TagsFlag{}, "Adds `<tag>` to log entry")
	setConfig = miniflag.Flag("config", "c", ConfigFlag{}, "Changes config setting with `<key=value>`")
)

var (
	ErrKeyNeedsValueF      = func(k string) error { return fmt.Errorf("key needs a value (ex. %s=<value>)", k) }
	ErrExpectedOneArgument = func(n int) error { return fmt.Errorf("expected 1 argument, got %d", n) }
	ErrWriteLog            = func(e error) error { return fmt.Errorf("failed to write log - %w", e) }
)

type TagsFlag []string

func (s *TagsFlag) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *TagsFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

type ConfigFlag map[string]string

func (c *ConfigFlag) String() string {
	return fmt.Sprintf("%s", *c)
}

func (c *ConfigFlag) Set(value string) error {
	pair := strings.Split(value, "=")

	k := pair[0]

	if len(pair) < 2 {
		return ErrKeyNeedsValueF(k)
	}

	v := pair[1]

	(*c)[k] = v

	return nil
}

func Run() error {
	out := miniflag.CommandLine.Output()
	miniflag.Parse()

	if *dir {
		// Return current repository path
		fmt.Fprintln(out, config.Config.CurrentWorkspace)
		return nil
	}

	if *workspace != "" {
		if exists := config.Config.Workspaces.Has(*workspace); !exists {
			return config.ErrWorkspaceIsNotValid(*workspace, config.Config.Workspaces)
		}

		if err := config.Write(map[string]string{config.CurrentWorkspaceKey: *workspace}); err != nil {
			return err
		}

		fmt.Fprintf(out, "workspace changed to \"%s\"", *workspace)

		return nil
	}

	// Config flags were used needs to do configuration change
	if len(*setConfig) > 0 {
		if err := config.Write(*setConfig); err != nil {
			return err
		}

		for k, v := range *setConfig {
			fmt.Fprintf(out, "config \"%s\" set as \"%s\"\n", k, v)
		}
		return nil
	}

	return writeLog(out)
}

func writeLog(out io.Writer) error {
	args := miniflag.Args()
	argN := len(args)

	if argN > 1 {
		return ErrWriteLog(ErrExpectedOneArgument(argN))
	}

	if argN == 0 {
		input, err := core.CaptureEditorInput()
		if err != nil {
			return ErrWriteLog(err)
		}

		meta := core.Meta{Date: time.Now(), Page: *page}

		if err := core.WriteLog(out, core.NewLog(meta, string(input), *tags)); err != nil {
			return ErrWriteLog(err)
		}

		return nil
	}

	meta := core.Meta{Date: time.Now(), Page: *page}

	if err := core.WriteLog(out, core.NewLog(meta, strings.Join(args, "\n"), *tags)); err != nil {
		return ErrWriteLog(err)
	}

	return nil
}
