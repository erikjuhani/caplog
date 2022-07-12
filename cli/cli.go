package cli

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/erikjuhani/caplog/config"
	"github.com/erikjuhani/caplog/core"
)

var (
	setConfig = flag.Bool("c", false, "Change config setting with key and value")
	getDir    = flag.Bool("g", false, "Returns the local repository directory")
	page      = flag.String("p", "", "Save log entry to sub-directory (=page)")
	tags      = flagCustom("t", &StringsFlag{}, "Add tags to log entry")
)

type StringsFlag []string

func (s *StringsFlag) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *StringsFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func flagCustom[T flag.Value](name string, value T, usage string) T {
	flag.Var(value, name, usage)
	return value
}

func Run() error {
	flag.Parse()

	if *getDir {
		// Return current repository path
		fmt.Println(config.Config.Git.LocalRepository)
		return nil
	}

	if *setConfig {
		return setConfigValue(flag.Args())
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

func setConfigValue(args []string) error {
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
