package cmd

import (
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/jasonwashburn/memectl/internal/imgflip"
	"github.com/spf13/cobra"
)

type templateRetriever interface {
	Templates(context.Context) ([]imgflip.Template, error)
}

func newGetCmd(client templateRetriever) *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Display Imgflip resources",
	}
	getCmd.AddCommand(newTemplatesCmd(client))
	return getCmd
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
