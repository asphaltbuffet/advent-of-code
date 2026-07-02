package exercises

import (
	"fmt"
	"sort"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 21.
type Exercise struct {
	common.BaseExercise
}

// food is one line: its ingredient list and its declared allergens.
type food struct {
	ingredients []string
	allergens   []string
}

func parse(instr string) []food {
	var foods []food
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		ingPart := line
		var allergens []string
		if i := strings.Index(line, " (contains "); i >= 0 {
			ingPart = line[:i]
			list := strings.TrimSuffix(line[i+len(" (contains "):], ")")
			for _, a := range strings.Split(list, ", ") {
				allergens = append(allergens, a)
			}
		}
		foods = append(foods, food{ingredients: strings.Fields(ingPart), allergens: allergens})
	}
	return foods
}

// candidates maps each allergen to the set of ingredients that could contain it:
// the intersection of ingredient lists across every food declaring that allergen.
func candidates(foods []food) map[string]map[string]bool {
	cand := map[string]map[string]bool{}
	for _, f := range foods {
		ingSet := map[string]bool{}
		for _, ing := range f.ingredients {
			ingSet[ing] = true
		}
		for _, a := range f.allergens {
			if cand[a] == nil {
				// first sighting: copy the ingredient set
				cand[a] = map[string]bool{}
				for ing := range ingSet {
					cand[a][ing] = true
				}
				continue
			}
			// intersect with this food's ingredients
			for ing := range cand[a] {
				if !ingSet[ing] {
					delete(cand[a], ing)
				}
			}
		}
	}
	return cand
}

// One counts appearances of ingredients that cannot be any allergen.
func (e Exercise) One(instr string) (any, error) {
	foods := parse(instr)
	cand := candidates(foods)

	// Ingredients that are a candidate for at least one allergen.
	suspect := map[string]bool{}
	for _, ings := range cand {
		for ing := range ings {
			suspect[ing] = true
		}
	}

	count := 0
	for _, f := range foods {
		for _, ing := range f.ingredients {
			if !suspect[ing] {
				count++
			}
		}
	}

	return fmt.Sprintf("%d", count), nil
}

// Two resolves each allergen to its single ingredient and returns the canonical
// dangerous list: those ingredients ordered by allergen name, comma-joined.
func (e Exercise) Two(instr string) (any, error) {
	cand := candidates(parse(instr))

	// Constraint propagation: repeatedly fix allergens with one candidate.
	resolved := map[string]string{} // allergen -> ingredient
	for len(resolved) < len(cand) {
		progress := false
		for a, ings := range cand {
			if _, done := resolved[a]; done || len(ings) != 1 {
				continue
			}
			var ing string
			for i := range ings {
				ing = i
			}
			resolved[a] = ing
			progress = true
			for other := range cand {
				delete(cand[other], ing)
			}
		}
		if !progress {
			return nil, fmt.Errorf("could not resolve allergens")
		}
	}

	allergens := make([]string, 0, len(resolved))
	for a := range resolved {
		allergens = append(allergens, a)
	}
	sort.Strings(allergens)

	ings := make([]string, len(allergens))
	for i, a := range allergens {
		ings[i] = resolved[a]
	}
	return strings.Join(ings, ","), nil
}
