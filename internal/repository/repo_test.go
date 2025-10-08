package repository

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitialize(t *testing.T) {
	tempDir := t.TempDir()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	if err := Initialize(); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	if _, err := os.Stat(gotDir); os.IsNotExist(err) {
		t.Errorf(".got directory was not created")
	}

	objectsPath := filepath.Join(gotDir, objectsDir)
	if _, err := os.Stat(objectsPath); os.IsNotExist(err) {
		t.Errorf(".got/objects directory was not created")
	}

	headsPath := filepath.Join(gotDir, refsDir, headsDir)
	if _, err := os.Stat(headsPath); os.IsNotExist(err) {
		t.Errorf(".got/refs/heads directory was not created")
	}

	headPath := filepath.Join(gotDir, headFile)
	content, err := os.ReadFile(headPath)
	if err != nil {
		t.Errorf("Failed to read HEAD file: %v", err)
	}
	expectedContent := "ref: " + mainBranch
	if string(content) != expectedContent {
		t.Errorf("HEAD content = %q, want %q", string(content), expectedContent)
	}

	indexPath := filepath.Join(gotDir, indexFile)
	indexContent, err := os.ReadFile(indexPath)
	if err != nil {
		t.Errorf("Failed to read index file: %v", err)
	}
	if len(indexContent) != 0 {
		t.Errorf("index file should be empty, got %d bytes", len(indexContent))
	}
}

func TestInitializeAlreadyExists(t *testing.T) {
	tempDir := t.TempDir()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	if err := Initialize(); err != nil {
		t.Fatalf("First Initialize() failed: %v", err)
	}

	err = Initialize()
	if err == nil {
		t.Error("Initialize() should return error when repository already exists")
	}
}
