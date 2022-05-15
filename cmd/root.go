package cmd

import (
	"os"
	"time"

	"github.com/erikjuhani/caplog/core"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "caplog",
	Short: "A Journaling System",
	Args:  cobra.MaximumNArgs(1),
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			capture, err := core.CaptureEditorInput()
			if err != nil {
				return err
			}
			return core.WriteLog(core.CreateLog(time.Now(), string(capture)))
		}
		return core.WriteLog(core.CreateLog(time.Now(), args[0]))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
