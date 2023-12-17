package exercises

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed input.txt
var input string

func TestExercise_Two(t *testing.T) {
	type args struct {
		instr string
	}
	tests := []struct {
		name      string
		e         Exercise
		args      args
		want      any
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "solution",
			args: args{
				instr: input,
			},
			want:      nil,
			assertion: assert.NoError,
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
