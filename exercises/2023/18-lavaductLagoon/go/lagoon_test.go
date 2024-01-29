package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_shoelace(t *testing.T) {
	type args struct {
		points []Point
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "wiki example",
			args: args{
				points: []Point{
					{1, 6}, {3, 1}, {7, 2}, {4, 4}, {8, 5},
				},
			},
			want: 16.5,
		},
		{
			name: "part1 example",
			args: args{
				points: []Point{
					{10, 10}, {16, 10}, {16, 5}, {14, 5}, {14, 3}, {16, 3}, {16, 1}, {11, 1}, {11, 3}, {10, 3}, {10, 5}, {12, 5}, {12, 8}, {10, 8},
				},
			},
			want: 62.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, shoelace(tt.args.points), 0.0001)
		})
	}
}
