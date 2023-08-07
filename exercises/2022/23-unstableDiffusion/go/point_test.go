package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_point_add(t *testing.T) {
	tests := []struct {
		name string
		p1   point
		p2   point
		want point
	}{
		{
			name: "add 0 to both",
			p1:   point{2, 3},
			p2:   point{0, 0},
			want: point{2, 3},
		},
		{
			name: "add 1 to both",
			p1:   point{2, 3},
			p2:   point{1, 1},
			want: point{3, 4},
		},
		{
			name: "add only to x",
			p1:   point{2, 3},
			p2:   point{1, 0},
			want: point{3, 3},
		},
		{
			name: "add only to y",
			p1:   point{2, 3},
			p2:   point{0, 1},
			want: point{2, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p1.add(tt.p2)

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_point_normalize(t *testing.T) {
	type args struct {
		maxX int
		maxY int
	}

	tests := []struct {
		name string
		p    point
		args args
	}{
		{
			name: "no change",
			p:    point{1, 2},
			args: args{
				maxX: 3,
				maxY: 3,
			},
		},
		{
			name: "x change",
			p:    point{5, 2},
			args: args{
				maxX: 3,
				maxY: 3,
			},
		},
		{
			name: "y change",
			p:    point{1, 5},
			args: args{
				maxX: 3,
				maxY: 3,
			},
		},
		{
			name: "both change",
			p:    point{5, 5},
			args: args{
				maxX: 3,
				maxY: 3,
			},
		},
		{
			name: "both change",
			p:    point{-1, -2},
			args: args{
				maxX: 3,
				maxY: 3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.normalize(tt.args.maxX, tt.args.maxY)

			assert.Less(t, tt.p.x, tt.args.maxX)
			assert.Less(t, tt.p.y, tt.args.maxY)

			assert.Positive(t, tt.p.x)
			assert.Positive(t, tt.p.y)
		})
	}
}
