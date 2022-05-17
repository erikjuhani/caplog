package config

import (
	"fmt"
	"os"

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

func setDefaults(localPath string) {
	viper.SetDefault(GitLocalRepositoryKey, fmt.Sprintf("%s/%s", localPath, defaultLogDir))
	viper.SetDefault(EditorKey, "vi")
}

func Get(key ConfigKey) string {
	return viper.GetString(key)
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
		return viper.ReadInConfig()
	}

	return err
}
