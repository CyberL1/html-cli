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

	buildCmd.Flags().StringVarP(&constants.Config.Build.Directory, "out", "o", "build", "Output directory")
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build html",
	Run: func(cmd *cobra.Command, args []string) {
		utils.LoadConfig(".")

		dirsToProcess := []string{"."}
		var currentDir string

		for len(dirsToProcess) > 0 {
			currentDir = dirsToProcess[len(dirsToProcess)-1]
			dirsToProcess = dirsToProcess[:len(dirsToProcess)-1]

			dirContents, err := os.ReadDir(currentDir)
			if err != nil {
				fmt.Println("Build failed:", err)
				return
			}

			for _, file := range dirContents {
				if file.Name() == "html-cli.json" || (file.IsDir() && file.Name() == constants.Config.Build.Directory) {
					continue
				}

				if file.IsDir() {
					dirsToProcess = append(dirsToProcess, filepath.Join(currentDir, file.Name()))
					continue
				}

				fmt.Println("Building:", filepath.Join(currentDir, file.Name()))
				fileContents, err := os.ReadFile(filepath.Join(currentDir, file.Name()))
				if err != nil {
					fmt.Println("Build failed:", err)
					return
				}

				if filepath.Ext(file.Name()) == ".html" {
					fileContents = utils.ApplyBoilerplate(fileContents, false)
					fileContents = utils.Minify(fileContents)
				}

				os.MkdirAll(filepath.Join(constants.Config.Build.Directory, currentDir), 0775)
				os.WriteFile(filepath.Join(constants.Config.Build.Directory, currentDir, file.Name()), fileContents, 0644)
			}
		}
		fmt.Println("\nBuild complete")
	},
}
