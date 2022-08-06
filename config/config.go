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

func Contains(key ConfigKey) bool {
	keys := map[ConfigKey]struct{}{GitLocalRepositoryKey: {}, EditorKey: {}}
	_, exists := keys[key]
	return exists
}

func Write(c map[ConfigKey]string) error {
	for k, v := range c {
		viper.Set(k, v)
	}
	return viper.WriteConfig()
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

// Viable config paths
// .config/caplog.toml
// ~/.caplog.toml

// Reading config toml
// 1. initialize default configuration values
// 2. search for existing config file
// 3. if found use replace default values with existing user configurations
// 4. unmarshal to config struct from toml

// Writing config toml
// 1. check that the key is a correct configuration key
// 2. check if configuration or existing configuration already exists
// 3. read to struct
// 4. write filled keys to configuration file

// TODO: Add package github.com/adrg/xdg for better configuration support

type C struct {
	Git struct {
		LocalRepository string
	}
	Editor string
}
