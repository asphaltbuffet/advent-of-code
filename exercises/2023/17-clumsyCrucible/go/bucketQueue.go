package exercises

import (
	"container/list"
	"math"
)

// BucketQueue represents a bucket queue data structure
type BucketQueue struct {
	buckets map[int]*list.List // Buckets of elements
	min     int                // Minimum priority in the queue
}

// NewBucketQueue initializes a new BucketQueue with a given maxPriority
func NewBucketQueue() *BucketQueue {
	bq := &BucketQueue{
		buckets: map[int]*list.List{},
		min:     math.MaxInt,
	}

	return bq
}

// Enqueue adds a new element to the queue with a specified priority
func (bq *BucketQueue) Enqueue(priority int, value any) {
	if _, ok := bq.buckets[priority]; !ok {
		bq.buckets[priority] = list.New()

		if priority < bq.min {
			bq.min = priority
		}
	}

	bq.buckets[priority].PushBack(value)
}

// Dequeue removes and returns the element with the highest priority (lowest number)
func (bq *BucketQueue) Dequeue() (any, bool) {
	b, ok := bq.buckets[bq.min]
	if !ok || b.Len() == 0 {
		return nil, false
	}

	e := b.Front()
	b.Remove(e)

	if b.Len() == 0 {
		bq.update()
	}

	return e.Value, true
}

// update sets min to the lowest priority and removes empty buckets.
func (bq *BucketQueue) update() {
	bq.min = math.MaxInt

	for pri, b := range bq.buckets {
		if b.Len() > 0 && pri < bq.min {
			bq.min = pri
		} else if b.Len() == 0 {
			delete(bq.buckets, pri)
		}
	}
}

// IsEmpty checks if the queue is empty
func (bq *BucketQueue) IsEmpty() bool {
	return len(bq.buckets) == 0
}
