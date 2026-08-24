// Package cmd contains memectl's Cobra commands.
package cmd

import (
	"os"

	"github.com/jasonwashburn/memectl/internal/imgflip"
	"github.com/jasonwashburn/memectl/internal/inventory"
	"github.com/spf13/cobra"
)

var (
	// These values may be set with -ldflags during a release build.
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func newRootCmd(store memeStore) *cobra.Command {
	client := imgflip.NewClient(nil)
	root := &cobra.Command{
		Use:           "memectl",
		Short:         "Generate memes through Imgflip",
		Long:          "memectl is a command-line tool for generating memes through Imgflip.",
		Version:       version + " (commit: " + commit + ", built: " + date + ")",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.AddCommand(newCreateCmd(client, store, defaultGetenv))
	root.AddCommand(newGetCmd(client, store))
	return root
}

// Execute runs the root command.
func Execute() error {
	storePath, err := inventory.ResolvePath(defaultGetenv, os.UserHomeDir)
	if err != nil {
		return err
	}
	return newRootCmd(inventory.New(storePath)).Execute()
}
