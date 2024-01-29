package exercises

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type Vertex struct {
	X int
	Y int
	Z int
}

// Brick is the AABB of a rectangular prism.
type Brick struct {
	min Vertex
	max Vertex
}

func countDisintegrated(bricksBelow, bricksAbove map[int][]int, i int) int {
	disintegrated := map[int]bool{i: true} // we'll need to subtract 1 from the final count
	q := []int{i}

	for len(q) > 0 {
		check := q[0]
		q = q[1:]

		// for all bricks above, check if all bricks below are disintegrated
		for _, a := range bricksAbove[check] {
			removed := 0

			for _, b := range bricksBelow[a] {
				if _, poof := disintegrated[b]; poof {
					removed++
				}
			}

			if len(bricksBelow[a]) == removed {
				// all supports are gone, add it to the list
				q = append(q, a)
				disintegrated[a] = true
			}
		}
	}

	return len(disintegrated) - 1 // subtract 1 for the original brick
}

func canMoveDown(b Brick, below []Brick) (Brick, []int, bool) {
	b.min.Z--
	b.max.Z--

	touching := []int{}

	for i, bb := range below {
		if b.Overlaps(bb) {
			touching = append(touching, i)
		}
	}

	return b, touching, len(touching) == 0
}

func (b Brick) Overlaps(a Brick) bool {
	return (a.min.X <= b.max.X && a.max.X >= b.min.X) &&
		(a.min.Y <= b.max.Y && a.max.Y >= b.min.Y) &&
		(a.min.Z <= b.max.Z && a.max.Z >= b.min.Z)
}

func parseInput(instr string) []Brick {
	lines := strings.Split(instr, "\n")
	bricks := make([]Brick, 0, len(lines))

	for i, line := range lines {
		var x1, y1, z1, x2, y2, z2 int

		_, err := fmt.Sscanf(line, "%d,%d,%d~%d,%d,%d", &x1, &y1, &z1, &x2, &y2, &z2)
		if err != nil {
			fmt.Printf("line %d, error parsing %s: %v\n", i, line, err)
			panic(err)
		}
		bricks = append(bricks, Brick{Vertex{x1, y1, z1}, Vertex{x2, y2, z2}})
	}

	sort.Slice(bricks, func(i, j int) bool {
		return bricks[i].min.Z < bricks[j].min.Z
	})

	return bricks
}

func getBrickOrder(bricks []Brick) (map[int][]int, map[int][]int) {
	bricksBelow := make(map[int][]int, len(bricks))
	bricksAbove := make(map[int][]int, len(bricks))

	for i := 0; i < len(bricks); {
		brick := bricks[i]

		if brick.min.Z == 1 {
			i++ // Skip processing for bricks that start at Z=1
			continue
		}

		b, supporting, ok := canMoveDown(brick, bricks[:i])
		if ok {
			bricks[i] = b
			// we're shifting the brick so don't increment i yet
		} else {
			updateBrickRelations(i, supporting, bricksBelow, bricksAbove)
			i++
		}
	}

	return bricksBelow, bricksAbove
}

// updateBrickRelations updates the bricksBelow and bricksAbove maps for the given brick.
func updateBrickRelations(curIdx int, supportBricks []int, bricksBelow, bricksAbove map[int][]int) {
	bricksBelow[curIdx] = supportBricks

	for _, supportIdx := range supportBricks {
		bricksAbove[supportIdx] = append(bricksAbove[supportIdx], curIdx)
	}
}

func getNonSupportingBricks(bricks []Brick, bricksBelow, bricksAbove map[int][]int) map[int]bool {
	nonSupportingBricks := map[int]bool{}

	for brick := 0; brick < len(bricks); brick++ {
		nonSupporting := true

		if list, found := bricksAbove[brick]; found {
			// for each brick above current brick
			for _, idx := range list {
				// if that bricks only has one brick below it (i.e. the current brick), it would fall
				if len(bricksBelow[idx]) <= 1 {
					nonSupporting = false
					break
				}
			}
		}

		// If nonSupporting is still true, it means the brick either doesn't have bricks above,
		// or none of the bricks above it are solely relying on it for support.
		if nonSupporting {
			nonSupportingBricks[brick] = true
		}
	}

	return nonSupportingBricks
}

func countDisintegratable(bricks []Brick, bricksBelow, bricksAbove map[int][]int, canDisintegrate map[int]bool) int {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	wg := &sync.WaitGroup{}
	ch := make(chan int, len(bricks))

	wg.Add(len(bricks))

	for i := 0; i < len(bricks); i++ {
		go func(i int, ch chan int) {
			defer wg.Done()

			if _, found := canDisintegrate[i]; !found {
				select {
				case ch <- countDisintegrated(bricksBelow, bricksAbove, i):
					// sent result successfully
				case <-ctx.Done():
					// timed out or cancelled
					ch <- 0
				}
			}
		}(i, ch)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// all goroutines finished
	case <-ctx.Done():
		// timed out or cancelled
		return -1
	}

	close(ch)

	total := 0

	for count := range ch {
		total += count
	}

	return total
}
