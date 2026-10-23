package solution

import "testing"

func TestGetTable(t *testing.T) {
	m := NewTimeMap()
	m.Set("foo", "bar", 1)
	m.Set("foo", "baz", 4)
	m.Set("other", "x", 2)
	tests := []struct {
		name string
		key  string
		ts   int
		want string
	}{
		{"exact first", "foo", 1, "bar"},
		{"between", "foo", 3, "bar"},
		{"exact second", "foo", 4, "baz"},
		{"after last", "foo", 100, "baz"},
		{"before first", "foo", 0, ""},
		{"unknown key", "nope", 9, ""},
		{"other key isolated", "other", 3, "x"},
		{"other key before", "other", 1, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.Get(tt.key, tt.ts); got != tt.want {
				t.Errorf("Get(%q, %d) = %q, want %q", tt.key, tt.ts, got, tt.want)
			}
		})
	}
}

func TestEmptyStore(t *testing.T) {
	m := NewTimeMap()
	if got := m.Get("a", 5); got != "" {
		t.Errorf("Get on empty store = %q, want \"\"", got)
	}
}

func TestOutOfOrderSet(t *testing.T) {
	m := NewTimeMap()
	m.Set("k", "ten", 10)
	m.Set("k", "five", 5)
	m.Set("k", "twenty", 20)
	m.Set("k", "seven", 7)
	tests := []struct {
		ts   int
		want string
	}{
		{4, ""},
		{5, "five"},
		{6, "five"},
		{7, "seven"},
		{9, "seven"},
		{10, "ten"},
		{19, "ten"},
		{20, "twenty"},
	}
	for _, tt := range tests {
		if got := m.Get("k", tt.ts); got != tt.want {
			t.Errorf("Get(k, %d) = %q, want %q", tt.ts, got, tt.want)
		}
	}
}

func TestOverwriteSameTimestamp(t *testing.T) {
	m := NewTimeMap()
	m.Set("k", "old", 3)
	m.Set("k", "new", 3)
	if got := m.Get("k", 3); got != "new" {
		t.Errorf("Get(k, 3) = %q, want \"new\"", got)
	}
	if got := m.Get("k", 8); got != "new" {
		t.Errorf("Get(k, 8) = %q, want \"new\"", got)
	}
}

func TestZeroAndNegativeTimestamps(t *testing.T) {
	m := NewTimeMap()
	m.Set("k", "neg", -5)
	m.Set("k", "zero", 0)
	tests := []struct {
		ts   int
		want string
	}{
		{-6, ""},
		{-5, "neg"},
		{-1, "neg"},
		{0, "zero"},
		{1, "zero"},
	}
	for _, tt := range tests {
		if got := m.Get("k", tt.ts); got != tt.want {
			t.Errorf("Get(k, %d) = %q, want %q", tt.ts, got, tt.want)
		}
	}
}

func TestEmptyValueIsStored(t *testing.T) {
	m := NewTimeMap()
	m.Set("k", "v", 1)
	m.Set("k", "", 2)
	if got := m.Get("k", 2); got != "" {
		t.Errorf("Get(k, 2) = %q, want \"\"", got)
	}
	if got := m.Get("k", 1); got != "v" {
		t.Errorf("Get(k, 1) = %q, want \"v\"", got)
	}
}

func TestManyEntries(t *testing.T) {
	m := NewTimeMap()
	// Insert even timestamps in descending order.
	for ts := 2000; ts >= 0; ts -= 2 {
		m.Set("k", string(rune('a'+ts%26)), ts)
	}
	for ts := 0; ts <= 2001; ts++ {
		base := ts - ts%2
		want := string(rune('a' + base%26))
		if got := m.Get("k", ts); got != want {
			t.Fatalf("Get(k, %d) = %q, want %q", ts, got, want)
		}
	}
}
