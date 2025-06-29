package commands

import (
	"fmt"
	"os"

	"github.com/muratmirgun/socketeer/internal/spec"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var validateFile string

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate wsapi.yaml file",
	Long:  `Validates the structure and content of a wsapi.yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if validateFile == "" {
			validateFile = "wsdocs/wsapi.yaml"
		}

		// Check if file exists
		if _, err := os.Stat(validateFile); os.IsNotExist(err) {
			fmt.Printf("Error: File %s does not exist\n", validateFile)
			return
		}

		// Read and parse the file
		data, err := os.ReadFile(validateFile)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		var s spec.Spec
		if err := yaml.Unmarshal(data, &s); err != nil {
			fmt.Printf("Error parsing YAML: %v\n", err)
			return
		}

		// Validate the spec
		errors := validateSpec(&s)
		if len(errors) == 0 {
			fmt.Printf("✅ %s is valid\n", validateFile)
		} else {
			fmt.Printf("❌ %s has validation errors:\n", validateFile)
			for _, err := range errors {
				fmt.Printf("  - %s\n", err)
			}
		}
	},
}

func validateSpec(s *spec.Spec) []string {
	var errors []string

	// Validate info
	if s.Info.Title == "" {
		errors = append(errors, "Info.Title is required")
	}
	if s.Info.Version == "" {
		errors = append(errors, "Info.Version is required")
	}

	// Validate sockets
	if len(s.Sockets) == 0 {
		errors = append(errors, "At least one WebSocket endpoint is required")
	}

	for i, socket := range s.Sockets {
		if socket.Name == "" {
			errors = append(errors, fmt.Sprintf("Socket[%d].Name is required", i))
		}
		if socket.URL == "" {
			errors = append(errors, fmt.Sprintf("Socket[%d].URL is required", i))
		}

		// Validate grouped messages
		for j, groupedMsg := range socket.GroupedMessages {
			if groupedMsg.Type == "" {
				errors = append(errors, fmt.Sprintf("Socket[%d].GroupedMessages[%d].Type is required", i, j))
			}
			if groupedMsg.Send == nil && groupedMsg.Receive == nil {
				errors = append(errors, fmt.Sprintf("Socket[%d].GroupedMessages[%d] must have at least one Send or Receive message", i, j))
			}
		}

		// Validate legacy messages (for backward compatibility)
		for j, msg := range socket.Messages {
			if msg.Type == "" {
				errors = append(errors, fmt.Sprintf("Socket[%d].Messages[%d].Type is required", i, j))
			}
			if msg.Direction == "" {
				errors = append(errors, fmt.Sprintf("Socket[%d].Messages[%d].Direction is required", i, j))
			}
		}
	}

	return errors
}

func init() {
	validateCmd.Flags().StringVar(&validateFile, "file", "wsdocs/wsapi.yaml", "File to validate")
	rootCmd.AddCommand(validateCmd)
} 