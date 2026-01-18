package exercises

import (
	"fmt"
	"strconv"
	"strings"
)

type Farm struct {
	presents map[int]int
	regions  []Region
}

type Region struct {
	area     int
	presents []int
}

func Parse(s string) *Farm {
	lines := strings.Split(s, "\n\n")

	f := &Farm{
		presents: make(map[int]int, len(lines)),
	}

	for i, l := range lines[:len(lines)-1] {
		f.presents[i] = strings.Count(l, "#")
	}

	n := len(f.presents)
	rl := strings.Split(lines[len(lines)-1], "\n")
	f.regions = make([]Region, len(rl))

	for i, r := range rl {
		tok := strings.Fields(r)
		var l, w int
		fmt.Sscanf(tok[0], "%dx%d:", &l, &w)

		shapes := make([]int, n)
		for j, c := range tok[1:] {
			d, _ := strconv.Atoi(c)

			shapes[j] = d
		}

		f.regions[i] = Region{
			area:     l * w,
			presents: shapes,
		}
	}

	return f
}
