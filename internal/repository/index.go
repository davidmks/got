package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type IndexEntry struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

type Index struct {
	Entries []IndexEntry `json:"entries"`
}

// LoadIndex reads the index file from gotDir.
func LoadIndex(gotDir string) (*Index, error) {
	indexPath := filepath.Join(gotDir, indexFile)

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read index file: %w", err)
	}

	// Empty file is valid (created by init)
	if len(data) == 0 {
		return &Index{Entries: []IndexEntry{}}, nil
	}

	var index Index
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("failed to parse index file: %w", err)
	}

	return &index, nil
}

// Save writes the index to disk.
func (idx *Index) Save(gotDir string) error {
	indexPath := filepath.Join(gotDir, indexFile)

	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize index: %w", err)
	}

	if err := os.WriteFile(indexPath, data, filePerm); err != nil {
		return fmt.Errorf("failed to write index file: %w", err)
	}

	return nil
}

// Add adds or updates an entry in the index.
func (idx *Index) Add(path, hash string) {
	for i, entry := range idx.Entries {
		if entry.Path == path {
			idx.Entries[i].Hash = hash
			return
		}
	}
	idx.Entries = append(idx.Entries, IndexEntry{Path: path, Hash: hash})
}

// Get returns the hash for a path, or ("", false) if not found.
func (idx *Index) Get(path string) (string, bool) {
	for _, entry := range idx.Entries {
		if entry.Path == path {
			return entry.Hash, true
		}
	}
	return "", false
}

// Remove removes an entry by path, returning true if it existed.
func (idx *Index) Remove(path string) bool {
	for i, entry := range idx.Entries {
		if entry.Path == path {
			idx.Entries[i] = idx.Entries[len(idx.Entries)-1]
			idx.Entries = idx.Entries[:len(idx.Entries)-1]
			return true
		}
	}
	return false
}
