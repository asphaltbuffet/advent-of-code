package exercises

import (
	"fmt"

	util "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

type Block struct {
	HeatLoss int
	Position Point
	City     *City
}

func (b *Block) GetNeighbors() []Pather {
	if nm == nil {
		fmt.Println("nodemap is nil")
		return nil
	}

	// fmt.Printf("getting neighbors for %v\n", b.Position)
	neighbors := []Pather{}

	for _, m := range []Point{right, down, up, left} {
		next := b.City.Blocks[b.Position.add(m)]
		if next == nil {
			// fmt.Printf("no block at %v\n", b.Position.add(m))
			continue
		}

		// fmt.Printf("checking %v -> %v\n", b.Position, next.Position)
		// fmt.Printf("%v parent: %+v\n", b.Position, nm[b].parent)

		if nm[b].parent != nil {
			if next == nm[b].parent.pather.(*Block) {
				// cannot move backwards
				continue
			}

			// check if we need to change direction
			if !b.hasTurnedWithin(next.Position.X, next.Position.Y, 3) { //nolint:gomnd // don't care right now
				continue
			}
		}

		// fmt.Printf("adding %v\n as neighbor to %v\n", next.Position, b.Position)

		neighbors = append(neighbors, next)
	}

	return neighbors
}

func (b *Block) hasTurnedWithin(nextX, nextY, n int) bool {
	// fmt.Printf("checking history for %+v\n", nm[b])
	dx := util.AbsInt(b.Position.X - nextX)
	dy := util.AbsInt(b.Position.Y - nextY)

	pn := nm[b].parent

	for i := 0; i < n; i++ {
		// fmt.Printf("iteration %d, pn=%+v\n", i, pn)
		if pn == nil {
			// fmt.Printf("no parent for %v\n", b.Position)
			// we've reached the start
			return true
		}

		prev, ok := pn.pather.(*Block)
		if !ok {
			fmt.Printf("failed to cast %v to *Block\n", pn.pather)
			panic("failed to cast parent to *Block")
		}

		// fmt.Printf("prev=%v\n", prev)

		dx += util.AbsInt(b.Position.X - prev.Position.X)
		dy += util.AbsInt(b.Position.Y - prev.Position.Y)

		pn = pn.parent
	}

	return dx != 0 && dy != 0
}

func (b *Block) PathNeighborCost(to Pather) float64 {
	return float64(to.(*Block).HeatLoss)
}

func (b *Block) PathEstimatedCost(to Pather) float64 {
	dx := abs(b.Position.X - to.(*Block).Position.X)
	dy := abs(b.Position.Y - to.(*Block).Position.Y)

	return float64(dx + dy)
}
