package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/jasonwashburn/memectl/internal/imgflip"
	"github.com/spf13/cobra"
)

type memeCreator interface {
	CaptionImage(context.Context, imgflip.CaptionImageRequest) (imgflip.CaptionImageResult, error)
}

func newCreateCmd(client memeCreator, getenv func(string) string) *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create Imgflip resources",
	}
	createCmd.AddCommand(newMemeCmd(client, getenv))
	return createCmd
}

func newMemeCmd(client memeCreator, getenv func(string) string) *cobra.Command {
	var texts []string

	memeCmd := &cobra.Command{
		Use:   "meme <template-id>",
		Short: "Create a captioned meme",
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("create meme: a template ID is required")
			}
			if len(args) > 1 {
				return fmt.Errorf("create meme: accepts exactly one template ID")
			}
			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(texts) == 0 {
				return fmt.Errorf("create meme: at least one --text value is required")
			}

			username := getenv("IMGFLIP_USERNAME")
			if username == "" {
				return fmt.Errorf("create meme: IMGFLIP_USERNAME must be set")
			}
			password := getenv("IMGFLIP_PASSWORD")
			if password == "" {
				return fmt.Errorf("create meme: IMGFLIP_PASSWORD must be set")
			}

			result, err := client.CaptionImage(cmd.Context(), imgflip.CaptionImageRequest{
				TemplateID: args[0],
				Username:   username,
				Password:   password,
				Texts:      texts,
			})
			if err != nil {
				return fmt.Errorf("create meme: %w", err)
			}

			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Created meme from template %s.\nImage URL: %s\nImgflip page URL: %s\n", args[0], result.ImageURL, result.PageURL); err != nil {
				return fmt.Errorf("write created meme: %w", err)
			}
			return nil
		},
	}
	memeCmd.Flags().StringSliceVar(&texts, "text", nil, "Caption text (repeatable; required)")
	return memeCmd
}

func defaultGetenv(key string) string {
	return os.Getenv(key)
}
