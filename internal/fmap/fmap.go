//spellchecker:words fmap
package fmap

//spellchecker:words strings
import (
	"strings"
)

// FMap represents a map that identifies keys via unicode case-folding.
// The original key case is maintained.
type FMap[K ~string, V any] map[K]V

// Add adds a new element to this f-map.
// If the element was previously contained, does not change the value.
// Returns true if the element was not contained previously.
func (fmap FMap[K, V]) Add(key K, value V) (new bool) {
	if fmap.Has(key) {
		return false
	}
	fmap[key] = value
	return true
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
	_, ok := fmap.Key(key)
	return ok
}

// Key returns the case of key contained in the map.
// The boolean indicates if the key was found.
//
// If the key was not found, the provided key argument is returned.
func (fmap FMap[K, V]) Key(key K) (K, bool) {
	// fast path: exact string is contained
	if _, ok := fmap[key]; ok {
		return key, true
	}

	// slow path: check every possibly matchinbg key
	for elem := range fmap {
		if strings.EqualFold(string(elem), string(key)) {
			return elem, true
		}
	}
	return key, false
}

// Get returns the key and value for a given element
func (fmap FMap[K, V]) Get(key K) (K, V, bool) {
	key, ok := fmap.Key(key)
	if !ok {
		return key, *new(V), false
	}
	return key, fmap[key], true
}
