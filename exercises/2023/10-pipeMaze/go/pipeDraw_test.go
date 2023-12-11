package exercises

import (
	_ "embed"
	"testing"
)

//go:embed input.txt
var input string

func Test_drawPipes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				s: input,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			drawPipes(tt.args.s)
		})
	}
}
