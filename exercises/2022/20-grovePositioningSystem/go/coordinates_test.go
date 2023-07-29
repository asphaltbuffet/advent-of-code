package exercises

import (
	"testing"

	"github.com/asphaltbuffet/advent-of-code/pkg/ring"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMoveSingleRingValue(t *testing.T) {
	tests := []struct {
		name       string
		ringValues []int
		start      int
		move       int
		expected   []int
	}{
		{"Move 0", []int{0, 1, 2, 3, 4}, 1, 0, []int{0, 1, 2, 3, 4}},
		{"Move 1, offset start", []int{3, 4, 0, 1, 2}, 3, 1, []int{3, 4, 0, 2, 1}},
		{"Move 3", []int{0, 1, 2, 3, 4}, 1, 3, []int{0, 2, 3, 4, 1}},
		{"Move 53", []int{0, 1, 2, 3, 4}, 1, 53, []int{0, 2, 3, 4, 1}},
		{"Move 5", []int{0, 1, 2, 3, 4}, 1, 5, []int{0, 1, 2, 3, 4}},
		{"Move -1", []int{0, 1, 2, 3, 4}, 1, -1, []int{0, 2, 3, 4, 1}},
		{"Move -1051", []int{0, 1, 2, 3, 4}, 1, -1051, []int{0, 2, 3, 4, 1}},
		{"Move -2", []int{0, 1, 2, 3, 4}, 1, -2, []int{0, 2, 3, 1, 4}},
		{"Move 8", []int{0, 1, 2, 3, 4}, 1, 8, []int{0, 2, 3, 4, 1}},
		{"Move -7", []int{0, 1, 2, 3, 4}, 1, -7, []int{0, 2, 3, 1, 4}},
		{"example", []int{0, -2, 5, 6, 7, 8, 9}, 1, -2, []int{0, 5, 6, 7, 8, -2, 9}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := makeRing(t, tt.ringValues)

			require.Equal(t, tt.ringValues[0], r.Element.value, "The first value should be 0")

			shift(r.Move(tt.start), tt.move)

			got := []int{}
			r.Do(func(v digit) {
				got = append(got, v.value)
			})

			assert.Equal(t, tt.expected, got, "They should be equal")
		})
	}
}

func makeRing(t *testing.T, values []int) *ring.Ring[digit] {
	t.Helper()

	r := ring.New[digit](len(values))
	for i := 0; i < r.Len(); i++ {
		r.Element.value = values[i]
		r = r.Next()
	}

	return r
}

func TestPrint(t *testing.T) {
	tests := []struct {
		name       string
		ringValues string
		expected   string
	}{
		{"start at 0", "0\n1\n2\n3\n4", "0 1 2 3 4"},
		{"mid zero", "3\n4\n0\n1\n2", "0 1 2 3 4"},
		{"one item", "0", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf, err := parse(tt.ringValues)
			require.NoError(t, err, "parse failed")

			got := cf.decryptedToString()

			assert.Equal(t, tt.expected, got, "They should be equal")
		})
	}
}

func TestDecrypt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"example - orig", "1\n2\n-3\n3\n-2\n0\n4", "0 3 -2 1 2 -3 4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf, err := parse(tt.input)
			require.NoError(t, err, "parse failed")

			err = cf.decrypt()
			assert.NoError(t, err, "decrypt should not fail")

			got := cf.decryptedToString()

			assert.Equal(t, tt.expected, got, "They should be equal")
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"example from 0", "0\n4\n1\n2\n-3\n3\n-2", "0 4 1 2 -3 3 -2"},
		{"example - orig", "1\n2\n-3\n3\n-2\n0\n4", "0 4 1 2 -3 3 -2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf, err := parse(tt.input)
			require.NoError(t, err, "parse failed")

			got := cf.decryptedToString()

			assert.Equal(t, cf.zero.Element.value, 0)
			assert.Equal(t, tt.expected, got, "They should be equal")
		})
	}
}

func TestGetCoordinates(t *testing.T) {
	type expected struct {
		one, two, three int
	}

	tests := []struct {
		name       string
		ringValues string
		expected   expected
	}{
		{"example", "1\n2\n-3\n4\n0\n3\n-2", expected{one: 4, two: -3, three: 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf, err := parse(tt.ringValues)
			require.NoError(t, err, "parse failed")
			require.Equal(t, cf.zero.Element.value, 0)

			got1, got2, got3 := cf.getCoordinates()

			assert.Equal(t, tt.expected.one, got1)
			assert.Equal(t, tt.expected.two, got2)
			assert.Equal(t, tt.expected.three, got3)
		})
	}
}
