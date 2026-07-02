package exercises

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 4.
type Exercise struct {
	common.BaseExercise
}

// requiredFields are the keys every valid passport must carry ("cid" is optional).
var requiredFields = []string{"byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid"}

// parse splits the batch into passports, each a map of key:value fields. Records
// are separated by blank lines; fields may be spread across lines or spaces.
func parse(instr string) []map[string]string {
	blocks := strings.Split(strings.TrimSpace(instr), "\n\n")
	passports := make([]map[string]string, 0, len(blocks))
	for _, block := range blocks {
		fields := map[string]string{}
		for _, tok := range strings.Fields(block) {
			k, v, ok := strings.Cut(tok, ":")
			if ok {
				fields[k] = v
			}
		}
		passports = append(passports, fields)
	}
	return passports
}

// hasRequired reports whether every required field is present.
func hasRequired(p map[string]string) bool {
	for _, k := range requiredFields {
		if _, ok := p[k]; !ok {
			return false
		}
	}
	return true
}

// One counts passports that carry all required fields.
func (e Exercise) One(instr string) (any, error) {
	valid := 0
	for _, p := range parse(instr) {
		if hasRequired(p) {
			valid++
		}
	}
	return fmt.Sprintf("%d", valid), nil
}

var (
	hclRe = regexp.MustCompile(`^#[0-9a-f]{6}$`)
	pidRe = regexp.MustCompile(`^[0-9]{9}$`)
	eclOK = map[string]bool{
		"amb": true, "blu": true, "brn": true, "gry": true,
		"grn": true, "hzl": true, "oth": true,
	}
)

// yearInRange reports whether s is a 4-digit year within [lo,hi].
func yearInRange(s string, lo, hi int) bool {
	if len(s) != 4 {
		return false
	}
	n, err := strconv.Atoi(s)
	return err == nil && n >= lo && n <= hi
}

// heightOK reports whether hgt is 150-193cm or 59-76in.
func heightOK(s string) bool {
	switch {
	case strings.HasSuffix(s, "cm"):
		n, err := strconv.Atoi(strings.TrimSuffix(s, "cm"))
		return err == nil && n >= 150 && n <= 193
	case strings.HasSuffix(s, "in"):
		n, err := strconv.Atoi(strings.TrimSuffix(s, "in"))
		return err == nil && n >= 59 && n <= 76
	}
	return false
}

// isValid reports whether a passport has all required fields and each obeys its
// value rule.
func isValid(p map[string]string) bool {
	if !hasRequired(p) {
		return false
	}
	return yearInRange(p["byr"], 1920, 2002) &&
		yearInRange(p["iyr"], 2010, 2020) &&
		yearInRange(p["eyr"], 2020, 2030) &&
		heightOK(p["hgt"]) &&
		hclRe.MatchString(p["hcl"]) &&
		eclOK[p["ecl"]] &&
		pidRe.MatchString(p["pid"])
}

// Two counts passports whose required fields are all present and valid.
func (e Exercise) Two(instr string) (any, error) {
	valid := 0
	for _, p := range parse(instr) {
		if isValid(p) {
			valid++
		}
	}
	return fmt.Sprintf("%d", valid), nil
}
