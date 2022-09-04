package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type ConfigKey = string

// Valid configuration keys
const (
	GitLocalRepositoryKey ConfigKey = "git.local_repository"
	EditorKey             ConfigKey = "editor"
)

// Default path location constants
const (
	defaultConfigLocation = "~/.caplog.toml"
	defaultRepositoryPath = "%s/.caplog/capbook"
)

// Config is exported to give access to existing configuration values
var Config = config{}

var (
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

type gitConfig struct {
	LocalRepository string `toml:"local_repository,omitempty"`
}

type config struct {
	Git    *gitConfig `toml:"git,omitempty"`
	Editor string     `toml:"editor,omitempty"`
}

// Load initializes configuration to memory either with default values
// or read from a `caplog.toml` configuration file if such exists.
func Load(homeDir string) error {
	return load(homeDir, &Config)
}

// Write writes given map structure to a configuration file
// located in configuration path
func Write(c map[ConfigKey]string) error {
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
	setDefaults(fmt.Sprintf(defaultRepositoryPath, homeDir), config)

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

	return nil
}

func mergeMapToConfig(c map[ConfigKey]string, config *config) error {
	for k, v := range c {
		switch k {
		case GitLocalRepositoryKey:
			config.Git = &gitConfig{LocalRepository: v}
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

func setDefaults(rpath string, c *config) {
	c.Git = &gitConfig{}
	c.Git.LocalRepository = rpath
	c.Editor = "vi"
}

func writeTo(configPath string, c map[ConfigKey]string) error {
	var localConfig config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
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

	return os.WriteFile(configPath, config[:len(config)-1], 0644)
}
