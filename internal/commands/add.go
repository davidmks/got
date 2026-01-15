package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/davidmks/got/internal/objects"
	"github.com/davidmks/got/internal/repository"
)

func Add(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no files specified")
	}

	gotDir, err := repository.FindGotDir()
	if err != nil {
		return err
	}

	repoRoot := repository.GetRepoRoot(gotDir)

	index, err := repository.LoadIndex(gotDir)
	if err != nil {
		return err
	}

	for _, arg := range args {
		if err := addFile(gotDir, repoRoot, index, arg); err != nil {
			return err
		}
	}

	if err := index.Save(gotDir); err != nil {
		return err
	}

	return nil
}

func addFile(gotDir, repoRoot string, index *repository.Index, filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("invalid path %s: %w", filePath, err)
	}

	// Ensure file is within repository
	if !strings.HasPrefix(absPath, repoRoot+string(filepath.Separator)) {
		return fmt.Errorf("file is outside repository: %s", filePath)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filePath)
		}
		return fmt.Errorf("cannot access %s: %w", filePath, err)
	}

	if info.IsDir() {
		return fmt.Errorf("cannot add directory: %s", filePath)
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", filePath, err)
	}

	hash, err := objects.Store(gotDir, content)
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(repoRoot, absPath)
	if err != nil {
		return fmt.Errorf("cannot compute relative path: %w", err)
	}

	index.Add(relPath, hash)
	fmt.Printf("staged: %s\n", relPath)

	return nil
}
