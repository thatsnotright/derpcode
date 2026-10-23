# Autocomplete Trie

Implement a prefix tree (trie) that stores words and can suggest completions for a prefix.

## Methods to implement

```go
// NewAutocomplete returns an empty trie.
func NewAutocomplete() *Autocomplete

// Insert adds word to the trie. Inserting a word twice is harmless.
func (a *Autocomplete) Insert(word string)

// Search reports whether word was inserted exactly (a prefix of an inserted
// word does not count).
func (a *Autocomplete) Search(word string) bool

// StartsWith reports whether any inserted word begins with prefix.
func (a *Autocomplete) StartsWith(prefix string) bool

// Suggest returns up to k inserted words that begin with prefix, in
// lexicographic (byte-wise) order. It returns an empty slice if k <= 0 or
// nothing matches. A word equal to prefix is itself a valid suggestion.
func (a *Autocomplete) Suggest(prefix string, k int) []string
```

## Example

```go
a := NewAutocomplete()
for _, w := range []string{"car", "card", "care", "cat", "dog"} {
	a.Insert(w)
}
a.Search("car")         // true
a.Search("ca")          // false
a.StartsWith("ca")      // true
a.StartsWith("cow")     // false
a.Suggest("car", 2)     // ["car", "card"]
a.Suggest("ca", 10)     // ["car", "card", "care", "cat"]
a.Suggest("z", 3)       // []
```

## Constraints

- Words are arbitrary strings (treat them as bytes or runes, but be consistent; tests use lowercase ASCII plus a few non-ASCII words).
- The empty string is a valid prefix matching every word. The empty word may be inserted.
- `Suggest` should not visit more of the trie than needed: stop once `k` words are found.

## Follow-ups

1. Add `Delete(word string)` that also frees nodes no longer needed.
2. Rank suggestions by insertion frequency instead of lexicographic order. How do you keep `Suggest` fast?
3. Memory is dominated by child maps. What alternatives (arrays, compressed/radix trie) would you consider?
