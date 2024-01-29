package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_hash(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"HASH", "HASH", 52},
		{"example step 1", "rn=1", 30},
		{"example step 2", "cm-", 253},
		{"example step 3", "qp=3", 97},
		{"example step 4", "cm=2", 47},
		{"example step 5", "qp-", 14},
		{"example step 6", "pc=4", 180},
		{"example step 7", "ot=9", 9},
		{"example step 8", "ab=5", 197},
		{"example step 9", "pc-", 48},
		{"example step 10", "pc=6", 214},
		{"example step 11", "ot=7", 231},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, hash(tt.input))
		})
	}
}

func Test_parseStep(t *testing.T) {
	type args struct {
		step string
	}
	tests := []struct {
		name      string
		args      args
		want      *Op
		assertion require.ErrorAssertionFunc
	}{
		{
			name: "set operation",
			args: args{
				step: "HASH=1",
			},
			want: &Op{
				Label:       "HASH",
				Box:         52,
				Action:      '=',
				FocalLength: 1,
			},
			assertion: require.NoError,
		},
		{
			name: "remove operation",
			args: args{
				step: "HASH-",
			},
			want: &Op{
				Label:       "HASH",
				Box:         52,
				Action:      '-',
				FocalLength: 0,
			},
			assertion: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseStep(tt.args.step)

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
