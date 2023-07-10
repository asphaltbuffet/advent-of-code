//go:build test
// +build test

package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalculateMovement(t *testing.T) {
	type args struct {
		h, t Point
	}

	tests := []struct {
		name string
		args args
		want Point
	}{
		{"Same points", args{Point{X: 1, Y: 1}, Point{X: 1, Y: 1}}, Point{X: 0, Y: 0}},
		{"Diagonal", args{Point{X: 4, Y: 1}, Point{X: 3, Y: 0}}, Point{X: 0, Y: 0}},
		{"Right", args{Point{X: 2, Y: 2}, Point{X: 0, Y: 2}}, Point{X: 1, Y: 0}},
		{"Left", args{Point{X: 2, Y: 2}, Point{X: 4, Y: 2}}, Point{X: -1, Y: 0}},
		{"Up", args{Point{X: 4, Y: 3}, Point{X: 4, Y: 1}}, Point{X: 0, Y: 1}},
		{"Down", args{Point{X: 4, Y: 3}, Point{X: 4, Y: 5}}, Point{X: 0, Y: -1}},
		{"Up and Right", args{Point{X: 4, Y: 2}, Point{X: 3, Y: 0}}, Point{X: 1, Y: 1}},
		{"Down and Right", args{Point{X: 4, Y: 0}, Point{X: 2, Y: 1}}, Point{X: 1, Y: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateMovement(tt.args.h, tt.args.t)
			assert.Equal(t, tt.want, got)
		})
	}
}
