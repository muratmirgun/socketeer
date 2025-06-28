package commands

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var dir string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve docs and playground from a directory",
	Long:  `Starts a local HTTP server to serve index.html, wsapi.yaml, logo.png, etc. from the specified directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		if dir == "" {
			dir = "wsdocs"
		}
		port := "8080"
		if p := os.Getenv("SOCKETEER_PORT"); p != "" {
			port = p
		}
		fmt.Printf("Serving %s at http://localhost:%s ...\n", dir, port)
		http.Handle("/", http.FileServer(http.Dir(dir)))
		http.ListenAndServe(":"+port, nil)
	},
}

func init() {
	serveCmd.Flags().StringVar(&dir, "dir", "wsdocs", "Directory to serve static files from")
	rootCmd.AddCommand(serveCmd)
}
