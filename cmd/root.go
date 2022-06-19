package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/erikjuhani/caplog/config"
	"github.com/erikjuhani/caplog/core"
	"github.com/spf13/cobra"
)

var (
	flagTags   []string
	flagPage   string
	flagConfig bool
	flagGetDir bool
)

var rootCmd = &cobra.Command{
	Use:   "caplog",
	Short: "A Journaling System",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagGetDir {
			fmt.Println(config.Config.Git.LocalRepository)
			return nil
		}

		if flagConfig {
			argAmount := len(args)

			if argAmount%2 > 0 {
				return fmt.Errorf("expected %d arguments, got %d", argAmount+1, argAmount)
			}

			m := make(map[config.ConfigKey]string)

			for i := 0; i < argAmount; i += 2 {
				key := args[i]
				if ok := config.Contains(key); !ok {
					return fmt.Errorf("%s is not a valid configuration key", key)
				}
				m[key] = args[i+1]
			}

			return config.Write(m)
		}

		if len(args) == 0 {
			input, err := core.CaptureEditorInput()
			if err != nil {
				return err
			}

			meta := core.Meta{Date: time.Now(), Page: flagPage}

			return core.WriteLog(core.NewLog(meta, string(input), flagTags))
		}

		meta := core.Meta{Date: time.Now(), Page: flagPage}

		return core.WriteLog(core.NewLog(meta, strings.Join(args, "\n"), flagTags))
	},
}

func Execute() {
	rootCmd.Flags().StringSliceVarP(&flagTags, "tag", "t", []string{}, "Add tags to log entry")
	rootCmd.Flags().StringVarP(&flagPage, "page", "p", "", "Save log entry to sub-directory (=page)")
	rootCmd.Flags().BoolVarP(&flagConfig, "config", "c", false, "Change config setting with key and value")
	rootCmd.Flags().BoolVarP(&flagGetDir, "get-dir", "g", false, "Returns the local repository directory")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
