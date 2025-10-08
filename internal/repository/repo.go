package repository

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	gotDir     = ".got"
	objectsDir = "objects"
	refsDir    = "refs"
	headsDir   = "heads"
	headFile   = "HEAD"
	indexFile  = "index"
	mainBranch = "refs/heads/main"

	dirPerm  = 0755 // rwxr-xr-x: owner can read/write/execute, others can read/execute
	filePerm = 0644 // rw-r--r--: owner can read/write, others can read only
)

func Initialize() error {
	if _, err := os.Stat(gotDir); err == nil {
		return fmt.Errorf("repository already initialized")
	}

	if err := os.Mkdir(gotDir, dirPerm); err != nil {
		return fmt.Errorf("failed to create .got directory: %w", err)
	}

	objectsPath := filepath.Join(gotDir, objectsDir)
	if err := os.Mkdir(objectsPath, dirPerm); err != nil {
		return fmt.Errorf("failed to create objects directory: %w", err)
	}

	headsPath := filepath.Join(gotDir, refsDir, headsDir)
	if err := os.MkdirAll(headsPath, dirPerm); err != nil {
		return fmt.Errorf("failed to create refs/heads directory: %w", err)
	}

	headPath := filepath.Join(gotDir, headFile)
	headContent := []byte("ref: " + mainBranch)
	if err := os.WriteFile(headPath, headContent, filePerm); err != nil {
		return fmt.Errorf("failed to create HEAD file: %w", err)
	}

	indexPath := filepath.Join(gotDir, indexFile)
	if err := os.WriteFile(indexPath, []byte{}, filePerm); err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}

	return nil
}
