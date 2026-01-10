package objects

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	objectsDir = "objects"
	filePerm   = 0644
)

func objectPath(gotDir, hash string) string {
	return filepath.Join(gotDir, objectsDir, hash)
}

// Store writes content to the object store and returns its hash.
// If an object with the same hash already exists, it returns the hash
// without writing (content-addressable deduplication).
func Store(gotDir string, content []byte) (string, error) {
	hash := Hash(content)
	path := objectPath(gotDir, hash)

	if _, err := os.Stat(path); err == nil {
		return hash, nil
	}

	if err := os.WriteFile(path, content, filePerm); err != nil {
		return "", fmt.Errorf("failed to write object %s: %w", hash, err)
	}

	return hash, nil
}

// Load retrieves content from the object store by its hash.
// Returns an error if the object does not exist.
func Load(gotDir string, hash string) ([]byte, error) {
	path := objectPath(gotDir, hash)

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("object not found: %s", hash)
		}
		return nil, fmt.Errorf("failed to read object %s: %w", hash, err)
	}

	return content, nil
}

// Exists checks whether an object with the given hash exists.
func Exists(gotDir string, hash string) bool {
	path := objectPath(gotDir, hash)
	_, err := os.Stat(path)
	return err == nil
}
