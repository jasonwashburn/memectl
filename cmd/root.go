// Package cmd contains memectl's Cobra commands.
package cmd

import (
	"fmt"
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
	storePath, err := inventory.ResolvePath(defaultGetenv, os.UserHomeDir)
	if err != nil {
		panic(fmt.Sprintf("resolve meme inventory: %v", err))
	}
	store := inventory.New(storePath)
	root.AddCommand(newCreateCmd(client, store, defaultGetenv))
	root.AddCommand(newGetCmd(client, store))
	return root
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
