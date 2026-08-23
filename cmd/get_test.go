package cmd

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/jasonwashburn/memectl/internal/imgflip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTemplatesHelp(t *testing.T) {
	root := newRootCmd()
	var output bytes.Buffer
	root.SetOut(&output)
	root.SetArgs([]string{"get", "templates", "--help"})

	require.NoError(t, root.Execute())
	assert.Contains(t, output.String(), "memectl get templates")
}

func TestGetTemplates(t *testing.T) {
	tests := []struct {
		name    string
		client  templateRetriever
		want    string
		wantErr string
	}{
		{
			name: "success",
			client: fakeTemplateClient{templates: []imgflip.Template{
				{ID: "1", Name: "One Does Not Simply", BoxCount: 2, Width: 568, Height: 335},
			}},
			want: "ID  NAME                 BOXES  DIMENSIONS\n1   One Does Not Simply  2      568x335\n",
		},
		{
			name:    "client failure",
			client:  fakeTemplateClient{err: errors.New("network unavailable")},
			wantErr: "get templates: network unavailable",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command := newTemplatesCmd(test.client)
			var output bytes.Buffer
			command.SetOut(&output)
			command.SetArgs([]string{})

			err := command.Execute()
			if test.wantErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.wantErr)
				assert.Empty(t, output.String())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, output.String())
		})
	}
}

func TestGetTemplatesOutputFlag(t *testing.T) {
	for _, args := range [][]string{{"--output", "wide"}, {"-o", "wide"}} {
		command := newTemplatesCmd(fakeTemplateClient{templates: []imgflip.Template{
			{ID: "1", Name: "One Does Not Simply", BoxCount: 2, Width: 568, Height: 335, URL: "https://i.imgflip.com/1.jpg"},
		}})
		var output bytes.Buffer
		command.SetOut(&output)
		command.SetArgs(args)

		require.NoError(t, command.Execute(), "args %q", args)
		assert.Equal(t, "ID  NAME                 BOXES  DIMENSIONS  URL\n1   One Does Not Simply  2      568x335     https://i.imgflip.com/1.jpg\n", output.String())
	}
}

func TestGetTemplatesRejectsUnsupportedOutput(t *testing.T) {
	command := newTemplatesCmd(fakeTemplateClient{err: errors.New("client should not be called")})
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetArgs([]string{"--output", "json"})

	err := command.Execute()
	require.Error(t, err)
	assert.ErrorContains(t, err, `unsupported output format "json"`)
	assert.Empty(t, output.String())
}

type fakeTemplateClient struct {
	templates []imgflip.Template
	err       error
}

func (c fakeTemplateClient) Templates(context.Context) ([]imgflip.Template, error) {
	return c.templates, c.err
}
