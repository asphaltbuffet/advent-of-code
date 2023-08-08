// Package set provides a basic set implementation for comparable types in Go.
package set

// Set represents a collection of elements with no duplicates.
// The underlying structure is a map with the set item type as the key.
// The value is an empty struct, ensuring minimal memory usage.
type Set[T comparable] map[T]struct{}

// Add inserts the provided items into the set. If an item already exists in the set,
// it will not create a duplicate.
func (set *Set[T]) Add(items ...T) {
	for _, item := range items {
		(*set)[item] = struct{}{}
	}
}

// Delete removes the specified item from the set. If the item doesn't exist in the set,
// the operation is a no-op.
func (set *Set[T]) Delete(item T) {
	delete(*set, item)
}

// Has checks if the specified item exists in the set and returns a boolean value.
func (set *Set[T]) Has(item T) bool {
	_, ok := (*set)[item]

	return ok
}
