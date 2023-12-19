package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	exampleOneInputOne = "#.##..##.\n..#.##.#.\n##......#\n##......#\n..#.##.#.\n..##..##.\n#.#.##.#."
	exampleOneInputTwo = "#...##..#\n#....#..#\n..##..###\n#####.##.\n#####.##.\n..##..###\n#....#..#"

	exampleOneInput = exampleOneInputOne + "\n\n" + exampleOneInputTwo
)

func Test_parsePattern(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *Pattern
	}{
		{
			"example one",
			args{s: exampleOneInputOne},
			&Pattern{
				Row: []string{"#.##..##.", "..#.##.#.", "##......#", "##......#", "..#.##.#.", "..##..##.", "#.#.##.#."},
				Col: []string{"#.##..#", "..##...", "##..###", "#....#.", ".#..#.#", ".#..#.#", "#....#.", "##..###", "..##..."},
			},
		},
		{
			"example two", args{s: exampleOneInputTwo}, &Pattern{
				Row: []string{"#...##..#", "#....#..#", "..##..###", "#####.##.", "#####.##.", "..##..###", "#....#..#"},
				Col: []string{"##.##.#", "...##..", "..####.", "..####.", "#..##..", "##....#", "..####.", "..####.", "###..##"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parsePattern(tt.args.s))
		})
	}
}

func Test_getPatterns(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want []*Pattern
	}{
		{
			name: "example",
			args: args{input: exampleOneInput},
			want: []*Pattern{
				{
					Row: []string{
						"#.##..##.",
						"..#.##.#.",
						"##......#",
						"##......#",
						"..#.##.#.",
						"..##..##.",
						"#.#.##.#.",
					},
					Col: []string{
						"#.##..#",
						"..##...",
						"##..###",
						"#....#.",
						".#..#.#",
						".#..#.#",
						"#....#.",
						"##..###",
						"..##...",
					},
				},
				{
					Row: []string{
						"#...##..#",
						"#....#..#",
						"..##..###",
						"#####.##.",
						"#####.##.",
						"..##..###",
						"#....#..#",
					},
					Col: []string{
						"##.##.#",
						"...##..",
						"..####.",
						"..####.",
						"#..##..",
						"##....#",
						"..####.",
						"..####.",
						"###..##",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getPatterns(tt.args.input))
		})
	}
}

func TestPattern_getHorizontalPlane(t *testing.T) {
	tests := []struct {
		name       string
		p          *Pattern
		hasSmudge  bool
		wantRow    int
		wantResult int
	}{
		{
			name: "example pattern one",
			p: &Pattern{
				Row: []string{"#.##..##.", "..#.##.#.", "##......#", "##......#", "..#.##.#.", "..##..##.", "#.#.##.#."},
				Col: []string{"#.##..#", "..##...", "##..###", "#....#.", ".#..#.#", ".#..#.#", "#....#.", "##..###", "..##..."},
			},
			hasSmudge:  false,
			wantRow:    0,
			wantResult: -1,
		},
		{
			name: "example pattern two",
			p: &Pattern{
				Row: []string{"#...##..#", "#....#..#", "..##..###", "#####.##.", "#####.##.", "..##..###", "#....#..#"},
				Col: []string{"##.##.#", "...##..", "..####.", "..####.", "#..##..", "##....#", "..####.", "..####.", "###..##"},
			},
			hasSmudge:  false,
			wantRow:    4,
			wantResult: 0,
		},
		{
			name: "smudge pattern one",
			p: &Pattern{
				Row: []string{"#.##..##.", "..#.##.#.", "##......#", "##......#", "..#.##.#.", "..##..##.", "#.#.##.#."},
				Col: []string{"#.##..#", "..##...", "##..###", "#....#.", ".#..#.#", ".#..#.#", "#....#.", "##..###", "..##..."},
			},
			hasSmudge:  true,
			wantRow:    3,
			wantResult: 1,
		},
		{
			name: "smudge pattern two",
			p: &Pattern{
				Row: []string{"#...##..#", "#....#..#", "..##..###", "#####.##.", "#####.##.", "..##..###", "#....#..#"},
				Col: []string{"##.##.#", "...##..", "..####.", "..####.", "#..##..", "##....#", "..####.", "..####.", "###..##"},
			},
			hasSmudge:  true,
			wantRow:    1,
			wantResult: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRow, gotResult := findMirror(tt.p.Row, tt.hasSmudge)

			assert.Equal(t, tt.wantRow, gotRow)
			assert.Equal(t, tt.wantResult, gotResult)
		})
	}
}

func TestPattern_getVerticalPlane(t *testing.T) {
	tests := []struct {
		name       string
		p          *Pattern
		hasSmudge  bool
		wantCol    int
		wantResult int
	}{
		{
			name: "example pattern one",
			p: &Pattern{
				Row: []string{"#.##..##.", "..#.##.#.", "##......#", "##......#", "..#.##.#.", "..##..##.", "#.#.##.#."},
				Col: []string{"#.##..#", "..##...", "##..###", "#....#.", ".#..#.#", ".#..#.#", "#....#.", "##..###", "..##..."},
			},
			hasSmudge:  false,
			wantCol:    5,
			wantResult: 0,
		},
		{
			name: "example pattern two",
			p: &Pattern{
				Row: []string{"#...##..#", "#....#..#", "..##..###", "#####.##.", "#####.##.", "..##..###", "#....#..#"},
				Col: []string{"##.##.#", "...##..", "..####.", "..####.", "#..##..", "##....#", "..####.", "..####.", "###..##"},
			},
			hasSmudge:  false,
			wantCol:    0,
			wantResult: -1,
		},
		{
			name: "smudge pattern one",
			p: &Pattern{
				Row: []string{"#.##..##.", "..#.##.#.", "##......#", "##......#", "..#.##.#.", "..##..##.", "#.#.##.#."},
				Col: []string{"#.##..#", "..##...", "##..###", "#....#.", ".#..#.#", ".#..#.#", "#....#.", "##..###", "..##..."},
			},
			hasSmudge:  true,
			wantCol:    0,
			wantResult: -1,
		},
		{
			name: "smudge pattern two",
			p: &Pattern{
				Row: []string{"#...##..#", "#....#..#", "..##..###", "#####.##.", "#####.##.", "..##..###", "#....#..#"},
				Col: []string{"##.##.#", "...##..", "..####.", "..####.", "#..##..", "##....#", "..####.", "..####.", "###..##"},
			},
			hasSmudge:  true,
			wantCol:    0,
			wantResult: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCol, gotResult := findMirror(tt.p.Col, tt.hasSmudge)

			assert.Equal(t, tt.wantCol, gotCol)
			assert.Equal(t, tt.wantResult, gotResult)
		})
	}
}
