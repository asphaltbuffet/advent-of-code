package exercises

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseInput(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name      string
		args      args
		want      *Dish
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "",
			args: args{
				input: "O....#....\nO.OO#....#\n.....##...\nOO.#O....O\n.O.....O#.\nO.#..O.#.#\n..O..#O..O\n.......O..\n#....###..\n#OO..#....",
			},
			want: &Dish{
				Rocks: [][]byte{
					[]byte("O....#...."),
					[]byte("O.OO#....#"),
					[]byte(".....##..."),
					[]byte("OO.#O....O"),
					[]byte(".O.....O#."),
					[]byte("O.#..O.#.#"),
					[]byte("..O..#O..O"),
					[]byte(".......O.."),
					[]byte("#....###.."),
					[]byte("#OO..#...."),
				},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseInput(tt.args.input)

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_cmpRock(t *testing.T) {
	type args struct {
		a byte
		b byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"out of order", args{a: '.', b: 'O'}, 1},
		{"equal round", args{a: 'O', b: 'O'}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, cmpRock(tt.args.a, tt.args.b))
		})
	}
}

func Test_sortRocks(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"no cubes", args{in: []byte("OO..O")}, []byte("OOO..")},
		{"cube between dot", args{in: []byte("OO.#.O")}, []byte("OO.#O.")},
		{"cube beginning dot", args{in: []byte("..#OO.#.O")}, []byte("..#OO.#O.")},
		{"cube ending dot", args{in: []byte("OO.#.O#..")}, []byte("OO.#O.#..")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slices.SortFunc(tt.args.in, cmpRock)

			assert.Equal(t, tt.want, tt.args.in)
		})
	}
}

func Test_transpose(t *testing.T) {
	type args struct {
		d Dish
	}
	tests := []struct {
		name string
		args args
		want Dish
	}{
		{
			name: "square",
			args: args{
				d: Dish{
					Rocks: [][]byte{
						[]byte("abc"),
						[]byte("def"),
						[]byte("hij"),
					},
				},
			},
			want: Dish{
				Rocks: [][]byte{
					[]byte("adh"),
					[]byte("bei"),
					[]byte("cfj"),
				},
			},
		},
		{
			name: "wide",
			args: args{
				d: Dish{
					Rocks: [][]byte{
						[]byte("abcdef"),
						[]byte("ghijkl"),
						[]byte("mnopqr"),
					},
				},
			},
			want: Dish{
				Rocks: [][]byte{
					[]byte("agm"),
					[]byte("bhn"),
					[]byte("cio"),
					[]byte("djp"),
					[]byte("ekq"),
					[]byte("flr"),
				},
			},
		},
		{
			name: "tall",
			args: args{
				d: Dish{
					Rocks: [][]byte{
						[]byte("abc"),
						[]byte("def"),
						[]byte("ghi"),
						[]byte("jkl"),
						[]byte("mno"),
						[]byte("pqr"),
					},
				},
			},
			want: Dish{
				Rocks: [][]byte{
					[]byte("adgjmp"),
					[]byte("behknq"),
					[]byte("cfilor"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.d.transpose()
			assert.Equal(t, tt.want, tt.args.d)
		})
	}
}

func TestDish_RotateCCW(t *testing.T) {
	tests := []struct {
		name string
		d    *Dish
		want *Dish
	}{
		{
			name: "first rotation",
			d: &Dish{
				Rocks: [][]byte{
					[]byte("abc"),
					[]byte("def"),
					[]byte("ghi"),
					[]byte("jkl"),
					[]byte("mno"),
					[]byte("pqr"),
				},
			},
			want: &Dish{
				Rocks: [][]byte{
					[]byte("cfilor"),
					[]byte("behknq"),
					[]byte("adgjmp"),
				},
			},
		},
		{
			name: "second rotation",
			d: &Dish{
				Rocks: [][]byte{
					[]byte("cfilor"),
					[]byte("behknq"),
					[]byte("adgjmp"),
				},
			},
			want: &Dish{
				Rocks: [][]byte{
					[]byte("rqp"),
					[]byte("onm"),
					[]byte("lkj"),
					[]byte("ihg"),
					[]byte("fed"),
					[]byte("cba"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.RotateCCW()

			assert.Equal(t, tt.want, tt.d)
		})
	}
}

func TestDish_spin(t *testing.T) {
	tests := []struct {
		name string
		d    *Dish
		want *Dish
	}{
		{
			name: "one cycle",
			d: &Dish{
				Rocks: [][]byte{
					[]byte("O....#...."),
					[]byte("O.OO#....#"),
					[]byte(".....##..."),
					[]byte("OO.#O....O"),
					[]byte(".O.....O#."),
					[]byte("O.#..O.#.#"),
					[]byte("..O..#O..O"),
					[]byte(".......O.."),
					[]byte("#....###.."),
					[]byte("#OO..#...."),
				},
			},
			want: &Dish{
				Rocks: [][]byte{
					[]byte(".....#...."),
					[]byte("....#...O#"),
					[]byte("...OO##..."),
					[]byte(".OO#......"),
					[]byte(".....OOO#."),
					[]byte(".O#...O#.#"),
					[]byte("....O#...."),
					[]byte("......OOOO"),
					[]byte("#...O###.."),
					[]byte("#..OO#...."),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.transpose()
			tt.d.spin()
			tt.d.transpose()

			assert.Equal(t, tt.want, tt.d)
		})
	}
}
