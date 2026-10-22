// Package listing provides a paginated, searchable in-memory item catalogue.
package listing

import (
	"slices"
	"strings"
)

// Item is a catalogue entry.
type Item struct {
	ID   int
	Name string
}

// Service serves pages of items. It is not safe for concurrent mutation.
type Service struct {
	items []Item
	tags  map[string][]int
}

// NewService returns a Service backed by a copy of items.
func NewService(items []Item) *Service {
	return &Service{
		items: slices.Clone(items),
		tags:  make(map[string][]int),
	}
}

// Items returns a copy of all items.
func (s *Service) Items() []Item {
	return slices.Clone(s.items)
}

// Page returns a copy of the given 1-based page of items. Invalid arguments
// (page or size below one) and pages past the end yield an empty slice.
func (s *Service) Page(page, size int) []Item {
	if page < 1 || size < 1 {
		return []Item{}
	}
	// Checked this way round so (page-1)*size cannot overflow.
	if page-1 > len(s.items)/size {
		return []Item{}
	}
	start := (page - 1) * size
	if start >= len(s.items) {
		return []Item{}
	}
	end := len(s.items)
	if size < end-start {
		end = start + size
	}
	return slices.Clone(s.items[start:end])
}

// TotalPages returns how many pages of the given size exist. It returns 0
// when size is below one.
func (s *Service) TotalPages(size int) int {
	if size < 1 {
		return 0
	}
	return (len(s.items) + size - 1) / size
}

// Search returns the items whose name starts with prefix.
func (s *Service) Search(prefix string) []Item {
	matches := []Item{}
	for _, it := range s.items {
		if strings.HasPrefix(it.Name, prefix) {
			matches = append(matches, it)
		}
	}
	return matches
}

// Tag attaches tag to the item with the given id.
func (s *Service) Tag(id int, tag string) {
	s.tags[tag] = append(s.tags[tag], id)
}

// ByTag returns the items carrying tag.
func (s *Service) ByTag(tag string) []Item {
	out := []Item{}
	for _, id := range s.tags[tag] {
		for _, it := range s.items {
			if it.ID == id {
				out = append(out, it)
			}
		}
	}
	return out
}
