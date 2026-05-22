package multifile

import (
	"fmt"
	"path/filepath"
	"sort"
)

// GlobFiles expands one or more glob patterns into an ordered list of File
// entries. Duplicates are removed; files are sorted lexicographically.
func GlobFiles(patterns []string) ([]File, error) {
	if len(patterns) == 0 {
		return nil, fmt.Errorf("multifile: no patterns provided")
	}

	seen := make(map[string]struct{})
	var files []File

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("multifile: invalid glob pattern %q: %w", pattern, err)
		}
		for _, m := range matches {
			if _, dup := seen[m]; dup {
				continue
			}
			seen[m] = struct{}{}
			files = append(files, File{Path: m})
		}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("multifile: no files matched patterns %v", patterns)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	return files, nil
}
