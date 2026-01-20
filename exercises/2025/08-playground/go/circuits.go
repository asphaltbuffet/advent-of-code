package exercises

import (
	"cmp"
	"slices"
	"strconv"
	"strings"

	dsu "github.com/arunksaha/gdsu/compact"
)

type Vector struct {
	X int
	Y int
	Z int
}

type Distance struct {
	A int
	B int
	D int
}

type Junctions []Vector

func NewJunctions(s string) (*Junctions, error) {
	lines := strings.Split(s, "\n")
	vv := make(Junctions, len(lines))
	for i, l := range lines {
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

		vv[i] = Vector{X: x, Y: y, Z: z}
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

func (jj Junctions) AllDists() []Distance {
	dists := []Distance{}
	for i, a := range jj {
		for j, b := range slices.Backward(jj) {
			if i == j {
				break
			}

			d := Dist(a, b)
			dists = append(dists, Distance{A: i, B: j, D: d})
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

	ds := dsu.New(len(jj))

	for _, d := range dists[:wires] {
		ds.Union(d.A, d.B)
	}

	// determine sizes of all sets (circuits)
	sizes := make([]int, len(jj))
	for i := range jj {
		sizes[ds.Find(i)] += 1
	}

	slices.SortFunc(sizes, func(a, b int) int {
		return cmp.Compare(b, a)
	})

	// return 3 largest sizes
	return sizes[:3]
}

func (jj Junctions) EndCircuits() int {
	// get shortest <wires> distances
	dists := jj.AllDists()

	ds := dsu.New(len(jj))

	connections := 0
	for _, d := range dists {
		if ds.Union(d.A, d.B) {
			connections++

			if connections == len(jj)-1 {
				return jj[d.A].X * jj[d.B].X
			}
		}
	}

	// we shouldn't get here
	return -1
}
