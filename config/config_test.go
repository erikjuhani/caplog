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
		input    map[string]string
		actual   config
		expected config
	}{
		{},
		{
			input: map[string]string{"not_a_correct_key": "_"},
		},
		{
			input: map[string]string{WorkspacesKey: ""},
		},
		{
			actual:   config{Workspaces: []Workspace{{Name: "test", Path: "~/test"}}},
			input:    map[string]string{CurrentWorkspaceKey: "_"},
			expected: config{Workspaces: []Workspace{{Name: "test", Path: "~/test"}}},
		},
		{
			actual:   config{Workspaces: []Workspace{{Name: "test", Path: "~/test"}}},
			input:    map[string]string{CurrentWorkspaceKey: "test"},
			expected: config{CurrentWorkspace: "test", Workspaces: []Workspace{{Name: "test", Path: "~/test"}}},
		},
		{
			input:    map[string]string{EditorKey: "vim"},
			expected: config{Editor: "vim"},
		},
		{
			input:    map[string]string{WorkspacesKey: "test:~/test"},
			expected: config{Workspaces: []Workspace{{Name: "test", Path: "~/test"}}},
		},
		{
			input: map[string]string{WorkspacesKey: "test"},
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

func TestWorkspacePath(t *testing.T) {
	homeDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(homeDir)

	tests := []struct {
		config   config
		expected string
	}{
		{
			config: config{CurrentWorkspace: "_", Workspaces: []Workspace{{Name: "test", Path: "~/test"}}},
		},
		{
			config:   config{CurrentWorkspace: "test", Workspaces: []Workspace{{Name: "test", Path: "~/test"}}},
			expected: homeDir + "/test",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			actual := workspacePath(homeDir, &tt.config)

			if tt.expected != actual {
				t.Fatalf("workspace path did not match expected %s, got %s", tt.expected, actual)
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
		input    map[string]string
		expected string
	}{
		{},
		{
			input: map[string]string{"not_a_correct_key": "_"},
		},
		{
			input:    map[string]string{EditorKey: "vim"},
			expected: "editor = 'vim'",
		},
		{
			input:    map[string]string{EditorKey: "vim", WorkspacesKey: "test:~/test"},
			expected: "workspaces = [{name = 'test', path = '~/test'}]\neditor = 'vim'",
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
				t.Fatalf("config toml did not match expected %s, got %s", tt.expected, actual)
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
			expected: config{CurrentWorkspace: "default", Editor: "vi", Workspaces: []Workspace{{Name: "default", Path: homeDir + "/.caplog/capbook"}}},
		},
		{
			content:  "editor = 'vim'",
			expected: config{CurrentWorkspace: "default", Editor: "vim", Workspaces: []Workspace{{Name: "default", Path: homeDir + "/.caplog/capbook"}}},
		},
		{
			content:  "faulty_toml'_'",
			expected: config{CurrentWorkspace: "default", Editor: "vi"},
		},
		{
			content: "editor = 'vim'\nworkspaces = [{Name = 'test', Path = '~/test'},{Name = 'test0', Path = '~/test0'}]",
			expected: config{CurrentWorkspace: "default", Editor: "vim", Workspaces: []Workspace{
				{
					Name: "test",
					Path: "~/test",
				},
				{
					Name: "test0",
					Path: "~/test0",
				},
				{
					Name: "default",
					Path: homeDir + "/.caplog/capbook",
				},
			}},
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
