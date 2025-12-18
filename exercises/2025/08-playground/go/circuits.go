package exercises

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
)

type Vector struct {
	X int
	Y int
	Z int
}

type Junctions []Vector

func NewJunctions(s string) (*Junctions, error) {
	vv := Junctions{}
	for _, l := range strings.Split(s, "\n") {
		tokens := strings.Split(l, ",")

		x, err := strconv.Atoi(tokens[0])
		if err != nil {
			return nil, err
		}
		y, err := strconv.Atoi(tokens[1])
		if err != nil {
			return nil, err
		}
		z, err := strconv.Atoi(tokens[2])
		if err != nil {
			return nil, err
		}

		vv = append(vv, Vector{X: x, Y: y, Z: z})
	}

	return &vv, nil
}

func Dist(a Vector, b Vector) int {
	return pow(a.X-b.X, 2) + pow(a.Y-b.Y, 2) + pow(a.Z-b.Z, 2)
}

func pow(a, b int) int {
	if b == 0 {
		return 1
	}

	n := a
	for _ = range b - 1 {
		n *= a
	}

	return n
}

func find_set(parent map[int]int, v int) int {
	if v == parent[v] {
		return v
	}

	parent[v] = find_set(parent, parent[v])
	return parent[v]
}

func merge_sets(parent map[int]int, a, b int) {
	parent[find_set(parent, parent[b])] = find_set(parent, a)
}

type Distance struct {
	A int
	B int
	D int
}

func (jj Junctions) AllDists() []Distance {
	dists := []Distance{}
	for i, a := range jj {
		for j, b := range jj {
			if i < j {
				d := Dist(a, b)
				dists = append(dists, Distance{A: i, B: j, D: d})
			}
		}
	}

	slices.SortFunc(dists, func(m, n Distance) int {
		return cmp.Compare(m.D, n.D)
	})

	return dists
}

func (jj Junctions) CreateCircuits(wires int) []int {
	// get shortest <wires> distances
	dists := jj.AllDists()

	// merge_sets for number of wires
	parents := make(map[int]int, len(jj))
	for i := range jj {
		parents[i] = i
	}

	for _, d := range dists[:wires] {
		merge_sets(parents, d.A, d.B)
	}

	// determine sizes of all sets (circuits)
	sizes := make([]int, len(jj))
	for i := range jj {
		sizes[find_set(parents, i)] += 1
	}

	// fmt.Println(sizes)
	slices.Sort(sizes)
	slices.Reverse(sizes)

	return sizes
}

func (jj Junctions) EndCircuits() int {
	// get shortest <wires> distances
	dists := jj.AllDists()

	// merge_sets for number of wires
	parents := make(map[int]int, len(jj))
	for i := range jj {
		parents[i] = i
	}

	connections := 0
	for _, d := range dists {
		a := find_set(parents, d.A)
		b := find_set(parents, d.B)

		if a != b {
			connections++

			merge_sets(parents, d.A, d.B)

			if connections == len(jj)-1 {
				return jj[d.A].X * jj[d.B].X
			}
		}
	}

	// we shouldn't get here
	return -1
}
