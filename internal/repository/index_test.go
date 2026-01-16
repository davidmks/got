package repository

import (
	"os"
	"path/filepath"
	"testing"
)

func createTestIndex(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, indexFile)
	if err := os.WriteFile(path, []byte(content), filePerm); err != nil {
		t.Fatalf("Failed to create test index: %v", err)
	}
	return dir
}

func TestLoadIndex_Empty(t *testing.T) {
	dir := createTestIndex(t, "")

	index, err := LoadIndex(dir)
	if err != nil {
		t.Fatalf("LoadIndex() error = %v", err)
	}

	if len(index.Entries) != 0 {
		t.Errorf("Expected empty entries, got %d", len(index.Entries))
	}
}

func TestLoadIndex_WithEntries(t *testing.T) {
	dir := createTestIndex(t, `{"entries":[{"path":"test.txt","hash":"abc123"}]}`)

	index, err := LoadIndex(dir)
	if err != nil {
		t.Fatalf("LoadIndex() error = %v", err)
	}

	if len(index.Entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(index.Entries))
	}
	if index.Entries[0].Path != "test.txt" || index.Entries[0].Hash != "abc123" {
		t.Errorf("Entry = %+v, want {test.txt, abc123}", index.Entries[0])
	}
}

func TestIndex_SaveAndLoad(t *testing.T) {
	dir := createTestIndex(t, "")

	original := &Index{
		Entries: []IndexEntry{
			{Path: "file1.txt", Hash: "hash1"},
			{Path: "file2.txt", Hash: "hash2"},
		},
	}

	if err := original.Save(dir); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := LoadIndex(dir)
	if err != nil {
		t.Fatalf("LoadIndex() error = %v", err)
	}

	if len(loaded.Entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Path != "file1.txt" || loaded.Entries[0].Hash != "hash1" {
		t.Errorf("First entry = %+v, want {file1.txt, hash1}", loaded.Entries[0])
	}
}

func TestIndex_Add(t *testing.T) {
	t.Run("new entry", func(t *testing.T) {
		index := &Index{}
		index.Add("test.txt", "hash123")

		if len(index.Entries) != 1 {
			t.Fatalf("Expected 1 entry, got %d", len(index.Entries))
		}
		if index.Entries[0].Path != "test.txt" || index.Entries[0].Hash != "hash123" {
			t.Errorf("Entry = %+v, want {test.txt, hash123}", index.Entries[0])
		}
	})

	t.Run("update existing", func(t *testing.T) {
		index := &Index{
			Entries: []IndexEntry{{Path: "test.txt", Hash: "old"}},
		}
		index.Add("test.txt", "new")

		if len(index.Entries) != 1 {
			t.Fatalf("Expected 1 entry, got %d", len(index.Entries))
		}
		if index.Entries[0].Hash != "new" {
			t.Errorf("Hash = %q, want %q", index.Entries[0].Hash, "new")
		}
	})
}

func TestIndex_Get(t *testing.T) {
	index := &Index{
		Entries: []IndexEntry{{Path: "test.txt", Hash: "hash123"}},
	}

	hash, exists := index.Get("test.txt")
	if !exists || hash != "hash123" {
		t.Errorf("Get(test.txt) = (%q, %v), want (hash123, true)", hash, exists)
	}

	hash, exists = index.Get("nonexistent.txt")
	if exists {
		t.Errorf("Get(nonexistent.txt) = (%q, %v), want ('', false)", hash, exists)
	}
}

func TestIndex_Remove(t *testing.T) {
	index := &Index{
		Entries: []IndexEntry{
			{Path: "file1.txt", Hash: "hash1"},
			{Path: "file2.txt", Hash: "hash2"},
		},
	}

	if !index.Remove("file1.txt") {
		t.Error("Remove(file1.txt) = false, want true")
	}
	if len(index.Entries) != 1 || index.Entries[0].Path != "file2.txt" {
		t.Errorf("After remove: entries = %+v, want [{file2.txt, hash2}]", index.Entries)
	}

	if index.Remove("nonexistent.txt") {
		t.Error("Remove(nonexistent.txt) = true, want false")
	}
}
