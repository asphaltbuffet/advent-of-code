package aoc22_14_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_14"
)

func Test_ParseToken(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  aoc22_14.Point
	}{
		{"No trim", "498,4", aoc22_14.Point{498, 4}},
		{"Trim right", "498,4 ", aoc22_14.Point{498, 4}},
		{"Trim left", " 498,4", aoc22_14.Point{498, 4}},
		{"Full trim", " 498,4 ", aoc22_14.Point{498, 4}},
		{"No separator", "4984", aoc22_14.Point{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := aoc22_14.ParseToken(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_GetPointsBetween(t *testing.T) {
	type args struct {
		left  aoc22_14.Point
		right aoc22_14.Point
	}

	tests := []struct {
		name  string
		input args
		want  []aoc22_14.Point
	}{
		{
			"Horiz - None between",
			args{aoc22_14.Point{500, 10}, aoc22_14.Point{501, 10}},
			[]aoc22_14.Point{
				{500, 10},
				{501, 10},
			},
		},
		{
			"Horiz - Three between",
			args{aoc22_14.Point{500, 10}, aoc22_14.Point{504, 10}},
			[]aoc22_14.Point{
				{500, 10},
				{501, 10},
				{502, 10},
				{503, 10},
				{504, 10},
			},
		},
		{
			"Horiz Bckwd - Three between",
			args{aoc22_14.Point{504, 10}, aoc22_14.Point{500, 10}},
			[]aoc22_14.Point{
				{504, 10},
				{503, 10},
				{502, 10},
				{501, 10},
				{500, 10},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := aoc22_14.GetPointsBetween(tt.input.left, tt.input.right)
			assert.Equal(t, tt.want, got)
		})
	}
}
