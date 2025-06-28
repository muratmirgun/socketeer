package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "v0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show wsdoc version",
	Long:  `Prints the current version of wsdoc CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wsdoc version", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
