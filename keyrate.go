// Package keyrate implements a non-locking rate-limiter for a
// slice of Key:Value [Thing], safe for concurrent use.
//
// This simple module processes a slice of [Thing] where items
// (Thing.Value) with common keys (Thing.Key) are provided via the [Get]
// function channel at an interval at KeyRate (or longer) time
// intervals.
//
// # Details
//
// Rate limiting utilises a time.Ticker to put each Thing.Value in a
// Thing.Key group (or "bunch") on the result channel until exhausted.
// The ticker is reset after each put to the channel to account for
// blocking delays by the consumer or contention from processing other
// bunches. This is not a bucket-type rate-limiting strategy with
// replenishing tokens.
//
// This simple solution requires items to be known upfront which are
// then provided back to the consumer via the [Get] function, on a
// per-bunch, rate-limited basis which is concurrent safe without the
// need for locking.
//
// # Other Solutions
//
// Other solutions, such as https://github.com/sethvargo/go-limiter and
// https://golang.org/x/time/rate provide more features such as
// alternate storage backends or burstability and work well for
// open-ended scenarios. The first offers a `Take` function which
// returns a reset unix nanosecond -- when a wait is required --
// utilising a refilling bucket model. The second offers a `Wait`
// function implemented as a token bucket which also refills. Both use
// mutex locking.
//
// # Example usage
//
// Example with "things" which have a string key and string value:
//
//	KeyRate = time.Millisecond * 10  // ~100/sec
//	a := time.Now()
//	getter := Get([]Thing[string, string]{
//		Thing[string, string]{"a", "b"},
//		Thing[string, string]{"b", "c"},
//		Thing[string, string]{"a", "c"},
//	})
//
//	// getter can be consumed by multiple goroutines
//	for t := range getter {
//		fmt.Println(t)
//	}
//	fmt.Println(time.Now().Sub(a)) // will run in just over 10ms
//
// Another example with an int key and type X value:
//
//	type X struct {
//		a string
//	}
//
//	KeyRate = time.Millisecond * 100  // ~10/sec
//	getter := Get([]Thing[int, X]{
//		Thing[int, X]{0, X{"b"}},
//		Thing[int, X]{1, X{"c"}},
//		Thing[int, X]{0, X{"c"}},
//	})
//
//	for t := range getter {
//		fmt.Printf("%v\n", t)
//	}
//
// # Licence
//
// Provided under the MIT License.
package keyrate

import (
	"sync"
	"time"
)

// Thing is a Key/Value type where the Key is used for bunching items by
// the Get function.
type Thing[T comparable, S any] struct {
	Key   T
	Value S
}

// KeyRate is the fastest rate at which bunched Thing items should be
// provided. Note that this implementation may provide Thing items by
// bunch at a slower rate due to contention on the Get result chan.
var KeyRate = time.Duration(1 * time.Second)

// Get provides a chan of Thing.Value where bunches of Thing having the
// same Key are rate limited to at least KeyRate.
func Get[T comparable, S any](things []Thing[T, S]) <-chan S {

	keyVals := map[T][]S{}

	thingValChan := make(chan S)
	if len(things) < 1 {
		close(thingValChan)
		return thingValChan
	}

	// bunch Thing values by Key
	for _, t := range things {
		if _, ok := keyVals[t.Key]; !ok {
			keyVals[t.Key] = []S{t.Value}
		} else {
			keyVals[t.Key] = append(keyVals[t.Key], t.Value)
		}
	}

	// Process each bunch of values by Key in its own goroutine. Where
	// there is more than one Value for a bunch, ensure that values are
	// only pushed onto the result chan at timedVal intervals.
	var wg sync.WaitGroup
	for _, v := range keyVals {
		wg.Add(1)
		go func(vals []S) {
			defer wg.Done()
			thingValChan <- vals[0]
			if len(vals) == 1 {
				return
			}
			tick := time.NewTicker(KeyRate)
			defer tick.Stop()
			for _, timedVal := range vals[1:] {
				<-tick.C // wait for tick
				thingValChan <- timedVal
				// Reset ticker since thingValChan may have blocked.
				// While backlogged tickers might be fine to use,
				// they may have the effect of bursting which might
				// have unintended consequences. So instead of
				// using a bucket strategy, simply resetting the
				// keyRate is a conservative (and simple) solution.
				tick.Reset(KeyRate)
			}
		}(v)
	}
	go func() {
		wg.Wait()
		close(thingValChan)
	}()
	return thingValChan
}
