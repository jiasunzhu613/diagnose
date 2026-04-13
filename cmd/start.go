package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	startCmd = &cobra.Command{
		Use: "start",
		Short: "start a diagnose session",
		Long: "used to start your diagnose session", 
		Run: func(cmd *cobra.Command, arg []string) {
			fmt.Println("Started diagnose sessions")
		},
	}
)

func init() {
	rootCmd.AddCommand(startCmd)
}

