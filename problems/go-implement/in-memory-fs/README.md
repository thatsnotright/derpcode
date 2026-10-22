# In-Memory File System

Implement a tiny in-memory file system with directories and text files.

Paths are absolute, slash-separated, and start with `/`. `/` is the root directory and always exists. Repeated or trailing slashes are ignored (`/a//b/` is the same as `/a/b`). A path that does not start with `/` (including the empty string) is invalid.

## Methods to implement

```go
var (
	ErrNotFound    = errors.New("no such file or directory")
	ErrNotDir      = errors.New("not a directory")
	ErrIsDir       = errors.New("is a directory")
	ErrInvalidPath = errors.New("invalid path")
)

func NewFileSystem() *FileSystem

// Mkdir creates the directory at path, creating missing parents (like mkdir -p).
// It succeeds if the directory already exists. Returns ErrNotDir if any
// component along the path is a file.
func (fs *FileSystem) Mkdir(path string) error

// WriteFile creates or overwrites the file at path. The parent directory must
// already exist. Returns ErrNotFound if the parent is missing, ErrNotDir if
// the parent (or any ancestor) is a file, and ErrIsDir if path is a directory.
func (fs *FileSystem) WriteFile(path, content string) error

// ReadFile returns the contents of the file at path. Returns ErrNotFound if
// missing, ErrIsDir if path is a directory, and ErrNotDir if an ancestor is a file.
func (fs *FileSystem) ReadFile(path string) (string, error)

// List returns the sorted names (not full paths) of the entries in the
// directory at path. If path is a file, it returns just that file's name.
func (fs *FileSystem) List(path string) ([]string, error)

// Remove deletes the file or directory at path; directories are removed
// recursively. The root cannot be removed (ErrInvalidPath).
func (fs *FileSystem) Remove(path string) error
```

## Example

```go
fs := NewFileSystem()
fs.Mkdir("/a/b")
fs.WriteFile("/a/b/hello.txt", "hi")
fs.WriteFile("/a/readme", "x")
fs.List("/a")                    // ["b", "readme"], nil
fs.List("/a/b/hello.txt")        // ["hello.txt"], nil
fs.ReadFile("/a/b/hello.txt")    // "hi", nil
fs.WriteFile("/nope/x", "y")     // ErrNotFound
fs.Remove("/a")                  // removes everything under /a
```

## Constraints

- Operations on a path should take time proportional to its depth (plus listing size for `List`).
- Single-threaded use is fine.
- Return the sentinel errors above directly or wrapped; tests use `errors.Is`.

## Follow-ups

1. Add `Move(src, dst string) error`. What edge cases come up (moving a directory into itself)?
2. Add `Find(root, name string) []string` returning full paths of matches, sorted.
3. Add a `Size(path string) int` that returns the total bytes under a path in O(1) amortized. How do you keep it updated on write and remove?
