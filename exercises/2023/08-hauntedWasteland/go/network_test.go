package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseMovement(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "empty",
			args: args{s: ""},
			want: 0,
		},
		{
			name: "left right",
			args: args{s: "LR"},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseMovement(tt.args.s)
			assert.Equal(t, tt.want, got.Len())
		})
	}
}

func Test_parseMap(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want map[string]map[Move]string
	}{
		{
			name: "single line",
			args: args{s: "AAA = (BBB, CCC)"},
			want: map[string]map[Move]string{"AAA": {MoveLeft: "BBB", MoveRight: "CCC"}},
		},
		{
			name: "three lines",
			args: args{s: "AAA = (BBB, CCC)\nBBB = (DDD, EEE)\nCCC = (ZZZ, GGG)"},
			want: map[string]map[Move]string{
				"AAA": {MoveLeft: "BBB", MoveRight: "CCC"},
				"BBB": {MoveLeft: "DDD", MoveRight: "EEE"},
				"CCC": {MoveLeft: "ZZZ", MoveRight: "GGG"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parseMap(tt.args.s))
		})
	}
}
