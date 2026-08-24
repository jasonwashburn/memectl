package cmd

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jasonwashburn/memectl/internal/imgflip"
	"github.com/jasonwashburn/memectl/internal/inventory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateMemeHelp(t *testing.T) {
	root := newRootCmd(&fakeMemeStore{})
	var output bytes.Buffer
	root.SetOut(&output)
	root.SetArgs([]string{"create", "meme", "--help"})

	require.NoError(t, root.Execute())
	assert.Contains(t, output.String(), "memectl create meme <name> --template <template-id>")
}

func TestCreateMeme(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		getenv    func(string) string
		client    memeCreator
		store     *fakeMemeStore
		want      string
		wantErr   string
		wantInput *imgflip.CaptionImageRequest
		wantCalls int
	}{
		{
			name:      "success",
			args:      []string{"first-meme", "--template", "181913649", "--text", "first, still first", "--text", "second"},
			getenv:    credentials,
			client:    &fakeMemeClient{result: imgflip.CaptionImageResult{ImageURL: "https://i.imgflip.com/image.jpg", PageURL: "https://imgflip.com/i/page"}},
			store:     &fakeMemeStore{},
			want:      "Created meme \"first-meme\" from template 181913649.\nImage URL: https://i.imgflip.com/image.jpg\nImgflip page URL: https://imgflip.com/i/page\n",
			wantInput: &imgflip.CaptionImageRequest{TemplateID: "181913649", Username: "meme-user", Password: "meme-password", Texts: []string{"first, still first", "second"}},
			wantCalls: 1,
		},
		{
			name:    "missing text",
			args:    []string{"meme", "--template", "181913649"},
			getenv:  credentials,
			client:  &fakeMemeClient{},
			wantErr: "failed to create meme: at least one --text value is required",
		},
		{
			name:    "missing local name",
			args:    []string{"--template", "181913649", "--text", "caption"},
			getenv:  credentials,
			client:  &fakeMemeClient{},
			wantErr: "failed to create meme: a local meme name is required",
		},
		{
			name:    "too many names",
			args:    []string{"one", "two", "--template", "181913649", "--text", "caption"},
			getenv:  credentials,
			client:  &fakeMemeClient{},
			wantErr: "failed to create meme: accepts exactly one local meme name",
		},
		{
			name:    "missing template",
			args:    []string{"meme", "--text", "caption"},
			getenv:  credentials,
			client:  &fakeMemeClient{},
			wantErr: "failed to create meme: --template is required",
		},
		{
			name:    "invalid name",
			args:    []string{"Not-valid", "--template", "181913649", "--text", "caption"},
			getenv:  credentials,
			client:  &fakeMemeClient{},
			wantErr: "failed to create meme: name \"Not-valid\" must be a DNS-label-like value",
		},
		{
			name:    "duplicate name",
			args:    []string{"meme", "--template", "181913649", "--text", "caption"},
			getenv:  credentials,
			client:  &fakeMemeClient{},
			store:   &fakeMemeStore{memes: []inventory.Meme{{Name: "meme"}}},
			wantErr: `failed to create meme: meme "meme" already exists`,
		},
		{
			name:    "missing username",
			args:    []string{"meme", "--template", "181913649", "--text", "caption"},
			getenv:  func(string) string { return "" },
			client:  &fakeMemeClient{},
			wantErr: "failed to create meme: IMGFLIP_USERNAME must be set",
		},
		{
			name: "missing password",
			args: []string{"meme", "--template", "181913649", "--text", "caption"},
			getenv: func(key string) string {
				if key == "IMGFLIP_USERNAME" {
					return "meme-user"
				}
				return ""
			},
			client:  &fakeMemeClient{},
			wantErr: "failed to create meme: IMGFLIP_PASSWORD must be set",
		},
		{
			name:      "client failure",
			args:      []string{"meme", "--template", "181913649", "--text", "caption"},
			getenv:    credentials,
			client:    &fakeMemeClient{err: errors.New("Imgflip rejected request")},
			wantErr:   "failed to create meme: Imgflip rejected request",
			wantCalls: 1,
		},
		{
			name:      "local persistence failure",
			args:      []string{"meme", "--template", "181913649", "--text", "caption"},
			getenv:    credentials,
			client:    &fakeMemeClient{result: imgflip.CaptionImageResult{ImageURL: "https://i.imgflip.com/image.jpg", PageURL: "https://imgflip.com/i/page"}},
			store:     &fakeMemeStore{addErr: errors.New("disk full")},
			want:      "Meme \"meme\" was created remotely but was not recorded locally.\nImage URL: https://i.imgflip.com/image.jpg\nImgflip page URL: https://imgflip.com/i/page\n",
			wantErr:   "failed to create meme: remote meme \"meme\" was not recorded locally",
			wantCalls: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.store == nil {
				test.store = &fakeMemeStore{}
			}
			command := newMemeCmd(test.client, test.store, test.getenv)
			var output bytes.Buffer
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.Execute()
			if test.wantErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.wantErr)
				assert.Equal(t, test.want, output.String())
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
			if test.wantCalls == 1 && test.wantErr == "" {
				require.Len(t, test.store.memes, 1)
				assert.Equal(t, "first-meme", test.store.memes[0].Name)
				assert.Equal(t, []string{"first, still first", "second"}, test.store.memes[0].Texts)
				assert.WithinDuration(t, time.Now().UTC(), test.store.memes[0].CreatedAt, time.Second)
			}
		})
	}
}

type fakeMemeStore struct {
	memes   []inventory.Meme
	loadErr error
	addErr  error
}

func (s *fakeMemeStore) Load() ([]inventory.Meme, error) { return s.memes, s.loadErr }
func (s *fakeMemeStore) Add(meme inventory.Meme) error {
	if s.addErr != nil {
		return s.addErr
	}
	s.memes = append(s.memes, meme)
	return nil
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
