package exercises

import (
	"cmp"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Hand struct {
	Cards    string
	Value    int
	Strength int
	Bid      int
}

type Cards []Card

type Card int

// Hash returns a numeric representation of a dedup reduction of cards.
//
// Example:
// "23456" -> 00005
// "12344" -> 00013
// "33333" -> 10000
func Hash(s string, useWildJ bool) int {
	if len(s) != 5 {
		return 0
	}

	seen := map[rune]int{}

	for _, card := range s {
		seen[card]++
	}

	if useWildJ {
		w := seen['J']
		seen['J'] = 0

		// find card with highest frequency
		var max int
		var maxCard rune
		for card, freq := range seen {
			if freq > max {
				max = freq
				maxCard = card
			}
		}

		seen[maxCard] += w

	}

	var n float64
	for _, freq := range seen {
		n += math.Pow10(freq - 1)
	}

	return int(n)
}

func ToHex(cards string, useWildJ bool) int {
	var hexMap map[string]string

	if useWildJ {
		hexMap = map[string]string{"T": "a", "J": "1", "Q": "c", "K": "d", "A": "e"}
	} else {
		hexMap = map[string]string{"T": "a", "J": "b", "Q": "c", "K": "d", "A": "e"}
	}

	for f, t := range hexMap {
		cards = strings.ReplaceAll(cards, f, t)
	}

	n, err := strconv.ParseInt(cards, 16, 64)
	if err != nil {
		return -1
	}

	return int(n)
}

func parseLine(s string) (string, int) {
	var hand string
	var bid int

	n, err := fmt.Sscan(s, &hand, &bid)
	if err != nil || n != 2 {
		panic("invalid input")
	}

	return hand, bid
}

func parseHand(s string, useWildJ bool) Hand {
	h, b := parseLine(s)

	return Hand{
		Cards:    s,
		Value:    Hash(h, useWildJ),
		Strength: ToHex(h, useWildJ),
		Bid:      b,
	}
}

func handSort(a, b Hand) int {
	if n := cmp.Compare(a.Value, b.Value); n != 0 {
		return n
	}

	return cmp.Compare(a.Strength, b.Strength)
}
