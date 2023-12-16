package exercises

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Record struct {
	Condition string
	Checksum  []int
}

func parseLine(s string) (*Record, error) {
	data, second, ok := strings.Cut(s, " ")
	if !ok {
		return nil, fmt.Errorf("invalid record: %s", s)
	}

	var d string
	for _, c := range data {
		switch c {
		case '#':
			d += "d"

		case '?':
			d += "u"

		case '.':
			d += "."

		default:
			return nil, fmt.Errorf("invalid character: %c", c)
		}
	}

	checks := []int{}
	for _, i := range strings.Split(second, ",") {
		v, err := strconv.Atoi(i)
		if err != nil {
			return nil, fmt.Errorf("invalid checksum: %s", i)
		}

		checks = append(checks, v)
	}

	return &Record{
		Condition: d,
		Checksum:  checks,
	}, nil
}

func expandAndParseLine(s string) (*Record, error) {
	left, right, _ := strings.Cut(s, " ")
	expandedLeft := make([]string, 5)
	expandedRight := make([]string, 5)

	for i := 0; i < 5; i++ {
		expandedLeft[i] = left
		expandedRight[i] = right
	}

	newLine := strings.Trim(strings.Join(expandedLeft, "?"), ".") + " " + strings.Join(expandedRight, ",")

	return parseLine(newLine)
}

func generateRegex(sizes []int) (*regexp.Regexp, error) {
	// if len(sizes) == 0 {
	// 	return nil, fmt.Errorf("no contiguous values provided")
	// }

	reSections := []string{}
	for _, c := range sizes {
		// if c < 1 {
		// 	return nil, fmt.Errorf("invalid contiguous section value: %d", c)
		// }

		reSections = append(reSections, fmt.Sprintf("[ud]{%d}", c))
	}

	genRegex := `^[u\.]*` + strings.Join(reSections, `[u\.]+`) + `[u\.]*$`

	return regexp.MustCompile(genRegex), nil
}

// Generates all possible combinations of the string where 'u' is replaced by 'd'.
func (r *Record) countCombinations() (int, error) {
	// if len(r.Checksum) == 0 {
	// 	return -1, fmt.Errorf("no contiguous values provided")
	// }

	re, _ := generateRegex(r.Checksum)
	// if err != nil {
	// 	return -1, err
	// }
	var n int
	for _, i := range r.Checksum {
		n += i
	}

	// fmt.Printf("starting with %q\n", r.Condition)
	return countHelper(r.Condition, n, re), nil
}

func countHelper(s string, n int, re *regexp.Regexp) int {
	var sum int
	idx := strings.Index(s, "u")

	if idx == -1 || n == strings.Count(s, "d") {
		if re.MatchString(s) {
			// fmt.Printf("nothing to replace, %q matches pattern\n", s)
			return 1
		}

		return 0
	}

	// unknown, replace with 'd' and '.' and recurse
	dSub := s[:idx] + "d" + s[idx+1:]
	nSub := s[:idx] + "." + s[idx+1:]

	// reDot := regexp.MustCompile(`\.{2,}`)
	// nSub = reDot.ReplaceAllLiteralString(nSub, ".")

	// fmt.Printf("from %q checking %q and %q\n", s, dSub, nSub)
	if re.MatchString(dSub) {
		// fmt.Printf("recurse with %q\n", dSub)
		sum += countHelper(dSub, n, re)
	}

	if re.MatchString(nSub) {
		// fmt.Printf("recurse with %q\n", nSub)
		sum += countHelper(nSub, n, re)
	}

	return sum
}
