package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/davidmks/got/internal/objects"
	"github.com/davidmks/got/internal/repository"
)

func setupTestRepo(t *testing.T) (repoRoot string, gotDir string) {
	t.Helper()

	repoRoot = t.TempDir()
	gotDir = filepath.Join(repoRoot, ".got")
	objectsDir := filepath.Join(gotDir, "objects")

	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		t.Fatalf("failed to create .got/objects: %v", err)
	}

	indexPath := filepath.Join(gotDir, "index")
	if err := os.WriteFile(indexPath, []byte(""), 0644); err != nil {
		t.Fatalf("failed to create index: %v", err)
	}

	return repoRoot, gotDir
}

func createFile(t *testing.T, path, content string) {
	t.Helper()

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create directory %s: %v", dir, err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create file %s: %v", path, err)
	}
}

func TestCollectWorkingFiles(t *testing.T) {
	t.Run("empty repo", func(t *testing.T) {
		repoRoot, _ := setupTestRepo(t)

		files, err := collectWorkingFiles(repoRoot)
		if err != nil {
			t.Fatalf("collectWorkingFiles() error = %v", err)
		}

		if len(files) != 0 {
			t.Errorf("expected 0 files, got %d: %v", len(files), files)
		}
	})

	t.Run("excludes .got directory", func(t *testing.T) {
		repoRoot, gotDir := setupTestRepo(t)

		// Create a file inside .got that should be ignored
		createFile(t, filepath.Join(gotDir, "objects", "somehash"), "object content")
		// Create a regular file that should be included
		createFile(t, filepath.Join(repoRoot, "file.txt"), "hello")

		files, err := collectWorkingFiles(repoRoot)
		if err != nil {
			t.Fatalf("collectWorkingFiles() error = %v", err)
		}

		if len(files) != 1 {
			t.Fatalf("expected 1 file, got %d: %v", len(files), files)
		}
		if files[0] != "file.txt" {
			t.Errorf("expected file.txt, got %s", files[0])
		}
	})

	t.Run("finds nested files", func(t *testing.T) {
		repoRoot, _ := setupTestRepo(t)

		createFile(t, filepath.Join(repoRoot, "root.txt"), "root")
		createFile(t, filepath.Join(repoRoot, "src", "main.go"), "package main")
		createFile(t, filepath.Join(repoRoot, "src", "lib", "util.go"), "package lib")

		files, err := collectWorkingFiles(repoRoot)
		if err != nil {
			t.Fatalf("collectWorkingFiles() error = %v", err)
		}

		if len(files) != 3 {
			t.Fatalf("expected 3 files, got %d: %v", len(files), files)
		}

		expected := map[string]bool{
			"root.txt":                             true,
			filepath.Join("src", "main.go"):        true,
			filepath.Join("src", "lib", "util.go"): true,
		}
		for _, f := range files {
			if !expected[f] {
				t.Errorf("unexpected file: %s", f)
			}
		}
	})
}

func TestCategorizeFiles(t *testing.T) {
	t.Run("untracked files", func(t *testing.T) {
		repoRoot, gotDir := setupTestRepo(t)
		createFile(t, filepath.Join(repoRoot, "new.txt"), "new content")

		index, err := repository.LoadIndex(gotDir)
		if err != nil {
			t.Fatalf("LoadIndex() error = %v", err)
		}
		workingFiles := []string{"new.txt"}

		result, err := categorizeFiles(repoRoot, index, workingFiles)
		if err != nil {
			t.Fatalf("categorizeFiles() error = %v", err)
		}

		if len(result.Untracked) != 1 || result.Untracked[0] != "new.txt" {
			t.Errorf("Untracked = %v, want [new.txt]", result.Untracked)
		}
		if len(result.Staged) != 0 {
			t.Errorf("Staged = %v, want []", result.Staged)
		}
	})

	t.Run("staged files", func(t *testing.T) {
		repoRoot, gotDir := setupTestRepo(t)
		content := "staged content"
		createFile(t, filepath.Join(repoRoot, "staged.txt"), content)

		index, err := repository.LoadIndex(gotDir)
		if err != nil {
			t.Fatalf("LoadIndex() error = %v", err)
		}
		index.Add("staged.txt", objects.Hash([]byte(content)))

		workingFiles := []string{"staged.txt"}

		result, err := categorizeFiles(repoRoot, index, workingFiles)
		if err != nil {
			t.Fatalf("categorizeFiles() error = %v", err)
		}

		if len(result.Staged) != 1 || result.Staged[0] != "staged.txt" {
			t.Errorf("Staged = %v, want [staged.txt]", result.Staged)
		}
		if len(result.Modified) != 0 {
			t.Errorf("Modified = %v, want []", result.Modified)
		}
	})

	t.Run("modified files", func(t *testing.T) {
		repoRoot, gotDir := setupTestRepo(t)
		createFile(t, filepath.Join(repoRoot, "file.txt"), "modified content")

		index, err := repository.LoadIndex(gotDir)
		if err != nil {
			t.Fatalf("LoadIndex() error = %v", err)
		}
		index.Add("file.txt", objects.Hash([]byte("original content")))

		workingFiles := []string{"file.txt"}

		result, err := categorizeFiles(repoRoot, index, workingFiles)
		if err != nil {
			t.Fatalf("categorizeFiles() error = %v", err)
		}

		if len(result.Modified) != 1 || result.Modified[0] != "file.txt" {
			t.Errorf("Modified = %v, want [file.txt]", result.Modified)
		}
		// File is still in index, so it's also staged
		if len(result.Staged) != 1 {
			t.Errorf("Staged = %v, want [file.txt]", result.Staged)
		}
	})

	t.Run("deleted files", func(t *testing.T) {
		repoRoot, gotDir := setupTestRepo(t)

		index, err := repository.LoadIndex(gotDir)
		if err != nil {
			t.Fatalf("LoadIndex() error = %v", err)
		}
		index.Add("deleted.txt", "somehash")

		// File is in index but not in working directory
		workingFiles := []string{}

		result, err := categorizeFiles(repoRoot, index, workingFiles)
		if err != nil {
			t.Fatalf("categorizeFiles() error = %v", err)
		}

		if len(result.Deleted) != 1 || result.Deleted[0] != "deleted.txt" {
			t.Errorf("Deleted = %v, want [deleted.txt]", result.Deleted)
		}
	})

	t.Run("mixed state", func(t *testing.T) {
		repoRoot, gotDir := setupTestRepo(t)

		stagedContent := "staged"
		modifiedContent := "modified now"

		createFile(t, filepath.Join(repoRoot, "staged.txt"), stagedContent)
		createFile(t, filepath.Join(repoRoot, "modified.txt"), modifiedContent)
		createFile(t, filepath.Join(repoRoot, "untracked.txt"), "untracked")

		index, err := repository.LoadIndex(gotDir)
		if err != nil {
			t.Fatalf("LoadIndex() error = %v", err)
		}
		index.Add("staged.txt", objects.Hash([]byte(stagedContent)))
		index.Add("modified.txt", objects.Hash([]byte("original")))
		index.Add("deleted.txt", "somehash")

		workingFiles := []string{"staged.txt", "modified.txt", "untracked.txt"}

		result, err := categorizeFiles(repoRoot, index, workingFiles)
		if err != nil {
			t.Fatalf("categorizeFiles() error = %v", err)
		}

		if len(result.Staged) != 3 {
			t.Errorf("Staged count = %d, want 3", len(result.Staged))
		}
		if len(result.Modified) != 1 || result.Modified[0] != "modified.txt" {
			t.Errorf("Modified = %v, want [modified.txt]", result.Modified)
		}
		if len(result.Deleted) != 1 || result.Deleted[0] != "deleted.txt" {
			t.Errorf("Deleted = %v, want [deleted.txt]", result.Deleted)
		}
		if len(result.Untracked) != 1 || result.Untracked[0] != "untracked.txt" {
			t.Errorf("Untracked = %v, want [untracked.txt]", result.Untracked)
		}
	})
}
