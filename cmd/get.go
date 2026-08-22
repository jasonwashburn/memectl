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
	return &cobra.Command{
		Use:           "templates",
		Short:         "List public meme templates",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			templates, err := client.Templates(cmd.Context())
			if err != nil {
				return fmt.Errorf("get templates: %w", err)
			}

			writer := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(writer, "ID\tNAME\tBOXES\tDIMENSIONS")
			for _, template := range templates {
				fmt.Fprintf(writer, "%s\t%s\t%d\t%dx%d\n", template.ID, template.Name, template.BoxCount, template.Width, template.Height)
			}
			return writer.Flush()
		},
	}
}
