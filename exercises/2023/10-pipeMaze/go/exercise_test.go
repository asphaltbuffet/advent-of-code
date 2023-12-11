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
			// .....
			// .S-7.
			// .|.|.
			// .L-J.
			// .....
			name: "easy example",
			args: args{
				instr: ".....\n.S-7.\n.|.|.\n.L-J.\n.....",
			},
			want:      4,
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
