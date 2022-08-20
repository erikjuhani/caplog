package flag

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Flag struct {
	Longhand  string
	Shorthand string
	Usage     string
}

type FlagSet[T any] struct {
	*flag.FlagSet
	flags []Flag
}

var commandLine = newFlagSet[any](os.Args[0], flag.ExitOnError)

func newFlagSet[T any](name string, errorHandling flag.ErrorHandling) *FlagSet[T] {
	return &FlagSet[T]{flag.NewFlagSet(name, errorHandling), []Flag{}}
}

func args(fs *FlagSet[any]) []string {
	args := fs.Args()

	pArgs := []string{}
	for i, arg := range args {
		if arg[0] == '-' {
			continue
		}
		if i > 0 && args[i-1][0] == '-' {
			f := fs.Lookup(strings.ReplaceAll(args[i-1], "-", ""))

			if f != nil && reflect.TypeOf(f.Value).Elem().Kind() != reflect.Bool {
				continue
			}
		}
		pArgs = append(pArgs, arg)
	}

	return pArgs
}

func Args() []string {
	return args(commandLine)
}

func Usage(name string) {
	commandLine.Usage = usageFn(commandLine, name)
}

func Parse() error {
	return commandLine.Parse(os.Args[1:])
}

func val[T any](fs *FlagSet[any], name string, shorthand string, value T, usage string) *T {
	setFlagUsage(&fs.flags, name, shorthand, usage)
	if name == shorthand {
		shorthand = ""
	}
	switch v := any(value).(type) {
	case bool:
		return any((boolVar(fs, name, shorthand, v, usage))).(*T)
	case string:
		return any(stringVar(fs, name, shorthand, v, usage)).(*T)
	case int:
		return any(intVar(fs, name, shorthand, v, usage)).(*T)
	case int64:
		return any(int64Var(fs, name, shorthand, v, usage)).(*T)
	case uint:
		return any(uintVar(fs, name, shorthand, v, usage)).(*T)
	case uint64:
		return any(uint64Var(fs, name, shorthand, v, usage)).(*T)
	case float64:
		return any(float64Var(fs, name, shorthand, v, usage)).(*T)
	case T:
		return valueVar(fs, name, shorthand, v, usage)
	}
	return &value
}

func Val[T any](name string, shorthand string, value T, usage string) *T {
	return val(commandLine, name, shorthand, value, usage)
}

func usageFn(fs *FlagSet[any], name string) func() {
	var s, u strings.Builder

	s.WriteString("usage: " + name)

	p := s.Len()

	// append help flag?
	for i, f := range fs.flags {
		var c strings.Builder

		if f.Shorthand == "" && f.Longhand == "" {
			continue
		}

		if f.Shorthand != "" {
			fmt.Fprintf(&c, "-%s", f.Shorthand)
		}

		if f.Longhand != "" {
			if f.Shorthand != "" {
				c.WriteRune(' ')
			}
			fmt.Fprintf(&c, "--%s", f.Longhand)
		}

		compound := c.String()

		fmt.Fprintf(&s, " [%s]", compound)

		if (i+1)%4 == 0 {
			fmt.Fprintf(&s, "\n%*s", p, "")
		}

		fmt.Fprintf(
			&u,
			"%*s%*s\n",
			len(compound)+4,
			compound,
			len(f.Usage)-len(compound)+16,
			f.Usage,
		)
	}

	return func() {
		fmt.Fprint(fs.Output(), s.String(), "\n", u.String())
	}
}

func setFlagUsage(flags *[]Flag, name string, shorthand string, usage string) {
	if len(name) == 1 {
		shorthand = name
		name = ""
	}

	*flags = append(*flags, Flag{Longhand: name, Shorthand: shorthand, Usage: usage})
}

func boolVar(fs *FlagSet[any], name string, shorthand string, value bool, usage string) interface{} {
	fs.BoolVar(&value, name, value, usage)
	if shorthand != "" {
		fs.BoolVar(&value, shorthand, value, usage)
	}
	return &value
}

func stringVar(fs *FlagSet[any], name string, shorthand string, value string, usage string) *string {
	fs.StringVar(&value, name, value, usage)
	if shorthand != "" {
		fs.StringVar(&value, shorthand, value, usage)
	}
	return &value
}

func intVar(fs *FlagSet[any], name string, shorthand string, value int, usage string) *int {
	fs.IntVar(&value, name, value, usage)
	if shorthand != "" {
		fs.IntVar(&value, shorthand, value, usage)
	}
	return &value
}

func int64Var(fs *FlagSet[any], name string, shorthand string, value int64, usage string) *int64 {
	fs.Int64Var(&value, name, value, usage)
	if shorthand != "" {
		fs.Int64Var(&value, shorthand, value, usage)
	}
	return &value
}

func uintVar(fs *FlagSet[any], name string, shorthand string, value uint, usage string) *uint {
	fs.UintVar(&value, name, value, usage)
	if shorthand != "" {
		fs.UintVar(&value, shorthand, value, usage)
	}
	return &value
}

func uint64Var(fs *FlagSet[any], name string, shorthand string, value uint64, usage string) *uint64 {
	fs.Uint64Var(&value, name, value, usage)
	if shorthand != "" {
		fs.Uint64Var(&value, shorthand, value, usage)
	}
	return &value
}

func float64Var(fs *FlagSet[any], name string, shorthand string, value float64, usage string) *float64 {
	fs.Float64Var(&value, name, value, usage)
	if shorthand != "" {
		fs.Float64Var(&value, shorthand, value, usage)
	}
	return &value
}

func valueVar[T any](fs *FlagSet[any], name string, shorthand string, value T, usage string) *T {
	fs.Var(any(&value).(flag.Value), name, usage)
	if shorthand != "" {
		fs.Var(any(&value).(flag.Value), shorthand, usage)
	}
	return &value
}
