package solution

import "errors"

// ErrNoTransaction is returned by Commit and Rollback when no transaction is open.
var ErrNoTransaction = errors.New("no transaction in progress")

// Store is an in-memory key-value store with nested transactions.
type Store struct {
	// TODO: add fields
}

// NewStore returns an empty store.
func NewStore() *Store {
	// TODO: implement
	return &Store{}
}

// Set stores value under key.
func (s *Store) Set(key, value string) {
	// TODO: implement
}

// Get returns the current value for key and whether it exists.
func (s *Store) Get(key string) (string, bool) {
	// TODO: implement
	return "", false
}

// Delete removes key. Deleting a missing key is a no-op.
func (s *Store) Delete(key string) {
	// TODO: implement
}

// Begin opens a new (possibly nested) transaction.
func (s *Store) Begin() {
	// TODO: implement
}

// Commit merges the innermost transaction into its parent (or into the store
// itself if it is the outermost). Returns ErrNoTransaction if none is open.
func (s *Store) Commit() error {
	// TODO: implement
	return nil
}

// Rollback discards the innermost transaction. Returns ErrNoTransaction if
// none is open.
func (s *Store) Rollback() error {
	// TODO: implement
	return nil
}
