package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/jasonwashburn/memectl/internal/inventory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteMemeHelp(t *testing.T) {
	root := newRootCmd(&fakeMemeStore{})
	var output bytes.Buffer
	root.SetOut(&output)
	root.SetArgs([]string{"delete", "meme", "--help"})

	require.NoError(t, root.Execute())
	assert.Contains(t, output.String(), "memectl delete meme <name> [<name>...]")
	assert.Contains(t, output.String(), "hosted Imgflip images remain unchanged")
}

func TestDeleteMeme(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		store     *fakeMemeStore
		want      string
		wantErr   string
		wantMemes []inventory.Meme
		wantCalls int
	}{
		{name: "single record without credentials", args: []string{"first"}, store: &fakeMemeStore{memes: []inventory.Meme{{Name: "first"}}}, want: "Meme \"first\" deleted.\n", wantMemes: []inventory.Meme{}, wantCalls: 1},
		{name: "multiple records preserve others", args: []string{"first", "third"}, store: &fakeMemeStore{memes: []inventory.Meme{{Name: "first"}, {Name: "second"}, {Name: "third"}}}, want: "Meme \"first\" deleted.\nMeme \"third\" deleted.\n", wantMemes: []inventory.Meme{{Name: "second"}}, wantCalls: 1},
		{name: "final records", args: []string{"first", "second"}, store: &fakeMemeStore{memes: []inventory.Meme{{Name: "first"}, {Name: "second"}}}, want: "Meme \"first\" deleted.\nMeme \"second\" deleted.\n", wantMemes: []inventory.Meme{}, wantCalls: 1},
		{name: "mixed present and absent", args: []string{"first", "missing"}, store: &fakeMemeStore{memes: []inventory.Meme{{Name: "first"}, {Name: "second"}}}, want: "Meme \"first\" deleted.\n", wantErr: "meme \"missing\" not found", wantMemes: []inventory.Meme{{Name: "second"}}, wantCalls: 1},
		{name: "duplicate name", args: []string{"first", "first"}, store: &fakeMemeStore{memes: []inventory.Meme{{Name: "first"}}}, want: "Meme \"first\" deleted.\n", wantErr: "meme \"first\" not found", wantMemes: []inventory.Meme{}, wantCalls: 1},
		{name: "all absent", args: []string{"missing", "unknown"}, store: &fakeMemeStore{memes: []inventory.Meme{{Name: "first"}}}, wantErr: "meme \"missing\" not found; meme \"unknown\" not found", wantMemes: []inventory.Meme{{Name: "first"}}, wantCalls: 1},
		{name: "missing name", store: &fakeMemeStore{memes: []inventory.Meme{{Name: "first"}}}, wantErr: "at least one local meme name is required", wantMemes: []inventory.Meme{{Name: "first"}}},
		{name: "invalid name", args: []string{"Not-valid"}, store: &fakeMemeStore{memes: []inventory.Meme{{Name: "first"}}}, wantErr: "name \"Not-valid\" must be a DNS-label-like value", wantMemes: []inventory.Meme{{Name: "first"}}},
		{name: "storage failure", args: []string{"first"}, store: &fakeMemeStore{removeErr: errors.New("disk full")}, wantErr: "failed to delete meme: disk full", wantCalls: 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command := newDeleteMemeCmd(test.store)
			var output bytes.Buffer
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.Execute()
			if test.wantErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, test.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.want, output.String())
			assert.Equal(t, test.wantMemes, test.store.memes)
			assert.Equal(t, test.wantCalls, test.store.removeCalls)
		})
	}
}
