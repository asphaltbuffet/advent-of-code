package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_lineToIntSlice(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "all positive",
			args: args{
				line: "0 3 6 9 12 15",
			},
			want: []int{0, 3, 6, 9, 12, 15},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, lineToIntSlice(tt.args.line))
		})
	}
}

func Test_reduceToDiffs(t *testing.T) {
	type args struct {
		in []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "example 1",
			args: args{
				in: []int{0, 3, 6, 9, 12, 15},
			},
			want: []int{3, 3, 3, 3, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, reduceToDiffs(tt.args.in))
		})
	}
}

func Test_calculateReductions(t *testing.T) {
	type args struct {
		in []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "example 1",
			args: args{
				in: []int{0, 3, 6, 9, 12, 15},
			},
			want: 18,
		},
		{
			name: "example 2",
			args: args{
				in: []int{1, 3, 6, 10, 15, 21},
			},
			want: 28,
		},
		{
			name: "example 3",
			args: args{
				in: []int{10, 13, 16, 21, 30, 45},
			},
			want: 68,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, calculateReductions(tt.args.in))
		})
	}
}
