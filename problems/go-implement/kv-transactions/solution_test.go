package solution

import (
	"errors"
	"testing"
)

func mustGet(t *testing.T, s *Store, key, want string) {
	t.Helper()
	got, ok := s.Get(key)
	if !ok || got != want {
		t.Errorf("Get(%q) = (%q, %v), want (%q, true)", key, got, ok, want)
	}
}

func mustMiss(t *testing.T, s *Store, key string) {
	t.Helper()
	if got, ok := s.Get(key); ok {
		t.Errorf("Get(%q) = (%q, true), want missing", key, got)
	}
}

func TestSetGetDelete(t *testing.T) {
	s := NewStore()
	mustMiss(t, s, "a")
	s.Set("a", "1")
	mustGet(t, s, "a", "1")
	s.Set("a", "2")
	mustGet(t, s, "a", "2")
	s.Delete("a")
	mustMiss(t, s, "a")
	s.Delete("never-existed") // must not panic
}

func TestEmptyStringValueExists(t *testing.T) {
	s := NewStore()
	s.Set("k", "")
	mustGet(t, s, "k", "")
}

func TestNoTransactionErrors(t *testing.T) {
	s := NewStore()
	if err := s.Commit(); !errors.Is(err, ErrNoTransaction) {
		t.Errorf("Commit() = %v, want ErrNoTransaction", err)
	}
	if err := s.Rollback(); !errors.Is(err, ErrNoTransaction) {
		t.Errorf("Rollback() = %v, want ErrNoTransaction", err)
	}
}

func TestReadsSeeUncommittedChangesInsideTransaction(t *testing.T) {
	s := NewStore()
	s.Set("a", "1")
	s.Begin()
	s.Set("a", "2")
	s.Set("b", "x")
	s.Delete("a")
	mustMiss(t, s, "a")
	mustGet(t, s, "b", "x")
}

func TestRollbackRestoresState(t *testing.T) {
	s := NewStore()
	s.Set("a", "1")
	s.Set("b", "2")
	s.Begin()
	s.Set("a", "changed")
	s.Delete("b")
	s.Set("c", "new")
	if err := s.Rollback(); err != nil {
		t.Fatalf("Rollback() = %v", err)
	}
	mustGet(t, s, "a", "1")
	mustGet(t, s, "b", "2")
	mustMiss(t, s, "c")
}

func TestCommitPersists(t *testing.T) {
	s := NewStore()
	s.Set("a", "1")
	s.Set("b", "2")
	s.Begin()
	s.Set("a", "changed")
	s.Delete("b")
	s.Set("c", "new")
	if err := s.Commit(); err != nil {
		t.Fatalf("Commit() = %v", err)
	}
	mustGet(t, s, "a", "changed")
	mustMiss(t, s, "b")
	mustGet(t, s, "c", "new")
	if err := s.Rollback(); !errors.Is(err, ErrNoTransaction) {
		t.Errorf("Rollback() after commit = %v, want ErrNoTransaction", err)
	}
}

func TestNestedRollbackOnlyDiscardsInner(t *testing.T) {
	s := NewStore()
	s.Set("a", "0")
	s.Begin()
	s.Set("a", "1")
	s.Begin()
	s.Set("a", "2")
	s.Set("b", "inner")
	mustGet(t, s, "a", "2")
	if err := s.Rollback(); err != nil {
		t.Fatalf("inner Rollback() = %v", err)
	}
	mustGet(t, s, "a", "1")
	mustMiss(t, s, "b")
	if err := s.Commit(); err != nil {
		t.Fatalf("outer Commit() = %v", err)
	}
	mustGet(t, s, "a", "1")
}

func TestInnerCommitThenOuterRollbackDiscardsEverything(t *testing.T) {
	s := NewStore()
	s.Set("a", "0")
	s.Begin()
	s.Set("a", "1")
	s.Begin()
	s.Set("a", "2")
	s.Set("b", "x")
	if err := s.Commit(); err != nil { // inner into outer
		t.Fatalf("inner Commit() = %v", err)
	}
	mustGet(t, s, "a", "2")
	mustGet(t, s, "b", "x")
	if err := s.Rollback(); err != nil { // outer
		t.Fatalf("outer Rollback() = %v", err)
	}
	mustGet(t, s, "a", "0")
	mustMiss(t, s, "b")
}

func TestDeleteThenSetInTransaction(t *testing.T) {
	s := NewStore()
	s.Set("a", "1")
	s.Begin()
	s.Delete("a")
	mustMiss(t, s, "a")
	s.Set("a", "2")
	mustGet(t, s, "a", "2")
	s.Rollback()
	mustGet(t, s, "a", "1")
}

func TestInnerDeleteCommittedIntoOuterHidesBaseValue(t *testing.T) {
	s := NewStore()
	s.Set("a", "1")
	s.Begin()
	s.Begin()
	s.Delete("a")
	s.Commit()
	mustMiss(t, s, "a")
	s.Commit()
	mustMiss(t, s, "a")
}
