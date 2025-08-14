package cmd

import (
	"fmt"
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

	devCmd.Flags().StringP("root", "r", ".", "Directory to watch")
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Run a dev server for html",
	Run: func(cmd *cobra.Command, args []string) {
		root, _ := cmd.Flags().GetString("root")

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Println("Error setting up file watcher:", err)
			return
		}
		defer watcher.Close()

		err = watcher.Add(root)
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

				if filepath.Ext(fileName) == "" {
					fileName = r.URL.Path[1:] + ".html"
				}
			}

			fileContents, err := os.ReadFile(filepath.Join(root, fileName))
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}

			fileType := mime.TypeByExtension(filepath.Ext(fileName))
			if strings.HasPrefix(fileType, "text/html") {
				hotReloadScript := "<script>(()=>{const sse = new EventSource('/_html/hot-reload');sse.onopen=()=>{console.log('* HTML hot-reload enabled *')};sse.onmessage=e=>{if(e.data==='reload'){console.log('* HTML hot-reload triggered *');location.reload()}}})()</script>"

				fileContents = []byte(string(fileContents) + hotReloadScript)
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

		fmt.Printf("HTML Dev Server ready\n\nAddress: %s\nRoot: %s\n", "http://localhost:8080", root)
		http.ListenAndServe(":8080", handler)
	},
}
