package solution

import (
	"errors"
	"reflect"
	"testing"
)

func TestMkdirAndList(t *testing.T) {
	fs := NewFileSystem()
	if err := fs.Mkdir("/a/b/c"); err != nil {
		t.Fatalf("Mkdir = %v", err)
	}
	tests := []struct {
		path string
		want []string
	}{
		{"/", []string{"a"}},
		{"/a", []string{"b"}},
		{"/a/b", []string{"c"}},
		{"/a/b/c", []string{}},
	}
	for _, tt := range tests {
		got, err := fs.List(tt.path)
		if err != nil {
			t.Errorf("List(%q) error = %v", tt.path, err)
			continue
		}
		if len(got) != len(tt.want) || (len(got) > 0 && !reflect.DeepEqual(got, tt.want)) {
			t.Errorf("List(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestMkdirIdempotent(t *testing.T) {
	fs := NewFileSystem()
	fs.Mkdir("/a")
	fs.WriteFile("/a/f", "keep")
	if err := fs.Mkdir("/a"); err != nil {
		t.Errorf("Mkdir of existing dir = %v, want nil", err)
	}
	if got, _ := fs.ReadFile("/a/f"); got != "keep" {
		t.Errorf("existing contents lost: %q", got)
	}
	if err := fs.Mkdir("/"); err != nil {
		t.Errorf("Mkdir(/) = %v, want nil", err)
	}
}

func TestWriteAndReadFile(t *testing.T) {
	fs := NewFileSystem()
	fs.Mkdir("/docs")
	if err := fs.WriteFile("/docs/a.txt", "hello"); err != nil {
		t.Fatalf("WriteFile = %v", err)
	}
	if got, err := fs.ReadFile("/docs/a.txt"); err != nil || got != "hello" {
		t.Errorf("ReadFile = (%q, %v), want (hello, nil)", got, err)
	}
	if err := fs.WriteFile("/docs/a.txt", "bye"); err != nil {
		t.Fatalf("overwrite = %v", err)
	}
	if got, _ := fs.ReadFile("/docs/a.txt"); got != "bye" {
		t.Errorf("after overwrite got %q, want bye", got)
	}
	if err := fs.WriteFile("/top.txt", ""); err != nil {
		t.Errorf("write at root = %v", err)
	}
	if got, err := fs.ReadFile("/top.txt"); err != nil || got != "" {
		t.Errorf("empty file = (%q, %v)", got, err)
	}
}

func TestWriteFileErrors(t *testing.T) {
	fs := NewFileSystem()
	fs.Mkdir("/d")
	fs.WriteFile("/d/f", "x")
	tests := []struct {
		name string
		path string
		want error
	}{
		{"missing parent", "/nope/x", ErrNotFound},
		{"parent is file", "/d/f/x", ErrNotDir},
		{"path is dir", "/d", ErrIsDir},
		{"root", "/", ErrIsDir},
		{"relative", "d/x", ErrInvalidPath},
		{"empty", "", ErrInvalidPath},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fs.WriteFile(tt.path, "data"); !errors.Is(err, tt.want) {
				t.Errorf("WriteFile(%q) = %v, want %v", tt.path, err, tt.want)
			}
		})
	}
}

func TestReadFileErrors(t *testing.T) {
	fs := NewFileSystem()
	fs.Mkdir("/d")
	fs.WriteFile("/d/f", "x")
	tests := []struct {
		name string
		path string
		want error
	}{
		{"missing", "/d/missing", ErrNotFound},
		{"missing parent", "/zzz/f", ErrNotFound},
		{"directory", "/d", ErrIsDir},
		{"through file", "/d/f/x", ErrNotDir},
		{"relative", "d/f", ErrInvalidPath},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := fs.ReadFile(tt.path); !errors.Is(err, tt.want) {
				t.Errorf("ReadFile(%q) = %v, want %v", tt.path, err, tt.want)
			}
		})
	}
}

func TestMkdirThroughFile(t *testing.T) {
	fs := NewFileSystem()
	fs.WriteFile("/f", "x")
	if err := fs.Mkdir("/f/sub"); !errors.Is(err, ErrNotDir) {
		t.Errorf("Mkdir through file = %v, want ErrNotDir", err)
	}
	if err := fs.Mkdir("/f"); !errors.Is(err, ErrNotDir) {
		t.Errorf("Mkdir over file = %v, want ErrNotDir", err)
	}
}

func TestListSortedAndFileAndErrors(t *testing.T) {
	fs := NewFileSystem()
	fs.Mkdir("/d/zeta")
	fs.WriteFile("/d/beta.txt", "1")
	fs.WriteFile("/d/alpha.txt", "2")
	got, err := fs.List("/d")
	want := []string{"alpha.txt", "beta.txt", "zeta"}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("List(/d) = (%v, %v), want %v", got, err, want)
	}
	got, err = fs.List("/d/beta.txt")
	if err != nil || !reflect.DeepEqual(got, []string{"beta.txt"}) {
		t.Errorf("List(file) = (%v, %v), want [beta.txt]", got, err)
	}
	if _, err := fs.List("/missing"); !errors.Is(err, ErrNotFound) {
		t.Errorf("List(missing) = %v, want ErrNotFound", err)
	}
}

func TestPathNormalization(t *testing.T) {
	fs := NewFileSystem()
	fs.Mkdir("/a//b/")
	if err := fs.WriteFile("//a/b//f.txt/", "v"); err != nil {
		t.Fatalf("WriteFile = %v", err)
	}
	if got, err := fs.ReadFile("/a/b/f.txt"); err != nil || got != "v" {
		t.Errorf("ReadFile = (%q, %v), want (v, nil)", got, err)
	}
	if got, _ := fs.List("/a/b/"); !reflect.DeepEqual(got, []string{"f.txt"}) {
		t.Errorf("List = %v", got)
	}
}

func TestRemove(t *testing.T) {
	fs := NewFileSystem()
	fs.Mkdir("/a/b")
	fs.WriteFile("/a/b/f", "x")
	fs.WriteFile("/a/g", "y")

	if err := fs.Remove("/a/g"); err != nil {
		t.Fatalf("Remove file = %v", err)
	}
	if _, err := fs.ReadFile("/a/g"); !errors.Is(err, ErrNotFound) {
		t.Errorf("removed file still readable: %v", err)
	}
	if err := fs.Remove("/a"); err != nil {
		t.Fatalf("Remove dir = %v", err)
	}
	if _, err := fs.List("/a"); !errors.Is(err, ErrNotFound) {
		t.Errorf("removed dir still listed: %v", err)
	}
	if _, err := fs.ReadFile("/a/b/f"); !errors.Is(err, ErrNotFound) {
		t.Errorf("descendant still readable: %v", err)
	}
	if got, _ := fs.List("/"); len(got) != 0 {
		t.Errorf("root should be empty, got %v", got)
	}
}

func TestRemoveErrors(t *testing.T) {
	fs := NewFileSystem()
	fs.WriteFile("/f", "x")
	if err := fs.Remove("/"); !errors.Is(err, ErrInvalidPath) {
		t.Errorf("Remove(/) = %v, want ErrInvalidPath", err)
	}
	if err := fs.Remove("/missing"); !errors.Is(err, ErrNotFound) {
		t.Errorf("Remove(missing) = %v, want ErrNotFound", err)
	}
	if err := fs.Remove("/f/x"); !errors.Is(err, ErrNotDir) {
		t.Errorf("Remove(through file) = %v, want ErrNotDir", err)
	}
}

func TestRecreateAfterRemove(t *testing.T) {
	fs := NewFileSystem()
	fs.Mkdir("/a")
	fs.WriteFile("/a/f", "old")
	fs.Remove("/a")
	if err := fs.Mkdir("/a"); err != nil {
		t.Fatal(err)
	}
	if _, err := fs.ReadFile("/a/f"); !errors.Is(err, ErrNotFound) {
		t.Errorf("old file resurfaced: %v", err)
	}
}
