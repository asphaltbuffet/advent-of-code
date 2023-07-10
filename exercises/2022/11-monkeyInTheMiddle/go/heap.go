package exercises

import "container/heap"

func getHeap(mh []*Monkey) *MonkeyHeap {
	h := &MonkeyHeap{}
	heap.Init(h)

	for _, v := range mh {
		v := *v
		heap.Push(h, v)
	}

	return h
}

// MonkeyHeap implements heap.Interface and holds Monkeys.
// Ref https://golang.org/pkg/container/heap/
type MonkeyHeap []Monkey

// Less is greater-than here so we can pop *larger* items.
func (h MonkeyHeap) Less(i, j int) bool { return h[i].Count > h[j].Count }

// Swap swaps the elements with indexes i and j.
func (h MonkeyHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Len is the number of elements in the collection.
func (h MonkeyHeap) Len() int { return len(h) }

// Push pushes the element x onto the heap.
func (h *MonkeyHeap) Push(x interface{}) {
	*h = append(*h, x.(Monkey))
}

// Pop removes and returns the maximum element (according to Less) from the heap.
func (h *MonkeyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}
