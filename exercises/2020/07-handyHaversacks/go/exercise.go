package exercises

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 7.
type Exercise struct {
	common.BaseExercise
}

const target = "shiny gold"

// content is one "N adjective color" entry inside a bag rule.
type content struct {
	count int
	color string
}

// contentRe matches each "3 pale aqua bag(s)" clause on the right of a rule.
var contentRe = regexp.MustCompile(`(\d+) (\w+ \w+) bags?`)

// parse maps each bag color to the bags it directly contains.
func parse(instr string) map[string][]content {
	rules := map[string][]content{}
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		outer, rest, ok := strings.Cut(line, " bags contain ")
		if !ok {
			continue
		}
		var contents []content
		for _, m := range contentRe.FindAllStringSubmatch(rest, -1) {
			n, _ := strconv.Atoi(m[1])
			contents = append(contents, content{n, m[2]})
		}
		rules[outer] = contents
	}
	return rules
}

// One counts bag colors that can eventually contain a shiny gold bag.
func (e Exercise) One(instr string) (any, error) {
	rules := parse(instr)

	// Invert: child color -> parents that directly contain it.
	parents := map[string][]string{}
	for outer, contents := range rules {
		for _, c := range contents {
			parents[c.color] = append(parents[c.color], outer)
		}
	}

	// BFS out from the target over the inverted graph.
	seen := map[string]bool{}
	queue := []string{target}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, p := range parents[cur] {
			if !seen[p] {
				seen[p] = true
				queue = append(queue, p)
			}
		}
	}

	return fmt.Sprintf("%d", len(seen)), nil
}

// countInside returns how many bags are contained within one bag of the given
// color, memoized over colors.
func countInside(color string, rules map[string][]content, memo map[string]int) int {
	if v, ok := memo[color]; ok {
		return v
	}
	total := 0
	for _, c := range rules[color] {
		total += c.count * (1 + countInside(c.color, rules, memo))
	}
	memo[color] = total
	return total
}

// Two counts the total number of bags required inside a single shiny gold bag.
func (e Exercise) Two(instr string) (any, error) {
	rules := parse(instr)
	return fmt.Sprintf("%d", countInside(target, rules, map[string]int{})), nil
}
