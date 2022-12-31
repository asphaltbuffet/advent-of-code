package aoc22_09_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_09"
)

func Test_Day9Part1(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 1 Example", []string{
			"R 4",
			"U 4",
			"L 3",
			"D 1",
			"R 4",
			"D 1",
			"L 5",
			"R 2",
		}, "13"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_09.D9P1(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Day9Part2(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"Part 2 Example", []string{
			"R 4",
			"U 4",
			"L 3",
			"D 1",
			"R 4",
			"D 1",
			"L 5",
			"R 2",
		}, "1"},
		{"Part 2 Example - Long", []string{
			"R 5",
			"U 8",
			"L 8",
			"D 3",
			"R 17",
			"D 10",
			"L 25",
			"U 20",
		}, "36"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_09.D9P2(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_CalculateMovement(t *testing.T) {
	type args struct {
		h, t aoc22_09.Point
	}

	tests := []struct {
		name string
		args args
		want aoc22_09.Point
	}{
		{"Same points", args{aoc22_09.Point{X: 1, Y: 1}, aoc22_09.Point{X: 1, Y: 1}}, aoc22_09.Point{X: 0, Y: 0}},
		{"Diagonal", args{aoc22_09.Point{X: 4, Y: 1}, aoc22_09.Point{X: 3, Y: 0}}, aoc22_09.Point{X: 0, Y: 0}},
		{"Right", args{aoc22_09.Point{X: 2, Y: 2}, aoc22_09.Point{X: 0, Y: 2}}, aoc22_09.Point{X: 1, Y: 0}},
		{"Left", args{aoc22_09.Point{X: 2, Y: 2}, aoc22_09.Point{X: 4, Y: 2}}, aoc22_09.Point{X: -1, Y: 0}},
		{"Up", args{aoc22_09.Point{X: 4, Y: 3}, aoc22_09.Point{X: 4, Y: 1}}, aoc22_09.Point{X: 0, Y: 1}},
		{"Down", args{aoc22_09.Point{X: 4, Y: 3}, aoc22_09.Point{X: 4, Y: 5}}, aoc22_09.Point{X: 0, Y: -1}},
		{"Up and Right", args{aoc22_09.Point{X: 4, Y: 2}, aoc22_09.Point{X: 3, Y: 0}}, aoc22_09.Point{X: 1, Y: 1}},
		{"Down and Right", args{aoc22_09.Point{X: 4, Y: 0}, aoc22_09.Point{X: 2, Y: 1}}, aoc22_09.Point{X: 1, Y: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aoc22_09.CalculateMovement(tt.args.h, tt.args.t)
			assert.Equal(t, tt.want, got)
		})
	}
}
