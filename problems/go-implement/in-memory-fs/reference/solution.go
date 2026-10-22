package solution

import (
	"errors"
	"sort"
	"strings"
)

var (
	ErrNotFound    = errors.New("no such file or directory")
	ErrNotDir      = errors.New("not a directory")
	ErrIsDir       = errors.New("is a directory")
	ErrInvalidPath = errors.New("invalid path")
)

type node struct {
	children map[string]*node // nil for files
	content  string
}

func newDir() *node { return &node{children: make(map[string]*node)} }

func (n *node) isDir() bool { return n.children != nil }

// FileSystem is an in-memory file system of directories and text files.
type FileSystem struct {
	root *node
}

// NewFileSystem returns a file system containing only the root directory.
func NewFileSystem() *FileSystem {
	return &FileSystem{root: newDir()}
}

// split validates path and returns its non-empty components.
func split(path string) ([]string, error) {
	if !strings.HasPrefix(path, "/") {
		return nil, ErrInvalidPath
	}
	var parts []string
	for _, p := range strings.Split(path, "/") {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts, nil
}

// walk follows parts from the root. It returns ErrNotFound for a missing
// component and ErrNotDir if it must descend through a file.
func (fs *FileSystem) walk(parts []string) (*node, error) {
	cur := fs.root
	for _, name := range parts {
		if !cur.isDir() {
			return nil, ErrNotDir
		}
		next, ok := cur.children[name]
		if !ok {
			return nil, ErrNotFound
		}
		cur = next
	}
	return cur, nil
}

// Mkdir creates the directory at path, creating missing parents (like mkdir -p).
// It succeeds if the directory already exists. Returns ErrNotDir if any
// component along the path is a file.
func (fs *FileSystem) Mkdir(path string) error {
	parts, err := split(path)
	if err != nil {
		return err
	}
	cur := fs.root
	for _, name := range parts {
		next, ok := cur.children[name]
		if !ok {
			next = newDir()
			cur.children[name] = next
		}
		if !next.isDir() {
			return ErrNotDir
		}
		cur = next
	}
	return nil
}

// WriteFile creates or overwrites the file at path. The parent directory must
// already exist. Returns ErrNotFound if the parent is missing, ErrNotDir if
// the parent (or any ancestor) is a file, and ErrIsDir if path is a directory.
func (fs *FileSystem) WriteFile(path, content string) error {
	parts, err := split(path)
	if err != nil {
		return err
	}
	if len(parts) == 0 {
		return ErrIsDir
	}
	parent, err := fs.walk(parts[:len(parts)-1])
	if err != nil {
		return err
	}
	if !parent.isDir() {
		return ErrNotDir
	}
	name := parts[len(parts)-1]
	if existing, ok := parent.children[name]; ok && existing.isDir() {
		return ErrIsDir
	}
	parent.children[name] = &node{content: content}
	return nil
}

// ReadFile returns the contents of the file at path. Returns ErrNotFound if
// missing, ErrIsDir if path is a directory, and ErrNotDir if an ancestor is a file.
func (fs *FileSystem) ReadFile(path string) (string, error) {
	parts, err := split(path)
	if err != nil {
		return "", err
	}
	n, err := fs.walk(parts)
	if err != nil {
		return "", err
	}
	if n.isDir() {
		return "", ErrIsDir
	}
	return n.content, nil
}

// List returns the sorted names (not full paths) of the entries in the
// directory at path. If path is a file, it returns just that file's name.
func (fs *FileSystem) List(path string) ([]string, error) {
	parts, err := split(path)
	if err != nil {
		return nil, err
	}
	n, err := fs.walk(parts)
	if err != nil {
		return nil, err
	}
	if !n.isDir() {
		return []string{parts[len(parts)-1]}, nil
	}
	names := make([]string, 0, len(n.children))
	for name := range n.children {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

// Remove deletes the file or directory at path; directories are removed
// recursively. The root cannot be removed (ErrInvalidPath).
func (fs *FileSystem) Remove(path string) error {
	parts, err := split(path)
	if err != nil {
		return err
	}
	if len(parts) == 0 {
		return ErrInvalidPath
	}
	parent, err := fs.walk(parts[:len(parts)-1])
	if err != nil {
		return err
	}
	if !parent.isDir() {
		return ErrNotDir
	}
	name := parts[len(parts)-1]
	if _, ok := parent.children[name]; !ok {
		return ErrNotFound
	}
	delete(parent.children, name)
	return nil
}
