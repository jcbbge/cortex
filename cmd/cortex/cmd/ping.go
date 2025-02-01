package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// pingCmd represents a simple test command
var pingCmd = &cobra.Command{
	Use:    "ping",
	Short:  "Test if Cortex is responding",
	Hidden: true, // This keeps it from showing up in help
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pong! 🏓")
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
