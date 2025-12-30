package exercises

import (
	"fmt"
	"iter"
	"math"
	"strconv"
	"strings"
)

type Machine struct {
	Lights  int
	Buttons []int
	Joltage []int
}

func ParseMachines(s string) []*Machine {
	lines := strings.Split(s, "\n")
	mm := make([]*Machine, len(lines))

	for i, l := range lines {
		tok := strings.Fields(l)

		lights := lightsToInt(tok[0])
		btns := buttonsToInts(tok[1 : len(tok)-1])

		// fmt.Println("light:", lights)
		// fmt.Println("buttons:", btns)

		mm[i] = &Machine{
			Lights:  lights,
			Buttons: btns,
			Joltage: []int{},
		}
	}

	return mm
}

func lightsToInt(s string) int {
	b := strings.Map(func(r rune) rune {
		switch r {
		case '.':
			return '0'
		case '#':
			return '1'
		default:
			return r
		}
	}, strings.Trim(s, "[]"))

	b = reverse(b)

	var n int
	fmt.Sscanf(b, "%b", &n)

	return n
}

func buttonsToInts(s []string) []int {
	buttons := make([]int, len(s))

	for i, b := range s {
		sum := 0
		for _, n := range strings.Split(strings.Trim(b, "()"), ",") {
			digit, _ := strconv.Atoi(n)
			sum += int(math.Pow(2, float64(digit)))
		}

		buttons[i] = sum
	}

	return buttons
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func XORComb(xs []int, k int) iter.Seq[int] {
	return func(yield func(int) bool) {
		n := len(xs)
		if k <= 0 || k > n {
			return
		}

		// Initial index combination: [0,1,2,...,k-1]
		idx := make([]int, k)
		for i := 0; i < k; i++ {
			idx[i] = i
		}

		for {
			// Compute OR for this combination
			orv := 0
			for _, j := range idx {
				orv ^= xs[j]
			}
			if !yield(orv) {
				return
			}

			// Generate next lexicographic combination
			i := k - 1
			for i >= 0 && idx[i] == n-k+i {
				i--
			}
			if i < 0 {
				break
			}
			idx[i]++
			for j := i + 1; j < k; j++ {
				idx[j] = idx[j-1] + 1
			}
		}
	}
}

func (m *Machine) GetButtonPresses() int {
	for i := 1; i < len(m.Buttons); i++ {
		for c := range XORComb(m.Buttons, i) {
			if m.Lights == c {
				return i
			}
		}
	}

	return -1
}
