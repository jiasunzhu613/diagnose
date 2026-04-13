package cmd

import (
	"fmt"
	"log"

	"github.com/jiasunzhu613/diagnose/cmd/lib"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "diagnose",
	Short: "A simple to use CLI doctor to resolve all your errors",
	Long: "This is a much longer description for diagnose doctor!!",
	Run: func(cmd *cobra.Command, arg []string) { // for some reason this can be called through "--"
		if len(arg) == 0 {
			log.Fatal("not enough arguments to root command")
		}
	
		lib.ExecWorkflow(arg)
		if lib.StderrBuf.Len() > 0 {
			fmt.Println("Found some errors and probably should run LLM workflow")
		}

		// This should not run because exec overrides the entire program
		fmt.Println(arg)
		fmt.Println("Ran root command")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	fmt.Println("Init for root.go")
}