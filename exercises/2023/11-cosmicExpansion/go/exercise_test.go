package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExercise_One(t *testing.T) {
	type args struct {
		instr string
	}
	tests := []struct {
		name      string
		args      args
		want      any
		assertion require.ErrorAssertionFunc
	}{
		{
			name: "example",
			args: args{
				instr: "...#......\n.......#..\n#.........\n..........\n......#...\n.#........\n.........#\n..........\n.......#..\n#...#.....",
			},
			want:      374,
			assertion: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Exercise{}

			got, err := e.One(tt.args.instr)

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExercise_Two(t *testing.T) {
	type args struct {
		instr string
	}
	tests := []struct {
		name      string
		args      args
		want      any
		assertion require.ErrorAssertionFunc
	}{
		{
			name: "example",
			args: args{
				instr: "...#......\n.......#..\n#.........\n..........\n......#...\n.#........\n.........#\n..........\n.......#..\n#...#.....",
			},
			want:      82000210, // 8410 -> 100x; 1030 -> 10x
			assertion: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Exercise{}

			got, err := e.Two(tt.args.instr)

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
