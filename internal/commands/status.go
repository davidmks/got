package commands

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/davidmks/got/internal/objects"
	"github.com/davidmks/got/internal/repository"
)

// StatusResult holds the categorized state of files in the working directory.
type StatusResult struct {
	Staged    []string
	Modified  []string
	Deleted   []string
	Untracked []string
}

// Status displays the working tree status compared to the index.
func Status() error {
	gotDir, err := repository.FindGotDir()
	if err != nil {
		return err
	}

	repoRoot := repository.GetRepoRoot(gotDir)

	index, err := repository.LoadIndex(gotDir)
	if err != nil {
		return err
	}

	workingFiles, err := collectWorkingFiles(repoRoot)
	if err != nil {
		return err
	}

	result, err := categorizeFiles(repoRoot, index, workingFiles)
	if err != nil {
		return err
	}

	printStatus(result)
	return nil
}

func collectWorkingFiles(repoRoot string) ([]string, error) {
	var files []string
	gotPath := filepath.Join(repoRoot, ".got")

	err := filepath.WalkDir(repoRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == gotPath {
			return filepath.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}

		files = append(files, relPath)
		return nil
	})

	return files, err
}

func categorizeFiles(repoRoot string, index *repository.Index, workingFiles []string) (*StatusResult, error) {
	result := &StatusResult{}

	indexMap := make(map[string]string)
	for _, entry := range index.Entries {
		indexMap[entry.Path] = entry.Hash
	}

	workingSet := make(map[string]bool)
	for _, f := range workingFiles {
		workingSet[f] = true
	}

	for _, file := range workingFiles {
		indexHash, inIndex := indexMap[file]
		if !inIndex {
			result.Untracked = append(result.Untracked, file)
			continue
		}

		currentHash, err := hashWorkingFile(repoRoot, file)
		if err != nil {
			return nil, fmt.Errorf("cannot hash %s: %w", file, err)
		}

		if currentHash != indexHash {
			result.Modified = append(result.Modified, file)
		}
	}

	for _, entry := range index.Entries {
		if !workingSet[entry.Path] {
			result.Deleted = append(result.Deleted, entry.Path)
		}
	}

	// All indexed files are "staged" since commits don't exist yet
	for _, entry := range index.Entries {
		result.Staged = append(result.Staged, entry.Path)
	}

	sort.Strings(result.Staged)
	sort.Strings(result.Modified)
	sort.Strings(result.Deleted)
	sort.Strings(result.Untracked)

	return result, nil
}

func hashWorkingFile(repoRoot, relPath string) (string, error) {
	absPath := filepath.Join(repoRoot, relPath)
	content, err := os.ReadFile(absPath)
	if err != nil {
		return "", err
	}
	return objects.Hash(content), nil
}

func printStatus(result *StatusResult) {
	fmt.Println("On branch main")
	fmt.Println()
	fmt.Println("No commits yet")
	fmt.Println()

	hasOutput := false

	if len(result.Staged) > 0 {
		hasOutput = true
		fmt.Println("Changes to be committed:")
		fmt.Println("  (use \"got rm --cached <file>...\" to unstage)")
		fmt.Println()
		for _, file := range result.Staged {
			fmt.Printf("\tnew file:   %s\n", file)
		}
		fmt.Println()
	}

	if len(result.Modified) > 0 || len(result.Deleted) > 0 {
		hasOutput = true
		fmt.Println("Changes not staged for commit:")
		fmt.Println("  (use \"got add <file>...\" to update what will be committed)")
		fmt.Println()
		for _, file := range result.Modified {
			fmt.Printf("\tmodified:   %s\n", file)
		}
		for _, file := range result.Deleted {
			fmt.Printf("\tdeleted:    %s\n", file)
		}
		fmt.Println()
	}

	if len(result.Untracked) > 0 {
		hasOutput = true
		fmt.Println("Untracked files:")
		fmt.Println("  (use \"got add <file>...\" to include in what will be committed)")
		fmt.Println()
		for _, file := range result.Untracked {
			fmt.Printf("\t%s\n", file)
		}
		fmt.Println()
	}

	if !hasOutput {
		fmt.Println("nothing to commit (create/copy files and use \"got add\" to track)")
	}
}
