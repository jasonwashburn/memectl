package cmd

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/jasonwashburn/memectl/internal/inventory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDescribeMemeHelp(t *testing.T) {
	root := newRootCmd(&fakeMemeStore{})
	var output bytes.Buffer
	root.SetOut(&output)
	root.SetArgs([]string{"describe", "meme", "--help"})

	require.NoError(t, root.Execute())
	assert.Contains(t, output.String(), "memectl describe meme <name>")

	root = newRootCmd(&fakeMemeStore{})
	output.Reset()
	root.SetOut(&output)
	root.SetArgs([]string{"desc", "meme", "--help"})

	require.NoError(t, root.Execute())
	assert.Contains(t, output.String(), "memectl describe meme <name>")
}

func TestDescribeMeme(t *testing.T) {
	createdAt := time.Date(2026, time.August, 24, 12, 0, 0, 123456789, time.UTC)
	tests := []struct {
		name    string
		args    []string
		store   *fakeMemeStore
		want    string
		wantErr string
	}{
		{
			name: "success",
			args: []string{"writing-memes"},
			store: &fakeMemeStore{memes: []inventory.Meme{{
				Name: "writing-memes", TemplateID: "181913649", Texts: []string{"Writing memes manually", "Using memectl"}, ImageURL: "https://i.imgflip.com/image.jpg", PageURL: "https://imgflip.com/i/page", CreatedAt: createdAt,
			}}},
			want: "Name: writing-memes\nTemplate ID: 181913649\nTexts:\n  0: Writing memes manually\n  1: Using memectl\nImage URL: https://i.imgflip.com/image.jpg\nImgflip page URL: https://imgflip.com/i/page\nCreated at: 2026-08-24T12:00:00.123456789Z\n",
		},
		{
			name: "empty caption",
			args: []string{"meme"},
			store: &fakeMemeStore{memes: []inventory.Meme{{
				Name: "meme", TemplateID: "1", Texts: []string{"first", ""}, ImageURL: "https://image", PageURL: "https://page", CreatedAt: createdAt,
			}}},
			want: "Name: meme\nTemplate ID: 1\nTexts:\n  0: first\n  1: \nImage URL: https://image\nImgflip page URL: https://page\nCreated at: 2026-08-24T12:00:00.123456789Z\n",
		},
		{name: "missing name", store: &fakeMemeStore{}, wantErr: "a local meme name is required"},
		{name: "extra name", args: []string{"first", "second"}, store: &fakeMemeStore{}, wantErr: "accepts exactly one local meme name"},
		{name: "invalid name", args: []string{"Not-valid"}, store: &fakeMemeStore{}, wantErr: "name \"Not-valid\" must be a DNS-label-like value"},
		{name: "missing record", args: []string{"missing"}, store: &fakeMemeStore{}, wantErr: "meme \"missing\" not found"},
		{name: "inventory read failure", args: []string{"meme"}, store: &fakeMemeStore{loadErr: errors.New("corrupt inventory")}, wantErr: "corrupt inventory"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command := newDescribeMemeCmd(test.store)
			var output bytes.Buffer
			command.SetOut(&output)
			command.SetArgs(test.args)

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
