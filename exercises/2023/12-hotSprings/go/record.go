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

	newLine := strings.Join(expandedLeft, "?") + " " + strings.Join(expandedRight, ",")

	return parseLine(newLine)
}

func generateRegex(sizes []int) (*regexp.Regexp, error) {
	reSections := []string{}
	for _, c := range sizes {
		reSections = append(reSections, fmt.Sprintf("[ud]{%d}", c))
	}

	genRegex := `^[u\.]*` + strings.Join(reSections, `[u\.]+`) + `[u\.]*$`
	genRegex = strings.ReplaceAll(genRegex, `{1}`, "")

	return regexp.MustCompile(genRegex), nil
}

func (r *Record) countCombinations() (int, error) {
	if len(r.Checksum) == 0 {
		return -1, fmt.Errorf("no contiguous values provided")
	}

	// determine how much 'd' we need (max)
	var needed int
	for _, size := range r.Checksum {
		needed += size
	}

	re, err := generateRegex(r.Checksum)
	if err != nil {
		return -1, err
	}

	memo := make(map[string]int)
	return countHelper(r.Condition, needed, re, memo), nil
}

func countHelper(s string, n int, re *regexp.Regexp, memo map[string]int) int {
	if v, ok := memo[s]; ok {
		return v
	}

	var sum int
	left, right, canReplace := strings.Cut(s, "u")
	diff := n - strings.Count(s, "d")

	if diff < 0 { // too much 'd'
		return 0
	} else if diff == 0 { // just enough 'd'
		if re.MatchString(s) {
			return 1
		}

		return 0

	} else if !canReplace { // not enough 'd', can't get more
		return 0
	} else if !re.MatchString(s) { // is it worth continuing?
		return 0
	}

	// unknown, replace with 'd' and '.' and recurse
	dSub := replaceDotsAndJoin(left, "d", right)
	nSub := replaceDotsAndJoin(left, ".", right)

	sum += countHelper(dSub, n, re, memo)
	sum += countHelper(nSub, n, re, memo)

	memo[s] = sum

	return sum
}

func replaceDotsAndJoin(left, middle, right string) string {
	var result strings.Builder
	consecutiveDots := 0

	for _, str := range []string{left, middle, right} {
		for _, ch := range str {
			if ch == '.' {
				consecutiveDots++
				if consecutiveDots == 1 {
					result.WriteRune(ch)
				}
			} else {
				consecutiveDots = 0
				result.WriteRune(ch)
			}
		}
	}

	return result.String()
}
