package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2022 day 21.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer:
func (e Exercise) One(instr string) (any, error) {
	raw := parse(instr)

	result, err := calc("root", raw, make(map[string]int))
	if err != nil {
		return nil, fmt.Errorf("calculating part 1 answer: %w", err)
	}

	return result, nil
}

// Two returns the answer to the second part of the exercise.
// wrong: 7010269744524
// answer: 3558714869436
func (e Exercise) Two(instr string) (any, error) {
	raw := parse(instr)
	results := map[string]int{}

	// force error when calculating "humn"
	raw["humn"] = "die_die_die"

	// change root equation to be  left / right = 1
	reversed := map[string]string{"root": "1"}

	orig := strings.Split(raw["root"], " ")
	orig[1] = "/"
	raw["root"] = strings.Join(orig, " ")
	currentName := "root"

	for currentName != "humn" {
		var left, operator, right string
		_, _ = fmt.Sscanf(raw[currentName], "%s %s %s", &left, &operator, &right)

		leftVal, errLeft := calc(left, raw, results)
		if errLeft == nil {
			reversed[left] = strconv.Itoa(leftVal)
		}

		rightVal, errRight := calc(right, raw, results)
		if errRight == nil {
			reversed[right] = strconv.Itoa(rightVal)
		}

		switch operator {
		case "+":
			switch {
			case errLeft != nil:
				reversed[left] = fmt.Sprintf("%s - %s", currentName, right)
				currentName = left
			case errRight != nil:
				reversed[right] = fmt.Sprintf("%s - %s", currentName, left)
				currentName = right
			default:
				return nil, fmt.Errorf("no error path found in %q: %s", currentName, raw[currentName])
			}
		case "-":
			switch {
			case errLeft != nil:
				reversed[left] = fmt.Sprintf("%s + %s", currentName, right)
				currentName = left
			case errRight != nil:
				reversed[right] = fmt.Sprintf("%s - %s", left, currentName)
				currentName = right
			default:
				return nil, fmt.Errorf("no error path found in %q: %s", currentName, raw[currentName])
			}
		case "*":
			switch {
			case errLeft != nil:
				reversed[left] = fmt.Sprintf("%s / %s", currentName, right)
				currentName = left
			case errRight != nil:
				reversed[right] = fmt.Sprintf("%s / %s", currentName, left)
				currentName = right
			default:
				return nil, fmt.Errorf("no error path found in %q: %s", currentName, raw[currentName])
			}
		case "/":
			switch {
			case errLeft != nil:
				reversed[left] = fmt.Sprintf("%s * %s", currentName, right)
				currentName = left
			case errRight != nil:
				reversed[right] = fmt.Sprintf("%s / %s", left, currentName)
				currentName = right
			default:
				return nil, fmt.Errorf("no error path found in %q: %s", currentName, raw[currentName])
			}
		default:
			return nil, fmt.Errorf("%q has invalid operator in %q", currentName, raw[currentName])
		}
	}

	return calc("humn", reversed, map[string]int{})
}
