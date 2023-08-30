package permutation

// perm := permutation.NewPermutation( slice )
//
//	for next := perm.Next(); next != nil; next = perm.Next() {
//	   ...
//	}
//
// or
//
// permCh := permutation.NewPermutationChan( ["foo", "bar", "baz"] )
//
//	for next := range ch {
//	    ...
//	}

type Permutation[P any] struct {
	orig []P
	perm []int
}

func NewPermutation[T any](orig []T) *Permutation[T] {
	p := Permutation[T]{orig: orig, perm: make([]int, len(orig))}
	return &p
}

func NewPermutationChan[T any](orig []T) chan []T {
	p := NewPermutation(orig)
	ch := make(chan []T)

	go func() {
		defer close(ch)

		for n := p.Next(); n != nil; n = p.Next() {
			ch <- n
		}
	}()

	return ch
}

func (p *Permutation[T]) nextPerm() {
	for i := len(p.perm) - 1; i >= 0; i-- {
		if i == 0 || p.perm[i] < len(p.perm)-i-1 {
			p.perm[i]++

			return
		}

		p.perm[i] = 0
	}
}

func (p *Permutation[T]) Next() []T {
	defer p.nextPerm()

	if len(p.perm) == 0 || p.perm[0] >= len(p.perm) {
		return nil
	}

	result := append([]T{}, p.orig...)

	for i, v := range p.perm {
		result[i], result[i+v] = result[i+v], result[i]
	}

	return result
}
