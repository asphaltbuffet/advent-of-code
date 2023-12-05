package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *Card
	}{
		{
			name: "example one",
			args: args{
				s: "Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53",
			},
			want: &Card{
				ID:        1,
				Winning:   []string{"41", "48", "83", "86", "17"},
				Scratched: []string{"83", "86", "6", "31", "17", "9", "48", "53"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.args.s)

			if assert.Len(t, c.Winning, 5) {
				assert.ElementsMatch(t, tt.want.Winning, c.Winning)
			}

			if assert.Len(t, c.Scratched, 8) {
				assert.ElementsMatch(t, tt.want.Scratched, c.Scratched)
			}
		})
	}
}

func TestCard_Score(t *testing.T) {
	tests := []struct {
		name string
		c    *Card
		want int
	}{
		{
			name: "five winning numbers",
			c: &Card{
				ID:        1,
				Winning:   []string{"41", "48", "83", "86", "17"},
				Scratched: []string{"83", "86", "6", "31", "17", "9", "48", "53"},
			},
			want: 8,
		},
		{
			name: "two winning numbers",
			c: &Card{
				ID:        2,
				Winning:   []string{"13", "32", "20", "16", "61"},
				Scratched: []string{"61", "30", "68", "82", "17", "32", "24", "19"},
			},
			want: 2,
		},
		{
			name: "no winning numbers",
			c: &Card{
				ID:        6,
				Winning:   []string{"31", "18", "13", "56", "72"},
				Scratched: []string{"74", "77", "10", "23", "35", "67", "36", "11"},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.c.Score())
		})
	}
}

func Test_addWonCards(t *testing.T) {
	type args struct {
		id    int
		count int
		mult  map[int]int
	}

	tests := []struct {
		name string
		args args
		want map[int]int
	}{
		{
			name: "increase next four by one",
			args: args{
				id:    1,
				count: 4,
				mult:  map[int]int{1: 1, 2: 1, 3: 1, 4: 1, 5: 1, 6: 1},
			},
			want: map[int]int{1: 1, 2: 2, 3: 2, 4: 2, 5: 2, 6: 1},
		},
		{
			name: "increase next two by two",
			args: args{
				id:    2,
				count: 2,
				mult:  map[int]int{1: 1, 2: 2, 3: 2, 4: 2, 5: 2, 6: 1},
			},
			want: map[int]int{1: 1, 2: 2, 3: 4, 4: 4, 5: 2, 6: 1},
		},
		{
			name: "increase next two by four",
			args: args{
				id:    3,
				count: 2,
				mult:  map[int]int{1: 1, 2: 2, 3: 4, 4: 4, 5: 2, 6: 1},
			},
			want: map[int]int{1: 1, 2: 2, 3: 4, 4: 8, 5: 6, 6: 1},
		},
		{
			name: "increase next one by eight",
			args: args{
				id:    4,
				count: 1,
				mult:  map[int]int{1: 1, 2: 2, 3: 4, 4: 8, 5: 6, 6: 1},
			},
			want: map[int]int{1: 1, 2: 2, 3: 4, 4: 8, 5: 14, 6: 1},
		},
		{
			name: "increase none by fourteen",
			args: args{
				id:    5,
				count: 0,
				mult:  map[int]int{1: 1, 2: 2, 3: 4, 4: 8, 5: 14, 6: 1},
			},
			want: map[int]int{1: 1, 2: 2, 3: 4, 4: 8, 5: 14, 6: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addWonCards(tt.args.id, tt.args.count, tt.args.mult)

			assert.Equal(t, tt.want, tt.args.mult)
		})
	}
}

func Test_countTotalCards(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "no copies",
			args: args{
				lines: []string{
					"Card 1: 31 18 13 56 72 | 74 77 10 23 35 67 36 11",
					"Card 2: 31 18 13 56 72 | 74 77 10 23 35 67 36 11",
					"Card 3: 31 18 13 56 72 | 74 77 10 23 35 67 36 11",
				},
			},
			want: 3,
		},
		{
			name: "one copy",
			args: args{
				lines: []string{
					"Card 1: 41 92 73 84 69 | 59 84 76 51 58  5 54 83",
					"Card 2: 31 18 13 56 72 | 74 77 10 23 35 67 36 11",
					"Card 3: 31 18 13 56 72 | 74 77 10 23 35 67 36 11",
				},
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, countTotalCards(tt.args.lines))
		})
	}
}
