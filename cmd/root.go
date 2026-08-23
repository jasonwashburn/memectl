// Package cmd contains memectl's Cobra commands.
package cmd

import (
	"github.com/jasonwashburn/memectl/internal/imgflip"
	"github.com/spf13/cobra"
)

var (
	// These values may be set with -ldflags during a release build.
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = newRootCmd()

func newRootCmd() *cobra.Command {
	client := imgflip.NewClient(nil)
	root := &cobra.Command{
		Use:           "memectl",
		Short:         "Generate memes through Imgflip",
		Long:          "memectl is a command-line tool for generating memes through Imgflip.",
		Version:       version + " (commit: " + commit + ", built: " + date + ")",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.AddCommand(newCreateCmd(client, defaultGetenv))
	root.AddCommand(newGetCmd(client))
	return root
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
