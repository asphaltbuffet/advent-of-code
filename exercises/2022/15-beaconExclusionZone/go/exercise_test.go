package exercises

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Sensor
	}{
		{
			name:  "Single line",
			input: "Sensor at x=2, y=18: closest beacon is at x=-2, y=15",
			want:  []Sensor{{Location: image.Point{X: 2, Y: 18}, Beacon: image.Point{X: -2, Y: 15}, Dist: 7}},
		},
		{
			name:  "Single line",
			input: "Sensor at x=9, y=16: closest beacon is at x=10, y=16\nSensor at x=13, y=2: closest beacon is at x=15, y=3",
			want: []Sensor{
				{Location: image.Point{X: 9, Y: 16}, Beacon: image.Point{X: 10, Y: 16}, Dist: 1},
				{Location: image.Point{X: 13, Y: 2}, Beacon: image.Point{X: 15, Y: 3}, Dist: 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
