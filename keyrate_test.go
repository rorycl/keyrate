package keyrate

import (
	"fmt"
	"testing"
	"time"
)

func TestGetStringString(t *testing.T) {

	tests := []struct {
		keyRate     time.Duration
		things      []Thing[string, string]
		maxDuration time.Duration
		results     int
	}{
		{
			keyRate:     10 * time.Millisecond,
			things:      []Thing[string, string]{},
			maxDuration: time.Millisecond * 5,
			results:     0,
		},
		{
			keyRate: 10 * time.Millisecond,
			things: []Thing[string, string]{
				Thing[string, string]{"a", "b"},
				Thing[string, string]{"b", "c"},
				Thing[string, string]{"a", "c"},
				Thing[string, string]{"a", "d"},
				Thing[string, string]{"b", "d"},
				Thing[string, string]{"z", "n"},
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

	tests := []struct {
		keyRate     time.Duration
		things      []Thing[int, string]
		maxDuration time.Duration
		results     int
	}{
		{
			keyRate: 10 * time.Millisecond,
			things: []Thing[int, string]{
				Thing[int, string]{0, "b"},
				Thing[int, string]{1, "c"},
				Thing[int, string]{0, "c"},
				Thing[int, string]{1, "d"},
				Thing[int, string]{9, "n"},
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

	tests := []struct {
		keyRate     time.Duration
		things      []Thing[int, X]
		maxDuration time.Duration
		results     int
	}{
		{
			keyRate: 10 * time.Millisecond,
			things: []Thing[int, X]{
				Thing[int, X]{0, X{"b"}},
				Thing[int, X]{1, X{"c"}},
				Thing[int, X]{0, X{"c"}},
				Thing[int, X]{1, X{"d"}},
				Thing[int, X]{9, X{"n"}},
			},
			maxDuration: time.Millisecond * 12,
			results:     5,
		},
		{
			keyRate: 10 * time.Millisecond,
			things: []Thing[int, X]{
				Thing[int, X]{0, X{"b"}},
				Thing[int, X]{1, X{"c"}},
				Thing[int, X]{0, X{"c"}},
				Thing[int, X]{1, X{"d"}},
				Thing[int, X]{9, X{"a"}},
				Thing[int, X]{9, X{"b"}},
				Thing[int, X]{9, X{"c"}},
				Thing[int, X]{9, X{"d"}},
				Thing[int, X]{9, X{"e"}},
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
