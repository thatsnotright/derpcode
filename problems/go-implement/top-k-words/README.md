# Top K Frequent Words

Implement a streaming word counter that can report the `k` most frequent words seen so far.

Text is fed in chunks via `Add`. Tokenization rules:

- Words are case-insensitive; store and report them lowercased.
- A word is a maximal run of Unicode letters and digits. Every other character (spaces, punctuation, apostrophes, hyphens) is a separator. So `"Don't stop"` yields `don`, `t`, `stop`.

## Methods to implement

```go
// WordCount pairs a word with how many times it occurred.
type WordCount struct {
	Word  string
	Count int
}

func NewWordCounter() *WordCounter

// Add tokenizes text and adds its words to the running counts.
func (w *WordCounter) Add(text string)

// Count returns how many times word has been seen (case-insensitive).
func (w *WordCounter) Count(word string) int

// TopK returns up to k words ordered by Count descending; ties are broken
// alphabetically (ascending). If k <= 0 it returns an empty slice. If fewer
// than k distinct words exist, it returns all of them.
func (w *WordCounter) TopK(k int) []WordCount
```

## Example

```go
w := NewWordCounter()
w.Add("the quick brown fox jumps over the lazy dog")
w.Add("The dog barks; the fox runs.")
w.Count("THE")  // 4
w.TopK(3)       // [{the 4} {dog 2} {fox 2}]
```

## Constraints

- Input may be large; do not store the raw text.
- `TopK` is called rarely compared to `Add`, but `k` is usually much smaller than the number of distinct words.

## Follow-ups

1. Sorting all distinct words costs O(n log n). Do it in O(n log k) with a size-k heap. What comparator does the heap need?
2. Support counting over a sliding window of the last N `Add` calls.
3. The vocabulary no longer fits in memory. What approximate or external approach would you use (Count-Min sketch, Misra-Gries, sharded counting)?
