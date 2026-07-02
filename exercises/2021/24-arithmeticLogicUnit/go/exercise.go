package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 24.
type Exercise struct {
	common.BaseExercise
}

// The MONAD program is 14 near-identical blocks, one per input digit. Each block
// differs only in three parameters extracted from fixed lines:
//
//	div z A   (A is 1 → "push", or 26 → "pop")
//	add x B
//	add y C
//
// A push block does z = z*26 + (w+C). A pop block requires w == z%26 + B to
// cancel a previous push (leaving z reducible toward 0); B is always negative for
// pops. Pairing each pop with the push it cancels yields a linear constraint
// between the two digits, and the model number is built by choosing digits within
// 1..9 to satisfy them.
type block struct {
	div, addX, addY int
}

func parseBlocks(instr string) ([]block, error) {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	var blocks []block
	for i := 0; i+18 <= len(lines); i += 18 {
		if strings.TrimSpace(lines[i]) != "inp w" {
			return nil, fmt.Errorf("expected 'inp w' at line %d, got %q", i, lines[i])
		}
		div, err1 := lastInt(lines[i+4]) // div z A
		addX, err2 := lastInt(lines[i+5]) // add x B
		addY, err3 := lastInt(lines[i+15]) // add y C
		if err1 != nil || err2 != nil || err3 != nil {
			return nil, fmt.Errorf("parsing block at line %d", i)
		}
		blocks = append(blocks, block{div, addX, addY})
	}
	if len(blocks) != 14 {
		return nil, fmt.Errorf("expected 14 blocks, got %d", len(blocks))
	}
	return blocks, nil
}

func lastInt(line string) (int, error) {
	fields := strings.Fields(line)
	return strconv.Atoi(fields[len(fields)-1])
}

// pairing holds a constraint w[j] = w[i] + delta between two digit positions.
type pairing struct {
	i, j, delta int
}

// constraints pairs each pop block with the push it cancels and returns the
// resulting digit relations.
func constraints(blocks []block) []pairing {
	type entry struct {
		idx, addY int
	}
	var stack []entry
	var pairs []pairing
	for k, b := range blocks {
		if b.div == 1 {
			stack = append(stack, entry{k, b.addY})
		} else {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			// w[k] = w[top] + (top.addY + b.addX)
			pairs = append(pairs, pairing{top.idx, k, top.addY + b.addX})
		}
	}
	return pairs
}

// bestModel fills a 14-digit model number satisfying every pairing, choosing the
// largest digits when maximize is true and the smallest otherwise.
func bestModel(pairs []pairing, maximize bool) string {
	digits := make([]int, 14)
	for _, p := range pairs {
		// w[j] = w[i] + delta, both in 1..9.
		// Pick the extreme value of w[i] that keeps w[j] in range.
		for wi := 1; wi <= 9; wi++ {
			wj := wi + p.delta
			if wj < 1 || wj > 9 {
				continue
			}
			digits[p.i], digits[p.j] = wi, wj
			if maximize {
				// keep scanning; last valid (largest wi) wins
			} else {
				break // first valid (smallest wi) wins
			}
		}
	}
	var b strings.Builder
	for _, d := range digits {
		b.WriteByte(byte('0' + d))
	}
	return b.String()
}

// One returns the largest valid 14-digit model number.
func (e Exercise) One(instr string) (any, error) {
	blocks, err := parseBlocks(instr)
	if err != nil {
		return nil, err
	}
	return bestModel(constraints(blocks), true), nil
}

// Two returns the smallest valid 14-digit model number.
func (e Exercise) Two(instr string) (any, error) {
	blocks, err := parseBlocks(instr)
	if err != nil {
		return nil, err
	}
	return bestModel(constraints(blocks), false), nil
}
