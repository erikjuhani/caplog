package config

import (
	"os"
	"reflect"
	"testing"
)

func TestFindExistingConfigFile(t *testing.T) {
	homeDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(homeDir)
	tests := []struct {
		dir      string
		filename string
		expected string
	}{
		{
			dir:      "",
			filename: "",
			expected: homeDir + "/.caplog.toml",
		},
		{
			dir:      "",
			filename: ".caplog.toml",
			expected: homeDir + "/.caplog.toml",
		},
		{
			dir:      "/.config",
			filename: "caplog.toml",
			expected: homeDir + "/.config/caplog.toml",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if tt.dir != "" {
				if err := os.MkdirAll(homeDir+tt.dir, 0644); err != nil {
					t.Fatal(err)
				}
			}

			if tt.filename != "" {
				if err := os.WriteFile(homeDir+tt.dir+tt.filename, []byte{}, 0644); err != nil {
					t.Fatal(err)
				}
			}

			actual := findExistingConfigFile(homeDir)

			if tt.expected != actual {
				t.Fatalf("found config path did not match expected %s, got %s", tt.expected, actual)
			}
		})
	}
}

func TestMergeMapToConfig(t *testing.T) {
	tests := []struct {
		input    map[ConfigKey]string
		actual   config
		expected config
	}{
		{},
		{
			input: map[ConfigKey]string{"not_a_correct_key": "_"},
		},
		{
			actual:   config{Editor: "vim"},
			expected: config{Editor: "vim"},
		},
		{
			input:    map[ConfigKey]string{EditorKey: "vim"},
			expected: config{Editor: "vim"},
		},
		{
			input:    map[ConfigKey]string{EditorKey: "vim", GitLocalRepositoryKey: "~/test"},
			expected: config{Editor: "vim", Git: &gitConfig{LocalRepository: "~/test"}},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			mergeMapToConfig(tt.input, &tt.actual)

			if !reflect.DeepEqual(tt.expected, tt.actual) {
				t.Fatalf("config did not match expected %+v, got %+v", tt.expected, tt.actual)
			}
		})
	}
}

func TestReplaceTilde(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "~",
			expected: "/home",
		},
		{
			input:    "~/.config",
			expected: "/home/.config",
		},
		{
			input:    "~/.caplog.toml",
			expected: "/home/.caplog.toml",
		},
		{
			input:    "~/~/test",
			expected: "/home/~/test",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			actual := replaceTilde(tt.input, "/home")

			if tt.expected != actual {
				t.Fatalf("replaced string did not match expected %s, got %s", tt.expected, actual)
			}
		})
	}
}

func TestWriteTo(t *testing.T) {
	homeDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(homeDir)
	configPath := homeDir + "/.caplog.toml"
	tests := []struct {
		input    map[ConfigKey]string
		expected string
	}{
		{},
		{
			input: map[ConfigKey]string{"not_a_correct_key": "_"},
		},
		{
			input:    map[ConfigKey]string{EditorKey: "vim"},
			expected: "editor = 'vim'",
		},
		{
			input:    map[ConfigKey]string{EditorKey: "vim", GitLocalRepositoryKey: "~/test"},
			expected: "editor = 'vim'\n[git]\nlocal_repository = '~/test'\n",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			writeTo(configPath, tt.input)
			actual, err := os.ReadFile(configPath)
			if err != nil {
				t.Fatal(err)
			}

			if tt.expected != string(actual) {
				t.Fatalf("found config path did not match expected %s, got %s", tt.expected, actual)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	homeDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(homeDir)
	tests := []struct {
		content  string
		expected config
	}{
		{
			expected: config{Editor: "vi", Git: &gitConfig{LocalRepository: homeDir + "/.caplog/capbook"}},
		},
		{
			content:  "editor = 'vim'",
			expected: config{Editor: "vim", Git: &gitConfig{LocalRepository: homeDir + "/.caplog/capbook"}},
		},
		{
			content:  "faulty_toml'_'",
			expected: config{Editor: "vi", Git: &gitConfig{LocalRepository: homeDir + "/.caplog/capbook"}},
		},
		{
			content:  "editor = 'vim'\n[git]\nlocal_repository = '~/test'\n",
			expected: config{Editor: "vim", Git: &gitConfig{LocalRepository: "~/test"}},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if tt.content != "" {
				if err := os.WriteFile(homeDir+"/.caplog.toml", []byte(tt.content), 0644); err != nil {
					t.Fatal(err)
				}
			}

			actual := config{}
			load(homeDir, &actual)

			if !reflect.DeepEqual(tt.expected, actual) {
				t.Fatalf("config did not match expected %+v, got %+v", tt.expected, actual)
			}
		})
	}
}
