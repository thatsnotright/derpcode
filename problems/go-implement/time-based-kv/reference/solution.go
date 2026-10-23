package solution

import "sort"

type entry struct {
	timestamp int
	value     string
}

// TimeMap is a key-value store that remembers every value over time.
type TimeMap struct {
	history map[string][]entry // each slice is sorted by timestamp
}

// NewTimeMap returns an empty store.
func NewTimeMap() *TimeMap {
	return &TimeMap{history: make(map[string][]entry)}
}

// Set records that key had value at the given timestamp. Setting the same
// (key, timestamp) again overwrites the earlier value.
func (m *TimeMap) Set(key, value string, timestamp int) {
	entries := m.history[key]
	i := sort.Search(len(entries), func(i int) bool { return entries[i].timestamp >= timestamp })
	if i < len(entries) && entries[i].timestamp == timestamp {
		entries[i].value = value
		return
	}
	entries = append(entries, entry{})
	copy(entries[i+1:], entries[i:])
	entries[i] = entry{timestamp: timestamp, value: value}
	m.history[key] = entries
}

// Get returns the value whose timestamp is the largest one <= timestamp for
// key, or "" if there is none.
func (m *TimeMap) Get(key string, timestamp int) string {
	entries := m.history[key]
	// First entry strictly after timestamp; the one before it is the answer.
	i := sort.Search(len(entries), func(i int) bool { return entries[i].timestamp > timestamp })
	if i == 0 {
		return ""
	}
	return entries[i-1].value
}
