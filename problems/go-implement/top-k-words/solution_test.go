package solution

import (
	"reflect"
	"testing"
)

func TestCountCaseInsensitive(t *testing.T) {
	w := NewWordCounter()
	w.Add("Go go GO gopher")
	tests := []struct {
		word string
		want int
	}{
		{"go", 3},
		{"GO", 3},
		{"Gopher", 1},
		{"missing", 0},
		{"", 0},
	}
	for _, tt := range tests {
		if got := w.Count(tt.word); got != tt.want {
			t.Errorf("Count(%q) = %d, want %d", tt.word, got, tt.want)
		}
	}
}

func TestTokenizationSeparators(t *testing.T) {
	w := NewWordCounter()
	w.Add("Hello, world! Hello... world?? don't re-run; x86\t42\nend")
	tests := []struct {
		word string
		want int
	}{
		{"hello", 2},
		{"world", 2},
		{"don", 1},
		{"t", 1},
		{"re", 1},
		{"run", 1},
		{"x86", 1},
		{"42", 1},
		{"end", 1},
	}
	for _, tt := range tests {
		if got := w.Count(tt.word); got != tt.want {
			t.Errorf("Count(%q) = %d, want %d", tt.word, got, tt.want)
		}
	}
}

func TestUnicodeWords(t *testing.T) {
	w := NewWordCounter()
	w.Add("Café café CAFÉ naïve")
	if got := w.Count("café"); got != 3 {
		t.Errorf("Count(café) = %d, want 3", got)
	}
	if got := w.Count("naïve"); got != 1 {
		t.Errorf("Count(naïve) = %d, want 1", got)
	}
}

func TestAddAccumulatesAcrossCalls(t *testing.T) {
	w := NewWordCounter()
	w.Add("apple banana")
	w.Add("apple")
	w.Add("")
	w.Add("   ,,, ")
	if got := w.Count("apple"); got != 2 {
		t.Errorf("Count(apple) = %d, want 2", got)
	}
	if got := w.Count("banana"); got != 1 {
		t.Errorf("Count(banana) = %d, want 1", got)
	}
}

func TestTopKOrdering(t *testing.T) {
	w := NewWordCounter()
	w.Add("the quick brown fox jumps over the lazy dog")
	w.Add("The dog barks; the fox runs.")
	want := []WordCount{{"the", 4}, {"dog", 2}, {"fox", 2}}
	if got := w.TopK(3); !reflect.DeepEqual(got, want) {
		t.Errorf("TopK(3) = %v, want %v", got, want)
	}
}

func TestTopKTieBreaksAlphabetically(t *testing.T) {
	w := NewWordCounter()
	w.Add("pear apple mango apple pear mango kiwi")
	want := []WordCount{{"apple", 2}, {"mango", 2}, {"pear", 2}, {"kiwi", 1}}
	if got := w.TopK(4); !reflect.DeepEqual(got, want) {
		t.Errorf("TopK(4) = %v, want %v", got, want)
	}
}

func TestTopKLargerThanDistinct(t *testing.T) {
	w := NewWordCounter()
	w.Add("b a b")
	want := []WordCount{{"b", 2}, {"a", 1}}
	if got := w.TopK(10); !reflect.DeepEqual(got, want) {
		t.Errorf("TopK(10) = %v, want %v", got, want)
	}
}

func TestTopKNonPositiveAndEmpty(t *testing.T) {
	w := NewWordCounter()
	if got := w.TopK(5); len(got) != 0 {
		t.Errorf("TopK on empty counter = %v, want empty", got)
	}
	w.Add("a b c")
	for _, k := range []int{0, -3} {
		if got := w.TopK(k); len(got) != 0 {
			t.Errorf("TopK(%d) = %v, want empty", k, got)
		}
	}
}

func TestTopKDoesNotMutateCounts(t *testing.T) {
	w := NewWordCounter()
	w.Add("a a b")
	got := w.TopK(2)
	if len(got) != 2 {
		t.Fatalf("TopK(2) = %v, want 2 entries", got)
	}
	got[0].Count = 100
	if w.Count("a") != 2 {
		t.Error("modifying TopK result changed internal counts")
	}
	if again := w.TopK(1); len(again) != 1 || again[0].Count != 2 {
		t.Errorf("TopK(1) = %v, want count 2", again)
	}
}
