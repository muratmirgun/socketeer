package commands

import (
	"fmt"

	"github.com/muratmirgun/socketeer/internal/parser"
	"github.com/spf13/cobra"
)

var src string
var out string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate wsapi.yaml from Go source annotations",
	Long:  `Scans Go files for WebSocket annotations and generates wsapi.yaml spec.`,
	Run: func(cmd *cobra.Command, args []string) {
		if src == "" {
			src = "./"
		}
		if out == "" {
			out = "wsdocs/wsapi.yaml"
		}
		fmt.Printf("Parsing Go files in %s...\n", src)
		err := parser.ParseAndWriteSpec(src, out)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Spec written to %s\n", out)
	},
}

func init() {
	generateCmd.Flags().StringVar(&src, "src", "./", "Source directory to scan for Go files")
	generateCmd.Flags().StringVar(&out, "out", "wsdocs/wsapi.yaml", "Output spec file (YAML)")
	rootCmd.AddCommand(generateCmd)
}
