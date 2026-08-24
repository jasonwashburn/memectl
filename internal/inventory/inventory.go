// Package inventory manages the local meme inventory.
package inventory

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

const version = 1

var namePattern = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)

// Meme is a locally managed Imgflip creation.
type Meme struct {
	Name       string    `json:"name"`
	TemplateID string    `json:"templateID"`
	Texts      []string  `json:"texts"`
	ImageURL   string    `json:"imageURL"`
	PageURL    string    `json:"pageURL"`
	CreatedAt  time.Time `json:"createdAt"`
}

type document struct {
	Version int    `json:"version"`
	Memes   []Meme `json:"memes"`
}

// Store reads and appends managed memes.
type Store struct {
	path string
}

// New returns a Store for path.
func New(path string) *Store {
	return &Store{path: path}
}

// Path returns the inventory file path.
func (s *Store) Path() string {
	return s.path
}

// ResolvePath returns MEME_STORE when set, otherwise ~/.meme/memes.json.
func ResolvePath(getenv func(string) string, homeDir func() (string, error)) (string, error) {
	if path := getenv("MEME_STORE"); path != "" {
		return path, nil
	}
	home, err := homeDir()
	if err != nil {
		return "", fmt.Errorf("resolve meme inventory home: %w", err)
	}
	return filepath.Join(home, ".meme", "memes.json"), nil
}

// Load returns all stored memes. A missing inventory is empty.
func (s *Store) Load() ([]Meme, error) {
	contents, err := os.ReadFile(s.path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read meme inventory %q: %w", s.path, err)
	}

	var doc document
	if err := json.Unmarshal(contents, &doc); err != nil {
		return nil, fmt.Errorf("read meme inventory %q: malformed JSON: %w", s.path, err)
	}
	if doc.Version != version {
		return nil, fmt.Errorf("read meme inventory %q: unsupported version %d", s.path, doc.Version)
	}
	for _, meme := range doc.Memes {
		if err := validateMeme(meme); err != nil {
			return nil, fmt.Errorf("read meme inventory %q: invalid record: %w", s.path, err)
		}
	}
	return doc.Memes, nil
}

// Add appends meme unless its name is already present.
func (s *Store) Add(meme Meme) error {
	if err := validateMeme(meme); err != nil {
		return fmt.Errorf("save meme inventory %q: invalid record: %w", s.path, err)
	}
	memes, err := s.Load()
	if err != nil {
		return err
	}
	if contains(memes, meme.Name) {
		return fmt.Errorf("save meme inventory %q: meme %q already exists", s.path, meme.Name)
	}
	memes = append(memes, meme)
	return s.write(document{Version: version, Memes: memes})
}

func (s *Store) write(doc document) error {
	contents, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("serialize meme inventory %q: %w", s.path, err)
	}
	contents = append(contents, '\n')
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return fmt.Errorf("create meme inventory directory for %q: %w", s.path, err)
	}
	temporary, err := os.CreateTemp(filepath.Dir(s.path), ".memes-*")
	if err != nil {
		return fmt.Errorf("create temporary meme inventory for %q: %w", s.path, err)
	}
	temporaryPath := temporary.Name()
	defer func() { _ = os.Remove(temporaryPath) }()
	if err := temporary.Chmod(0o600); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("set temporary meme inventory permissions for %q: %w", s.path, err)
	}
	if _, err := temporary.Write(contents); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("write meme inventory %q: %w", s.path, err)
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("sync meme inventory %q: %w", s.path, err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close meme inventory %q: %w", s.path, err)
	}
	if err := os.Rename(temporaryPath, s.path); err != nil {
		return fmt.Errorf("replace meme inventory %q: %w", s.path, err)
	}
	return nil
}

// ValidName reports whether name is valid for a managed meme.
func ValidName(name string) bool {
	return len(name) <= 63 && namePattern.MatchString(name)
}

func validateMeme(meme Meme) error {
	if !ValidName(meme.Name) {
		return fmt.Errorf("name %q must be a DNS-label-like value", meme.Name)
	}
	if meme.TemplateID == "" || len(meme.Texts) == 0 || meme.ImageURL == "" || meme.PageURL == "" || meme.CreatedAt.IsZero() || meme.CreatedAt.Location() != time.UTC {
		return errors.New("name, template ID, texts, URLs, and UTC creation time are required")
	}
	return nil
}

// Contains reports whether memes contains name.
func Contains(memes []Meme, name string) bool {
	return contains(memes, name)
}

func contains(memes []Meme, name string) bool {
	for _, meme := range memes {
		if meme.Name == name {
			return true
		}
	}
	return false
}

// SortByName sorts memes by local name.
func SortByName(memes []Meme) {
	sort.Slice(memes, func(i, j int) bool { return memes[i].Name < memes[j].Name })
}
