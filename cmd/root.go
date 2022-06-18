package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/erikjuhani/caplog/config"
	"github.com/erikjuhani/caplog/core"
	"github.com/spf13/cobra"
)

var (
	tags   []string
	page   string
	getDir bool
)

var rootCmd = &cobra.Command{
	Use:   "caplog",
	Short: "A Journaling System",
	Args:  cobra.MaximumNArgs(1),
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if getDir {
			fmt.Println(config.Config.Git.LocalRepository)
			return nil
		}

		if len(args) == 0 {
			input, err := core.CaptureEditorInput()
			if err != nil {
				return err
			}

			meta := core.Meta{Date: time.Now(), Page: page}

			return core.WriteLog(core.NewLog(meta, string(input), tags))
		}

		meta := core.Meta{Date: time.Now(), Page: page}

		return core.WriteLog(core.NewLog(meta, args[0], tags))
	},
}

func Execute() {
	rootCmd.Flags().StringSliceVarP(&tags, "tag", "t", []string{}, "Add tags to log entry")
	rootCmd.Flags().StringVarP(&page, "page", "p", "", "Save log entry to sub-directory (=page)")
	rootCmd.Flags().BoolVarP(&getDir, "get-dir", "g", false, "Returns the local repository directory")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
