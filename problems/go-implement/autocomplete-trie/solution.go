package solution

// Autocomplete is a trie of words supporting prefix suggestions.
type Autocomplete struct {
	// TODO: add fields
}

// NewAutocomplete returns an empty trie.
func NewAutocomplete() *Autocomplete {
	// TODO: implement
	return &Autocomplete{}
}

// Insert adds word to the trie. Inserting a word twice is harmless.
func (a *Autocomplete) Insert(word string) {
	// TODO: implement
}

// Search reports whether word was inserted exactly (a prefix of an inserted
// word does not count).
func (a *Autocomplete) Search(word string) bool {
	// TODO: implement
	return false
}

// StartsWith reports whether any inserted word begins with prefix.
func (a *Autocomplete) StartsWith(prefix string) bool {
	// TODO: implement
	return false
}

// Suggest returns up to k inserted words that begin with prefix, in
// lexicographic (byte-wise) order. It returns an empty slice if k <= 0 or
// nothing matches. A word equal to prefix is itself a valid suggestion.
func (a *Autocomplete) Suggest(prefix string, k int) []string {
	// TODO: implement
	return nil
}
