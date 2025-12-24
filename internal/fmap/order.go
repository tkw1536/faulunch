package fmap

// Order turns a slice into an f-map with the index as the value.
func Order[T ~string](values ...T) FMap[T, int] {
	m := make(FMap[T, int], len(values))
	for index, item := range values {
		m.Add(item, index)
	}
	return m
}
