package commands

import (
	"fmt"
	"os"

	"github.com/muratmirgun/socketeer/internal/spec"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var fmtFile string
var fmtOutput string

var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Format wsapi.yaml file",
	Long:  `Formats and prettifies a wsapi.yaml file with consistent indentation and structure.`,
	Run: func(cmd *cobra.Command, args []string) {
		if fmtFile == "" {
			fmtFile = "wsdocs/wsapi.yaml"
		}
		if fmtOutput == "" {
			fmtOutput = fmtFile
		}

		// Check if file exists
		if _, err := os.Stat(fmtFile); os.IsNotExist(err) {
			fmt.Printf("Error: File %s does not exist\n", fmtFile)
			return
		}

		// Read and parse the file
		data, err := os.ReadFile(fmtFile)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		var s spec.Spec
		if err := yaml.Unmarshal(data, &s); err != nil {
			fmt.Printf("Error parsing YAML: %v\n", err)
			return
		}

		// Create output file
		outputFile, err := os.Create(fmtOutput)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			return
		}
		defer outputFile.Close()

		// Write formatted YAML
		encoder := yaml.NewEncoder(outputFile)
		encoder.SetIndent(2)
		if err := encoder.Encode(&s); err != nil {
			fmt.Printf("Error writing formatted YAML: %v\n", err)
			return
		}

		if fmtFile == fmtOutput {
			fmt.Printf("✅ Formatted %s\n", fmtFile)
		} else {
			fmt.Printf("✅ Formatted %s -> %s\n", fmtFile, fmtOutput)
		}
	},
}

func init() {
	fmtCmd.Flags().StringVar(&fmtFile, "file", "wsdocs/wsapi.yaml", "File to format")
	fmtCmd.Flags().StringVar(&fmtOutput, "output", "", "Output file (defaults to input file)")
	rootCmd.AddCommand(fmtCmd)
} 