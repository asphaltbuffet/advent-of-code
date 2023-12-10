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
				instr: "seeds: 79 14 55 13\n\nseed-to-soil map:\n50 98 2\n52 50 48\n\nsoil-to-fertilizer map:\n0 15 37\n37 52 2\n39 0 15\n\nfertilizer-to-water map:\n49 53 8\n0 11 42\n42 0 7\n57 7 4\n\nwater-to-light map:\n88 18 7\n18 25 70\n\nlight-to-temperature map:\n45 77 23\n81 45 19\n68 64 13\n\ntemperature-to-humidity map:\n0 69 1\n1 0 69\n\nhumidity-to-location map:\n60 56 37\n56 93 4",
			},
			want:      35,
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
