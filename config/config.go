package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// Valid configuration keys
const (
	CurrentWorkspaceKey = "current_workspace"
	WorkspacesKey       = "workspaces"
	EditorKey           = "editor"
)

// Default path location constants
const (
	defaultConfigLocation = "~/.caplog.toml"
	defaultRepositoryPath = "%s/.caplog/capbook"
)

var (
	// Users home directory
	HomeDir string

	// Config is exported to give access to existing configuration values
	Config = config{}

	// In-memory copy of the existing config file
	configFile []byte

	// The actual path to the nearest config file
	configPath string

	// Paths that are searched for configuration file
	validConfigPaths = []string{
		defaultConfigLocation,
		"~/.config/caplog.toml",
		"~/.config/caplog/caplog.toml",
	}
)

type Workspace struct {
	Name string `toml:"name"`
	Path string `toml:"path"`
}

type Workspaces []Workspace

func (w *Workspaces) Append(name string, path string) {
	*w = append(*w, Workspace{Name: name, Path: path})
}

func (w Workspaces) Has(workspace string) bool {
	for _, v := range w {
		if v.Name == workspace {
			return true
		}
	}

	return false
}

func (w Workspaces) Names() []string {
	var n []string
	for _, v := range w {
		n = append(n, v.Name)
	}

	return n
}

func WorkspacePath() string {
	return workspacePath(HomeDir, &Config)
}

func workspacePath(homeDir string, config *config) string {
	for _, v := range config.Workspaces {
		if v.Name == config.CurrentWorkspace {
			return replaceTilde(v.Path, homeDir)
		}
	}

	return ""
}

type config struct {
	CurrentWorkspace string     `toml:"current_workspace,omitempty"`
	Workspaces       Workspaces `toml:"workspaces,inline,omitempty"`
	Editor           string     `toml:"editor,omitempty"`
}

// Load initializes configuration to memory either with default values
// or read from a `caplog.toml` configuration file if such exists.
func Load() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	return load(homeDir, &Config)
}

// Write writes given map structure to a configuration file
// located in configuration path
func Write(c map[string]string) error {
	return writeTo(configPath, c)
}

func findExistingConfigFile(homeDir string) string {
	for _, v := range validConfigPaths {
		vv := replaceTilde(v, homeDir)
		if _, err := os.Stat(vv); os.IsNotExist(err) {
			continue
		}
		return vv
	}
	return replaceTilde(defaultConfigLocation, homeDir)
}

func load(homeDir string, config *config) error {
	HomeDir = homeDir

	defaultPath := fmt.Sprintf(defaultRepositoryPath, homeDir)

	// Set defaults

	config.CurrentWorkspace = "default"
	config.Editor = "vi"

	configPath = findExistingConfigFile(homeDir)

	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		// Store actual config file to memory if it exists
		// This is used when merging configurations together
		// so old configurations are not lost
		configFile, err = os.ReadFile(configPath)
		if err != nil {
			return err
		}
	}

	if err := toml.Unmarshal(configFile, config); err != nil {
		return err
	}

	// TODO: make this better, we need to append default workspace here
	// as that one needs to be always available
	config.Workspaces.Append("default", defaultPath)

	if exists := config.Workspaces.Has(config.CurrentWorkspace); !exists {
		return fmt.Errorf("given \"%s\" workspace is not a valid workspace\nvalid workspaces are: %v", config.CurrentWorkspace, config.Workspaces.Names())
	}

	return nil
}

func mergeMapToConfig(c map[string]string, config *config) error {
	for k, v := range c {
		switch k {
		case WorkspacesKey:
			var ws []Workspace
			wss := strings.Split(v, ",")

			// Accept no value to remove all workspace definitions
			if len(wss) == 1 && wss[0] == "" {
				config.Workspaces = ws
				return nil
			}

			for _, raw := range wss {
				rs := strings.SplitN(raw, ":", 2)
				if len(rs) < 2 {
					return fmt.Errorf("not enough arguments to set workspaces, set the value with double colon separator \"workspace:path\"")
				}
				ws = append(ws, Workspace{Name: rs[0], Path: rs[1]})
			}

			config.Workspaces = ws
		case CurrentWorkspaceKey:
			if exists := config.Workspaces.Has(v); !exists && v != "default" {
				return fmt.Errorf("given \"%s\" workspace is not a valid workspace\nvalid workspaces are: %v", v, config.Workspaces.Names())
			}
			config.CurrentWorkspace = v
		case EditorKey:
			config.Editor = v
		default:
			return fmt.Errorf("\"%s\" is not a valid configuration key", k)
		}
	}
	return nil
}

func replaceTilde(s, r string) string {
	return strings.Replace(s, "~", r, 1)
}

func writeTo(configPath string, c map[string]string) error {
	var localConfig config
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		if err := toml.Unmarshal(configFile, &localConfig); err != nil {
			return err
		}
	}

	if err := mergeMapToConfig(c, &localConfig); err != nil {
		return err
	}

	config, err := toml.Marshal(&localConfig)
	if err != nil {
		return err
	}

	trimmed := strings.Trim(string(config), "\n")

	return os.WriteFile(configPath, []byte(trimmed), 0644)
}
