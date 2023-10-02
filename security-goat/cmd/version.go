package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("🐐 security-goat version: %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
