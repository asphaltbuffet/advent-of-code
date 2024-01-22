package exercises

import (
	"fmt"
	"math"
)

type State struct {
	point, vector Point
	dist          int
}

func countHeatLoss(city map[Point]int, start, end Point, minDist, maxDist int) int {
	minHeatLoss := math.MaxInt32

	bq := NewBucketQueue()

	bq.Enqueue(0, State{start, Point{1, 0}, 0})
	bq.Enqueue(0, State{start, Point{0, 1}, 0})

	visited := map[State]int{{start, Point{0, 0}, 0}: 0}

	for !bq.IsEmpty() {
		current := bq.Pop()

		if current.point == end {
			if current.dist >= minDist {
				minHeatLoss = min(minHeatLoss, visited[current])
				break
			}

			continue
		}

		var skipLeft, skipRight bool
		nextPoint := Point{current.point.X + current.vector.X, current.point.Y + current.vector.Y}

		// check if we can go straight
		if _, ok := city[nextPoint]; ok {
			if current.dist < maxDist {
				nextState := State{nextPoint, current.vector, current.dist + 1}
				totalHeatLoss := visited[current] + city[nextState.point]

				if val, found := visited[nextState]; !found || val > totalHeatLoss {
					visited[nextState] = totalHeatLoss
					bq.Enqueue(totalHeatLoss, nextState)
				}
			}
		} else {
			if current.vector.X == 1 || current.vector.Y == -1 {
				skipLeft = true
			} else if current.vector.X == -1 || current.vector.Y == 1 {
				skipRight = true
			}
		}

		if current.dist < minDist {
			// we can't turn yet
			continue
		}

		// check if we can turn left
		if !skipLeft {
			leftVector := turnLeft(current.vector)
			nextPoint = Point{current.point.X + leftVector.X, current.point.Y + leftVector.Y}

			if _, ok := city[nextPoint]; ok {
				nextState := State{nextPoint, leftVector, 1}
				totalHeatLoss := visited[current] + city[nextState.point]

				if val, found := visited[nextState]; !found || val > totalHeatLoss {
					visited[nextState] = totalHeatLoss
					bq.Enqueue(totalHeatLoss, nextState)
				}
			}
		}

		// check if we can turn right
		if !skipRight {
			rightVector := turnRight(current.vector)
			nextPoint = Point{current.point.X + rightVector.X, current.point.Y + rightVector.Y}

			if _, ok := city[nextPoint]; ok {
				nextState := State{nextPoint, rightVector, 1}
				totalHeatLoss := visited[current] + city[nextState.point]

				if val, found := visited[nextState]; !found || val > totalHeatLoss {
					visited[nextState] = totalHeatLoss
					bq.Enqueue(totalHeatLoss, nextState)
				}
			}
		}
	}

	return minHeatLoss
}

func (bq *BucketQueue) Pop() State {
	element, popped := bq.Dequeue()
	if !popped {
		fmt.Println("no minimum element found")
		panic("no minimum element found")
	}

	state, ok := element.(State)
	if !ok {
		fmt.Println("element is not a HeatState")
		panic("element is not a HeatState")
	}

	return state
}

func turnLeft(p Point) Point {
	return Point{p.Y, -p.X}
}

func turnRight(p Point) Point {
	return Point{-p.Y, p.X}
}
