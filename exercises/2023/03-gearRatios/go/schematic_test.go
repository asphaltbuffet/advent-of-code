package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_schematic_parseLine(t *testing.T) {
	type args struct {
		line string
	}

	type wants struct {
		symbols map[point]bool
		parts   map[int]part
	}

	tests := []struct {
		name      string
		args      args
		wants     wants
		assertion require.ErrorAssertionFunc
	}{
		{
			name: "two numbers",
			args: args{
				line: "467..114..",
			},
			wants: wants{
				parts: map[int]part{
					0: {value: "467", points: []point{}, id: 0},
					1: {value: "114", points: []point{}, id: 1},
				},
				symbols: map[point]bool{},
			},
			assertion: require.NoError,
		},
		{
			name: "symbol after number",
			args: args{
				line: "467*......",
			},
			wants: wants{
				parts: map[int]part{
					0: {value: "467", points: []point{}, id: 0},
				},
				symbols: map[point]bool{
					{3, 0}: true,
				},
			},
			assertion: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schm := &schematic{
				parts:    make(map[int]part),
				isSymbol: make(map[point]bool),
			}

			tt.assertion(t, schm.parseLine(tt.args.line, 0))
			assert.Equal(t, tt.wants.parts, schm.parts)
			assert.Equal(t, tt.wants.symbols, schm.isSymbol)
		})
	}
}

func Test_schematic_addPart(t *testing.T) {
	type args struct {
		in     string
		origin point
	}

	type wants struct {
		parts map[int]part
		next  int
	}

	tests := []struct {
		name  string
		s     *schematic
		args  args
		wants wants
	}{
		{
			name: "first part",
			s: &schematic{
				parts: make(map[int]part),
			},
			args: args{
				in:     "467..114..",
				origin: point{0, 0},
			},
			wants: wants{
				parts: map[int]part{
					0: {
						value:  "467",
						points: []point{{0, 0}, {1, 0}, {2, 0}},
						bounds: []point{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {2, -1}, {2, 1}, {3, -1}, {3, 0}, {3, 1}},
						id:     0,
					},
				},
				next: 3,
			},
		},
		{
			name: "another part",
			s: &schematic{
				parts: map[int]part{0: {
					value:  "467",
					points: []point{{0, 0}, {1, 0}, {2, 0}},
					bounds: []point{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {2, -1}, {2, 1}, {3, -1}, {3, 0}, {3, 1}},
					id:     0,
				}},
			},
			args: args{in: "114..", origin: point{5, 0}},
			wants: wants{
				parts: map[int]part{
					0: {
						value:  "467",
						points: []point{{0, 0}, {1, 0}, {2, 0}},
						bounds: []point{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {2, -1}, {2, 1}, {3, -1}, {3, 0}, {3, 1}},
						id:     0,
					},
					1: {
						value:  "114",
						points: []point{{5, 0}, {6, 0}, {7, 0}},
						bounds: []point{{4, -1}, {4, 0}, {4, 1}, {5, -1}, {5, 1}, {6, -1}, {6, 1}, {7, -1}, {7, 1}, {8, -1}, {8, 0}, {8, 1}},
						id:     1,
					},
				},

				next: 3,
			},
		},
		{
			name: "symbol at start",
			s: &schematic{
				// engine: make(map[point]string),
				parts: map[int]part{0: {
					value:  "467",
					points: []point{{0, 0}, {1, 0}, {2, 0}},
					bounds: []point{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {2, -1}, {2, 1}, {3, -1}, {3, 0}, {3, 1}},
					id:     0,
				}},
			},
			args: args{in: "*114..", origin: point{4, 0}},
			wants: wants{
				parts: map[int]part{0: {
					value:  "467",
					points: []point{{0, 0}, {1, 0}, {2, 0}},
					bounds: []point{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {2, -1}, {2, 1}, {3, -1}, {3, 0}, {3, 1}},
					id:     0,
				}},
				next: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.addPart(tt.args.in, tt.args.origin)

			require.Len(t, tt.s.parts, len(tt.wants.parts), "unexpected number of parts")
			for i := 0; i < len(tt.s.parts); i++ {
				assert.ElementsMatchf(t, tt.wants.parts[i].points, tt.s.parts[i].points, "part[%d] points do not match", i)
				assert.ElementsMatchf(t, tt.wants.parts[i].bounds, tt.s.parts[i].bounds, "part[%d] bounds do not match", i)
				assert.Equalf(t, tt.wants.parts[i].value, tt.s.parts[i].value, "part[%d] values does not match", i)
				assert.Equalf(t, tt.wants.parts[i].id, tt.s.parts[i].id, "part[%d] ids does not match", i)
			}

			assert.Equal(t, tt.wants.next, got)
		})
	}
}
