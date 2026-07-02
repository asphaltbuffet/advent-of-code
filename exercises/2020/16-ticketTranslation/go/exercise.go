package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 16.
type Exercise struct {
	common.BaseExercise
}

// field is a named validity constraint: a value is valid if it falls in either
// [lo1,hi1] or [lo2,hi2].
type field struct {
	name                   string
	lo1, hi1, lo2, hi2 int
}

func (f field) valid(v int) bool {
	return (v >= f.lo1 && v <= f.hi1) || (v >= f.lo2 && v <= f.hi2)
}

type notes struct {
	fields  []field
	mine    []int
	tickets [][]int
}

func parseInts(csv string) ([]int, error) {
	parts := strings.Split(strings.TrimSpace(csv), ",")
	out := make([]int, len(parts))
	for i, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}
		out[i] = n
	}
	return out, nil
}

func parse(instr string) (notes, error) {
	blocks := strings.Split(strings.TrimSpace(instr), "\n\n")
	if len(blocks) != 3 {
		return notes{}, fmt.Errorf("expected 3 blocks, got %d", len(blocks))
	}

	var n notes
	for _, line := range strings.Split(blocks[0], "\n") {
		name, ranges, ok := strings.Cut(line, ": ")
		if !ok {
			continue
		}
		var f field
		f.name = name
		if _, err := fmt.Sscanf(ranges, "%d-%d or %d-%d", &f.lo1, &f.hi1, &f.lo2, &f.hi2); err != nil {
			return notes{}, fmt.Errorf("parsing ranges %q: %w", ranges, err)
		}
		n.fields = append(n.fields, f)
	}

	mineLines := strings.Split(blocks[1], "\n")
	mine, err := parseInts(mineLines[len(mineLines)-1])
	if err != nil {
		return notes{}, fmt.Errorf("parsing your ticket: %w", err)
	}
	n.mine = mine

	for _, line := range strings.Split(blocks[2], "\n")[1:] {
		t, err := parseInts(line)
		if err != nil {
			return notes{}, fmt.Errorf("parsing nearby ticket %q: %w", line, err)
		}
		n.tickets = append(n.tickets, t)
	}

	return n, nil
}

// validForAny reports whether v satisfies at least one field's ranges.
func validForAny(fields []field, v int) bool {
	for _, f := range fields {
		if f.valid(v) {
			return true
		}
	}
	return false
}

// One sums every nearby-ticket value that is valid for no field.
func (e Exercise) One(instr string) (any, error) {
	n, err := parse(instr)
	if err != nil {
		return nil, err
	}

	rate := 0
	for _, t := range n.tickets {
		for _, v := range t {
			if !validForAny(n.fields, v) {
				rate += v
			}
		}
	}

	return fmt.Sprintf("%d", rate), nil
}

// assignFields deduces which field belongs to each column. It builds the set of
// candidate fields per column (those whose ranges accept every valid ticket's
// value in that column), then repeatedly locks in columns that have a single
// candidate, removing it from the others. Returns column index -> field name.
func assignFields(n notes) map[int]string {
	cols := len(n.mine)

	// valid tickets only (all values valid for some field).
	valid := [][]int{n.mine}
	for _, t := range n.tickets {
		ok := true
		for _, v := range t {
			if !validForAny(n.fields, v) {
				ok = false
				break
			}
		}
		if ok {
			valid = append(valid, t)
		}
	}

	// candidates[c] = set of field indices possible for column c.
	candidates := make([]map[int]bool, cols)
	for c := 0; c < cols; c++ {
		candidates[c] = map[int]bool{}
		for fi, f := range n.fields {
			fits := true
			for _, t := range valid {
				if !f.valid(t[c]) {
					fits = false
					break
				}
			}
			if fits {
				candidates[c][fi] = true
			}
		}
	}

	assigned := map[int]string{}
	for len(assigned) < cols {
		progress := false
		for c := 0; c < cols; c++ {
			if len(candidates[c]) != 1 {
				continue
			}
			var fi int
			for k := range candidates[c] {
				fi = k
			}
			assigned[c] = n.fields[fi].name
			progress = true
			for cc := 0; cc < cols; cc++ {
				delete(candidates[cc], fi)
			}
		}
		if !progress {
			break
		}
	}
	return assigned
}

// Two multiplies your ticket's values for the six "departure" fields.
func (e Exercise) Two(instr string) (any, error) {
	n, err := parse(instr)
	if err != nil {
		return nil, err
	}

	assigned := assignFields(n)
	product := 1
	for c, name := range assigned {
		if strings.HasPrefix(name, "departure") {
			product *= n.mine[c]
		}
	}

	return fmt.Sprintf("%d", product), nil
}
