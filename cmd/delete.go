package cmd

import (
	"fmt"
	"strings"

	"github.com/jasonwashburn/memectl/internal/inventory"
	"github.com/spf13/cobra"
)

func newDeleteCmd(store memeStore) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete locally managed resources",
	}
	deleteCmd.AddCommand(newDeleteMemeCmd(store))
	return deleteCmd
}

func newDeleteMemeCmd(store memeStore) *cobra.Command {
	return &cobra.Command{
		Use:           "meme <name> [<name>...]",
		Short:         "Delete locally managed memes",
		Long:          "Delete locally managed memes. This removes only local inventory metadata; hosted Imgflip images remain unchanged.",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("failed to delete meme: at least one local meme name is required")
			}
			for _, name := range args {
				if !inventory.ValidName(name) {
					return fmt.Errorf("failed to delete meme: name %q must be a DNS-label-like value", name)
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			absent, err := store.Remove(args)
			if err != nil {
				return fmt.Errorf("failed to delete meme: %w", err)
			}
			absentSet := make(map[string]bool, len(absent))
			for _, name := range absent {
				absentSet[name] = true
			}
			for _, name := range args {
				if absentSet[name] {
					continue
				}
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Meme %q deleted.\n", name); err != nil {
					return fmt.Errorf("write deleted meme: %w", err)
				}
			}
			if len(absent) == 0 {
				return nil
			}
			notFound := make([]string, len(absent))
			for i, name := range absent {
				notFound[i] = fmt.Sprintf("meme %q not found", name)
			}
			return fmt.Errorf("failed to delete meme: %s", strings.Join(notFound, "; "))
		},
	}
}
