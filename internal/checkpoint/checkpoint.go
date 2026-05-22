// Package checkpoint provides resume support for large log file processing
// by persisting the last successfully processed byte offset to disk.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// ErrNotFound is returned when no checkpoint file exists for the given key.
var ErrNotFound = errors.New("checkpoint: no checkpoint found")

// Record holds the persisted state for a log file processing run.
type Record struct {
	FilePath  string    `json:"file_path"`
	Offset    int64     `json:"offset"`
	LineCount int64     `json:"line_count"`
	SavedAt   time.Time `json:"saved_at"`
}

// Store persists and retrieves checkpoint records.
type Store struct {
	dir string
}

// NewStore creates a Store that writes checkpoint files under dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}

// Save writes rec to disk, keyed by rec.FilePath.
func (s *Store) Save(rec Record) error {
	rec.SavedAt = time.Now().UTC()
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(rec.FilePath), data, 0o644)
}

// Load retrieves the checkpoint for filePath.
// Returns ErrNotFound if no checkpoint exists.
func (s *Store) Load(filePath string) (Record, error) {
	data, err := os.ReadFile(s.path(filePath))
	if errors.Is(err, os.ErrNotExist) {
		return Record{}, ErrNotFound
	}
	if err != nil {
		return Record{}, err
	}
	var rec Record
	if err := json.Unmarshal(data, &rec); err != nil {
		return Record{}, err
	}
	return rec, nil
}

// Delete removes the checkpoint for filePath, if it exists.
func (s *Store) Delete(filePath string) error {
	err := os.Remove(s.path(filePath))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// path returns the checkpoint file path for a given log file path.
func (s *Store) path(filePath string) string {
	// Use a simple hash-like encoding: replace path separators.
	safe := sanitize(filePath)
	return s.dir + "/" + safe + ".json"
}

// sanitize replaces characters that are unsafe in filenames.
func sanitize(p string) string {
	out := make([]byte, len(p))
	for i := range len(p) {
		c := p[i]
		if c == '/' || c == '\\' || c == ':' {
			out[i] = '_'
		} else {
			out[i] = c
		}
	}
	return string(out)
}
