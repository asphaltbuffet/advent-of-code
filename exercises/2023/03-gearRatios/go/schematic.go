package exercises

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type schematic struct {
	// engine   map[point]string
	parts    map[int]part
	isSymbol map[point]bool
}

type part struct {
	value  string
	number int
	points []point
	bounds []point
	id     int
}

func New(s string) *schematic {
	lines := strings.Split(s, "\n")

	schm := &schematic{
		parts:    make(map[int]part),
		isSymbol: make(map[point]bool),
	}

	for i, line := range lines {
		if err := schm.parseLine(line, i); err != nil {
			fmt.Printf("error parsing line %d: %s\n", i, err)
			panic(err)
		}
	}

	return schm
}

func (schm *schematic) parseLine(line string, y int) error {
	for i := 0; i < len(line); i++ {
		// fmt.Printf("checking position %d: %s\n", i, string(line[i]))

		c := rune(line[i])
		p := point{x: i, y: y}

		if unicode.IsDigit(c) {
			i += schm.addPart(line[i:], p) - 1
		} else if c != '.' {
			schm.isSymbol[p] = true
		}

	}

	return nil
}

func (schm *schematic) addPart(s string, origin point) int {
	if !unicode.IsDigit(rune(s[0])) {
		return 1
	}

	p := part{
		id:     len(schm.parts),
		value:  "",
		points: []point{},
		bounds: []point{
			{x: origin.x - 1, y: origin.y - 1},
			{x: origin.x - 1, y: origin.y},
			{x: origin.x - 1, y: origin.y + 1},
		},
	}

	for i, c := range s {
		if !unicode.IsDigit(c) {
			break
		}

		// add point to part
		p.points = append(p.points, point{x: origin.x + i, y: origin.y})

		// add surrounding area to part
		p.bounds = append(p.bounds,
			point{x: origin.x + i, y: origin.y - 1},
			point{x: origin.x + i, y: origin.y + 1},
		)

		p.value += string(c)
	}

	plen := len(p.value)

	// add ending bounds
	p.bounds = append(p.bounds,
		point{x: origin.x + plen, y: origin.y - 1},
		point{x: origin.x + plen, y: origin.y},
		point{x: origin.x + plen, y: origin.y + 1},
	)

	n, err := strconv.Atoi(p.value)
	if err != nil {
		fmt.Printf("error converting part[%d] value %q to int: %s\n", p.id, p.value, err)
		// panic(err)
	}

	p.number = n

	schm.parts[p.id] = p

	return plen
}
