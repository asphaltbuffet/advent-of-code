package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_loadInput(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      ModConfig
		assertion require.ErrorAssertionFunc
	}{
		{
			name:  "simple example",
			input: "broadcaster -> a, b, c\n%a -> b\n%b -> c\n%c -> inv\n&inv -> a",
			want: ModConfig{
				"broadcaster": &Broadcast{destinations: []string{"a", "b", "c"}},
				"a":           &FlipFlop{ID: "a", State: Off, destinations: []string{"b"}},
				"b":           &FlipFlop{ID: "b", State: Off, destinations: []string{"c"}},
				"c":           &FlipFlop{ID: "c", State: Off, destinations: []string{"inv"}},
				"inv":         &Conjunction{ID: "inv", destinations: []string{"a"}, memory: map[string]Pulse{"c": Low}},
			},
			assertion: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadInput(tt.input)

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
