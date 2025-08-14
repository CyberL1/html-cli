package cmd

import (
	"fmt"
	"html-cli/constants"
	"html-cli/utils"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays your cli version",
	Run: func(cmd *cobra.Command, args []string) {
		remoteVersion, err := utils.GetLatestRelease()
		if err != nil {
			fmt.Println("Failed to get latest release:", err)
			return
		}

		fmt.Printf("Your version: %s\nRemote version: %s\n", constants.Version, remoteVersion.TagName)

		if remoteVersion.TagName > constants.Version {
			fmt.Println("\nNew version available, do 'html-cli upgrade' to download")
		}
	},
}
