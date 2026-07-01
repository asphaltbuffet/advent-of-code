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

	// Standard springs DP: recurse on (position in condition, group index),
	// consuming one group of damaged springs or one operational spring at a
	// time. Memoised on those two small integers, this is linear in the record
	// length times the number of groups — no regex or string rebuilding.
	memo := map[[2]int]int{}
	var count func(pos, grp int) int
	count = func(pos, grp int) int {
		if pos == len(r.Condition) {
			if grp == len(r.Checksum) {
				return 1 // all groups placed, nothing left
			}
			return 0
		}
		key := [2]int{pos, grp}
		if v, ok := memo[key]; ok {
			return v
		}

		var total int
		c := r.Condition[pos]

		// Treat this cell as operational ('.'): skip it.
		if c == '.' || c == 'u' {
			total += count(pos+1, grp)
		}

		// Treat this cell as the start of the next damaged group ('#').
		if (c == 'd' || c == 'u') && grp < len(r.Checksum) {
			size := r.Checksum[grp]
			end := pos + size
			if end <= len(r.Condition) && noOperational(r.Condition, pos, end) &&
				(end == len(r.Condition) || r.Condition[end] != 'd') {
				if end == len(r.Condition) {
					total += count(end, grp+1)
				} else {
					total += count(end+1, grp+1) // consume the separator
				}
			}
		}

		memo[key] = total
		return total
	}

	return count(0, 0), nil
}

// noOperational reports whether Condition[start:end] contains no operational
// ('.') cells, i.e. every cell can be damaged.
func noOperational(cond string, start, end int) bool {
	for i := start; i < end; i++ {
		if cond[i] == '.' {
			return false
		}
	}
	return true
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
