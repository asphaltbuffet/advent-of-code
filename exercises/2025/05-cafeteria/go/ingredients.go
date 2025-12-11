package exercises

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Inventory struct {
	Lowest  int
	Highest int
	Ranges  []Range
}

type Range struct {
	Low  int
	High int
}

func LoadInventory(s string) (*Inventory, error) {
	lines := strings.Split(s, "\n")
	ranges := make([]Range, len(lines))

	lowest := math.MaxInt
	highest := 0

	for i, l := range lines {
		tokens := strings.Split(l, "-")
		if len(tokens) != 2 {
			return nil, fmt.Errorf("%d tokens: invalid range", len(tokens))
		}

		l, err := strconv.Atoi(tokens[0])
		if err != nil {
			return nil, fmt.Errorf("range %d, token=%q: %w", i, tokens[0], err)
		}

		r, err := strconv.Atoi(tokens[1])
		if err != nil {
			return nil, fmt.Errorf("range %d, token=%q: %w", i, tokens[0], err)
		}

		ranges[i] = Range{l, r}
		lowest = min(lowest, l)
		highest = max(highest, r)
	}

	return &Inventory{
		Lowest:  lowest,
		Highest: highest,
		Ranges:  ranges,
	}, nil
}

func LoadIngredients(s string) ([]int, error) {
	lines := strings.Split(s, "\n")
	ids := make([]int, len(lines))
	for i, l := range lines {
		id, err := strconv.Atoi(l)
		if err != nil {
			return nil, fmt.Errorf("load id=%d: %w", i, err)
		}

		ids[i] = id
	}

	return ids, nil
}

func (inv Inventory) IsFresh(id int) bool {
	if id < inv.Lowest || id > inv.Highest {
		return false
	}

	for _, r := range inv.Ranges {
		if id >= r.Low && id <= r.High {
			return true
		}
	}

	return false
}
