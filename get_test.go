package keyrate

import (
	"fmt"
	"slices"
	"time"
)

func Example() {
	KeyRate = time.Millisecond * 10
	a := time.Now()
	type ssT = Thing[string, string]
	getter := Get([]ssT{
		ssT{"a", "b"},
		ssT{"b", "c"},
		ssT{"a", "c"},
	})

	// getter is non-deterministic in terms of output order
	output := []string{}
	for t := range getter {
		output = append(output, t)
	}
	slices.Sort(output)

	fmt.Println(output)
	// will run in just over 10ms
	fmt.Printf("%dms\n", time.Since(a).Milliseconds())

	// Output:
	// [b c c]
	// 10ms

}
