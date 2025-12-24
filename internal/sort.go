//spellchecker:words internal
package internal

//spellchecker:words slices
import "slices"

// SortedKeysOf returns a slice of the given map, sorted using the sort function.
// If sort is nil, the keys are returned in any order.
func SortedKeysOf[K comparable, V any](m map[K]V, sort func(a, b K) int) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	if sort != nil {
		slices.SortStableFunc(keys, sort)
	}
	return keys
}
