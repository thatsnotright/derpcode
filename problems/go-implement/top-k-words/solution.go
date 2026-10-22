package solution

// WordCount pairs a word with how many times it occurred.
type WordCount struct {
	Word  string
	Count int
}

// WordCounter counts words across a stream of text.
type WordCounter struct {
	// TODO: add fields
}

// NewWordCounter returns an empty counter.
func NewWordCounter() *WordCounter {
	// TODO: implement
	return &WordCounter{}
}

// Add tokenizes text and adds its words to the running counts.
func (w *WordCounter) Add(text string) {
	// TODO: implement
}

// Count returns how many times word has been seen (case-insensitive).
func (w *WordCounter) Count(word string) int {
	// TODO: implement
	return 0
}

// TopK returns up to k words ordered by Count descending; ties are broken
// alphabetically (ascending). If k <= 0 it returns an empty slice. If fewer
// than k distinct words exist, it returns all of them.
func (w *WordCounter) TopK(k int) []WordCount {
	// TODO: implement
	return nil
}
