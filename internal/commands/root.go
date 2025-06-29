package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "socketeer",
	Short:   "socketeer - WebSocket API doc & playground generator",
	Long:    `socketeer is a modern, Swagger-like documentation and playground generator for WebSocket APIs in Go.`,
	Example: `  socketeer init\n  socketeer generate --src ./ --out ./wsdocs/wsapi.yaml\n  socketeer serve\n  socketeer validate\n  socketeer fmt\n  socketeer version`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
