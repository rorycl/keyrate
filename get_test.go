package keyrate

import (
	"fmt"
	"slices"
	"time"
)

// Example provides an example for using the keyrate package
func Example() {
	a := time.Now()

	// override the package default KeyRate
	KeyRate = time.Millisecond * 10

	// make a string string Thing type alias to reduce boilerplate
	type ssT = Thing[string, string]

	// define a channel from a slice of key/value pairs
	getter := Get([]ssT{
		ssT{"a", "b"},
		ssT{"b", "c"},
		ssT{"a", "c"},
	})

	// read off the getter chan
	// note that the order is non-deterministic; the only guarantee is
	// that items with the same key will be provided no sooner than
	// KeyRate.
	output := []string{}
	for t := range getter {
		output = append(output, t)
	}

	// sort output for this test (because the order of b/c can swap)
	slices.Sort(output)
	fmt.Println(output)

	// will run in just over 10ms
	fmt.Printf("%dms\n", time.Since(a).Milliseconds())

	// Output:
	// [b c c]
	// 10ms

}
