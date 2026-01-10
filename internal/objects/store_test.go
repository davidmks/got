package objects

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestGotDir(t *testing.T) string {
	t.Helper()

	tempDir := t.TempDir()
	gotDir := filepath.Join(tempDir, ".got")
	objectsPath := filepath.Join(gotDir, "objects")

	if err := os.MkdirAll(objectsPath, 0755); err != nil {
		t.Fatalf("failed to create test .got directory: %v", err)
	}

	return gotDir
}

func TestStore(t *testing.T) {
	gotDir := setupTestGotDir(t)
	content := []byte("hello world")

	hash, err := Store(gotDir, content)
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}

	if len(hash) != HashLength {
		t.Errorf("hash length = %d, want %d", len(hash), HashLength)
	}

	objectFile := filepath.Join(gotDir, "objects", hash)
	if _, err := os.Stat(objectFile); os.IsNotExist(err) {
		t.Errorf("object file was not created at %s", objectFile)
	}

	stored, err := os.ReadFile(objectFile)
	if err != nil {
		t.Fatalf("failed to read stored object: %v", err)
	}
	if string(stored) != string(content) {
		t.Errorf("stored content = %q, want %q", stored, content)
	}
}

func TestStoreDeduplication(t *testing.T) {
	gotDir := setupTestGotDir(t)
	content := []byte("duplicate content")

	hash1, err := Store(gotDir, content)
	if err != nil {
		t.Fatalf("first Store() error = %v", err)
	}

	hash2, err := Store(gotDir, content)
	if err != nil {
		t.Fatalf("second Store() error = %v", err)
	}

	if hash1 != hash2 {
		t.Errorf("hashes differ: %s != %s", hash1, hash2)
	}

	files, err := os.ReadDir(filepath.Join(gotDir, "objects"))
	if err != nil {
		t.Fatalf("failed to read objects dir: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 object file, got %d", len(files))
	}
}

func TestLoad(t *testing.T) {
	gotDir := setupTestGotDir(t)
	content := []byte("test content for loading")

	hash, err := Store(gotDir, content)
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}

	loaded, err := Load(gotDir, hash)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if string(loaded) != string(content) {
		t.Errorf("Load() = %q, want %q", loaded, content)
	}
}

func TestLoadNotFound(t *testing.T) {
	gotDir := setupTestGotDir(t)
	// A valid-looking but non-existent hash
	fakeHash := "0000000000000000000000000000000000000000000000000000000000000000"

	_, err := Load(gotDir, fakeHash)
	if err == nil {
		t.Error("Load() should return error for non-existent object")
	}
}

func TestExists(t *testing.T) {
	gotDir := setupTestGotDir(t)
	content := []byte("existence test")

	hash := Hash(content)
	if Exists(gotDir, hash) {
		t.Error("Exists() = true before storing, want false")
	}

	_, err := Store(gotDir, content)
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}

	if !Exists(gotDir, hash) {
		t.Error("Exists() = false after storing, want true")
	}
}

func TestStoreBinaryContent(t *testing.T) {
	gotDir := setupTestGotDir(t)
	// Binary content with null bytes and high bytes
	content := []byte{0x00, 0x01, 0x02, 0xFE, 0xFF}

	hash, err := Store(gotDir, content)
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}

	loaded, err := Load(gotDir, hash)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(loaded) != len(content) {
		t.Fatalf("loaded length = %d, want %d", len(loaded), len(content))
	}

	for i := range content {
		if loaded[i] != content[i] {
			t.Errorf("byte %d: got %x, want %x", i, loaded[i], content[i])
		}
	}
}
