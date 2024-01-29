package exercises

import (
	"fmt"
	"strings"
)

type Part struct {
	Rating map[Category]int
	Value  int
}

func parsePart(in string) *Part {
	s := strings.Trim(in, "{}")
	tokens := strings.Split(s, ",")

	p := &Part{
		Rating: make(map[Category]int),
	}

	for _, token := range tokens {
		var c rune
		var n int

		r, err := fmt.Sscanf(token, "%c=%d", &c, &n)
		if err != nil || r != 2 {
			fmt.Printf("read %d tokens: %v\n", r, err)
			return nil
		}

		p.Rating[getCategory(c)] = n
		p.Value += n
	}

	return p
}

func (p *Part) Process(w Workflows, flow string) Result {
	var next string

	curFlow, ok := w[flow]
	if !ok {
		fmt.Printf("workflow %q not found\n", flow)
		return Rejected
	}

	for {
		for _, t := range curFlow.Tests {
			testFunc := t.Test
			if testFunc(p.Rating[t.Category]) {
				next = t.Dest
				break
			}
		}

		switch next {
		case "R":
			// fmt.Printf("rejected part %v in workflow %q\n", p, flow)
			return Rejected
		case "A":
			return Accepted
		default:
			curFlow, ok = w[next]
			if !ok {
				fmt.Printf("next workflow %q not found\n", next)
				return Rejected
			}
		}
	}
}
