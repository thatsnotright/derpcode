package solution

import "testing"

func TestPutGetBasic(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(1, 10)
	c.Put(2, 20)
	tests := []struct {
		key  int
		want int
		ok   bool
	}{
		{1, 10, true},
		{2, 20, true},
		{3, 0, false},
	}
	for _, tt := range tests {
		got, ok := c.Get(tt.key)
		if got != tt.want || ok != tt.ok {
			t.Errorf("Get(%d) = (%d, %v), want (%d, %v)", tt.key, got, ok, tt.want, tt.ok)
		}
	}
}

func TestEvictsLeastRecentlyUsed(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(3, 3) // evicts 1
	if _, ok := c.Get(1); ok {
		t.Error("key 1 should have been evicted")
	}
	for _, k := range []int{2, 3} {
		if _, ok := c.Get(k); !ok {
			t.Errorf("key %d should be present", k)
		}
	}
}

func TestGetRefreshesRecency(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Get(1)    // 2 is now least recent
	c.Put(3, 3) // evicts 2
	if _, ok := c.Get(2); ok {
		t.Error("key 2 should have been evicted")
	}
	if v, ok := c.Get(1); !ok || v != 1 {
		t.Errorf("Get(1) = (%d, %v), want (1, true)", v, ok)
	}
}

func TestPutExistingUpdatesAndRefreshes(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(1, 100) // update, 2 becomes least recent
	if c.Len() != 2 {
		t.Fatalf("Len() = %d, want 2", c.Len())
	}
	c.Put(3, 3) // evicts 2
	if _, ok := c.Get(2); ok {
		t.Error("key 2 should have been evicted")
	}
	if v, _ := c.Get(1); v != 100 {
		t.Errorf("Get(1) = %d, want 100", v)
	}
}

func TestCapacityOne(t *testing.T) {
	c := NewLRUCache(1)
	c.Put(1, 1)
	c.Put(2, 2)
	if _, ok := c.Get(1); ok {
		t.Error("key 1 should have been evicted")
	}
	if v, ok := c.Get(2); !ok || v != 2 {
		t.Errorf("Get(2) = (%d, %v), want (2, true)", v, ok)
	}
	if c.Len() != 1 {
		t.Errorf("Len() = %d, want 1", c.Len())
	}
}

func TestCapacityZeroStoresNothing(t *testing.T) {
	c := NewLRUCache(0)
	c.Put(1, 1)
	if _, ok := c.Get(1); ok {
		t.Error("zero-capacity cache should not store entries")
	}
	if c.Len() != 0 {
		t.Errorf("Len() = %d, want 0", c.Len())
	}
}

func TestZeroValueIsDistinctFromMissing(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(7, 0)
	if v, ok := c.Get(7); !ok || v != 0 {
		t.Errorf("Get(7) = (%d, %v), want (0, true)", v, ok)
	}
}

func TestLenNeverExceedsCapacity(t *testing.T) {
	c := NewLRUCache(3)
	for i := 0; i < 100; i++ {
		c.Put(i, i)
		want := i + 1
		if want > 3 {
			want = 3
		}
		if c.Len() != want {
			t.Fatalf("after %d puts Len() = %d, want %d", i+1, c.Len(), want)
		}
	}
	for i := 97; i < 100; i++ {
		if v, ok := c.Get(i); !ok || v != i {
			t.Errorf("Get(%d) = (%d, %v), want (%d, true)", i, v, ok, i)
		}
	}
}
