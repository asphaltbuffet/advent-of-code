package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Vis shows the allergen-to-ingredient deduction as a candidate matrix: rows are
// allergens, columns are the suspect ingredients (those that could contain some
// allergen). A filled cell means that ingredient is a candidate for that allergen
// (it appears in every food declaring it); constraint propagation reduces the
// matrix to one ingredient per allergen, and those solved cells are outlined. The
// solved ingredients, read in allergen-name order, form the Part Two answer.
// Candidates and the solved matching are distinguished by outline and position as
// well as color, so the grid reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	foods := parse(instr)
	cand := candidates(foods)

	// Suspect ingredients (candidate for at least one allergen), sorted.
	suspectSet := map[string]bool{}
	for _, ings := range cand {
		for ing := range ings {
			suspectSet[ing] = true
		}
	}
	ingredients := make([]string, 0, len(suspectSet))
	for ing := range suspectSet {
		ingredients = append(ingredients, ing)
	}
	sort.Strings(ingredients)

	allergens := make([]string, 0, len(cand))
	for a := range cand {
		allergens = append(allergens, a)
	}
	sort.Strings(allergens)

	// Resolve (on a copy) to know the final assignment.
	work := map[string]map[string]bool{}
	for a, ings := range cand {
		work[a] = map[string]bool{}
		for ing := range ings {
			work[a][ing] = true
		}
	}
	solved := map[string]string{}
	for len(solved) < len(work) {
		for a, ings := range work {
			if _, done := solved[a]; done || len(ings) != 1 {
				continue
			}
			var ing string
			for i := range ings {
				ing = i
			}
			solved[a] = ing
			for other := range work {
				if other != a {
					delete(work[other], ing)
				}
			}
		}
	}

	const (
		cell   = 24
		labelW = 90
		colLbl = 90
		mL     = 10
	)
	mT := colLbl + 20
	nc := len(ingredients)
	nr := len(allergens)
	gridW := nc * cell
	W := mL + labelW + gridW + 20
	if W < 470 {
		W = 470 // ensure the title and footer fit
	}
	H := mT + nr*cell + 30

	candCol := "#56B4E9"  // sky blue: candidate
	solveCol := "#F0E442" // yellow: resolved

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="18" fill="#e8ecf4" font-size="14">Allergen Assessment: candidate ingredients per allergen</text>`, mL)

	gx := mL + labelW
	// Column (ingredient) labels, rotated.
	for c, ing := range ingredients {
		cx := gx + c*cell + cell/2
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="10" transform="rotate(-55 %d %d)" text-anchor="start">%s</text>`,
			cx, mT-6, cx, mT-6, ing)
	}

	for r, a := range allergens {
		y := mT + r*cell
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#c8d0dc" font-size="12" text-anchor="end">%s</text>`, mL+labelW-6, y+cell/2+4, a)
		for c, ing := range ingredients {
			cx := gx + c*cell
			isCand := cand[a][ing]
			isSolved := solved[a] == ing
			fillCol := "#1b1f26"
			if isSolved {
				fillCol = solveCol
			} else if isCand {
				fillCol = candCol
			}
			fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s"/>`, cx+1, y+1, cell-2, cell-2, fillCol)
			if isSolved {
				fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" fill="none" stroke="#e8ecf4" stroke-width="2"/>`, cx+1, y+1, cell-2, cell-2)
			}
		}
	}

	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#9aa4b2" font-size="11">filled = candidate; outlined = resolved (part 2, read in allergen order)</text>`, mL, H-10)

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "allergen-assessment.svg"), []byte(sb.String()), 0o600)
}
