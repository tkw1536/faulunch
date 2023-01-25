package main

import (
	"fmt"

	"github.com/tkw1536/faulunch"
)

func main() {
	res, err := faulunch.Fetch(faulunch.MensaLmp, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", res)
}
