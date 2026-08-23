package cmd

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/jasonwashburn/memectl/internal/imgflip"
)

func TestCreateMemeHelp(t *testing.T) {
	root := newRootCmd()
	var output bytes.Buffer
	root.SetOut(&output)
	root.SetArgs([]string{"create", "meme", "--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(output.String(), "memectl create meme <template-id>") {
		t.Fatalf("help output = %q, want meme command", output.String())
	}
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
				if err == nil || !strings.Contains(err.Error(), test.wantErr) {
					t.Fatalf("Execute() error = %v, want containing %q", err, test.wantErr)
				}
				if output.Len() != 0 {
					t.Fatalf("output = %q, want no output", output.String())
				}
				if calls := test.client.(*fakeMemeClient).calls; calls != test.wantCalls {
					t.Fatalf("CaptionImage() calls = %d, want %d", calls, test.wantCalls)
				}
				return
			}
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			if output.String() != test.want {
				t.Fatalf("output = %q, want %q", output.String(), test.want)
			}
			if test.wantInput != nil {
				input := test.client.(*fakeMemeClient).input
				if !reflect.DeepEqual(input, *test.wantInput) {
					t.Fatalf("CaptionImage() input = %#v, want %#v", input, *test.wantInput)
				}
			}
			if calls := test.client.(*fakeMemeClient).calls; calls != test.wantCalls {
				t.Fatalf("CaptionImage() calls = %d, want %d", calls, test.wantCalls)
			}
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
