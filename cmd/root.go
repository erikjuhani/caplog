package cmd

import (
	"os"
	"time"

	"github.com/erikjuhani/caplog/core"
	"github.com/spf13/cobra"
)

var tags []string

var rootCmd = &cobra.Command{
	Use:   "caplog",
	Short: "A Journaling System",
	Args:  cobra.MaximumNArgs(1),
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			input, err := core.CaptureEditorInput()
			if err != nil {
				return err
			}
			return core.WriteLog(core.CreateLog(time.Now(), string(input), tags))
		}
		return core.WriteLog(core.CreateLog(time.Now(), args[0], tags))
	},
}

func Execute() {
	rootCmd.Flags().StringSliceVarP(&tags, "tag", "t", []string{}, "Add tags to log entry")
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
