package exercises

import (
	"sort"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 20.
type Exercise struct {
	common.BaseExercise
}

type ipRange struct{ lo, hi int }

// parseRanges reads "lo-hi" blocked ranges, sorted ascending by lo.
func parseRanges(instr string) []ipRange {
	var ranges []ipRange
	for line := range strings.FieldsSeq(instr) {
		parts := strings.SplitN(line, "-", 2)
		lo, _ := strconv.Atoi(parts[0])
		hi, _ := strconv.Atoi(parts[1])
		ranges = append(ranges, ipRange{lo, hi})
	}
	sort.Slice(ranges, func(i, j int) bool { return ranges[i].lo < ranges[j].lo })
	return ranges
}

// maxIP returns the highest valid IP: 9 for the small example, else 2^32 - 1.
func maxIP(ranges []ipRange) int {
	for _, r := range ranges {
		if r.hi > 9 {
			return 4294967295
		}
	}
	return 9
}

// One returns the answer to the first part of the exercise: the lowest IP not
// covered by any blocked range.
func (e Exercise) One(instr string) (any, error) {
	ranges := parseRanges(instr)
	candidate := 0
	for _, r := range ranges {
		if r.lo > candidate {
			break // candidate falls in a gap before this range
		}
		if r.hi+1 > candidate {
			candidate = r.hi + 1
		}
	}
	return candidate, nil
}

// Two returns the answer to the second part of the exercise: the count of IPs
// not covered by any blocked range.
func (e Exercise) Two(instr string) (any, error) {
	ranges := parseRanges(instr)
	m := maxIP(ranges)

	allowed := 0
	next := 0 // lowest IP not yet accounted for
	for _, r := range ranges {
		if r.lo > next {
			allowed += r.lo - next // gap of unblocked IPs
		}
		if r.hi+1 > next {
			next = r.hi + 1
		}
	}
	if next <= m {
		allowed += m - next + 1
	}
	return allowed, nil
}
