package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

const (
	GitLocalRepositoryKey ConfigKey = "git.local_repository"
	EditorKey             ConfigKey = "editor"
)

const (
	configType = "toml"

	defaultDir    = ".caplog"
	defaultLogDir = "capbook"
)

var (
	Config     = config{}
	configFile []byte
	configPath string
)

type ConfigKey = string

type gitConfig struct {
	LocalRepository string `toml:"local_repository,omitempty"`
}

type config struct {
	Git    *gitConfig `toml:"git,omitempty"`
	Editor string     `toml:"editor,omitempty"`
}

func Load() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	localPath := fmt.Sprintf("%s/%s", homeDir, defaultDir)

	configPath, err = readConfigPath(localPath)
	if err != nil {
		return err
	}

	// Store actual config file to memory
	configFile, err = os.ReadFile(configPath)
	if err != nil {
		return err
	}

	setDefaults(&Config, localPath)

	if err := toml.Unmarshal(configFile, &Config); err != nil {
		return err
	}

	// Replace tilde to point to user home directory
	Config.Git.LocalRepository = strings.Replace(Config.Git.LocalRepository, "~", homeDir, 1)

	return nil
}

func Write(c map[ConfigKey]string) error {
	var localConfig config
	if err := toml.Unmarshal(configFile, &localConfig); err != nil {
		return err
	}

	if err := mergeMapToConfig(c, &localConfig); err != nil {
		return err
	}

	config, err := toml.Marshal(&localConfig)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, config, 0644)
}

func mergeMapToConfig(c map[ConfigKey]string, config *config) error {
	for k, v := range c {
		switch k {
		case GitLocalRepositoryKey:
			config.Git.LocalRepository = v
		case EditorKey:
			config.Editor = v
		default:
			return fmt.Errorf("\"%s\" is not a valid configuration key", k)
		}
	}
	return nil
}

func setDefaults(c *config, localPath string) {
	c.Git = &gitConfig{}
	c.Git.LocalRepository = fmt.Sprintf("%s/%s", localPath, defaultLogDir)
	c.Editor = "vi"
}

// TODO: Add package github.com/adrg/xdg for better configuration support
// TODO: Search for existing config file in all possible locations
// on first hit finish

// Viable config paths
// .config/caplog.toml
// ~/.caplog.toml

func readConfigPath(localPath string) (string, error) {
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		if err := os.Mkdir(localPath, os.ModePerm); err != nil {
			return "", fmt.Errorf("cannot create .caplog folder under home directory %w", err)
		}
	}

	return fmt.Sprintf("%s/%s", localPath, "config"), nil

}
