package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NextInvalid(t *testing.T) {

	tests := []struct {
		name string
		arg  int
		want int
	}{
		{"11", 11, 22},
		{"98", 98, 99},
		{"99", 99, 1010},
		{"111", 111, 1010},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NextInvalid(tt.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_IsInvalid(t *testing.T) {
	tests := []struct {
		name string
		arg  int
		want bool
	}{
		{"valid even", 10, false},
		{"valid odd", 101, false},
		{"invalid short", 11, true},
		{"invalid long", 420420, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsInvalid(tt.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}
