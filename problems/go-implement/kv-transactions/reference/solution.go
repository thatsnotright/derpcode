package solution

import "errors"

// ErrNoTransaction is returned by Commit and Rollback when no transaction is open.
var ErrNoTransaction = errors.New("no transaction in progress")

// change records a pending write; deleted marks a tombstone.
type change struct {
	value   string
	deleted bool
}

// Store is an in-memory key-value store with nested transactions.
type Store struct {
	data map[string]string
	txs  []map[string]change // innermost transaction is last
}

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{data: make(map[string]string)}
}

// Set stores value under key.
func (s *Store) Set(key, value string) {
	if tx := s.top(); tx != nil {
		tx[key] = change{value: value}
		return
	}
	s.data[key] = value
}

// Get returns the current value for key and whether it exists.
func (s *Store) Get(key string) (string, bool) {
	for i := len(s.txs) - 1; i >= 0; i-- {
		if c, ok := s.txs[i][key]; ok {
			if c.deleted {
				return "", false
			}
			return c.value, true
		}
	}
	v, ok := s.data[key]
	return v, ok
}

// Delete removes key. Deleting a missing key is a no-op.
func (s *Store) Delete(key string) {
	if tx := s.top(); tx != nil {
		tx[key] = change{deleted: true}
		return
	}
	delete(s.data, key)
}

// Begin opens a new (possibly nested) transaction.
func (s *Store) Begin() {
	s.txs = append(s.txs, make(map[string]change))
}

// Commit merges the innermost transaction into its parent (or into the store
// itself if it is the outermost). Returns ErrNoTransaction if none is open.
func (s *Store) Commit() error {
	tx := s.top()
	if tx == nil {
		return ErrNoTransaction
	}
	s.txs = s.txs[:len(s.txs)-1]
	parent := s.top()
	for k, c := range tx {
		switch {
		case parent != nil:
			parent[k] = c
		case c.deleted:
			delete(s.data, k)
		default:
			s.data[k] = c.value
		}
	}
	return nil
}

// Rollback discards the innermost transaction. Returns ErrNoTransaction if
// none is open.
func (s *Store) Rollback() error {
	if s.top() == nil {
		return ErrNoTransaction
	}
	s.txs = s.txs[:len(s.txs)-1]
	return nil
}

func (s *Store) top() map[string]change {
	if len(s.txs) == 0 {
		return nil
	}
	return s.txs[len(s.txs)-1]
}
