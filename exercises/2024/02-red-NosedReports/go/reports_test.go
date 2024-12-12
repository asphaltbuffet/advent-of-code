package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReport_isSafe(t *testing.T) {
	tests := []struct {
		name      string
		values    []int
		assertion assert.BoolAssertionFunc
	}{
		{"example 1", []int{7, 6, 4, 2, 1}, assert.True},
		{"example 2", []int{1, 2, 7, 8, 9}, assert.False},
		{"example 3", []int{9, 7, 6, 2, 1}, assert.False},
		{"example 4", []int{1, 3, 2, 4, 5}, assert.False},
		{"example 5", []int{8, 6, 4, 4, 1}, assert.False},
		{"example 6", []int{1, 3, 6, 7, 9}, assert.True},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Report{values: tt.values}

			tt.assertion(t, r.isSafe())
		})
	}
}

func Test_parse(t *testing.T) {
	tests := []struct {
		name      string
		s         string
		want      []Report
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "example all",
			s:    "7 6 4 2 1\n1 2 7 8 9\n9 7 6 2 1\n1 3 2 4 5\n8 6 4 4 1\n1 3 6 7 9",
			want: []Report{
				{
					values: []int{7, 6, 4, 2, 1},
				},
				{
					values: []int{1, 2, 7, 8, 9},
				},
				{
					values: []int{9, 7, 6, 2, 1},
				},
				{
					values: []int{1, 3, 2, 4, 5},
				},
				{
					values: []int{8, 6, 4, 4, 1},
				},
				{
					values: []int{1, 3, 6, 7, 9},
				},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(tt.s)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_removeAt(t *testing.T) {
	type args struct {
		r   []int
		idx int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "first",
			args: args{
				r:   []int{1, 2, 3},
				idx: 0,
			},
			want: []int{2, 3},
		},
		{
			name: "last",
			args: args{
				r:   []int{1, 2, 3},
				idx: 2,
			},
			want: []int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeAt(tt.args.r, tt.args.idx)
			assert.Len(t, got, len(tt.want))
			// assert.Equal(t, tt.want, removeAt(tt.args.r, tt.args.idx))
		})
	}
}
