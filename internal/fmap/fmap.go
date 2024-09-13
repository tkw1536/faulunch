package fmap

import "strings"

// FMap represents a map that identifies keys via unicode case-folding.
// The original key case is maintained.
type FMap[K ~string, V any] map[K]V

// Add adds a new element to this f-map.
// Returns true if the element was not contained previously.
func (fmap FMap[K, V]) Add(key K, value V) (new bool) {
	for elem := range fmap {
		if strings.EqualFold(string(elem), string(key)) {
			return false
		}
	}
	fmap[key] = value
	return
}

// Remove removes an element from this f-map.
// Returns true iff the element was actually contained.
func (fmap FMap[K, V]) Remove(key string) (ok bool) {
	for elem := range fmap {
		if strings.EqualFold(string(elem), string(key)) {
			delete(fmap, elem)
			return true
		}
	}
	return false
}

func (fmap FMap[K, V]) Has(key K) bool {
	_, _, ok := fmap.Get(key)
	return ok
}

// Get returns the key and value for a given element
func (fmap FMap[K, V]) Get(key K) (K, V, bool) {
	// fast path: exact string is contained
	if value, ok := fmap[key]; ok {
		return key, value, true
	}
	for elem, value := range fmap {
		if strings.EqualFold(string(elem), string(key)) {
			return elem, value, true
		}
	}
	return "", *new(V), false
}
