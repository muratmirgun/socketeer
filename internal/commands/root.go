package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "wsdoc",
	Short:   "wsdoc - WebSocket API doc & playground generator",
	Long:    `wsdoc is a modern, Swagger-like documentation and playground generator for WebSocket APIs in Go.`,
	Example: `  wsdoc init\n  wsdoc generate --src ./ --out ./wsdocs/wsapi.yaml\n  wsdoc serve\n  wsdoc validate\n  wsdoc fmt\n  wsdoc version`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
