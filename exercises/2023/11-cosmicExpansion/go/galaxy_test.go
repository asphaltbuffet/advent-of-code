package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getEmptyCols(t *testing.T) {
	type args struct {
		img []string
	}
	tests := []struct {
		name string
		args args
		want map[int]bool
	}{
		{
			name: "example",
			args: args{
				img: []string{
					"...#......",
					".......#..",
					"#.........",
					"..........",
					"......#...",
					".#........",
					".........#",
					"..........",
					".......#..",
					"#...#.....",
				},
			},
			want: map[int]bool{2: true, 5: true, 8: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getEmptyCols(tt.args.img))
		})
	}
}

func Test_expandImage(t *testing.T) {
	type args struct {
		img []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "no expansion",
			args: args{
				img: []string{"#...", ".#..", "..#.", "...#"},
			},
			want: []string{"#...", ".#..", "..#.", "...#"},
		},
		{
			name: "single horizontal expansion",
			args: args{
				img: []string{"#...", "..#.", "..#.", "...#"},
			},
			want: []string{"#|..", ".|#.", ".|#.", ".|.#"},
		},
		{
			name: "single vertical expansion",
			args: args{
				img: []string{"#...", "....", ".##.", "...#"},
			},
			want: []string{"#...", "-", ".##.", "...#"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, expandImage(tt.args.img))
		})
	}
}
