package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jasonwashburn/memectl/internal/imgflip"
	"github.com/jasonwashburn/memectl/internal/inventory"
	"github.com/spf13/cobra"
)

type memeCreator interface {
	CaptionImage(context.Context, imgflip.CaptionImageRequest) (imgflip.CaptionImageResult, error)
}

type memeStore interface {
	Load() ([]inventory.Meme, error)
	Add(inventory.Meme) error
	Remove([]string) ([]string, error)
}

func newCreateCmd(client memeCreator, store memeStore, getenv func(string) string) *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create Imgflip resources",
	}
	createCmd.AddCommand(newMemeCmd(client, store, getenv))
	return createCmd
}

func newMemeCmd(client memeCreator, store memeStore, getenv func(string) string) *cobra.Command {
	var texts []string
	var templateID string

	memeCmd := &cobra.Command{
		Use:   "meme <name> --template <template-id>",
		Short: "Create a captioned meme",
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("failed to create meme: a local meme name is required")
			}
			if len(args) > 1 {
				return fmt.Errorf("failed to create meme: accepts exactly one local meme name")
			}
			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(texts) == 0 {
				return fmt.Errorf("failed to create meme: at least one --text value is required")
			}
			if templateID == "" {
				return fmt.Errorf("failed to create meme: --template is required")
			}
			name := args[0]
			if !inventory.ValidName(name) {
				return fmt.Errorf("failed to create meme: name %q must be a DNS-label-like value", name)
			}
			memes, err := store.Load()
			if err != nil {
				return fmt.Errorf("failed to create meme: %w", err)
			}
			if inventory.Contains(memes, name) {
				return fmt.Errorf("failed to create meme: meme %q already exists", name)
			}

			username := getenv("IMGFLIP_USERNAME")
			if username == "" {
				return fmt.Errorf("failed to create meme: IMGFLIP_USERNAME must be set")
			}
			password := getenv("IMGFLIP_PASSWORD")
			if password == "" {
				return fmt.Errorf("failed to create meme: IMGFLIP_PASSWORD must be set")
			}

			result, err := client.CaptionImage(cmd.Context(), imgflip.CaptionImageRequest{
				TemplateID: templateID,
				Username:   username,
				Password:   password,
				Texts:      texts,
			})
			if err != nil {
				return fmt.Errorf("failed to create meme: %w", err)
			}

			record := inventory.Meme{Name: name, TemplateID: templateID, Texts: texts, ImageURL: result.ImageURL, PageURL: result.PageURL, CreatedAt: time.Now().UTC()}
			if err := store.Add(record); err != nil {
				if _, writeErr := fmt.Fprintf(cmd.OutOrStdout(), "Meme %q was created remotely but was not recorded locally.\nImage URL: %s\nImgflip page URL: %s\n", name, result.ImageURL, result.PageURL); writeErr != nil {
					return fmt.Errorf("write unrecorded meme: %w", writeErr)
				}
				return fmt.Errorf("failed to create meme: remote meme %q was not recorded locally: %w", name, err)
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Created meme %q from template %s.\nImage URL: %s\nImgflip page URL: %s\n", name, templateID, result.ImageURL, result.PageURL); err != nil {
				return fmt.Errorf("write created meme: %w", err)
			}
			return nil
		},
	}
	memeCmd.Flags().StringArrayVar(&texts, "text", nil, "Caption text (repeatable; required)")
	memeCmd.Flags().StringVar(&templateID, "template", "", "Imgflip template ID (required)")
	return memeCmd
}

func defaultGetenv(key string) string {
	return os.Getenv(key)
}
