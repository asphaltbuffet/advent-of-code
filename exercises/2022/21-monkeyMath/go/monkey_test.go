package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]string
	}{
		{"example - one line", "root: pppw + sjmn", map[string]string{"root": "pppw + sjmn"}},
		{
			"example - mult", "root: pppw + sjmn\ndbpl: 5\ncczh: sllz + lgvd\nzczc: 2",
			map[string]string{
				"root": "dbpl + cczh",
				"dbpl": "5",
				"cczh": "zczc + dbpl",
				"zczc": "2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parse(tt.input)

			assert.Equal(t, tt.want, got, "They should be equal")
		})
	}
}

func TestCalc(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]string
		want  int
	}{
		{"example - one line", map[string]string{"root": "5"}, 5},
		{
			"example - mult",
			map[string]string{
				"root": "dbpl + zczc",
				"dbpl": "5",
				"zczc": "2",
			},
			7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calc("root", tt.input, make(map[string]int))
			require.NoError(t, err)

			assert.Equal(t, tt.want, got, "They should be equal")
		})
	}
}
