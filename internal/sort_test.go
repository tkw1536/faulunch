package internal_test

import (
	"cmp"
	"fmt"

	"github.com/tkw1536/faulunch/internal"
)

func ExampleSortedKeysOf_withSort() {
	m := map[string]int{"banana": 1, "apple": 2, "cherry": 3}
	keys := internal.SortedKeysOf(m, cmp.Compare[string])
	fmt.Println(keys)
	// Output: [apple banana cherry]
}

func ExampleSortedKeysOf_withoutSort() {
	m := map[string]int{"only": 1}
	keys := internal.SortedKeysOf(m, nil)
	fmt.Println(keys)
	// Output: [only]
}
