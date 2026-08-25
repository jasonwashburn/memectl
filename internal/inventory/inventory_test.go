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

func TestAddSerializesConcurrentWriters(t *testing.T) {
	path := filepath.Join(t.TempDir(), "memes.json")
	first := Meme{Name: "first", TemplateID: "1", Texts: []string{"text"}, ImageURL: "https://first", PageURL: "https://first", CreatedAt: time.Now().UTC()}
	second := Meme{Name: "second", TemplateID: "2", Texts: []string{"text"}, ImageURL: "https://second", PageURL: "https://second", CreatedAt: time.Now().UTC()}
	errs := make(chan error, 2)
	start := make(chan struct{})
	for _, meme := range []Meme{first, second} {
		go func(meme Meme) {
			<-start
			errs <- New(path).Add(meme)
		}(meme)
	}
	close(start)
	require.NoError(t, <-errs)
	require.NoError(t, <-errs)
	memes, err := New(path).Load()
	require.NoError(t, err)
	SortByName(memes)
	assert.Equal(t, []Meme{first, second}, memes)
}

func TestRemove(t *testing.T) {
	first := testMeme("first")
	second := testMeme("second")
	third := testMeme("third")

	tests := []struct {
		name       string
		initial    []Meme
		remove     []string
		want       []Meme
		wantAbsent []string
	}{
		{name: "single record", initial: []Meme{first, second}, remove: []string{"first"}, want: []Meme{second}},
		{name: "multiple records", initial: []Meme{first, second, third}, remove: []string{"first", "third"}, want: []Meme{second}},
		{name: "final records", initial: []Meme{first, second}, remove: []string{"first", "second"}, want: []Meme{}},
		{name: "mixed present and absent", initial: []Meme{first, second}, remove: []string{"first", "missing"}, want: []Meme{second}, wantAbsent: []string{"missing"}},
		{name: "wholly absent", initial: []Meme{first}, remove: []string{"missing", "unknown"}, want: []Meme{first}, wantAbsent: []string{"missing", "unknown"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := New(filepath.Join(t.TempDir(), "memes.json"))
			for _, meme := range test.initial {
				require.NoError(t, store.Add(meme))
			}

			absent, err := store.Remove(test.remove)
			require.NoError(t, err)
			assert.Equal(t, test.wantAbsent, absent)
			got, err := store.Load()
			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestRemovePreservesInvalidState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "memes.json")
	contents := "{"
	require.NoError(t, os.WriteFile(path, []byte(contents), 0o600))

	_, err := New(path).Remove([]string{"meme"})
	require.Error(t, err)
	assert.ErrorContains(t, err, "malformed JSON")
	got, readErr := os.ReadFile(path)
	require.NoError(t, readErr)
	assert.Equal(t, contents, string(got))
}

func TestRemoveWriteFailurePreservesInventory(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "memes.json")
	store := New(path)
	meme := testMeme("meme")
	require.NoError(t, store.Add(meme))
	contents, err := os.ReadFile(path)
	require.NoError(t, err)
	require.NoError(t, os.Chmod(directory, 0o500))
	t.Cleanup(func() { _ = os.Chmod(directory, 0o700) })

	_, err = store.Remove([]string{"meme"})
	require.Error(t, err)
	got, readErr := os.ReadFile(path)
	require.NoError(t, readErr)
	assert.Equal(t, contents, got)
}

func testMeme(name string) Meme {
	return Meme{Name: name, TemplateID: "1", Texts: []string{"text"}, ImageURL: "https://image/" + name, PageURL: "https://page/" + name, CreatedAt: time.Date(2026, time.August, 24, 12, 0, 0, 0, time.UTC)}
}
