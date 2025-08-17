package cmd

import (
	"fmt"
	"html-cli/constants"
	"html-cli/utils"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
	utils.LoadConfig(".")

	buildCmd.Flags().StringVarP(&constants.Config.Build.Directory, "out", "o", "build", "Output directory")
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build html",
	Run: func(cmd *cobra.Command, args []string) {
		dirContents, err := os.ReadDir(".")
		if err != nil {
			fmt.Println("Build failed:", err)
			return
		}

		for _, file := range dirContents {
			if file.Name() == "html-cli.json" || (file.IsDir() && file.Name() == constants.Config.Build.Directory) {
				continue
			}

			if file.IsDir() {
				cmd.Run(cmd, []string{file.Name()})
			} else {
				fmt.Println("Building:", file.Name())

				fileContents, err := os.ReadFile(file.Name())
				if err != nil {
					fmt.Println("Build failed:", err)
					return
				}

				if filepath.Ext(file.Name()) == ".html" {
					fileContents = utils.ApplyBoilerplate(fileContents, false)
					fileContents = utils.Minify(fileContents)
				}

				os.MkdirAll(constants.Config.Build.Directory, 0775)
				os.WriteFile(filepath.Join(constants.Config.Build.Directory, file.Name()), fileContents, 0644)
			}
		}
		fmt.Println("\nBuild complete")
	},
}
