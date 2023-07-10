package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseToken(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Point
	}{
		{"No trim", "498,4", Point{498, 4}},
		{"Trim right", "498,4 ", Point{498, 4}},
		{"Trim left", " 498,4", Point{498, 4}},
		{"Full trim", " 498,4 ", Point{498, 4}},
		{"No separator", "4984", Point{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := ParseToken(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_GetPointsBetween(t *testing.T) {
	type args struct {
		left  Point
		right Point
	}

	tests := []struct {
		name  string
		input args
		want  []Point
	}{
		{
			"Horiz - None between",
			args{Point{500, 10}, Point{501, 10}},
			[]Point{
				{500, 10},
				{501, 10},
			},
		},
		{
			"Horiz - Three between",
			args{Point{500, 10}, Point{504, 10}},
			[]Point{
				{500, 10},
				{501, 10},
				{502, 10},
				{503, 10},
				{504, 10},
			},
		},
		{
			"Horiz Bckwd - Three between",
			args{Point{504, 10}, Point{500, 10}},
			[]Point{
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
			got, _ := GetPointsBetween(tt.input.left, tt.input.right)
			assert.Equal(t, tt.want, got)
		})
	}
}
