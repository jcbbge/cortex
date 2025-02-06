package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var (
    Version   = "0.5.0-dev"
    BuildTime = "unknown"
    GitCommit = "unknown"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print the version number of Cortex",
    Long:  `All software has versions. This is Cortex's.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Cortex v%s\n", Version)
        fmt.Printf("Build Time: %s\n", BuildTime)
        fmt.Printf("Git Commit: %s\n", GitCommit)
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
}