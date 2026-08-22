// Package cmd contains memectl's Cobra commands.
package cmd

import "github.com/spf13/cobra"

var (
	// These values may be set with -ldflags during a release build.
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "memectl",
	Short:   "Generate memes through Imgflip",
	Long:    "memectl is a command-line tool for generating memes through Imgflip.",
	Version: version + " (commit: " + commit + ", built: " + date + ")",
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
