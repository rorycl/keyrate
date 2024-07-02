# keyrate

Package keyrate implements a non-locking rate-limiter for a
slice of Key:Value [Thing], safe for concurrent use.

This simple module processes a slice of [Thing] where items
(Thing.Value) with common keys (Thing.Key) are provided via the [Get]
function channel at an interval at KeyRate (or longer) time
intervals.

## Details

Rate limiting utilises a time.Ticker to put each Thing.Value in a
Thing.Key group (or "bunch") on the result channel until exhausted.
The ticker is reset after each put to the channel to account for
blocking delays by the consumer or contention from processing other
bunches. This is not a bucket-type rate-limiting strategy with
replenishing tokens.

This simple solution requires items to be known upfront which are
then provided back to the consumer via the [Get] function, on a
per-bunch, rate-limited basis which is concurrent safe without the
need for locking.

## Other Solutions

Other solutions, such as https://github.com/sethvargo/go-limiter and
https://golang.org/x/time/rate provide more features such as
alternate storage backends or burstability and work well for
open-ended scenarios. The first offers a `Take` function which
returns a reset unix nanosecond -- when a wait is required --
utilising a refilling bucket model. The second offers a `Wait`
function implemented as a token bucket which also refills. Both use
mutex locking.

## Example usage

Example with "things" which have a string key and string value:

```go
KeyRate = time.Millisecond * 10  // ~100/sec
a := time.Now()
getter := Get([]Thing[string, string]{
	Thing[string, string]{"a", "b"},
	Thing[string, string]{"b", "c"},
	Thing[string, string]{"a", "c"},
})

// getter can be consumed by multiple goroutines
for t := range getter {
	fmt.Println(t)
}
fmt.Println(time.Now().Sub(a)) // will run in just over 10ms
```

Another example with an int key and type X value:

```go
type X struct {
	a string
}

KeyRate = time.Millisecond * 100  // ~10/sec
getter := Get([]Thing[int, X]{
	Thing[int, X]{0, X{"b"}},
	Thing[int, X]{1, X{"c"}},
	Thing[int, X]{0, X{"c"}},
})

for t := range getter {
	fmt.Printf("%v\n", t)
}
```

# Licence

Provided under the MIT Licence.