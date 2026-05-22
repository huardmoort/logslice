// Package multifile provides support for processing multiple log files
// in sequence, merging their output as if they were a single stream.
package multifile

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// File represents a single log file entry with optional ordering weight.
type File struct {
	Path   string
	Weight int // lower weight = processed first
}

// MultiReader reads from multiple files sequentially.
type MultiReader struct {
	files   []File
	current io.ReadCloser
	idx     int
}

// New creates a MultiReader from a list of file paths.
// Files are processed in the order provided.
func New(paths []string) (*MultiReader, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("multifile: no files provided")
	}
	files := make([]File, len(paths))
	for i, p := range paths {
		files[i] = File{Path: p, Weight: i}
	}
	return &MultiReader{files: files}, nil
}

// NewSorted creates a MultiReader with files sorted by weight then path.
func NewSorted(files []File) (*MultiReader, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("multifile: no files provided")
	}
	sorted := make([]File, len(files))
	copy(sorted, files)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Weight != sorted[j].Weight {
			return sorted[i].Weight < sorted[j].Weight
		}
		return sorted[i].Path < sorted[j].Path
	})
	return &MultiReader{files: sorted}, nil
}

// Read implements io.Reader, transparently moving to the next file on EOF.
func (m *MultiReader) Read(p []byte) (int, error) {
	for {
		if m.current == nil {
			if m.idx >= len(m.files) {
				return 0, io.EOF
			}
			f, err := os.Open(m.files[m.idx].Path)
			if err != nil {
				return 0, fmt.Errorf("multifile: open %q: %w", m.files[m.idx].Path, err)
			}
			m.current = f
			m.idx++
		}
		n, err := m.current.Read(p)
		if err == io.EOF {
			m.current.Close()
			m.current = nil
			if n > 0 {
				return n, nil
			}
			continue
		}
		return n, err
	}
}

// Close releases any open file handle.
func (m *MultiReader) Close() error {
	if m.current != nil {
		err := m.current.Close()
		m.current = nil
		return err
	}
	return nil
}

// CurrentPath returns the path of the file currently being read.
func (m *MultiReader) CurrentPath() string {
	if m.idx == 0 || m.idx > len(m.files) {
		return ""
	}
	return m.files[m.idx-1].Path
}
