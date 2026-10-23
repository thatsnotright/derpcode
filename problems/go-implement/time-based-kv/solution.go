package solution

// TimeMap is a key-value store that remembers every value over time.
type TimeMap struct {
	// TODO: add fields
}

// NewTimeMap returns an empty store.
func NewTimeMap() *TimeMap {
	// TODO: implement
	return &TimeMap{}
}

// Set records that key had value at the given timestamp. Setting the same
// (key, timestamp) again overwrites the earlier value.
func (m *TimeMap) Set(key, value string, timestamp int) {
	// TODO: implement
}

// Get returns the value whose timestamp is the largest one <= timestamp for
// key, or "" if there is none.
func (m *TimeMap) Get(key string, timestamp int) string {
	// TODO: implement
	return ""
}
