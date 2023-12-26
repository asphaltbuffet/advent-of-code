package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseInput(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name      string
		args      args
		want      map[Point]rune
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "small example",
			args: args{
				s: ".|...\n|.-.\\\n.....\n.....\n.....\n.....\n..../\n.-.-/\n.|...\n..//.",
			},
			want: map[Point]rune{
				{0, 0}: '.', {1, 0}: '|', {2, 0}: '.', {3, 0}: '.', {4, 0}: '.',
				{0, 1}: '|', {1, 1}: '.', {2, 1}: '-', {3, 1}: '.', {4, 1}: '\\',
				{0, 2}: '.', {1, 2}: '.', {2, 2}: '.', {3, 2}: '.', {4, 2}: '.',
				{0, 3}: '.', {1, 3}: '.', {2, 3}: '.', {3, 3}: '.', {4, 3}: '.',
				{0, 4}: '.', {1, 4}: '.', {2, 4}: '.', {3, 4}: '.', {4, 4}: '.',
				{0, 5}: '.', {1, 5}: '.', {2, 5}: '.', {3, 5}: '.', {4, 5}: '.',
				{0, 6}: '.', {1, 6}: '.', {2, 6}: '.', {3, 6}: '.', {4, 6}: '/',
				{0, 7}: '.', {1, 7}: '-', {2, 7}: '.', {3, 7}: '-', {4, 7}: '/',
				{0, 8}: '.', {1, 8}: '|', {2, 8}: '.', {3, 8}: '.', {4, 8}: '.',
				{0, 9}: '.', {1, 9}: '.', {2, 9}: '/', {3, 9}: '/', {4, 9}: '.',
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseInput(tt.args.s)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
