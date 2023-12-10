package exercises

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseRaces(t *testing.T) {
	t.Parallel()

	type args struct {
		s string
	}

	tests := []struct {
		name string
		args args
		want []Race
	}{
		{
			name: "empty string",
			args: args{s: ""},
			want: nil,
		},
		{
			name: "single race",
			args: args{s: "Time:      7\nDistance:  9"},
			want: []Race{{ID: 0, Time: 7, Distance: 9}},
		},
		{
			name: "multiple races",
			args: args{s: "Time:      7  15\nDistance:  9  40"},
			want: []Race{{ID: 0, Time: 7, Distance: 9}, {ID: 1, Time: 15, Distance: 40}},
		},
		{
			name: "invalid input",
			args: args{s: "Time:      7  foo   30\nDistance:  9  40  200"},
			want: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.ElementsMatch(t, tt.want, parseRaces(tt.args.s))
		})
	}
}

func TestRace_CalculateDistances(t *testing.T) {
	t.Parallel()

	type wants struct {
		n     int
		dists []int
	}

	tests := []struct {
		name string
		r    *Race
		want wants
	}{
		{
			name: "short",
			r:    &Race{ID: 0, Time: 7, Distance: 9},
			want: wants{n: 4, dists: []int{0, 6, 10, 12, 12, 10, 6, 0}},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotN, gotDists := tt.r.CalculateDistances()

			assert.Equal(t, tt.want.n, gotN)
			assert.Equal(t, tt.want.dists, gotDists)

			// slip in a test for CountFasterTimes() while we're here
			assert.Equal(t, tt.want.n, tt.r.CountFasterTimes())
		})
	}
}

func Test_parseBigRace(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *Race
	}{
		{
			name: "example",
			args: args{s: "Time:      7  15   30\nDistance:  9  40  200"},
			want: &Race{
				ID:       0,
				Time:     71530,
				Distance: 940200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parseBigRace(tt.args.s))
		})
	}
}

var exercise = Exercise{}

var table = []struct {
	s string
}{
	{s: "Time:      7  15   30\nDistance:  9  40  200"},
	{s: "Time:        60     80     86     76\nDistance:   601   1163   1559   1300"},
	{s: "Time:        41     77     70     96\nDistance:   249   1362   1127   1011"},
}

func BenchmarkPartOne(b *testing.B) {
	for nn, bb := range table {
		b.Run(fmt.Sprintf("input_%d", nn), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				exercise.One(bb.s)
			}
		})
	}
}

func BenchmarkPartTwo(b *testing.B) {
	for nn, bb := range table {
		b.Run(fmt.Sprintf("input_%d", nn), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				exercise.Two(bb.s)
			}
		})
	}
}

// var table = []struct {
// 	r Race
// }{
// 	{Race{ID: 0, Time: 7, Distance: 9}},
// 	{Race{ID: 1, Time: 15, Distance: 40}},
// 	{Race{ID: 2, Time: 30, Distance: 200}},
// }

// func BenchmarkCountWins(b *testing.B) {
// 	for _, bb := range table {
// 		b.Run(fmt.Sprintf("example_race_%d", bb.r.ID), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				bb.r.CountFasterTimes()
// 			}
// 		})
// 	}
// }

// func BenchmarkLoops(b *testing.B) {
// 	for _, bb := range table {
// 		b.Run(fmt.Sprintf("example_race_%d", bb.r.ID), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				bb.r.CalculateDistances()
// 			}
// 		})
// 	}
// }
