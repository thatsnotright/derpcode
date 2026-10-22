# Key-Value Store with Transactions

Implement an in-memory string key-value store that supports **nested transactions**.

- Outside a transaction, `Set` and `Delete` take effect immediately.
- `Begin` opens a new transaction nested inside the current one (if any).
- Changes made inside a transaction are visible to reads inside it, but are not final until committed.
- `Rollback` discards all changes made in the innermost open transaction.
- `Commit` merges the innermost transaction's changes into its parent transaction, or into the store itself if it was the outermost. Committing an inner transaction does not make its changes permanent if an outer transaction is later rolled back.

## Methods to implement

```go
var ErrNoTransaction = errors.New("no transaction in progress")

func NewStore() *Store

// Set stores value under key.
func (s *Store) Set(key, value string)

// Get returns the current value for key and whether it exists.
func (s *Store) Get(key string) (string, bool)

// Delete removes key. Deleting a missing key is a no-op.
func (s *Store) Delete(key string)

// Begin opens a new (possibly nested) transaction.
func (s *Store) Begin()

// Commit merges the innermost transaction into its parent.
// Returns ErrNoTransaction if none is open.
func (s *Store) Commit() error

// Rollback discards the innermost transaction.
// Returns ErrNoTransaction if none is open.
func (s *Store) Rollback() error
```

## Example

```go
s := NewStore()
s.Set("a", "1")
s.Begin()
s.Set("a", "2")
s.Delete("b")
s.Begin()
s.Set("a", "3")
s.Rollback()   // inner discarded
s.Get("a")     // "2", true
s.Commit()     // outer applied to the store
s.Get("a")     // "2", true
s.Rollback()   // ErrNoTransaction
```

## Constraints

- A key set to the empty string `""` exists; it is different from a deleted key.
- `Get`, `Set`, and `Delete` should not copy the whole store.
- Single-threaded use is fine.

## Follow-ups

1. Add `Count(value string) int` returning how many keys currently hold `value`, in O(1).
2. Make the store safe for concurrent use. What does isolation between two concurrent transactions mean here?
3. Add a `Keys()` method returning sorted keys. What is its cost with many nested transactions?
