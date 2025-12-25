package fmap_test

import (
	"fmt"

	"github.com/tkw1536/faulunch/internal/fmap"
)

func ExampleOrder() {
	m := fmap.Order("a", "b", "c")
	_, a, _ := m.Get("A")
	_, b, _ := m.Get("B")
	_, c, _ := m.Get("C")
	fmt.Println(a, b, c)
	// Output: 0 1 2
}
