package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize wsdoc project (creates wsdocs/ with base files)",
	Long:  `Creates a wsdocs/ directory with base wsapi.yaml and index.html for your WebSocket API docs project.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("wsdocs"); err == nil {
			fmt.Println("wsdocs/ already exists. Aborting.")
			return
		}
		os.Mkdir("wsdocs", 0755)
		copyFile("internal/templates/index.html", "wsdocs/index.html")
		copyFile("internal/templates/wsapi.yaml", "wsdocs/wsapi.yaml")
		fmt.Println("Initialized wsdocs/ with index.html and wsapi.yaml.")
	},
}

func copyFile(src, dst string) error {
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()
	dstF, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstF.Close()
	_, err = io.Copy(dstF, srcF)
	return err
}

func init() {
	rootCmd.AddCommand(initCmd)
}
