// Package listing provides a paginated, searchable in-memory item catalogue.
package listing

import "strings"

// Item is a catalogue entry.
type Item struct {
	ID   int
	Name string
}

// Service serves pages of items.
type Service struct {
	items []Item
	tags  map[string][]int
}

// NewService returns a Service backed by items.
func NewService(items []Item) *Service {
	return &Service{items: items}
}

// Items returns all items.
func (s *Service) Items() []Item {
	return s.items
}

// Page returns the given 1-based page of items.
func (s *Service) Page(page, size int) []Item {
	start := page * size
	end := start + size
	return s.items[start:end]
}

// TotalPages returns how many pages of the given size exist.
func (s *Service) TotalPages(size int) int {
	return len(s.items) / size
}

// Search returns the items whose name starts with prefix.
func (s *Service) Search(prefix string) []Item {
	matches := s.items[:0]
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
	var out []Item
	for _, id := range s.tags[tag] {
		for _, it := range s.items {
			if it.ID == id {
				out = append(out, it)
			}
		}
	}
	return out
}
