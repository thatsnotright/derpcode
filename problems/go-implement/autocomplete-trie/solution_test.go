package solution

import (
	"reflect"
	"testing"
)

func newWith(words ...string) *Autocomplete {
	a := NewAutocomplete()
	for _, w := range words {
		a.Insert(w)
	}
	return a
}

func TestSearchAndStartsWith(t *testing.T) {
	a := newWith("car", "card", "cat")
	tests := []struct {
		word       string
		search     bool
		startsWith bool
	}{
		{"car", true, true},
		{"card", true, true},
		{"ca", false, true},
		{"cards", false, false},
		{"dog", false, false},
		{"c", false, true},
	}
	for _, tt := range tests {
		if got := a.Search(tt.word); got != tt.search {
			t.Errorf("Search(%q) = %v, want %v", tt.word, got, tt.search)
		}
		if got := a.StartsWith(tt.word); got != tt.startsWith {
			t.Errorf("StartsWith(%q) = %v, want %v", tt.word, got, tt.startsWith)
		}
	}
}

func TestEmptyTrie(t *testing.T) {
	a := NewAutocomplete()
	if a.Search("") {
		t.Error("Search(\"\") on empty trie = true")
	}
	if a.StartsWith("a") {
		t.Error("StartsWith(a) on empty trie = true")
	}
	if got := a.Suggest("", 5); len(got) != 0 {
		t.Errorf("Suggest on empty trie = %v, want empty", got)
	}
}

func TestSuggestTable(t *testing.T) {
	a := newWith("car", "card", "care", "cat", "dog", "cab")
	tests := []struct {
		name   string
		prefix string
		k      int
		want   []string
	}{
		{"includes exact word", "car", 2, []string{"car", "card"}},
		{"all under prefix", "ca", 10, []string{"cab", "car", "card", "care", "cat"}},
		{"limit one", "ca", 1, []string{"cab"}},
		{"no match", "z", 3, []string{}},
		{"k zero", "ca", 0, []string{}},
		{"k negative", "ca", -1, []string{}},
		{"empty prefix", "", 3, []string{"cab", "car", "card"}},
		{"leaf", "dog", 5, []string{"dog"}},
		{"prefix longer than any word", "dogs", 5, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.Suggest(tt.prefix, tt.k)
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Suggest(%q, %d) = %v, want %v", tt.prefix, tt.k, got, tt.want)
			}
		})
	}
}

func TestSuggestLexicographicRegardlessOfInsertOrder(t *testing.T) {
	a := newWith("banana", "apple", "app", "apply", "ape", "band")
	want := []string{"ape", "app", "apple", "apply"}
	if got := a.Suggest("ap", 10); !reflect.DeepEqual(got, want) {
		t.Errorf("Suggest(ap, 10) = %v, want %v", got, want)
	}
	// "banana" < "band" lexicographically ('a' < 'd' at index 3).
	want = []string{"banana", "band"}
	if got := a.Suggest("ban", 10); !reflect.DeepEqual(got, want) {
		t.Errorf("Suggest(ban, 10) = %v, want %v", got, want)
	}
}

func TestDuplicateInsertNotDuplicatedInSuggestions(t *testing.T) {
	a := newWith("go", "go", "gopher")
	want := []string{"go", "gopher"}
	if got := a.Suggest("g", 10); !reflect.DeepEqual(got, want) {
		t.Errorf("Suggest(g, 10) = %v, want %v", got, want)
	}
}

func TestEmptyWord(t *testing.T) {
	a := newWith("", "a")
	if !a.Search("") {
		t.Error("Search(\"\") = false after inserting empty word")
	}
	want := []string{"", "a"}
	if got := a.Suggest("", 10); !reflect.DeepEqual(got, want) {
		t.Errorf("Suggest(\"\", 10) = %v, want %v", got, want)
	}
}

func TestNonASCII(t *testing.T) {
	a := newWith("héllo", "hello", "héros")
	if !a.Search("héllo") {
		t.Error("Search(héllo) = false")
	}
	if !a.StartsWith("hé") {
		t.Error("StartsWith(hé) = false")
	}
	want := []string{"héllo", "héros"}
	if got := a.Suggest("hé", 5); !reflect.DeepEqual(got, want) {
		t.Errorf("Suggest(hé, 5) = %v, want %v", got, want)
	}
}

func TestManyWords(t *testing.T) {
	a := NewAutocomplete()
	for c1 := 'a'; c1 <= 'z'; c1++ {
		for c2 := 'a'; c2 <= 'z'; c2++ {
			a.Insert(string([]rune{c1, c2}))
		}
	}
	if got := a.Suggest("mx", 1); len(got) != 1 || got[0] != "mx" {
		t.Errorf("Suggest(mx, 1) = %v", got)
	}
	got := a.Suggest("m", 26)
	if len(got) != 26 || got[0] != "ma" || got[25] != "mz" {
		t.Errorf("Suggest(m, 26) wrong: len=%d %v", len(got), got)
	}
}
