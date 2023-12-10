package exercises

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/pkg/ring"
)

type Move int

const (
	MoveLeft Move = iota + 1
	MoveRight
)

func parseMovement(s string) *ring.Ring[Move] {
	if s == "" {
		return nil
	}

	moves := ring.New[Move](len(s))

	for _, c := range s {
		switch c {
		case 'L':
			moves.Element = MoveLeft

		case 'R':
			moves.Element = MoveRight

		default:
			panic("invalid move")
		}

		moves = moves.Next()
	}

	return moves
}

func parseMap(s string) map[string]map[Move]string {
	lines := strings.Split(s, "\n")
	m := make(map[string]map[Move]string, len(lines))

	for _, line := range lines {
		re := regexp.MustCompile(`^(?P<src>\w+) = \((?P<left>\w+), (?P<right>\w+)\)$`)
		matches := findNamedMatches(re, line)
		if len(matches) != 3 {
			panic(fmt.Sprintf("invalid line: %s", line))
		}

		m[matches["src"]] = map[Move]string{MoveLeft: matches["left"], MoveRight: matches["right"]}
	}

	return m
}

func findNamedMatches(re *regexp.Regexp, s string) map[string]string {
	match := re.FindStringSubmatch(s)
	if len(match) == 0 {
		return nil
	}

	result := make(map[string]string)

	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return result
}
