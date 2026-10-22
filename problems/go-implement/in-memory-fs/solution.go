package solution

import "errors"

var (
	ErrNotFound    = errors.New("no such file or directory")
	ErrNotDir      = errors.New("not a directory")
	ErrIsDir       = errors.New("is a directory")
	ErrInvalidPath = errors.New("invalid path")
)

// FileSystem is an in-memory file system of directories and text files.
type FileSystem struct {
	// TODO: add fields
}

// NewFileSystem returns a file system containing only the root directory.
func NewFileSystem() *FileSystem {
	// TODO: implement
	return &FileSystem{}
}

// Mkdir creates the directory at path, creating missing parents (like mkdir -p).
// It succeeds if the directory already exists. Returns ErrNotDir if any
// component along the path is a file.
func (fs *FileSystem) Mkdir(path string) error {
	// TODO: implement
	return nil
}

// WriteFile creates or overwrites the file at path. The parent directory must
// already exist. Returns ErrNotFound if the parent is missing, ErrNotDir if
// the parent (or any ancestor) is a file, and ErrIsDir if path is a directory.
func (fs *FileSystem) WriteFile(path, content string) error {
	// TODO: implement
	return nil
}

// ReadFile returns the contents of the file at path. Returns ErrNotFound if
// missing, ErrIsDir if path is a directory, and ErrNotDir if an ancestor is a file.
func (fs *FileSystem) ReadFile(path string) (string, error) {
	// TODO: implement
	return "", nil
}

// List returns the sorted names (not full paths) of the entries in the
// directory at path. If path is a file, it returns just that file's name.
func (fs *FileSystem) List(path string) ([]string, error) {
	// TODO: implement
	return nil, nil
}

// Remove deletes the file or directory at path; directories are removed
// recursively. The root cannot be removed (ErrInvalidPath).
func (fs *FileSystem) Remove(path string) error {
	// TODO: implement
	return nil
}
