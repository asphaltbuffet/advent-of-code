package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"Zero", "0", 0},
		{"One", "1", 1},
		{"Two", "2", 2},
		{"Three", "1=", 3},
		{"Four", "1-", 4},
		{"Ten", "20", 10},
		{"Fifteen", "1=0", 15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Decode(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEncode(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  string
	}{
		{"Zero", 0, "0"},
		{"One", 1, "1"},
		{"Two", 2, "2"},
		{"Three", 3, "1="},
		{"Four", 4, "1-"},
		{"Ten", 10, "20"},
		{"Fifteen", 15, "1=0"},
		{"2022", 2022, "1=11-2"},
		{"12345", 12345, "1-0---0"},
		{"314159265", 314159265, "1121-1110-1=0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Encode(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func FuzzEncode(f *testing.F) {
	testcases := []int{0, 3, 15}

	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, n int) {
		if n < 0 {
			t.Skip("negative numbers are not supported")
		}

		e := Encode(n)
		d := Decode(e)
		assert.Equal(t, n, d)
	})
}
