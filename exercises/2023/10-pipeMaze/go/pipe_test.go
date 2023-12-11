package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetShape(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want PipeShape
	}{
		{"ground", args{'.'}, NoPipe},
		{"vertical", args{'|'}, VerticalPipe},
		{"horizontal", args{'-'}, HorizontalPipe},
		{"NE corner", args{'L'}, NECornerPipe},
		{"NW corner", args{'J'}, NWCornerPipe},
		{"SW corner", args{'7'}, SWCornerPipe},
		{"SE corner", args{'F'}, SECornerPipe},
		{"start", args{'S'}, StartPipe},
		{"invalid", args{'1'}, InvalidPipe},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetShape(tt.args.r))
		})
	}
}

func Test_parseInput(t *testing.T) {
	type args struct {
		s string
	}

	type wants struct {
		pipes     map[Point]Pipe
		start     Point
		assertion require.ErrorAssertionFunc
	}

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "no pipes",
			args: args{
				s: "...\n...\n...",
			},
			wants: wants{
				pipes: map[Point]Pipe{
					// {0, 0}: {Pos: Point{0, 0}, Shape: NoPipe},
					// {1, 0}: {Pos: Point{1, 0}, Shape: NoPipe},
					// {2, 0}: {Pos: Point{2, 0}, Shape: NoPipe},
					// {0, 1}: {Pos: Point{0, 1}, Shape: NoPipe},
					// {1, 1}: {Pos: Point{1, 1}, Shape: NoPipe},
					// {2, 1}: {Pos: Point{2, 1}, Shape: NoPipe},
					// {0, 2}: {Pos: Point{0, 2}, Shape: NoPipe},
					// {1, 2}: {Pos: Point{1, 2}, Shape: NoPipe},
					// {2, 2}: {Pos: Point{2, 2}, Shape: NoPipe},
				},
				start:     Point{},
				assertion: require.NoError,
			},
		},
		{
			name: "loop without start",
			args: args{
				s: "F-7\n|.|\nL-J",
			},
			wants: wants{
				pipes: map[Point]Pipe{
					{0, 0}: {Pos: Point{0, 0}, Shape: SECornerPipe, To: []Point{{0, 1}, {1, 0}}},
					{1, 0}: {Pos: Point{1, 0}, Shape: HorizontalPipe, To: []Point{{0, 0}, {2, 0}}},
					{2, 0}: {Pos: Point{2, 0}, Shape: SWCornerPipe, To: []Point{{1, 0}, {2, 1}}},
					{0, 1}: {Pos: Point{0, 1}, Shape: VerticalPipe, To: []Point{{0, 0}, {0, 2}}},
					{2, 1}: {Pos: Point{2, 1}, Shape: VerticalPipe, To: []Point{{2, 0}, {2, 2}}},
					{0, 2}: {Pos: Point{0, 2}, Shape: NECornerPipe, To: []Point{{0, 1}, {1, 2}}},
					{1, 2}: {Pos: Point{1, 2}, Shape: HorizontalPipe, To: []Point{{0, 2}, {2, 2}}},
					{2, 2}: {Pos: Point{2, 2}, Shape: NWCornerPipe, To: []Point{{1, 2}, {2, 1}}},
				},
				start:     Point{},
				assertion: require.NoError,
			},
		},
		{
			name: "loop with start",
			args: args{
				s: "FS7\n|.|\nL-J",
			},
			wants: wants{
				pipes: map[Point]Pipe{
					{0, 0}: {Pos: Point{0, 0}, Shape: SECornerPipe, To: []Point{{0, 1}, {1, 0}}},
					{1, 0}: {Pos: Point{1, 0}, Shape: StartPipe, To: []Point{{0, 0}, {1, 1}, {2, 0}, {1, -1}}},
					{2, 0}: {Pos: Point{2, 0}, Shape: SWCornerPipe, To: []Point{{1, 0}, {2, 1}}},
					{0, 1}: {Pos: Point{0, 1}, Shape: VerticalPipe, To: []Point{{0, 0}, {0, 2}}},
					{2, 1}: {Pos: Point{2, 1}, Shape: VerticalPipe, To: []Point{{2, 0}, {2, 2}}},
					{0, 2}: {Pos: Point{0, 2}, Shape: NECornerPipe, To: []Point{{0, 1}, {1, 2}}},
					{1, 2}: {Pos: Point{1, 2}, Shape: HorizontalPipe, To: []Point{{0, 2}, {2, 2}}},
					{2, 2}: {Pos: Point{2, 2}, Shape: NWCornerPipe, To: []Point{{1, 2}, {2, 1}}},
				},
				start:     Point{1, 0},
				assertion: require.NoError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPipes, gotStart, gotErr := parseInput(tt.args.s)

			tt.wants.assertion(t, gotErr)

			verifyPipes(t, tt.wants.pipes, gotPipes)
			assert.Equal(t, tt.wants.start, gotStart)
		})
	}
}

func verifyPipes(t *testing.T, expected map[Point]Pipe, actual map[Point]Pipe) {
	t.Helper()

	for k, pipe := range expected {
		require.NotEmpty(t, actual[k], "missing pipe at %v", k)

		assert.ElementsMatch(t, pipe.To, actual[k].To)
		assert.Equal(t, pipe.Shape, actual[k].Shape)
		assert.Equal(t, pipe.Pos, actual[k].Pos)
	}
}

func Test_getConnections(t *testing.T) {
	type args struct {
		pos   Point
		shape PipeShape
	}
	tests := []struct {
		name string
		args args
		want []Point
	}{
		{
			name: "all directions",
			args: args{
				pos:   Point{1, 1},
				shape: StartPipe,
			},
			want: []Point{{1, 0}, {1, 2}, {0, 1}, {2, 1}},
		},
		{
			name: "north and south",
			args: args{
				pos:   Point{1, 1},
				shape: VerticalPipe,
			},
			want: []Point{{1, 0}, {1, 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.want, getConnections(tt.args.pos, tt.args.shape))
		})
	}
}
