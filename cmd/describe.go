package cmd

import (
	"fmt"
	"time"

	"github.com/jasonwashburn/memectl/internal/inventory"
	"github.com/spf13/cobra"
)

func newDescribeCmd(store memeStore) *cobra.Command {
	describeCmd := &cobra.Command{
		Use:     "describe",
		Aliases: []string{"desc"},
		Short:   "Display managed meme details",
	}
	describeCmd.AddCommand(newDescribeMemeCmd(store))
	return describeCmd
}

func newDescribeMemeCmd(store memeStore) *cobra.Command {
	return &cobra.Command{
		Use:   "meme <name>",
		Short: "Display a locally managed meme",
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("failed to describe meme: a local meme name is required")
			}
			if len(args) > 1 {
				return fmt.Errorf("failed to describe meme: accepts exactly one local meme name")
			}
			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if !inventory.ValidName(name) {
				return fmt.Errorf("failed to describe meme: name %q must be a DNS-label-like value", name)
			}

			memes, err := store.Load()
			if err != nil {
				return fmt.Errorf("failed to describe meme: %w", err)
			}
			for _, meme := range memes {
				if meme.Name != name {
					continue
				}
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Name: %s\nTemplate ID: %s\nTexts:\n", meme.Name, meme.TemplateID); err != nil {
					return fmt.Errorf("write meme details: %w", err)
				}
				for index, text := range meme.Texts {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "  %d: %s\n", index, text); err != nil {
						return fmt.Errorf("write meme details: %w", err)
					}
				}
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Image URL: %s\nImgflip page URL: %s\nCreated at: %s\n", meme.ImageURL, meme.PageURL, meme.CreatedAt.UTC().Format(time.RFC3339Nano)); err != nil {
					return fmt.Errorf("write meme details: %w", err)
				}
				return nil
			}
			return fmt.Errorf("failed to describe meme: meme %q not found", name)
		},
	}
}
