package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 2.
type Exercise struct {
	common.BaseExercise
}

// entry is one password policy line: "lo-hi c: password".
type entry struct {
	lo, hi int
	ch     byte
	pass   string
}

func parse(instr string) ([]entry, error) {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	entries := make([]entry, 0, len(lines))
	for _, line := range lines {
		var e entry
		var ch string
		// e.g. "2-8 t: pncmjxlvckfbtrjh"
		_, err := fmt.Sscanf(line, "%d-%d %1s: %s", &e.lo, &e.hi, &ch, &e.pass)
		if err != nil {
			return nil, fmt.Errorf("parsing line %q: %w", line, err)
		}
		e.ch = ch[0]
		entries = append(entries, e)
	}
	return entries, nil
}

// One counts passwords where the policy character appears between lo and hi
// times (inclusive).
func (e Exercise) One(instr string) (any, error) {
	entries, err := parse(instr)
	if err != nil {
		return nil, err
	}

	valid := 0
	for _, en := range entries {
		n := strings.Count(en.pass, string(en.ch))
		if n >= en.lo && n <= en.hi {
			valid++
		}
	}

	return fmt.Sprintf("%d", valid), nil
}

// Two counts passwords where exactly one of the two 1-based positions lo and hi
// holds the policy character.
func (e Exercise) Two(instr string) (any, error) {
	entries, err := parse(instr)
	if err != nil {
		return nil, err
	}

	valid := 0
	for _, en := range entries {
		a := en.lo-1 < len(en.pass) && en.pass[en.lo-1] == en.ch
		b := en.hi-1 < len(en.pass) && en.pass[en.hi-1] == en.ch
		if a != b {
			valid++
		}
	}

	return fmt.Sprintf("%d", valid), nil
}
