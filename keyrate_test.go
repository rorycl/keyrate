package keyrate

import (
	"fmt"
	"testing"
	"time"
)

func TestGetStringString(t *testing.T) {

	type ssT = Thing[string, string]
	tests := []struct {
		keyRate     time.Duration
		things      []ssT
		maxDuration time.Duration
		results     int
	}{
		{
			keyRate:     10 * time.Millisecond,
			things:      []ssT{},
			maxDuration: time.Millisecond * 5,
			results:     0,
		},
		{
			keyRate: 10 * time.Millisecond,
			things: []ssT{
				ssT{"a", "b"},
				ssT{"b", "c"},
				ssT{"a", "c"},
				ssT{"a", "d"},
				ssT{"b", "d"},
				ssT{"z", "n"},
			},
			maxDuration: time.Millisecond * 33,
			results:     6,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			// change package global KeyRate
			KeyRate = tt.keyRate
			a := time.Now()
			counter := 0
			getter := Get(tt.things)
			for range getter {
				counter++
			}
			b := time.Since(a)
			if got, want := b, tt.maxDuration; got > want {
				t.Errorf("duration got %v want %v", got, want)
			}
			if got, want := counter, tt.results; got != want {
				t.Errorf("result count got %d want %d", got, want)
			}
		})
	}
}

func TestGetIntString(t *testing.T) {

	type isT = Thing[int, string]
	tests := []struct {
		keyRate     time.Duration
		things      []isT
		maxDuration time.Duration
		results     int
	}{
		{
			keyRate: 10 * time.Millisecond,
			things: []isT{
				isT{0, "b"},
				isT{1, "c"},
				isT{0, "c"},
				isT{1, "d"},
				isT{9, "n"},
			},
			maxDuration: time.Millisecond * 12,
			results:     5,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			// change package global KeyRate
			KeyRate = tt.keyRate
			a := time.Now()
			counter := 0
			getter := Get(tt.things)
			for range getter {
				counter++
			}
			b := time.Since(a)
			if got, want := b, tt.maxDuration; got > want {
				t.Errorf("duration got %v want %v", got, want)
			}
			if got, want := counter, tt.results; got != want {
				t.Errorf("result count got %d want %d", got, want)
			}
		})
	}
}

func TestGetIntStruct(t *testing.T) {

	type X struct {
		a string
	}

	type ixT = Thing[int, X]
	tests := []struct {
		keyRate     time.Duration
		things      []ixT
		maxDuration time.Duration
		results     int
	}{
		{
			keyRate: 10 * time.Millisecond,
			things: []ixT{
				ixT{0, X{"b"}},
				ixT{1, X{"c"}},
				ixT{0, X{"c"}},
				ixT{1, X{"d"}},
				ixT{9, X{"n"}},
			},
			maxDuration: time.Millisecond * 12,
			results:     5,
		},
		{
			keyRate: 10 * time.Millisecond,
			things: []ixT{
				ixT{0, X{"b"}},
				ixT{1, X{"c"}},
				ixT{0, X{"c"}},
				ixT{1, X{"d"}},
				ixT{9, X{"a"}},
				ixT{9, X{"b"}},
				ixT{9, X{"c"}},
				ixT{9, X{"d"}},
				ixT{9, X{"e"}},
			},
			maxDuration: time.Millisecond * 45, // drifts up to 5ms
			results:     9,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			// change package global KeyRate
			KeyRate = tt.keyRate
			a := time.Now()
			counter := 0
			getter := Get(tt.things)
			for range getter {
				counter++
			}
			b := time.Since(a)
			if got, want := b, tt.maxDuration; got > want {
				t.Errorf("duration got %v want %v", got, want)
			}
			if got, want := counter, tt.results; got != want {
				t.Errorf("result count got %d want %d", got, want)
			}
		})
	}
}
