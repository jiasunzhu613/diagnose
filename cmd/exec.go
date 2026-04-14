// TEMP: placeholder, may not be needed, could be used as alternative invocation of exec single command

package cmd

import (
	"github.com/jiasunzhu613/diagnose/cmd/lib"
	"github.com/spf13/cobra"
)

var (
	execCmd = &cobra.Command{
		Use: "exec",
		Short: "execute single command with diagnose",
		Long: "TODO: this is the longer explanation for exec subcommand",
		Run: func(cmd *cobra.Command, arg []string) {
			lib.ExecWorkflow(arg)
		},
	}
)

func init() {
	rootCmd.AddCommand(execCmd)
}