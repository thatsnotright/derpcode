package solution

import (
	"sort"
	"strings"
	"unicode"
)

// WordCount pairs a word with how many times it occurred.
type WordCount struct {
	Word  string
	Count int
}

// WordCounter counts words across a stream of text.
type WordCounter struct {
	counts map[string]int
}

// NewWordCounter returns an empty counter.
func NewWordCounter() *WordCounter {
	return &WordCounter{counts: make(map[string]int)}
}

// Add tokenizes text and adds its words to the running counts.
func (w *WordCounter) Add(text string) {
	words := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	for _, word := range words {
		w.counts[word]++
	}
}

// Count returns how many times word has been seen (case-insensitive).
func (w *WordCounter) Count(word string) int {
	return w.counts[strings.ToLower(word)]
}

// TopK returns up to k words ordered by Count descending; ties are broken
// alphabetically (ascending). If k <= 0 it returns an empty slice. If fewer
// than k distinct words exist, it returns all of them.
func (w *WordCounter) TopK(k int) []WordCount {
	if k <= 0 {
		return []WordCount{}
	}
	all := make([]WordCount, 0, len(w.counts))
	for word, n := range w.counts {
		all = append(all, WordCount{word, n})
	}
	sort.Slice(all, func(i, j int) bool {
		if all[i].Count != all[j].Count {
			return all[i].Count > all[j].Count
		}
		return all[i].Word < all[j].Word
	})
	if k < len(all) {
		all = all[:k]
	}
	return all
}
