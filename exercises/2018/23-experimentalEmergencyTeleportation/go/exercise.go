package exercises

import (
	"container/heap"
	"regexp"
	"strconv"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 23.
type Exercise struct {
	common.BaseExercise
}

type bot struct {
	x, y, z, r int
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	bots := parse(instr)

	strongest := bots[0]
	for _, b := range bots {
		if b.r > strongest.r {
			strongest = b
		}
	}
	count := 0
	for _, b := range bots {
		if abs(b.x-strongest.x)+abs(b.y-strongest.y)+abs(b.z-strongest.z) <= strongest.r {
			count++
		}
	}
	return count, nil
}

// Two returns the answer to the second part of the exercise — the distance from
// the origin to the point in range of the most nanobots (closest to origin on
// ties). It searches with an octree: cubes are examined best-first by how many
// bots could possibly reach them, subdividing until a single point wins.
func (e Exercise) Two(instr string) (any, error) {
	bots := parse(instr)

	// Cover all bots with a cube whose side is a power of two.
	var span int
	for _, b := range bots {
		span = max(span, abs(b.x), abs(b.y), abs(b.z))
	}
	size := 1
	for size < span {
		size *= 2
	}

	pq := &cubeQueue{}
	heap.Push(pq, newCube(bots, -size, -size, -size, 2*size))

	for pq.Len() > 0 {
		c := heap.Pop(pq).(cube)
		if c.size == 1 {
			// Best-first order guarantees this is the winning point.
			return c.distToOrigin, nil
		}
		half := c.size / 2
		for _, d := range [8][3]int{
			{0, 0, 0}, {half, 0, 0}, {0, half, 0}, {0, 0, half},
			{half, half, 0}, {half, 0, half}, {0, half, half}, {half, half, half},
		} {
			heap.Push(pq, newCube(bots, c.x+d[0], c.y+d[1], c.z+d[2], half))
		}
	}
	return -1, nil
}

// cube is an axis-aligned region [x, x+size) × … scored by how many bots' ranges
// reach it and how close its nearest corner is to the origin.
type cube struct {
	x, y, z, size int
	inRange       int
	distToOrigin  int
}

// newCube builds a cube and scores it: inRange counts bots whose Manhattan range
// intersects the cube, distToOrigin is the closest the cube gets to the origin.
func newCube(bots []bot, x, y, z, size int) cube {
	c := cube{x: x, y: y, z: z, size: size}
	for _, b := range bots {
		// Manhattan distance from the bot to the nearest point of the cube.
		d := axisDist(b.x, x, size) + axisDist(b.y, y, size) + axisDist(b.z, z, size)
		if d <= b.r {
			c.inRange++
		}
	}
	c.distToOrigin = axisDist(0, x, size) + axisDist(0, y, size) + axisDist(0, z, size)
	return c
}

// axisDist returns how far coordinate v lies outside the half-open interval
// [lo, lo+size) on one axis (0 if inside).
func axisDist(v, lo, size int) int {
	if v < lo {
		return lo - v
	}
	if v > lo+size-1 {
		return v - (lo + size - 1)
	}
	return 0
}

// parse reads the nanobots.
func parse(instr string) []bot {
	re := regexp.MustCompile(`-?\d+`)
	var bots []bot
	for _, line := range regexp.MustCompile(`\r?\n`).Split(instr, -1) {
		nums := re.FindAllString(line, -1)
		if len(nums) != 4 {
			continue
		}
		vals := make([]int, 4)
		for i, s := range nums {
			vals[i], _ = strconv.Atoi(s)
		}
		bots = append(bots, bot{vals[0], vals[1], vals[2], vals[3]})
	}
	return bots
}

// cubeQueue orders cubes best-first: most bots in range, then closest to origin,
// then smallest — so the first size-1 cube popped is the optimal point.
type cubeQueue []cube

func (q cubeQueue) Len() int { return len(q) }
func (q cubeQueue) Less(i, j int) bool {
	if q[i].inRange != q[j].inRange {
		return q[i].inRange > q[j].inRange
	}
	if q[i].distToOrigin != q[j].distToOrigin {
		return q[i].distToOrigin < q[j].distToOrigin
	}
	return q[i].size < q[j].size
}
func (q cubeQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i] }
func (q *cubeQueue) Push(x any)   { *q = append(*q, x.(cube)) }
func (q *cubeQueue) Pop() any {
	old := *q
	n := len(old)
	c := old[n-1]
	*q = old[:n-1]
	return c
}
