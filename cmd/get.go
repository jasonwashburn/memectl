package cmd

import (
	"context"
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/jasonwashburn/memectl/internal/imgflip"
	"github.com/jasonwashburn/memectl/internal/inventory"
	"github.com/spf13/cobra"
)

type templateRetriever interface {
	Templates(context.Context) ([]imgflip.Template, error)
}

func newGetCmd(client templateRetriever, store memeStore) *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Display Imgflip resources",
	}
	getCmd.AddCommand(newTemplatesCmd(client))
	getCmd.AddCommand(newMemesCmd(store))
	return getCmd
}

func newMemesCmd(store memeStore) *cobra.Command {
	return newMemesCmdAt(store, time.Now)
}

func newMemesCmdAt(store memeStore, now func() time.Time) *cobra.Command {
	var output string
	command := &cobra.Command{
		Use:           "memes",
		Short:         "List managed memes",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if output != "" && output != "wide" {
				return fmt.Errorf("unable to match a printer suitable for the output format %q", output)
			}
			memes, err := store.Load()
			if err != nil {
				return fmt.Errorf("get memes: %w", err)
			}
			if len(memes) == 0 {
				_, err := fmt.Fprintln(cmd.OutOrStdout(), "No resources found.")
				return err
			}
			inventory.SortByName(memes)
			writer := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			header := "NAME\tTEMPLATE ID\tAGE\tIMAGE URL"
			if output == "wide" {
				header += "\tPAGE URL"
			}
			if _, err := fmt.Fprintln(writer, header); err != nil {
				return fmt.Errorf("write meme header: %w", err)
			}
			for _, meme := range memes {
				age := now().Sub(meme.CreatedAt).Truncate(time.Second).String()
				if output == "wide" {
					_, err = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n", meme.Name, meme.TemplateID, age, meme.ImageURL, meme.PageURL)
				} else {
					_, err = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", meme.Name, meme.TemplateID, age, meme.ImageURL)
				}
				if err != nil {
					return fmt.Errorf("write meme row: %w", err)
				}
			}
			return writer.Flush()
		},
	}
	command.Flags().StringVarP(&output, "output", "o", "", "Output format: wide")
	return command
}

func newTemplatesCmd(client templateRetriever) *cobra.Command {
	var output string
	command := &cobra.Command{
		Use:           "templates",
		Short:         "List public meme templates",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if output != "" && output != "wide" {
				return fmt.Errorf("unsupported output format %q", output)
			}

			templates, err := client.Templates(cmd.Context())
			if err != nil {
				return fmt.Errorf("get templates: %w", err)
			}

			writer := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			header := "ID\tNAME\tBOXES\tDIMENSIONS"
			if output == "wide" {
				header += "\tURL"
			}
			if _, err := fmt.Fprintln(writer, header); err != nil {
				return fmt.Errorf("write template header: %w", err)
			}
			for _, template := range templates {
				if output == "wide" {
					if _, err := fmt.Fprintf(writer, "%s\t%s\t%d\t%dx%d\t%s\n", template.ID, template.Name, template.BoxCount, template.Width, template.Height, template.URL); err != nil {
						return fmt.Errorf("write template row: %w", err)
					}
					continue
				}
				if _, err := fmt.Fprintf(writer, "%s\t%s\t%d\t%dx%d\n", template.ID, template.Name, template.BoxCount, template.Width, template.Height); err != nil {
					return fmt.Errorf("write template row: %w", err)
				}
			}
			return writer.Flush()
		},
	}
	command.Flags().StringVarP(&output, "output", "o", "", "Output format: wide")
	return command
}
