//spellchecker:words internal
package internal

//spellchecker:words slices
import "slices"

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
