package importer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
)

const header = "id,name,email,age\n"

func TestImportHappyPath(t *testing.T) {
	s := NewStore()
	in := header + "1,Ada Lovelace,ada@example.com,36\n2,Alan Turing,alan@example.com,41\n"
	if err := Import(s, strings.NewReader(in)); err != nil {
		t.Fatal(err)
	}
	if s.Len() != 2 {
		t.Fatalf("Len = %d; want 2", s.Len())
	}
	got, ok := s.Get(2)
	want := User{ID: 2, Name: "Alan Turing", Email: "alan@example.com", Age: 41}
	if !ok || got != want {
		t.Fatalf("Get(2) = %+v, %v; want %+v", got, ok, want)
	}
}

func TestImportUpsertsExistingUsers(t *testing.T) {
	s := NewStore()
	s.Upsert(User{ID: 1, Name: "Old", Email: "old@example.com", Age: 1})
	if err := Import(s, strings.NewReader(header+"1,New,new@example.com,2\n")); err != nil {
		t.Fatal(err)
	}
	if u, _ := s.Get(1); u.Name != "New" || s.Len() != 1 {
		t.Fatalf("got %+v, Len %d; want the row to replace the old user", u, s.Len())
	}
}

func TestImportRejectsBadHeader(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{"no header", "1,Ada,ada@example.com,36\n2,Alan,alan@example.com,41\n"},
		{"wrong order", "name,id,email,age\nAda,1,ada@example.com,36\n"},
		{"missing column", "id,name,email\n1,Ada,ada@example.com\n"},
		{"empty input", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStore()
			err := Import(s, strings.NewReader(tt.in))
			if err == nil {
				t.Fatalf("Import succeeded; want a header error (store has %d users)", s.Len())
			}
			if !strings.Contains(err.Error(), "header") {
				t.Fatalf("err = %q; want it to mention the header", err)
			}
			if s.Len() != 0 {
				t.Fatalf("store has %d users after a rejected file", s.Len())
			}
		})
	}
}

func TestImportReportsLineNumbersAndWrapsCause(t *testing.T) {
	in := header +
		"1,Ada,ada@example.com,36\n" +
		"abc,Alan,alan@example.com,41\n"
	err := Import(NewStore(), strings.NewReader(in))
	if err == nil {
		t.Fatal("want an error for a non-numeric id")
	}
	if !strings.Contains(err.Error(), "line 3") {
		t.Fatalf("err = %q; want it to name line 3", err)
	}
	if !errors.Is(err, strconv.ErrSyntax) {
		t.Fatalf("err = %q; want it to wrap strconv.ErrSyntax", err)
	}

	in = header + "1,Ada,ada@example.com,old\n"
	err = Import(NewStore(), strings.NewReader(in))
	if err == nil || !strings.Contains(err.Error(), "line 2") {
		t.Fatalf("bad age: err = %v; want an error naming line 2", err)
	}
}

func TestImportHandlesQuotedFields(t *testing.T) {
	in := header + "1,\"Lovelace, Ada\",ada@example.com,36\n2,\"Multi\nLine\",m@example.com,5\n"
	s := NewStore()
	if err := Import(s, strings.NewReader(in)); err != nil {
		t.Fatal(err)
	}
	if u, _ := s.Get(1); u.Name != "Lovelace, Ada" {
		t.Fatalf("Name = %q; want the comma preserved", u.Name)
	}
	if u, _ := s.Get(2); u.Name != "Multi\nLine" {
		t.Fatalf("Name = %q; want the newline preserved", u.Name)
	}
}

func TestImportAcceptsLooseHeader(t *testing.T) {
	s := NewStore()
	if err := Import(s, strings.NewReader("ID, Name,EMAIL,Age\n1,Ada,ada@example.com,36\n")); err != nil {
		t.Fatal(err)
	}
	if s.Len() != 1 {
		t.Fatalf("Len = %d; want 1", s.Len())
	}
}

func TestImportReportsShortRows(t *testing.T) {
	in := header + "1,Ada,ada@example.com,36\n2,Alan,alan@example.com,41\n3,Grace\n"
	err := Import(NewStore(), strings.NewReader(in))
	if err == nil || !strings.Contains(err.Error(), "line 4") {
		t.Fatalf("err = %v; want an error naming line 4", err)
	}
}

func TestImportIsAllOrNothing(t *testing.T) {
	s := NewStore()
	s.Upsert(User{ID: 1, Name: "Before", Email: "b@example.com", Age: 9})
	in := header +
		"1,After,a@example.com,10\n" +
		"2,Alan,alan@example.com,41\n" +
		"3,Grace,grace@example.com,oops\n"
	if err := Import(s, strings.NewReader(in)); err == nil {
		t.Fatal("want an error")
	}
	if s.Len() != 1 {
		t.Fatalf("Len = %d; the failed import left rows behind", s.Len())
	}
	if u, _ := s.Get(1); u.Name != "Before" {
		t.Fatalf("user 1 = %+v; the failed import overwrote it", u)
	}
}

func TestImportHandlesVeryLongLines(t *testing.T) {
	bio := strings.Repeat("x", 200_000)
	in := header +
		"1,Ada,ada@example.com,36\n" +
		"2," + bio + ",big@example.com,41\n" +
		"3,Grace,grace@example.com,85\n"
	s := NewStore()
	if err := Import(s, strings.NewReader(in)); err != nil {
		t.Fatalf("Import: %v", err)
	}
	if s.Len() != 3 {
		t.Fatalf("Len = %d; want 3 (rows after the long line were silently dropped)", s.Len())
	}
	if u, _ := s.Get(2); len(u.Name) != len(bio) {
		t.Fatalf("long name was truncated to %d bytes", len(u.Name))
	}
}

type fakeFile struct {
	io.Reader
	fs *fakeFS
}

func (f *fakeFile) Close() error {
	f.fs.mu.Lock()
	defer f.fs.mu.Unlock()
	f.fs.open--
	return nil
}

type fakeFS struct {
	mu      sync.Mutex
	files   map[string]string
	open    int
	maxOpen int
}

func (fs *fakeFS) openFile(name string) (io.ReadCloser, error) {
	content, ok := fs.files[name]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
	}
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.open++
	if fs.open > fs.maxOpen {
		fs.maxOpen = fs.open
	}
	return &fakeFile{Reader: strings.NewReader(content), fs: fs}, nil
}

func useFakeFS(t *testing.T, files map[string]string) *fakeFS {
	t.Helper()
	fs := &fakeFS{files: files}
	old := openFile
	openFile = fs.openFile
	t.Cleanup(func() { openFile = old })
	return fs
}

func TestImportFilesClosesEachFileBeforeOpeningTheNext(t *testing.T) {
	files := map[string]string{}
	var paths []string
	for i := 1; i <= 20; i++ {
		name := fmt.Sprintf("users-%02d.csv", i)
		files[name] = header + fmt.Sprintf("%d,User %d,u%d@example.com,30\n", i, i, i)
		paths = append(paths, name)
	}
	fs := useFakeFS(t, files)

	s := NewStore()
	if err := ImportFiles(s, paths); err != nil {
		t.Fatal(err)
	}
	if s.Len() != 20 {
		t.Fatalf("Len = %d; want 20", s.Len())
	}
	if fs.maxOpen > 1 {
		t.Fatalf("up to %d files were open at once; want each closed before the next is opened", fs.maxOpen)
	}
	if fs.open != 0 {
		t.Fatalf("%d files still open after ImportFiles returned", fs.open)
	}
}

func TestImportFilesNamesTheBadFileAndAppliesNothing(t *testing.T) {
	useFakeFS(t, map[string]string{
		"a.csv": header + "1,Ada,ada@example.com,36\n",
		"b.csv": header + "2,Alan,alan@example.com,forty\n",
	})
	s := NewStore()
	err := ImportFiles(s, []string{"a.csv", "b.csv"})
	if err == nil {
		t.Fatal("want an error")
	}
	if !strings.Contains(err.Error(), "b.csv") || !strings.Contains(err.Error(), "line 2") {
		t.Fatalf("err = %q; want it to name b.csv and line 2", err)
	}
	if s.Len() != 0 {
		t.Fatalf("Len = %d; the good file was applied even though the batch failed", s.Len())
	}
}

func TestImportFilesMissingFile(t *testing.T) {
	useFakeFS(t, map[string]string{})
	err := ImportFiles(NewStore(), []string{"nope.csv"})
	if !errors.Is(err, os.ErrNotExist) || !strings.Contains(err.Error(), "nope.csv") {
		t.Fatalf("err = %v; want a wrapped ErrNotExist naming the file", err)
	}
}
