package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Vis shows the two-stage filter over the passport batch: how many pass Part One
// (all required fields present) and Part Two (all fields present and valid), then
// which field rule most often rejects the passports that had every field but
// still failed. The funnel bar plus a per-field rejection chart explains the drop
// from Part One to Part Two. Bars use colorblind-safe colors ordered by
// brightness and are labeled with counts, so the chart reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	passports := parse(instr)
	total := len(passports)

	var present, valid int
	// fieldFail counts, among passports that have every required field, how many
	// fail each specific field rule.
	fieldFail := map[string]int{}
	checks := []struct {
		name string
		ok   func(map[string]string) bool
	}{
		{"byr", func(p map[string]string) bool { return yearInRange(p["byr"], 1920, 2002) }},
		{"iyr", func(p map[string]string) bool { return yearInRange(p["iyr"], 2010, 2020) }},
		{"eyr", func(p map[string]string) bool { return yearInRange(p["eyr"], 2020, 2030) }},
		{"hgt", func(p map[string]string) bool { return heightOK(p["hgt"]) }},
		{"hcl", func(p map[string]string) bool { return hclRe.MatchString(p["hcl"]) }},
		{"ecl", func(p map[string]string) bool { return eclOK[p["ecl"]] }},
		{"pid", func(p map[string]string) bool { return pidRe.MatchString(p["pid"]) }},
	}

	for _, p := range passports {
		if !hasRequired(p) {
			continue
		}
		present++
		allGood := true
		for _, c := range checks {
			if !c.ok(p) {
				fieldFail[c.name]++
				allGood = false
			}
		}
		if allGood {
			valid++
		}
	}

	const (
		W    = 820
		H    = 460
		mX   = 60
		barW = 620
	)

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="40" fill="#e8ecf4" font-size="16">Passport Processing: %d passports through two filters</text>`, mX, total)

	// Funnel: total -> present (part one) -> valid (part two).
	stages := []struct {
		label string
		count int
		col   string
	}{
		{"all passports", total, "#7a869a"},
		{"fields present (part 1)", present, "#F0E442"},
		{"fields valid (part 2)", valid, "#009E73"},
	}
	for i, s := range stages {
		y := 70 + i*46
		w := barW * s.count / total
		fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="32" fill="%s"/>`, mX, y, w, s.col)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#0d0f18" font-size="14" font-weight="bold">%d</text>`, mX+8, y+22, s.count)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="13">%s</text>`, mX+w+10, y+22, s.label)
	}

	// Per-field rejection chart (among field-present passports).
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="14">Rejections by field rule (of the %d with all fields):</text>`, mX, 240, present)
	names := make([]string, 0, len(fieldFail))
	for k := range fieldFail {
		names = append(names, k)
	}
	sort.Slice(names, func(i, j int) bool { return fieldFail[names[i]] > fieldFail[names[j]] })
	maxFail := 1
	for _, n := range names {
		if fieldFail[n] > maxFail {
			maxFail = fieldFail[n]
		}
	}
	for i, n := range names {
		y := 262 + i*26
		w := 460 * fieldFail[n] / maxFail
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#c8d0dc" font-size="13" text-anchor="end">%s</text>`, mX+34, y+13, n)
		fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="17" fill="#0072B2"/>`, mX+44, y, w)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="12">%d</text>`, mX+44+w+6, y+13, fieldFail[n])
	}

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "passport-processing.svg"), []byte(sb.String()), 0o600)
}
