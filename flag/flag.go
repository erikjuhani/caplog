package flag

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Flag struct {
	LongHand  string
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

func Args() []string {
	return commandLine.Args()
}

func Usage(name string) {
	commandLine.Usage = usageFn(name)
}

func Parse() {
	commandLine.Parse(os.Args[1:])
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

func usageFn(name string) func() {
	return func() {
		// TODO: better version of this
		var _ strings.Builder

		summary := fmt.Sprintf("usage: %s", name)

		var flagUsage string

		for _, f := range commandLine.flags {
			var compound string

			if f.Shorthand != "" {
				compound = fmt.Sprintf("-%s", f.Shorthand)
			}
			if f.LongHand != "" {
				compound = fmt.Sprintf("%s|--%s", compound, f.LongHand)
			}

			summary = fmt.Sprintf("%s [%s]", summary, compound)

			flagUsage = fmt.Sprintf("%s  %s\t%s\n", flagUsage, compound, f.Usage)
		}

		fmt.Println(summary)
		fmt.Print(flagUsage)
	}
}

func setFlagUsage(flags *[]Flag, name string, shorthand string, usage string) {
	// TODO: better version of this
	var _ strings.Builder

	var f Flag

	if len(name) == 1 {
		shorthand = name
	}

	if shorthand != "" {
		f.Shorthand = shorthand
	}

	if len(name) > 1 {
		f.LongHand = name
	}

	f.Usage = usage

	*flags = append(*flags, f)
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
