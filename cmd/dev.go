package cmd

import (
	"fmt"
	"html-cli/constants"
	"html-cli/utils"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(devCmd)

	devCmd.Flags().StringVarP(&constants.Config.Dev.Root, "root", "r", ".", "Directory to watch")
	devCmd.Flags().Uint16VarP(&constants.Config.Dev.Port, "port", "p", 8080, "Port to run dev server on")
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Run a dev server for html",
	Run: func(cmd *cobra.Command, args []string) {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Println("Error setting up file watcher:", err)
			return
		}
		defer watcher.Close()

		err = watcher.Add(constants.Config.Dev.Root)
		if err != nil {
			fmt.Println("Error adding watcher:", err)
			return
		}

		sseClients := make([]chan struct{}, 0)
		handler := http.NewServeMux()

		handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			var fileName string

			if r.URL.Path == "/" {
				fileName = "index.html"
			} else {
				if strings.HasSuffix(r.URL.Path, ".html") {
					http.Redirect(w, r, strings.TrimSuffix(r.URL.Path, ".html"), http.StatusTemporaryRedirect)
					return
				}

				fileName = r.URL.Path[1:]
				if filepath.Ext(fileName) == "" {
					fileName = fileName + ".html"
				}
			}

			fileContents, err := os.ReadFile(filepath.Join(constants.Config.Dev.Root, fileName))
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}

			fileType := mime.TypeByExtension(filepath.Ext(fileName))
			if strings.HasPrefix(fileType, "text/html") {
				fileContents = utils.ApplyBoilerplate(fileContents, true)
			}

			w.Header().Set("Content-Type", fileType)
			w.Write(fileContents)
		})

		handler.HandleFunc("/_html/hot-reload", func(w http.ResponseWriter, r *http.Request) {
			reloadChan := make(chan struct{})
			sseClients = append(sseClients, reloadChan)

			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")

			w.WriteHeader(http.StatusOK)

			fmt.Fprintf(w, "data: connected\n\n")
			w.(http.Flusher).Flush()

			for {
				select {
				case <-r.Context().Done():
					for i, ch := range sseClients {
						if ch == reloadChan {
							sseClients = append(sseClients[:i], sseClients[i+1:]...)
							return
						}
					}
				case <-reloadChan:
					fmt.Fprintf(w, "data: reload\n\n")
					w.(http.Flusher).Flush()
				}
			}
		})

		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Op&fsnotify.Write == fsnotify.Write {
						for _, ch := range sseClients {
							ch <- struct{}{}
						}
					}
				case err := <-watcher.Errors:
					fmt.Println("Watcher error:", err)
				}
			}
		}()

		fmt.Printf("HTML Dev Server ready\n\nAddress: %s\nRoot: %s\n", "http://localhost:"+fmt.Sprint(constants.Config.Dev.Port), constants.Config.Dev.Root)
		http.ListenAndServe(":"+fmt.Sprint(constants.Config.Dev.Port), handler)
	},
}
