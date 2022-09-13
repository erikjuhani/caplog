package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/erikjuhani/caplog/config"
	"github.com/erikjuhani/caplog/core"
	"github.com/erikjuhani/caplog/flag"
)

var (
	dir       = flag.Val("getdir", "g", false, "Returns the local repository directory")
	page      = flag.Val("page", "p", "", "Saves log entry to <sub-directory>/<page>")
	workspace = flag.Val("workspace", "w", "", "Changes workspace to given <workspace> if it exists")
	tags      = flag.Val("tag", "t", TagsFlag{}, "Adds `<tag>` to log entry")
	setConfig = flag.Val("config", "c", ConfigFlag{}, "Changes config setting with `<key=value>`")
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
		return fmt.Errorf("key needs a value (ex. %s=<value>)", k)
	}

	v := pair[1]

	(*c)[k] = v

	return nil
}

func Run() error {
	flag.Usage("caplog")
	flag.Parse()

	if *dir {
		// Return current repository path
		fmt.Println(config.Config.CurrentWorkspace)
		return nil
	}

	if *workspace != "" {
		fmt.Printf("%+v", config.Config.Workspaces)
		if exists := config.Config.Workspaces.Has(*workspace); !exists {
			return fmt.Errorf("given \"%s\" workspace is not a valid workspace\nvalid workspaces are: %v", *workspace, config.Config.Workspaces.Names())
		}

		return config.Write(map[string]string{config.CurrentWorkspaceKey: *workspace})
	}

	// Config flags were used needs to do configuration change
	if len(*setConfig) > 0 {
		return config.Write(*setConfig)
	}

	return writeLog()
}

func writeLog() error {
	args := flag.Args()
	argN := len(args)

	if argN > 1 {
		return fmt.Errorf("expected 1 argument, got %d", argN)
	}

	if argN == 0 {
		input, err := core.CaptureEditorInput()
		if err != nil {
			return err
		}

		meta := core.Meta{Date: time.Now(), Page: *page}

		return core.WriteLog(core.NewLog(meta, string(input), *tags))
	}

	meta := core.Meta{Date: time.Now(), Page: *page}

	return core.WriteLog(core.NewLog(meta, strings.Join(args, "\n"), *tags))

}
