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

func TestCreateMemeHelp(t *testing.T) {
	root := newRootCmd()
	var output bytes.Buffer
	root.SetOut(&output)
	root.SetArgs([]string{"create", "meme", "--help"})

	require.NoError(t, root.Execute())
	assert.Contains(t, output.String(), "memectl create meme <template-id>")
}

func TestCreateMeme(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		getenv    func(string) string
		client    memeCreator
		want      string
		wantErr   string
		wantInput *imgflip.CaptionImageRequest
		wantCalls int
	}{
		{
			name:      "success",
			args:      []string{"181913649", "--text", "first, still first", "--text", "second"},
			getenv:    credentials,
			client:    &fakeMemeClient{result: imgflip.CaptionImageResult{ImageURL: "https://i.imgflip.com/image.jpg", PageURL: "https://imgflip.com/i/page"}},
			want:      "Created meme from template 181913649.\nImage URL: https://i.imgflip.com/image.jpg\nImgflip page URL: https://imgflip.com/i/page\n",
			wantInput: &imgflip.CaptionImageRequest{TemplateID: "181913649", Username: "meme-user", Password: "meme-password", Texts: []string{"first, still first", "second"}},
			wantCalls: 1,
		},
		{
			name:    "missing text",
			args:    []string{"181913649"},
			getenv:  credentials,
			client:  &fakeMemeClient{},
			wantErr: "at least one --text value is required",
		},
		{
			name:    "missing template ID",
			args:    []string{"--text", "caption"},
			getenv:  credentials,
			client:  &fakeMemeClient{},
			wantErr: "a template ID is required",
		},
		{
			name:    "too many template IDs",
			args:    []string{"one", "two", "--text", "caption"},
			getenv:  credentials,
			client:  &fakeMemeClient{},
			wantErr: "accepts exactly one template ID",
		},
		{
			name:    "missing username",
			args:    []string{"181913649", "--text", "caption"},
			getenv:  func(string) string { return "" },
			client:  &fakeMemeClient{},
			wantErr: "IMGFLIP_USERNAME must be set",
		},
		{
			name: "missing password",
			args: []string{"181913649", "--text", "caption"},
			getenv: func(key string) string {
				if key == "IMGFLIP_USERNAME" {
					return "meme-user"
				}
				return ""
			},
			client:  &fakeMemeClient{},
			wantErr: "IMGFLIP_PASSWORD must be set",
		},
		{
			name:      "client failure",
			args:      []string{"181913649", "--text", "caption"},
			getenv:    credentials,
			client:    &fakeMemeClient{err: errors.New("Imgflip rejected request")},
			wantErr:   "create meme: Imgflip rejected request",
			wantCalls: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command := newMemeCmd(test.client, test.getenv)
			var output bytes.Buffer
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.Execute()
			if test.wantErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.wantErr)
				assert.Empty(t, output.String())
				assert.Equal(t, test.wantCalls, test.client.(*fakeMemeClient).calls)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, output.String())
			if test.wantInput != nil {
				input := test.client.(*fakeMemeClient).input
				assert.Equal(t, *test.wantInput, input)
			}
			assert.Equal(t, test.wantCalls, test.client.(*fakeMemeClient).calls)
		})
	}
}

func credentials(key string) string {
	if key == "IMGFLIP_USERNAME" {
		return "meme-user"
	}
	return "meme-password"
}

type fakeMemeClient struct {
	result imgflip.CaptionImageResult
	err    error
	input  imgflip.CaptionImageRequest
	calls  int
}

func (c *fakeMemeClient) CaptionImage(_ context.Context, input imgflip.CaptionImageRequest) (imgflip.CaptionImageResult, error) {
	c.calls++
	c.input = input
	return c.result, c.err
}
