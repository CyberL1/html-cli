package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "html-cli",
	Short: "A cli for html",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	}}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
