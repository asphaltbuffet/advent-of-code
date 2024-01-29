package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseHailstones(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name      string
		args      args
		want      []Hailstone
		assertion require.ErrorAssertionFunc
	}{
		{
			name: "solitary good line",
			args: args{
				input: "19, 13, 30 @ -2,  1, -2",
			},
			want: []Hailstone{
				{
					pos: Vector3D{19.0, 13.0, 30.0},
					vel: Vector3D{-2.0, 1.0, -2.0},
				},
			},
			assertion: require.NoError,
		},
		{
			name: "multiple good lines",
			args: args{
				input: "19, 13, 30 @ -2,  1, -2\n18, 19, 22 @ -1, -1, -2",
			},
			want: []Hailstone{
				{
					pos: Vector3D{19.0, 13.0, 30.0},
					vel: Vector3D{-2.0, 1.0, -2.0},
				},
				{
					pos: Vector3D{18.0, 19.0, 22.0},
					vel: Vector3D{-1.0, -1.0, -2.0},
				},
			},
			assertion: require.NoError,
		},
		{
			name: "solitary line with bad position",
			args: args{
				input: "nineteen, 13, 30 @ -2,  1, -2",
			},
			want:      nil,
			assertion: require.Error,
		},
		{
			name: "solitary line with bad velocity",
			args: args{
				input: "19, 13, 30 @ -2,  one, -2",
			},
			want:      nil,
			assertion: require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// we don't care about the min/max values here
			got, err := parseInput(tt.args.input)

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
