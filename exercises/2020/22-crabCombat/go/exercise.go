package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 22.
type Exercise struct {
	common.BaseExercise
}

// parseDecks reads the two players' decks (top card first).
func parseDecks(instr string) ([]int, []int, error) {
	blocks := strings.Split(strings.TrimSpace(instr), "\n\n")
	if len(blocks) != 2 {
		return nil, nil, fmt.Errorf("expected two players")
	}
	deck := func(block string) ([]int, error) {
		lines := strings.Split(block, "\n")
		out := make([]int, 0, len(lines)-1)
		for _, l := range lines[1:] { // skip "Player N:"
			n, err := strconv.Atoi(strings.TrimSpace(l))
			if err != nil {
				return nil, err
			}
			out = append(out, n)
		}
		return out, nil
	}
	p1, err := deck(blocks[0])
	if err != nil {
		return nil, nil, err
	}
	p2, err := deck(blocks[1])
	if err != nil {
		return nil, nil, err
	}
	return p1, p2, nil
}

// score sums each card times its position counting up from the bottom.
func score(deck []int) int {
	total := 0
	n := len(deck)
	for i, c := range deck {
		total += c * (n - i)
	}
	return total
}

// playCombat plays the non-recursive game and returns the winning deck.
func playCombat(p1, p2 []int) []int {
	for len(p1) > 0 && len(p2) > 0 {
		a, b := p1[0], p2[0]
		p1, p2 = p1[1:], p2[1:]
		if a > b {
			p1 = append(p1, a, b)
		} else {
			p2 = append(p2, b, a)
		}
	}
	if len(p1) > 0 {
		return p1
	}
	return p2
}

// One plays Combat and returns the winner's score.
func (e Exercise) One(instr string) (any, error) {
	p1, p2, err := parseDecks(instr)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%d", score(playCombat(p1, p2))), nil
}

// stateKey encodes both decks so repeated game states can be detected.
func stateKey(p1, p2 []int) string {
	var b strings.Builder
	for _, c := range p1 {
		b.WriteByte(byte(c))
	}
	b.WriteByte('|')
	for _, c := range p2 {
		b.WriteByte(byte(c))
	}
	return b.String()
}

// playRecursive plays Recursive Combat. It returns true if player 1 wins, and the
// winning deck. A repeated deck state within a game is an instant player-1 win.
func playRecursive(p1, p2 []int) (bool, []int) {
	seen := map[string]bool{}
	for len(p1) > 0 && len(p2) > 0 {
		key := stateKey(p1, p2)
		if seen[key] {
			return true, p1
		}
		seen[key] = true

		a, b := p1[0], p2[0]
		p1, p2 = p1[1:], p2[1:]

		var p1Wins bool
		if len(p1) >= a && len(p2) >= b {
			// Recurse on copies of the next a / b cards.
			sub1 := append([]int(nil), p1[:a]...)
			sub2 := append([]int(nil), p2[:b]...)
			p1Wins, _ = playRecursive(sub1, sub2)
		} else {
			p1Wins = a > b
		}

		if p1Wins {
			p1 = append(p1, a, b)
		} else {
			p2 = append(p2, b, a)
		}
	}
	if len(p1) > 0 {
		return true, p1
	}
	return false, p2
}

// Two plays Recursive Combat and returns the winner's score.
func (e Exercise) Two(instr string) (any, error) {
	p1, p2, err := parseDecks(instr)
	if err != nil {
		return nil, err
	}
	_, win := playRecursive(p1, p2)
	return fmt.Sprintf("%d", score(win)), nil
}
