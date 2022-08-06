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
	getDir    = flag.Val("getdir", "g", false, "Returns the local repository directory")
	page      = flag.Val("page", "p", "", "Saves log entry to <sub-directory>/<page>")
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

	// Check that passed configuration keys are valid
	if ok := config.Contains(k); !ok {
		return fmt.Errorf("%s is not a valid configuration key", k)
	}

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

	if *getDir {
		// Return current repository path
		fmt.Println(config.Config.Git.LocalRepository)
		return nil
	}

	// Config flags were used needs to do configuration change
	if len(*setConfig) > 0 {
		return setConfigValue()
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

func oldSetConfigValue(args []string) error {
	argAmount := len(args)

	if argAmount == 0 {
		return fmt.Errorf("expected 2 arguments, got %d", argAmount)
	}

	if argAmount%2 > 0 {
		return fmt.Errorf("expected %d arguments, got %d", argAmount+1, argAmount)
	}

	m := make(map[config.ConfigKey]string)

	for i := 0; i < argAmount; i += 2 {
		key := args[i]
		if ok := config.Contains(key); !ok {
			return fmt.Errorf("%s is not a valid configuration key", key)
		}
		m[key] = args[i+1]
	}

	return config.Write(m)
}

func setConfigValue() error {
	return config.Write(*setConfig)
}
