package exercises

import (
	"testing"
)

var exercise = Exercise{}

var table = []struct {
	name  string
	input string
}{
	{name: "example", input: "32T3K 765\nT55J5 684\nKK677 28\nKTJJT 220\nQQQJA 483"},
}

func BenchmarkPartOne(b *testing.B) {
	for _, bb := range table {
		b.Run(bb.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				exercise.One(bb.input)
			}
		})
	}
}

func BenchmarkPartTwo(b *testing.B) {
	for _, bb := range table {
		b.Run(bb.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				exercise.Two(bb.input)
			}
		})
	}
}
