# got - Git Clone in Go - Learning Project Plan

## Project Overview

Build **got**, a simplified version-control system inspired by Git, implemented in Go. This project will teach fundamental Go concepts while creating a functional tool for tracking file changes and managing project history.

## Learning Objectives

- Master Go file I/O and filesystem operations
- Understand content-addressable storage and cryptographic hashing
- Practice working with Go's standard library
- Learn command-line application development
- Implement basic data structures and algorithms in Go

## Project Scope

### In Scope (MVP Features)

- Repository initialization
- File staging and committing
- Status inspection
- Commit history viewing
- Basic branching
- Switching between branches/commits

### Out of Scope (For This Project)

- Remote repository operations (push/pull/fetch)
- Merge operations
- Complex diff algorithms
- .gitignore functionality
- Advanced features (tags, stash, rebase, submodules)

## Architecture Design

### Directory Structure

```
.got/
├── objects/           # Content-addressable object storage
│   ├── [hash1]       # Blob objects (file contents)
│   ├── [hash2]       # Commit objects
│   └── ...
├── refs/
│   └── heads/        # Branch pointers (files containing commit hashes)
│       ├── main
│       └── [other-branches]
├── HEAD              # Pointer to current branch or commit
└── index             # Staging area (serialized file)
```

### Core Data Structures

**Blob Object**

- Represents file content
- Stored by content hash
- Format: raw file contents

**Commit Object**

- Contains metadata about a snapshot
- Fields:
    - Tree hash (root directory snapshot)
    - Parent commit hash(es)
    - Author name and timestamp
    - Commit message

**Tree Object** (Simplified)

- Maps filenames to blob hashes
- Can be a simple JSON structure for MVP

**Index/Staging Area**

- Tracks files staged for next commit
- Maps file paths to blob hashes

## Implementation Phases

### Phase 1: Foundation & Repository Initialization

**Goal:** Set up project structure and implement `init` command

**Tasks:**

1. Set up Go project with proper module structure
2. Implement basic CLI framework (using `flag` package or `cobra`)
3. Create `init` command that:
    - Creates `.got` directory structure
    - Initializes empty `objects/` and `refs/heads/` directories
    - Creates `HEAD` file pointing to `refs/heads/main`
    - Creates empty `index` file

**Key Go Concepts:**

- Project organization and package structure
- `os.Mkdir`, `os.MkdirAll` for directory creation
- File creation with `os.Create`
- Error handling patterns in Go

**Validation:**

- Running `got init` creates proper directory structure
- Can handle errors (e.g., already initialized repository)

---

### Phase 2: Object Storage & Hashing

**Goal:** Implement content-addressable storage system

**Tasks:**

1. Implement hash function (SHA-1 or SHA-256)
2. Create function to hash file contents and return hash string
3. Implement function to store content in `objects/` directory:
    - Read file content
    - Compute hash
    - Write content to `objects/[hash]`
4. Implement function to retrieve content by hash

**Key Go Concepts:**

- `crypto/sha1` or `crypto/sha256` packages
- `encoding/hex` for hash string conversion
- `io.ReadAll` and `os.ReadFile`
- `os.WriteFile` with proper permissions

**Validation:**

- Can hash content and get consistent hash values
- Can store and retrieve content by hash
- Same content produces same hash (deduplication)

---

### Phase 3: Staging Area (Index)

**Goal:** Implement `add` command to stage files

**Tasks:**

1. Design index file format (JSON or custom format)
2. Implement function to read current index
3. Implement function to write index back to disk
4. Create `add` command that:
    - Takes file path(s) as arguments
    - Hashes file content and stores as blob
    - Updates index with filename → blob hash mapping
5. Handle edge cases:
    - Non-existent files
    - Files outside repository
    - Multiple files at once

**Key Go Concepts:**

- `encoding/json` for serialization
- `filepath` package for path manipulation
- `filepath.Walk` for directory traversal (if adding directories)
- Map data structures

**Validation:**

- `got add file.txt` stages the file
- Index file correctly stores the mapping
- Can stage multiple files
- Proper error messages for invalid files

---

### Phase 4: Status Command

**Goal:** Show repository state (modified, staged, untracked files)

**Tasks:**

1. Implement working directory scanner
2. Compare working directory with:
    - Index (staged files)
    - Last commit (to find modifications)
3. Create `status` command that shows:
    - Files staged for commit (in index but changed from last commit)
    - Files modified but not staged (changed since added to index)
    - Untracked files (not in index at all)
4. Implement file comparison logic

**Key Go Concepts:**

- `os.ReadDir` for listing directories
- File comparison techniques
- String formatting and output
- Set operations (tracking what's in various states)

**Validation:**

- Shows correct categorization of files
- Updates properly after `add` operations
- Handles empty repository (no commits yet)

---

### Phase 5: Commit Creation

**Goal:** Implement `commit` command to save snapshots

**Tasks:**

1. Design commit object structure
2. Implement tree object creation:
    - Read index to get current staged files
    - Create tree structure (can be simplified as flat for MVP)
    - Hash and store tree object
3. Implement commit object creation:
    - Generate commit metadata (author, timestamp, message)
    - Reference tree hash and parent commit
    - Hash and store commit object
4. Create `commit` command that:
    - Takes commit message (e.g., `m "message"`)
    - Creates commit object from staged files
    - Updates current branch reference
    - Clears staging area
5. Handle edge cases:
    - No files staged
    - Empty commit message

**Key Go Concepts:**

- `time` package for timestamps
- Struct serialization
- File pointer updates (updating branch refs)
- Command-line flag parsing

**Validation:**

- Creates commit objects with correct structure
- Updates branch reference to new commit
- Subsequent commits have correct parent references
- Staging area is cleared after commit

---

### Phase 6: Commit History

**Goal:** Implement `log` command to view commit history

**Tasks:**

1. Implement function to read commit objects
2. Implement function to traverse commit history (follow parent links)
3. Create `log` command that:
    - Starts from current HEAD
    - Walks back through parent commits
    - Displays formatted commit information
4. Format output nicely (hash, author, date, message)

**Key Go Concepts:**

- Linked-list-like traversal
- String formatting (`fmt.Printf`)
- Struct deserialization
- Error handling for corrupt/missing objects

**Validation:**

- Shows commits in reverse chronological order
- Displays all commit metadata
- Stops at initial commit (no parent)

---

### Phase 7: Branch Management

**Goal:** Implement basic branch operations

**Tasks:**

1. Implement `branch` command (no args) to list branches:
    - Read all files in `refs/heads/`
    - Show current branch (from HEAD)
2. Implement `branch <name>` to create new branch:
    - Create new file in `refs/heads/<name>`
    - Point it to current commit
3. Implement branch deletion (optional enhancement)

**Key Go Concepts:**

- Directory listing and filtering
- File naming and path conventions
- Reading symbolic references

**Validation:**

- Can list all branches with current marked
- Creates new branches at current commit
- Branches are independent pointers

---

### Phase 8: Checkout (Branch Switching)

**Goal:** Implement `checkout` command to switch branches/commits

**Tasks:**

1. Implement working directory update function:
    - Clear current working directory (carefully!)
    - Restore files from commit's tree
2. Create `checkout <branch>` to switch branches:
    - Verify branch exists
    - Update HEAD to point to branch
    - Restore working directory to branch's commit
3. Create `checkout <commit-hash>` to view historical commits:
    - Put repository in "detached HEAD" state
    - Update HEAD to point directly to commit
    - Restore working directory
4. Safety checks:
    - Warn about uncommitted changes
    - Prevent data loss

**Key Go Concepts:**

- File deletion (`os.Remove`, `os.RemoveAll`)
- Working with file permissions
- State management
- User warnings and confirmations

**Validation:**

- Can switch between branches
- Working directory reflects branch state
- Can checkout specific commits
- HEAD file correctly updated

---

### Phase 9: Show Command

**Goal:** Display specific commit details

**Tasks:**

1. Implement `show <commit-hash>` command:
    - Parse commit object
    - Display commit metadata
    - Show files in that commit (from tree)
    - Optionally show basic diff from parent

**Key Go Concepts:**

- Object parsing and display
- Tree traversal
- Optional: simple diff algorithm

**Validation:**

- Shows commit information clearly
- Lists all files in commit
- Works with partial hash matching (optional)

---

### Phase 10: Polish & Documentation

**Goal:** Create a usable, documented tool

**Tasks:**

1. Add help text for all commands
2. Improve error messages
3. Write README with:
    - Installation instructions
    - Usage examples
    - Architecture explanation
4. Add tests for core functions
5. Handle edge cases gracefully

**Key Go Concepts:**

- Go testing (`testing` package)
- Documentation comments
- User experience considerations

---

## Technical Implementation Guidelines

### Recommended Go Packages

- **Standard Library:**
    - `os` - file operations
    - `io` - input/output utilities
    - `crypto/sha1` or `crypto/sha256` - hashing
    - `encoding/json` - serialization
    - `encoding/hex` - hash string conversion
    - `path/filepath` - path manipulation
    - `flag` - command-line parsing
    - `fmt` - formatted I/O
    - `time` - timestamps
    - `strings` - string utilities
- **Third-party (Optional):**
    - `github.com/spf13/cobra` - advanced CLI framework
    - `github.com/spf13/viper` - configuration management

### Code Organization Suggestion

```
got/
├── main.go                 # Entry point, CLI routing
├── cmd/
│   ├── init.go            # init command
│   ├── add.go             # add command
│   ├── commit.go          # commit command
│   ├── status.go          # status command
│   ├── log.go             # log command
│   ├── branch.go          # branch command
│   └── checkout.go        # checkout command
├── internal/
│   ├── objects/
│   │   ├── blob.go        # blob object handling
│   │   ├── commit.go      # commit object handling
│   │   └── tree.go        # tree object handling
│   ├── repository/
│   │   ├── repo.go        # repository operations
│   │   └── index.go       # staging area management
│   └── utils/
│       ├── hash.go        # hashing utilities
│       └── file.go        # file utilities
├── go.mod
└── README.md
```

### Best Practices

1. **Error Handling:** Always check errors, provide context
2. **File Permissions:** Use appropriate permissions (0644 for files, 0755 for directories)
3. **Path Handling:** Use `filepath.Join()` for cross-platform compatibility
4. **Testing:** Write tests for core logic (hashing, object storage)
5. **Documentation:** Add comments for exported functions
6. **Safety:** Validate user input and repository state

### Common Pitfalls to Avoid

- Don't modify working directory without safety checks
- Handle concurrent access carefully (file locking if needed)
- Properly close files after opening
- Use proper path separators (don't hardcode `/` or `\\`)
- Don't assume file encodings (handle binary files)

## Extension Ideas (After MVP)

Once core functionality works, consider:

1. **Config system** - user name, email configuration
2. **Diff viewer** - show changes between commits
3. **Simple merge** - fast-forward merges only
4. **File ignore patterns** - basic .gotignore
5. **Compression** - compress object storage
6. **Performance** - optimize for large repositories
7. **Remote operations** - basic clone/push/pull (advanced)

## Success Criteria

Your Git clone is successful when you can:

- [ ]  Initialize a repository
- [ ]  Stage files for commit
- [ ]  Create commits with messages
- [ ]  View repository status
- [ ]  See commit history
- [ ]  Create and list branches
- [ ]  Switch between branches
- [ ]  Restore old commits
- [ ]  Use it to track a real project's history

## Learning Resources

- **Go Documentation:** https://go.dev/doc/
- **Git Internals:** https://git-scm.com/book/en/v2/Git-Internals-Plumbing-and-Porcelain
- **Go by Example:** https://gobyexample.com/

## Final Notes

- Start simple, iterate gradually
- Test each phase before moving forward
- Don't worry about performance initially
- Focus on understanding the concepts
- It's okay to simplify (this isn't production Git!)
- Learn by doing - experiment and break things

Good luck with your project! 🚀
