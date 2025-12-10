package exercises

import (
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NextInvalid(t *testing.T) {

	tests := []struct {
		name string
		arg  int
		want int
	}{
		{"11", 11, 22},
		{"98", 98, 99},
		{"99", 99, 1010},
		{"111", 111, 1010},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NextInvalid(tt.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_IsInvalid(t *testing.T) {
	tests := []struct {
		name string
		arg  int
		want bool
	}{
		{"valid even", 10, false},
		{"valid odd", 101, false},
		{"invalid short", 11, true},
		{"invalid long", 420420, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsInvalid(tt.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_IsRepeating(t *testing.T) {
	tests := []struct {
		name      string
		arg       int
		assertion assert.BoolAssertionFunc
	}{
		{"valid single", 11, assert.True},
		{"valid long single", 4444444, assert.True},
		{"valid double", 1010, assert.True},
		{"valid triple", 321321321, assert.True},
		{"invalid single", 9, assert.False},
		{"invalid diff", 3456, assert.False},
		{"invalid mid-diff", 45345, assert.False},
		{"invalid long", 1221001221, assert.False},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, IsRepeating(tt.arg))
		})
	}
}

func Test_InvalidIds(t *testing.T) {
	got := slices.Collect(InvalidIds(11, 22))
	assert.Equal(t, []int{11, 22}, got)
}

func Test_Repeated(t *testing.T) {
	tests := []struct {
		name string
		pair Pair
		want []int
	}{
		{"shortest", Pair{10, 30, "10", "30"}, []int{11, 22}},
		{
			"single repeating",
			Pair{10, 700, "11", "700"},
			[]int{11, 22, 33, 44, 55, 66, 77, 88, 99, 111, 222, 333, 444, 555, 666}},
		{"4 digit", Pair{998, 1012, "998", "1012"}, []int{999, 1010}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pair.Repeated()
			if assert.Nil(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_IntLen(t *testing.T) {
	tests := []struct {
		name string
		arg  int
		want int
	}{
		{"single", 9, 1},
		{"double", 99, 2},
		{"long", 123515251, 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, intLen(tt.arg))
		})
	}
}

func FuzzRepeatedMatchesSlow(f *testing.F) {
	// A couple of seeds to help it start somewhere sane
	f.Add(1, 100)
	f.Add(10, 9999)
	f.Add(123, 12345)
	f.Add(523, 560)

	f.Fuzz(func(t *testing.T, lo, hi int) {
		// Normalize / constrain so fuzz doesn’t go wild
		if lo <= 0 {
			lo = 1
		}
		if hi <= 0 {
			hi = 1
		}
		if lo > hi {
			lo, hi = hi, lo
		}

		// Cap the range size so the slow version doesn't explode
		if hi-lo > 50000 {
			hi = lo + 50000
		}

		p := Pair{
			Lower:       lo,
			Upper:       hi,
			LowerString: strconv.Itoa(lo),
			UpperString: strconv.Itoa(hi),
		}

		got, err := p.Repeated()
		want := repeatedSlow(p)

		require.Nil(t, err)
		require.Equal(t, want, got,
			"mismatch for range [%d, %d]\nwant: %v\ngot:  %v",
			lo, hi, want, got,
		)
	})
}

func repeatedSlow(p Pair) []int {
	var out []int
	for n := p.Lower; n <= p.Upper; n++ {
		if isRepeatedPattern(n) {
			out = append(out, n)
		}
	}
	slices.Sort(out)
	return slices.Compact(out)
}
