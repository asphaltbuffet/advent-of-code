package exercises

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 12.
type Exercise struct {
	common.BaseExercise
}

var numRe = regexp.MustCompile(`-?\d+`)

// One returns the answer to the first part of the exercise: the sum of every
// number appearing in the document, regardless of structure.
func (e Exercise) One(instr string) (any, error) {
	sum := 0
	for _, m := range numRe.FindAllString(instr, -1) {
		n, _ := strconv.Atoi(m)
		sum += n
	}

	return sum, nil
}

// sumValue recursively sums numbers, skipping any object that has a value of
// exactly "red" (along with all of its descendants). Arrays are never skipped.
func sumValue(v any) int {
	switch x := v.(type) {
	case float64:
		return int(x)
	case []any:
		sum := 0
		for _, e := range x {
			sum += sumValue(e)
		}
		return sum
	case map[string]any:
		sum := 0
		for _, e := range x {
			if s, ok := e.(string); ok && s == "red" {
				return 0
			}
			sum += sumValue(e)
		}
		return sum
	default:
		return 0
	}
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	var doc any
	if err := json.Unmarshal([]byte(strings.TrimSpace(instr)), &doc); err != nil {
		return nil, err
	}

	return sumValue(doc), nil
}
