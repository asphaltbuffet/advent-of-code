package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/pkg/ring"
)

type coordinateFile struct {
	zero      *ring.Ring[digit]
	decrypted *ring.Ring[digit]
	encrypted []*ring.Ring[digit]
}

type digit struct {
	value   int
	visited bool
}

func parse(instr string) (*coordinateFile, error) {
	lines := strings.Split(instr, "\n")

	cf := &coordinateFile{
		decrypted: ring.New[digit](len(lines)),
		encrypted: []*ring.Ring[digit]{},
	}

	for i, line := range lines {
		var err error

		v, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("parsing line %d: %w", i, err)
		}

		cf.decrypted.Element = digit{v, false}
		cf.encrypted = append(cf.encrypted, cf.decrypted)

		// create an anchor to the zero value
		if v == 0 {
			cf.zero = cf.decrypted
		}

		cf.decrypted = cf.decrypted.Next()
	}

	return cf, nil
}

func (cf *coordinateFile) decryptedToString() string {
	var v []string

	cf.zero.Do(func(d digit) {
		v = append(v, strconv.Itoa(d.value))
	})

	return strings.Join(v, " ")
}

func (cf *coordinateFile) decrypt() error {
	// fmt.Println("initial\n", cf.Print())

	if len(cf.encrypted) != cf.decrypted.Len() {
		return fmt.Errorf("encrypted and decrypted rings are not the same length")
	}

	for _, e := range cf.encrypted {
		// fmt.Println(cf.PrintMove(e.Element.value))
		// fmt.Println(cf.Print())

		shift(e, e.Element.value)
	}

	return nil
}

// return the 1000th, 2000th, and 3000th numbers in the ring starting at value 0
func (cf *coordinateFile) getCoordinates() (int, int, int) {
	c := cf.zero.Move(1000)
	c1 := c.Element.value

	c = c.Move(1000)
	c2 := c.Element.value

	c = c.Move(1000)
	c3 := c.Element.value

	return c1, c2, c3
}

func shift(r *ring.Ring[digit], n int) {
	l := r.Len() - 1
	h := l >> 1

	// unlink starts at next element, so move back one
	i := r.Prev()
	removed := i.Unlink(1)

	// optimization, if moving more than half length, faster to go other way
	if n > h || n < -h {
		n %= l

		switch {
		case n > h:
			n -= l
		case n < -h:
			n += l
		}
	}

	// Move and link (insert) the removed element back.
	i.Move(n).Link(removed)
}
