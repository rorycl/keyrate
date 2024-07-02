package keyrate

import (
	"fmt"
	"slices"
	"time"
)

func ExampleKeyRate() {
	KeyRate = time.Millisecond * 10
	a := time.Now()
	getter := Get([]Thing[string, string]{
		Thing[string, string]{"a", "b"},
		Thing[string, string]{"b", "c"},
		Thing[string, string]{"a", "c"},
	})

	// getter is non-deterministic in terms of output order
	output := []string{}
	for t := range getter {
		output = append(output, t)
	}
	slices.Sort(output)

	fmt.Println(output)
	// will run in jut over 10ms
	fmt.Printf("%dms\n", time.Since(a).Milliseconds())

	// Output:
	// [b c c]
	// 10ms

}
