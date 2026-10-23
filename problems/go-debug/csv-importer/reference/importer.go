// Package importer loads user records from CSV files into an in-memory store.
package importer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

// columns is the expected CSV header.
var columns = []string{"id", "name", "email", "age"}

// User is a single imported user.
type User struct {
	ID    int
	Name  string
	Email string
	Age   int
}

// Store is a concurrency-safe in-memory user store.
type Store struct {
	mu    sync.Mutex
	users map[int]User
}

// NewStore returns an empty Store.
func NewStore() *Store {
	return &Store{users: map[int]User{}}
}

// Upsert inserts u or replaces the existing user with the same ID.
func (s *Store) Upsert(u User) {
	s.UpsertAll([]User{u})
}

// UpsertAll applies all users under a single lock, so readers never see a partial batch.
func (s *Store) UpsertAll(users []User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range users {
		s.users[u.ID] = u
	}
}

// Get returns the user with the given ID.
func (s *Store) Get(id int) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[id]
	return u, ok
}

// Len returns the number of stored users.
func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.users)
}

// openFile is a variable so tests can substitute an in-memory filesystem.
var openFile = func(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

// ImportFiles imports each file in paths into store, in order.
//
// The import is all-or-nothing: if any file is invalid, the store is left untouched and the error
// names the offending file and line.
func ImportFiles(store *Store, paths []string) error {
	var all []User
	for _, p := range paths {
		users, err := parseFile(p)
		if err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		all = append(all, users...)
	}
	store.UpsertAll(all)
	return nil
}

func parseFile(path string) ([]User, error) {
	f, err := openFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parse(f)
}

// Import reads CSV records (id,name,email,age) from r and upserts them into store.
// If any record is invalid, nothing is imported.
func Import(store *Store, r io.Reader) error {
	users, err := parse(r)
	if err != nil {
		return err
	}
	store.UpsertAll(users)
	return nil
}

// parse reads and validates every record in r without touching any store.
func parse(r io.Reader) ([]User, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1 // checked explicitly below so the header error can say "header"

	rec, err := cr.Read()
	if errors.Is(err, io.EOF) {
		return nil, errors.New("line 1: missing header")
	}
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	if !headerMatches(rec) {
		return nil, fmt.Errorf("line 1: unexpected header %q, want %q", rec, columns)
	}
	cr.FieldsPerRecord = len(columns)

	var users []User
	for {
		rec, err := cr.Read()
		if errors.Is(err, io.EOF) {
			return users, nil
		}
		if err != nil {
			return nil, err // *csv.ParseError already carries the line number
		}
		line, _ := cr.FieldPos(0)
		u, err := parseUser(rec)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}
		users = append(users, u)
	}
}

func headerMatches(rec []string) bool {
	if len(rec) != len(columns) {
		return false
	}
	for i, name := range rec {
		if strings.ToLower(strings.TrimSpace(name)) != columns[i] {
			return false
		}
	}
	return true
}

func parseUser(rec []string) (User, error) {
	id, err := strconv.Atoi(rec[0])
	if err != nil {
		return User{}, fmt.Errorf("id: %w", err)
	}
	age, err := strconv.Atoi(rec[3])
	if err != nil {
		return User{}, fmt.Errorf("age: %w", err)
	}
	return User{ID: id, Name: rec[1], Email: rec[2], Age: age}, nil
}
