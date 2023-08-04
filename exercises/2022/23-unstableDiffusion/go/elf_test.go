package exercises

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_unstableDiffusion(t *testing.T) {
	type args struct {
		elfCoords map[point]string
		part      int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := diffuse(tt.args.elfCoords, tt.args.part)
			if (err != nil) != tt.wantErr {
				t.Errorf("unstableDiffusion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("unstableDiffusion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashState(t *testing.T) {
	type args struct {
		elfCoords map[point]string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic",
			args: args{
				elfCoords: map[point]string{{0, 0}: "#", {1, 2}: "#"},
			},
			want: "7f2a18232af2401ec9fd3af7574833b3a0700e18",
		},
		{
			name: "basic - reversed",
			args: args{
				elfCoords: map[point]string{{1, 2}: "#", {0, 0}: "#"},
			},
			want: "7f2a18232af2401ec9fd3af7574833b3a0700e18",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hashState(tt.args.elfCoords)

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_planElfMoves(t *testing.T) {
	type args struct {
		elfCoords      map[point]string
		diffStartIndex int
	}

	tests := []struct {
		name  string
		args  args
		want  map[point]point
		want1 map[point]int
	}{
		{
			name: "small example",
			args: args{
				elfCoords:      map[point]string{{2, 1}: "#", {3, 1}: "#", {2, 2}: "#", {2, 4}: "#", {3, 4}: "#"},
				diffStartIndex: 0,
			},
			want:  map[point]point{{2, 1}: {2, 0}, {3, 1}: {3, 0}, {2, 2}: {2, 3}, {2, 4}: {2, 3}, {3, 4}: {3, 3}},
			want1: map[point]int{{2, 0}: 1, {3, 0}: 1, {2, 3}: 2, {3, 3}: 1},
		},
		{
			name: "non-moving",
			args: args{
				elfCoords:      map[point]string{{2, 0}: "#", {4, 1}: "#", {0, 2}: "#", {4, 3}: "#", {2, 5}: "#"},
				diffStartIndex: 0,
			},
			want:  map[point]point{{2, 0}: {2, 0}, {4, 1}: {4, 1}, {0, 2}: {0, 2}, {4, 3}: {4, 3}, {2, 5}: {2, 5}},
			want1: map[point]int{{2, 0}: 1, {4, 1}: 1, {0, 2}: 1, {4, 3}: 1, {2, 5}: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := planElfMoves(tt.args.elfCoords, tt.args.diffStartIndex)

			assert.Equal(t, tt.want, got, "planned moves")
			assert.Equal(t, tt.want1, got1, "targeting coord counts")
		})
	}
}

func Test_determineNextCoords(t *testing.T) {
	type args struct {
		elfCoords      map[point]string
		coords         point
		diffStartIndex int
	}

	tests := []struct {
		name string
		args args
		want point
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMove(tt.args.elfCoords, tt.args.coords, tt.args.diffStartIndex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("determineNextCoords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resetElfCoords(t *testing.T) {
	type args struct {
		plannedMoves   map[point]point
		targetingCoord map[point]int
	}

	tests := []struct {
		name string
		args args
		want map[point]string
	}{
		{
			name: "small example",
			args: args{
				plannedMoves:   map[point]point{{2, 1}: {2, 0}, {3, 1}: {3, 0}, {2, 2}: {2, 3}, {2, 4}: {2, 3}, {3, 4}: {3, 3}},
				targetingCoord: map[point]int{{2, 0}: 1, {3, 0}: 1, {2, 3}: 2, {3, 3}: 1},
			},
			want: map[point]string{{2, 0}: "#", {3, 0}: "#", {2, 2}: "#", {2, 4}: "#", {3, 3}: "#"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateElfLocations(tt.args.plannedMoves, tt.args.targetingCoord)

			assert.Equal(t, tt.want, got)
		})
	}
}
