package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type ConfigKey = string

const (
	GitLocalRepositoryKey ConfigKey = "git.local_repository"
	EditorKey             ConfigKey = "editor"
)

const (
	configType = "toml"

	defaultDir    = ".caplog"
	defaultLogDir = "capbook"
)

type config struct {
	Git struct {
		LocalRepository string `mapstructure:"local_repository"`
	}
	Editor string
}

var Config = config{}

func setDefaults(localPath string) {
	viper.SetDefault(GitLocalRepositoryKey, fmt.Sprintf("%s/%s", localPath, defaultLogDir))
	viper.SetDefault(EditorKey, "vi")
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot find user home directory %w", err)
	}

	localPath := fmt.Sprintf("%s/%s", homeDir, defaultDir)

	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		if err := os.Mkdir(localPath, os.ModePerm); err != nil {
			return "", fmt.Errorf("cannot create .caplog folder under home directory %w", err)
		}
	}

	return localPath, nil
}

func Load() error {
	path, err := getConfigPath()

	if err != nil {
		return err
	}

	configPath := fmt.Sprintf("%s/%s", path, "config")

	viper.SetConfigType(configType)
	viper.SetConfigFile(configPath)

	err = viper.SafeWriteConfigAs(configPath)

	// Set default values for configuration options
	setDefaults(path)

	if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		if err := viper.Unmarshal(&Config); err != nil {
			return err
		}

		// TODO: handle config decoding better
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cannot find user home directory %w", err)
		}

		// Replace tilde to point to user home directory
		Config.Git.LocalRepository = strings.Replace(Config.Git.LocalRepository, "~", homeDir, 1)

		return nil
	}

	return err
}
