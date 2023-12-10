package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHand_Hash(t *testing.T) {
	tests := []struct {
		name  string
		cards string
		wilds bool
		want  int
	}{
		{name: "one pair", cards: "32T3K", wilds: false, want: 13},
		{name: "3 of a kind", cards: "T55J5", wilds: false, want: 102},
		{name: "empty", cards: "", wilds: false, want: 0},

		{name: "pair and a wild", cards: "32J3K", wilds: true, want: 102},
		{name: "pair and two wild", cards: "QJJQ2", wilds: true, want: 1001},
		{name: "trip and two wild", cards: "222JJ", wilds: true, want: 10000},
		{name: "all wild", cards: "JJJJJ", wilds: true, want: 10000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Hash(tt.cards, tt.wilds))
		})
	}
}

func Test_ToHex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		wild bool
		want int
	}{
		{"less than 10", "22222", false, 139810},
		{"all 10s", "TTTTT", false, 699050},
		{"all 10s", "KKKKK", false, 908765},
		{"only 10s or higher", "TJQKA", false, 703710},
		{"mixed hex and numeric", "T2J9A", false, 666526},

		{"no wilds", "22222", true, 139810},
		{"some wilds", "T2J9A", true, 664222},
		{"all wilds", "JJJJJ", true, 69905},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, ToHex(tt.in, tt.wild))
		})
	}
}

func Test_parseLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantHand string
		wantBid  int
	}{
		{"example 1", "32T3K 765", "32T3K", 765},
		{"example 1", "T55J5 684", "T55J5", 684},
		{"example 1", "KK677 28", "KK677", 28},
		{"example 1", "KTJJT 220", "KTJJT", 220},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHand, gotBid := parseLine(tt.line)
			assert.Equal(t, tt.wantHand, gotHand)
			assert.Equal(t, tt.wantBid, gotBid)
		})
	}
}
