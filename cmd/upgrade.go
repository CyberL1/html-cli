package cmd

import (
	"archive/zip"
	"fmt"
	"html-cli/utils"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrades the CLI",
	Run: func(cmd *cobra.Command, args []string) {
		latestRelease, err := utils.GetLatestRelease()
		if err != nil {
			fmt.Println("Failed to get latest release:", err)
			return
		}

		var assetDownloadUrl string
		for _, asset := range latestRelease.Assets {
			os := strings.Split(asset.Name, "-")[2]
			arch := strings.Split(strings.TrimSuffix(asset.Name, ".zip"), "-")[3]

			if os == runtime.GOOS && arch == runtime.GOARCH {
				assetDownloadUrl = asset.BrowserDownloadUrl
			}
		}

		if assetDownloadUrl == "" {
			fmt.Println("No matching binary found")
			return
		}

		resp, err := http.Get(assetDownloadUrl)
		if err != nil {
			fmt.Println("Failed to get latest release", err)
			return
		}
		defer resp.Body.Close()

		tmpFile, err := os.CreateTemp("", "html-cli-*.zip")
		if err != nil {
			fmt.Println("Upgrade failed:", err)
			return
		}
		defer tmpFile.Close()
		defer os.Remove(tmpFile.Name())

		io.Copy(tmpFile, resp.Body)

		r, err := zip.OpenReader(tmpFile.Name())
		if err != nil {
			fmt.Println("Upgrade failed:", err)
			return
		}
		defer r.Close()

		for _, f := range r.File {
			rc, err := f.Open()
			if err != nil {
				fmt.Println("Upgrade failed:", err)
				return
			}
			defer rc.Close()

			execPath, _ := os.Executable()

			// Replace binary
			out, err := os.OpenFile(execPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if err != nil {
				fmt.Println("Upgrade failed:", err)
				return
			}
			defer out.Close()

			_, err = io.Copy(out, rc)
			if err != nil {
				fmt.Println("Upgrade failed:", err)
				return
			}
		}
	},
}
