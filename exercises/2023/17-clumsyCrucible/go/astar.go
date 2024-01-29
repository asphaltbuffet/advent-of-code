package exercises

import (
	"container/heap"
	"fmt"
	"image"
	"strings"
	"time"
)

var (
	nm nodeMap
	sm stateMap
)

type Pather interface {
	GetNeighbors() []Pather
	PathNeighborCost(to Pather) float64
	PathEstimatedCost(to Pather) float64
}

type Direction rune

const (
	UpDir    Direction = 'u'
	DownDir  Direction = 'd'
	RightDir Direction = 'r'
	LeftDir  Direction = 'l'
)

type node struct {
	pather  Pather
	history []Direction
	cost    float64
	rank    float64
	parent  *node
	open    bool
	closed  bool
	index   int
}

type nodeMap map[Pather]*node

func (nm nodeMap) get(p Pather) *node {
	n, ok := nm[p]
	if !ok {
		n = &node{pather: p}
		nm[p] = n
	}

	return n
}

type stateMap map[Pather]map[string]*node

func (sm stateMap) get(p Pather, hist ...Direction) *node {
	// not safe for concurrent use
	var sb strings.Builder
	for _, d := range hist {
		sb.WriteRune(rune(d))
	}

	hash := sb.String()

	if sm[p] == nil {
		sm[p] = make(map[string]*node)
	}

	n, ok := sm[p][hash]
	if !ok {
		n = &node{pather: p, history: hist}
		sm[p][hash] = n
	}

	return n
}

func getNodeHash(n *node) string {
	if n == nil {
		return ""
	}

	const hashLen = 4

	var hash string
	i := 0

	for v := n; v != nil && i < hashLen; v = v.parent {
		// fmt.Printf("v=%+v\n", v)
		if v.pather == nil {
			hash += "<nil>"
		} else {
			hash += fmt.Sprintf("<%d-%d>", v.pather.(*Block).Position.X, v.pather.(*Block).Position.Y)
		}

		i++
	}

	return hash
}

func (c *City) Path(start, end Pather) ([]Pather, float64, bool) {
	frames := []*image.RGBA{generateBackgroundImage(c, nil)}

	fmt.Printf("Path from %v to %v\n", start.(*Block).Position, end.(*Block).Position)

	nm = nodeMap{}
	// sm = stateMap{}
	pq := &priorityQueue{}

	heap.Init(pq)

	fromNode := nm.get(start)
	// fromNode := sm.get(start, RightDir)
	fromNode.open = true

	heap.Push(pq, fromNode)

	// loop until we get to goal or queue is empty
	for i := 0; ; i++ {
		if pq.Len() == 0 {
			// There's no path to the goal
			return nil, 0, false
		}

		current, ok := heap.Pop(pq).(*node)
		if !ok {
			fmt.Println("failed to pop from queue")
			return nil, 0, false
		}

		// hist := getNodeHash(current)

		// DEBUG: generate animation frame
		if i%500 == 0 {
			p := buildPath(current)
			frame, err := c.GenerateFrame(p, fmt.Sprintf("%s cost=%.0f", p[0].(*Block).Position.String(), nm[p[0]].cost))
			// frame, err := c.GenerateFrame(p, fmt.Sprintf("%s cost=%.0f", p[0].(*Block).Position.String(), sm[p[0]][hist].cost))
			if err == nil {
				frames = append(frames, frame)
			}
		}

		current.open = false
		current.closed = true

		if current == nm.get(end) {
			// if isEnd(current, end) {
			// DEBUG: generate animation
			fmt.Printf("generating animation (%d frames)\n", len(frames))
			outfile := fmt.Sprintf("./images/p1-%s.png", time.Now().Format("20060102_150405"))
			err := GenerateAPNG(frames, outfile)
			if err != nil {
				fmt.Println("failed to generate animation")
			}

			// Found a path to the goal.
			return buildPath(current), current.cost, true
		}

		addNeighborsToQueue(current, end, pq)
	}
}

func isEnd(n *node, end Pather) bool {
	if n == nil || n.pather == nil {
		return false
	}

	b, ok := n.pather.(*Block)
	if !ok {
		return false
	}

	return b.Position == end.(*Block).Position
}

func buildPath(n *node) []Pather {
	path := []Pather{}
	curr := n

	for curr != nil {
		path = append(path, curr.pather)
		curr = curr.parent
	}

	return path
}

func addNeighborsToQueue(current *node, to Pather, pq *priorityQueue) {
	for _, neighbor := range current.pather.GetNeighbors() {
		cost := current.cost + current.pather.PathNeighborCost(neighbor)
		neighborNode := nm.get(neighbor)

		// if we have a better path (lower cost), update the neighbor and make it active
		if cost < neighborNode.cost {
			if neighborNode.open { // FIX: it's more efficient to update the node than remove/add
				heap.Remove(pq, neighborNode.index)
			}

			neighborNode.open = false
			neighborNode.closed = false
		}

		if !neighborNode.open && !neighborNode.closed {
			neighborNode.cost = cost
			neighborNode.open = true
			neighborNode.rank = cost + neighbor.PathEstimatedCost(to)
			neighborNode.parent = current

			heap.Push(pq, neighborNode)
		}
	}
}
