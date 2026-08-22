package cmd

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jasonwashburn/memectl/internal/imgflip"
)

func TestGetTemplatesHelp(t *testing.T) {
	root := newRootCmd()
	var output bytes.Buffer
	root.SetOut(&output)
	root.SetArgs([]string{"get", "templates", "--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(output.String(), "memectl get templates") {
		t.Fatalf("help output = %q, want templates command", output.String())
	}
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
				if err == nil || !strings.Contains(err.Error(), test.wantErr) {
					t.Fatalf("Execute() error = %v, want containing %q", err, test.wantErr)
				}
				if output.Len() != 0 {
					t.Fatalf("output = %q, want no output", output.String())
				}
				return
			}
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			if output.String() != test.want {
				t.Fatalf("output = %q, want %q", output.String(), test.want)
			}
		})
	}
}

type fakeTemplateClient struct {
	templates []imgflip.Template
	err       error
}

func (c fakeTemplateClient) Templates(context.Context) ([]imgflip.Template, error) {
	return c.templates, c.err
}
