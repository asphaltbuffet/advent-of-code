package exercises

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 9.
type Exercise struct {
	common.BaseExercise
}

var intRe = regexp.MustCompile(`\d+`)

// parse reads the player count and the last marble's worth.
func parse(instr string) (int, int, error) {
	nums := intRe.FindAllString(instr, -1)
	if len(nums) != 2 {
		return 0, 0, fmt.Errorf("expected 2 numbers, got %d in %q", len(nums), instr)
	}

	players, _ := strconv.Atoi(nums[0])
	last, _ := strconv.Atoi(nums[1])

	return players, last, nil
}

// marble is a node in the circular doubly-linked ring of played marbles.
type marble struct {
	value      int
	prev, next *marble
}

// highScore plays the full game and returns the winning player's score.
func highScore(players, last int) int {
	scores := playGame(players, last)

	best := 0
	for _, s := range scores {
		if s > best {
			best = s
		}
	}

	return best
}

// playGame plays the full game and returns every player's score. The ring is a
// circular doubly-linked list, so both the normal insertion (two clockwise) and
// the scoring removal (seven counter-clockwise) are O(1).
func playGame(players, last int) []int {
	scores := make([]int, players)

	current := &marble{value: 0}
	current.prev, current.next = current, current

	for m := 1; m <= last; m++ {
		if m%23 == 0 {
			// Scoring move: take the marble seven counter-clockwise, add both it
			// and the current marble value to the active player's score.
			for range 7 {
				current = current.prev
			}
			scores[m%players] += m + current.value

			// Remove that marble; the one clockwise becomes current.
			current.prev.next = current.next
			current.next.prev = current.prev
			current = current.next
		} else {
			// Normal move: insert between the marbles one and two clockwise.
			left := current.next
			right := left.next
			inserted := &marble{value: m, prev: left, next: right}
			left.next = inserted
			right.prev = inserted
			current = inserted
		}
	}

	return scores
}

// One returns the answer to the first part of the exercise.
// Answer: 400493
func (e Exercise) One(instr string) (any, error) {
	players, last, err := parse(instr)
	if err != nil {
		return nil, err
	}

	return highScore(players, last), nil
}

// Two returns the answer to the second part of the exercise.
// Answer: 3338341690
func (e Exercise) Two(instr string) (any, error) {
	players, last, err := parse(instr)
	if err != nil {
		return nil, err
	}

	// Same game, but the last marble is worth 100 times as much.
	return highScore(players, last*100), nil
}
