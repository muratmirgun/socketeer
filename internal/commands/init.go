package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	baseURL = "https://raw.githubusercontent.com/muratmirgun/socketeer/main/internal/templates"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize wsdoc project (creates wsdocs/ with base files)",
	Long:  `Creates a wsdocs/ directory with base wsapi.yaml and index.html for your WebSocket API docs project.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create wsdocs directory if it doesn't exist
		if _, err := os.Stat("wsdocs"); os.IsNotExist(err) {
			os.Mkdir("wsdocs", 0755)
			fmt.Println("Created wsdocs/ directory.")
		} else {
			fmt.Println("wsdocs/ directory already exists. Updating files...")
		}
		
		// Download files from GitHub
		files := []string{"index.html", "wsapi.yaml", "logo.png"}
		
		for _, file := range files {
			fmt.Printf("Downloading %s...\n", file)
			url := baseURL + "/" + file
			err := downloadFile(url, filepath.Join("wsdocs", file))
			if err != nil {
				fmt.Printf("Error downloading %s: %v\n", file, err)
				return
			}
		}
		
		fmt.Println("Initialized wsdocs/ with index.html, wsapi.yaml, and logo.png.")
	},
}

func downloadFile(url, dst string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}
	
	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()
	
	_, err = io.Copy(file, resp.Body)
	return err
}

func init() {
	rootCmd.AddCommand(initCmd)
}
