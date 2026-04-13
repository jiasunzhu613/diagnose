package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	endCmd = &cobra.Command{
		Use: "end",
		Short: "end a diagnose session",
		Long: "used to end current diagnose session",
		Run: func(cmd *cobra.Command, arg []string) {
			fmt.Println("Ended diagnose sessions")
		},
	}
)

func init() {
	rootCmd.AddCommand(endCmd)
}

