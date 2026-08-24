package inventory

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "memes.json")
	meme := Meme{Name: "first-meme", TemplateID: "1", Texts: []string{"top", "bottom"}, ImageURL: "https://image", PageURL: "https://page", CreatedAt: time.Now().UTC().Truncate(time.Second)}
	store := New(path)
	require.NoError(t, store.Add(meme))
	got, err := store.Load()
	require.NoError(t, err)
	assert.Equal(t, []Meme{meme}, got)
	contents, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.NotContains(t, string(contents), "password")
}

func TestResolvePath(t *testing.T) {
	path, err := ResolvePath(func(string) string { return "" }, func() (string, error) { return "/home/meme", nil })
	require.NoError(t, err)
	assert.Equal(t, "/home/meme/.meme/memes.json", path)
	path, err = ResolvePath(func(string) string { return "/tmp/custom.json" }, os.UserHomeDir)
	require.NoError(t, err)
	assert.Equal(t, "/tmp/custom.json", path)
}

func TestLoadInvalidStatePreservesFile(t *testing.T) {
	for _, contents := range []string{"{", `{"version":2,"memes":[]}`, `{"version":1,"memes":[{"name":"BAD"}]}`} {
		t.Run(contents, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "memes.json")
			require.NoError(t, os.WriteFile(path, []byte(contents), 0o600))
			_, err := New(path).Load()
			require.Error(t, err)
			assert.ErrorContains(t, err, path)
			got, readErr := os.ReadFile(path)
			require.NoError(t, readErr)
			assert.Equal(t, contents, string(got))
		})
	}
}

func TestLoadMissingIsEmpty(t *testing.T) {
	memes, err := New(filepath.Join(t.TempDir(), "missing.json")).Load()
	require.NoError(t, err)
	assert.Empty(t, memes)
}

func TestLoadUnreadablePathReportsInventory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "inventory")
	require.NoError(t, os.Mkdir(path, 0o700))
	_, err := New(path).Load()
	require.Error(t, err)
	assert.ErrorContains(t, err, path)
}

func TestAddRejectsDuplicate(t *testing.T) {
	store := New(filepath.Join(t.TempDir(), "memes.json"))
	meme := Meme{Name: "meme", TemplateID: "1", Texts: []string{"text"}, ImageURL: "https://image", PageURL: "https://page", CreatedAt: time.Now().UTC()}
	require.NoError(t, store.Add(meme))
	require.Error(t, store.Add(meme))
	memes, err := store.Load()
	require.NoError(t, err)
	assert.Equal(t, []Meme{meme}, memes)
}
