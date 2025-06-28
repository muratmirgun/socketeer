package commands

import (
	"fmt"

	"github.com/muratmirgun/socketeer/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show socketeer version",
	Long:  `Prints the current version of socketeer CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("socketeer version %s\n", version.GetVersion())
		fmt.Printf("  commit: %s\n", version.GetCommit())
		fmt.Printf("  date: %s\n", version.GetDate())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
