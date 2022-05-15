package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type ConfigKey = string

const (
	// TODO: Add remote repository support for cloud storage
	// GitRemoteRepositoryKey ConfigKey = "git.remote_repository"

	// TODO: Not used yet
	GitLocalRepositoryKey ConfigKey = "git.local_repository"
	EditorKey             ConfigKey = "editor"
)

const (
	configPath = ".caplog/config"
	configType = "toml"
)

var (
	LocalRepositoryPath string
)

func setDefaults() {
	// TODO: Not used yet
	viper.SetDefault(GitLocalRepositoryKey, "")
	viper.SetDefault(EditorKey, "vi")
}

func loadConfig() error {
	// TODO: Ensure caplog directory exists
	err := viper.SafeWriteConfig()

	if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
		// Set default values for configuration options
		setDefaults()

		return viper.ReadInConfig()
	}

	return err
}

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("cannot find user home directory %s", err)
	}

	LocalRepositoryPath = fmt.Sprintf("%s/%s", homeDir, ".caplog")

	if err := os.MkdirAll(LocalRepositoryPath, os.ModePerm); err != nil {
		log.Fatalf("cannot create local repository %s", err)
	}

	viper.SetConfigType(configType)
	viper.SetConfigFile(LocalRepositoryPath)
}
