package solution

import "sort"

type node struct {
	children map[byte]*node
	isWord   bool
}

func newNode() *node {
	return &node{children: make(map[byte]*node)}
}

// Autocomplete is a trie of words supporting prefix suggestions.
type Autocomplete struct {
	root *node
}

// NewAutocomplete returns an empty trie.
func NewAutocomplete() *Autocomplete {
	return &Autocomplete{root: newNode()}
}

// Insert adds word to the trie. Inserting a word twice is harmless.
func (a *Autocomplete) Insert(word string) {
	n := a.root
	for i := 0; i < len(word); i++ {
		child, ok := n.children[word[i]]
		if !ok {
			child = newNode()
			n.children[word[i]] = child
		}
		n = child
	}
	n.isWord = true
}

// find returns the node reached by following prefix, or nil.
func (a *Autocomplete) find(prefix string) *node {
	n := a.root
	for i := 0; i < len(prefix); i++ {
		n = n.children[prefix[i]]
		if n == nil {
			return nil
		}
	}
	return n
}

// Search reports whether word was inserted exactly (a prefix of an inserted
// word does not count).
func (a *Autocomplete) Search(word string) bool {
	n := a.find(word)
	return n != nil && n.isWord
}

// StartsWith reports whether any inserted word begins with prefix.
func (a *Autocomplete) StartsWith(prefix string) bool {
	return a.find(prefix) != nil
}

// Suggest returns up to k inserted words that begin with prefix, in
// lexicographic (byte-wise) order. It returns an empty slice if k <= 0 or
// nothing matches. A word equal to prefix is itself a valid suggestion.
func (a *Autocomplete) Suggest(prefix string, k int) []string {
	out := []string{}
	n := a.find(prefix)
	if n == nil || k <= 0 {
		return out
	}
	buf := []byte(prefix)
	var walk func(n *node)
	walk = func(n *node) {
		if len(out) >= k {
			return
		}
		if n.isWord {
			out = append(out, string(buf))
		}
		keys := make([]byte, 0, len(n.children))
		for b := range n.children {
			keys = append(keys, b)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
		for _, b := range keys {
			if len(out) >= k {
				return
			}
			buf = append(buf, b)
			walk(n.children[b])
			buf = buf[:len(buf)-1]
		}
	}
	walk(n)
	return out
}
