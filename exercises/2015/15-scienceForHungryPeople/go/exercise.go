package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 15.
type Exercise struct {
	common.BaseExercise
}

const teaspoons = 100

type ingredient struct {
	capacity, durability, flavor, texture, calories int
}

// parse reads one ingredient per line. Only the integer fields matter, so we
// pull every signed number out of the line in order.
func parse(instr string) []ingredient {
	var is []ingredient

	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		nums := numbers(line)
		is = append(is, ingredient{nums[0], nums[1], nums[2], nums[3], nums[4]})
	}

	return is
}

// numbers extracts the signed integers from a line, ignoring punctuation.
func numbers(line string) []int {
	fields := strings.FieldsFunc(line, func(r rune) bool {
		return !(r == '-' || (r >= '0' && r <= '9'))
	})

	var out []int
	for _, f := range fields {
		if n, err := strconv.Atoi(f); err == nil {
			out = append(out, n)
		}
	}

	return out
}

// best walks every distribution of `teaspoons` over the ingredients. For each
// complete recipe it computes the score (product of clamped property totals);
// when calorieGoal >= 0 only recipes hitting that calorie count are scored.
func best(is []ingredient, calorieGoal int) int {
	amounts := make([]int, len(is))
	max := 0

	var rec func(i, remaining int)
	rec = func(i, remaining int) {
		// Last ingredient takes whatever teaspoons are left.
		if i == len(is)-1 {
			amounts[i] = remaining
			if s := score(is, amounts, calorieGoal); s > max {
				max = s
			}
			return
		}
		for a := 0; a <= remaining; a++ {
			amounts[i] = a
			rec(i+1, remaining-a)
		}
	}
	rec(0, teaspoons)

	return max
}

// score multiplies the per-property totals (negatives clamped to 0). If
// calorieGoal >= 0 and the recipe misses it, the recipe scores 0.
func score(is []ingredient, amounts []int, calorieGoal int) int {
	var capacity, dur, fla, tex, cal int
	for i, a := range amounts {
		capacity += is[i].capacity * a
		dur += is[i].durability * a
		fla += is[i].flavor * a
		tex += is[i].texture * a
		cal += is[i].calories * a
	}

	if calorieGoal >= 0 && cal != calorieGoal {
		return 0
	}

	for _, v := range []*int{&capacity, &dur, &fla, &tex} {
		if *v < 0 {
			*v = 0
		}
	}

	return capacity * dur * fla * tex
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	return best(parse(instr), -1), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return best(parse(instr), 500), nil
}
