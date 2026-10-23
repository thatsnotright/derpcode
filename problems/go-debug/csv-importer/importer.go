// Package importer loads user records from CSV files into an in-memory store.
package importer

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

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
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[u.ID] = u
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
func ImportFiles(store *Store, paths []string) error {
	for _, p := range paths {
		f, err := openFile(p)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := Import(store, f); err != nil {
			return err
		}
	}
	return nil
}

// Import reads CSV records (id,name,email,age) from r and upserts them into store.
func Import(store *Store, r io.Reader) error {
	sc := bufio.NewScanner(r)
	sc.Scan() // skip header
	for sc.Scan() {
		fields := strings.Split(sc.Text(), ",")
		if len(fields) != 4 {
			return errors.New("invalid row")
		}
		id, _ := strconv.Atoi(fields[0])
		age, _ := strconv.Atoi(fields[3])
		store.Upsert(User{ID: id, Name: fields[1], Email: fields[2], Age: age})
	}
	return nil
}
