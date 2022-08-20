package flag

import (
	"flag"
	"fmt"
	"testing"
)

func TestValNameIsShorthand(t *testing.T) {
	tests := []struct {
		args     []string
		expected bool
	}{
		{
			args:     []string{"-b"},
			expected: true,
		},
	}

	for _, tt := range tests {
		fs := newFlagSet[any]("", flag.ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := val(fs, "b", "b", false, "Test name is shorthand flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %t, got %t", tt.expected, *actual)
			}
		})
	}
}

func TestValBool(t *testing.T) {
	tests := []struct {
		args     []string
		expected bool
	}{
		{
			expected: false,
		},
		{
			args:     []string{"-b"},
			expected: true,
		},
		{
			args:     []string{"--bool"},
			expected: true,
		},
		{
			args:     []string{"--bool=false"},
			expected: false,
		},
	}

	for _, tt := range tests {
		fs := newFlagSet[any]("", flag.ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := val(fs, "bool", "b", false, "Test bool flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %t, got %t", tt.expected, *actual)
			}
		})
	}
}

func TestValString(t *testing.T) {
	tests := []struct {
		args     []string
		expected string
	}{
		{
			expected: "",
		},
		{
			args:     []string{"-s", "string"},
			expected: "string",
		},
		{
			args:     []string{"--string", "long", "string"},
			expected: "long",
		},
	}

	for _, tt := range tests {
		fs := newFlagSet[any]("", flag.ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := val(fs, "string", "s", "", "Test string flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %s, got %s", tt.expected, *actual)
			}
		})
	}
}

func TestValInt(t *testing.T) {
	tests := []struct {
		args     []string
		expected int
	}{
		{
			expected: 0,
		},
		{
			args:     []string{"-i", "1"},
			expected: 1,
		},
		{
			args:     []string{"--int", "-1"},
			expected: -1,
		},
	}

	for _, tt := range tests {
		fs := newFlagSet[any]("", flag.ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := val(fs, "int", "i", 0, "Test int flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %d, got %d", tt.expected, *actual)
			}
		})
	}
}

func TestValInt64(t *testing.T) {
	tests := []struct {
		args     []string
		expected int64
	}{
		{
			expected: 0,
		},
		{
			args:     []string{"-i", "1"},
			expected: 1,
		},
		{
			args:     []string{"--int64", "-1"},
			expected: -1,
		},
	}

	for _, tt := range tests {
		fs := newFlagSet[any]("", flag.ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := val(fs, "int64", "i", int64(0), "Test int64 flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %d, got %d", tt.expected, *actual)
			}
		})
	}
}

func TestValUint(t *testing.T) {
	tests := []struct {
		args     []string
		expected uint
	}{
		{
			expected: 0,
		},
		{
			args:     []string{"-u", "1"},
			expected: 1,
		},
		{
			args:     []string{"--uint", "10"},
			expected: 10,
		},
	}

	for _, tt := range tests {
		fs := newFlagSet[any]("", flag.ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := val(fs, "uint", "u", uint(0), "Test uint flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %d, got %d", tt.expected, *actual)
			}
		})
	}
}

func TestValUint64(t *testing.T) {
	tests := []struct {
		args     []string
		expected uint64
	}{
		{
			expected: 0,
		},
		{
			args:     []string{"-u", "1"},
			expected: 1,
		},
		{
			args:     []string{"--uint64", "10"},
			expected: 10,
		},
	}

	for _, tt := range tests {
		fs := newFlagSet[any]("", flag.ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := val(fs, "uint64", "u", uint64(0), "Test uint64 flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %d, got %d", tt.expected, *actual)
			}
		})
	}
}

func TestValFloat64(t *testing.T) {
	tests := []struct {
		args     []string
		expected float64
	}{
		{
			args:     []string{},
			expected: 0.000000,
		},
		{
			args:     []string{"-f", "1.0"},
			expected: 1.000000,
		},
		{
			args:     []string{"--float64", "10.000001"},
			expected: 10.000001,
		},
	}

	for _, tt := range tests {
		fs := newFlagSet[any]("", flag.ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := val(fs, "float64", "f", float64(0), "Test float64 flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if tt.expected != *actual {
				t.Fatalf("flag value did not match expected %f, got %f", tt.expected, *actual)
			}
		})
	}
}

type customSliceFlag []string

func (s *customSliceFlag) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *customSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func TestValCustom(t *testing.T) {
	tests := []struct {
		args     []string
		expected customSliceFlag
	}{
		{
			args:     []string{},
			expected: customSliceFlag{},
		},
		{
			args:     []string{"-s", "A"},
			expected: customSliceFlag{"A"},
		},
		{
			args:     []string{"--slice", "A", "--slice", "B"},
			expected: customSliceFlag{"A", "B"},
		},
	}

	for _, tt := range tests {
		fs := newFlagSet[any]("", flag.ContinueOnError)
		t.Run("", func(t *testing.T) {
			actual := val(fs, "slice", "s", customSliceFlag{}, "Test customSliceFlag flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			if len(tt.expected) != len(*actual) {
				t.Fatalf("flag value did not match expected %q, got %q", tt.expected, *actual)
			}
		})
	}
}

func TestArgs(t *testing.T) {
	tests := []struct {
		args     []string
		expected []string
	}{
		{
			args:     []string{},
			expected: []string{},
		},
		{
			args:     []string{"-s", "string"},
			expected: []string{},
		},
		{
			args:     []string{"--bool", "arg0"},
			expected: []string{"arg0"},
		},
		{
			args:     []string{"arg0", "--bool", "arg1", "-s", "string", "arg2"},
			expected: []string{"arg0", "arg1", "arg2"},
		},
	}

	for _, tt := range tests {
		fs := newFlagSet[any]("", flag.ContinueOnError)
		t.Run("", func(t *testing.T) {
			// setup flags
			// boolean is a special case so that needs to be tested
			val(fs, "bool", "b", false, "bool flag")
			val(fs, "string", "s", "", "string flag")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatal(err)
			}

			actual := args(fs)

			if len(tt.expected) != len(actual) {
				t.Fatalf("flag value did not match expected %q, got %q", tt.expected, actual)
			}
		})
	}
}
