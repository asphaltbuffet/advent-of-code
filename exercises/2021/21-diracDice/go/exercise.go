package exercises

import (
	"fmt"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 21.
type Exercise struct {
	common.BaseExercise
}

func parse(instr string) (int, int, error) {
	var p1, p2 int
	_, err := fmt.Sscanf(instr,
		"Player 1 starting position: %d\nPlayer 2 starting position: %d", &p1, &p2)
	if err != nil {
		return 0, 0, fmt.Errorf("parsing starting positions: %w", err)
	}
	return p1, p2, nil
}

// move advances a position by n on the 1..10 wrapping track.
func move(pos, n int) int {
	return (pos-1+n)%10 + 1
}

// One plays the deterministic game: a 100-sided die rolled three times per turn.
// The answer is the losing score times the number of die rolls.
func (e Exercise) One(instr string) (any, error) {
	p1, p2, err := parse(instr)
	if err != nil {
		return nil, err
	}
	positions := [2]int{p1, p2}
	scores := [2]int{0, 0}

	die := 0
	rolls := 0
	roll := func() int {
		die = die%100 + 1
		rolls++
		return die
	}

	turn := 0
	for {
		step := roll() + roll() + roll()
		positions[turn] = move(positions[turn], step)
		scores[turn] += positions[turn]
		if scores[turn] >= 1000 {
			loser := scores[1-turn]
			return fmt.Sprintf("%d", loser*rolls), nil
		}
		turn = 1 - turn
	}
}

// diracRolls holds the seven distinct three-roll sums of a 3-sided die and how
// many of the 27 universes produce each.
var diracRolls = [][2]int{{3, 1}, {4, 3}, {5, 6}, {6, 7}, {7, 6}, {8, 3}, {9, 1}}

type dstate struct {
	pos1, pos2, score1, score2 int
}

// countWins returns how many universes each player wins from this state, where
// player 1 is the one about to move.
func countWins(s dstate, memo map[dstate][2]int64) [2]int64 {
	if s.score2 >= 21 {
		return [2]int64{0, 1} // player 2 already won last turn
	}
	if v, ok := memo[s]; ok {
		return v
	}

	var wins [2]int64
	for _, r := range diracRolls {
		sum, mult := r[0], int64(r[1])
		np := move(s.pos1, sum)
		// After moving, roles swap: the opponent becomes the mover.
		sub := countWins(dstate{s.pos2, np, s.score2, s.score1 + np}, memo)
		wins[0] += mult * sub[1]
		wins[1] += mult * sub[0]
	}

	memo[s] = wins
	return wins
}

// Two counts the quantum game and returns the larger win total.
func (e Exercise) Two(instr string) (any, error) {
	p1, p2, err := parse(instr)
	if err != nil {
		return nil, err
	}

	wins := countWins(dstate{p1, p2, 0, 0}, map[dstate][2]int64{})
	best := wins[0]
	if wins[1] > best {
		best = wins[1]
	}

	return fmt.Sprintf("%d", best), nil
}
