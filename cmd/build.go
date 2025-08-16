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
		if err := build("."); err != nil {
			fmt.Println("Build failed:", err)
			return
		}

		fmt.Println("\nBuild complete")
	},
}

func build(directory string) error {
	dirContents, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	for _, file := range dirContents {
		if file.Name() == "html-cli.json" || (file.IsDir() && file.Name() == constants.Config.Build.Directory) {
			continue
		}

		if file.IsDir() {
			build(filepath.Join(directory, file.Name()))
		} else {
			fmt.Println("Building:", filepath.Join(directory, file.Name()))

			fileContents, err := os.ReadFile(filepath.Join(directory, file.Name()))
			if err != nil {
				return err
			}

			if filepath.Ext(file.Name()) == ".html" {
				fileContents = utils.ApplyBoilerplate(fileContents, false)
				fileContents = utils.Minify(fileContents)
			}

			os.MkdirAll(filepath.Join(constants.Config.Build.Directory, directory), 0775)
			os.WriteFile(filepath.Join(constants.Config.Build.Directory, directory, file.Name()), fileContents, 0644)
		}
	}
	return nil
}
