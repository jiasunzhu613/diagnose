package cmd

import (
	"github.com/jiasunzhu613/diagnose/cmd/lib"
	"github.com/spf13/cobra"
)

var (
	setupCmd = &cobra.Command{
		Use: "setup",
		Short: "setup your diagnose environment and choose your model provider",
		Long: "TODO: this is the longer explanation for setup subcommand",
		Run: func(cmd *cobra.Command, arg []string) {
			lib.SetupWorkflow()
		},
	}
)

func init() {
	rootCmd.AddCommand(setupCmd)
}