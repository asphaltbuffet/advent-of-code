package exercises

import (
	"cmp"
	"fmt"
	"math"
	"slices"
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
		Ranges:  Coalesce(ranges),
	}, nil
}

func Coalesce(rr []Range) []Range {
	slices.SortFunc(rr, func(a, b Range) int {
		lc := cmp.Compare(a.Low, b.Low)
		if lc == 0 {
			return cmp.Compare(a.High, b.High)
		}
		return lc
	})

	out := []Range{rr[0]}
	last := 0
	for i := 1; i < len(rr); i++ {
		if out[last].High >= rr[i].Low {
			out[last].High = max(out[last].High, rr[i].High)
		} else {
			out = append(out, rr[i])
			last++
		}
	}

	return out
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

func (inv Inventory) CountRanges() int {
	var count int
	for _, r := range inv.Ranges {
		count += r.High - r.Low + 1
	}

	return count
}
