package exercises

import (
	"fmt"
	"strconv"
	"strings"
)

type Box struct {
	ID     int
	Lenses []Lens
}

type Lens struct {
	FocalLength int
	Label       string
}

func sumAllSteps(input string) int {
	steps := strings.Split(input, ",")
	var sum int

	for _, step := range steps {
		sum += hash(step)
	}

	return sum
}

func hash(input string) int {
	var h int

	for _, c := range input {
		h += int(c)

		h *= 17
		h %= 256
	}

	return h
}

type Op struct {
	Label       string
	Box         int
	Action      rune
	FocalLength int
}

func parseStep(step string) (*Op, error) {
	var (
		name   string
		value  int
		action = '='
	)

	name, fl, ok := strings.Cut(step, "=")
	if ok {
		var err error
		value, err = strconv.Atoi(fl)
		if err != nil {
			return nil, err
		}
	} else {
		action = '-'
		name, ok = strings.CutSuffix(step, "-")
		if !ok {
			return nil, fmt.Errorf("invalid step: %q", step)
		}
	}

	return &Op{
		Label:       name,
		Box:         hash(name),
		Action:      action,
		FocalLength: value,
	}, nil
}

func calcFocusingPower(ops []*Op) int {
	boxes := make(map[int]*Box)

	for _, op := range ops {
		box, ok := boxes[op.Box]
		if !ok {
			box = &Box{ID: op.Box, Lenses: []Lens{}}
			boxes[op.Box] = box
		}

		var replaced bool

		for i, lens := range box.Lenses {
			if lens.Label == op.Label {
				if op.Action == '=' {
					box.Lenses[i].FocalLength = op.FocalLength
					replaced = true
					break
				}

				// action == '-'
				box.Lenses = append(box.Lenses[:i], box.Lenses[i+1:]...)
				break
			}
		}

		if !replaced && op.Action == '=' {
			box.Lenses = append(box.Lenses, Lens{
				FocalLength: op.FocalLength,
				Label:       op.Label,
			})
		}

		// probably not necessary, but keeping it clean for now
		if len(box.Lenses) == 0 {
			delete(boxes, op.Box)
		}

		// // debug print
		// fmt.Printf("after %q:\n", op.Orig)
		// for _, b := range boxes {
		// 	fmt.Println(b)
		// }
		// fmt.Println()
	}

	return calcBoxPower(boxes)
}

func calcBoxPower(boxes map[int]*Box) int {
	var power int
	lensPowers := make(map[string]int)

	for b, box := range boxes {
		for l, lens := range box.Lenses {
			lens := lens
			if _, ok := lensPowers[lens.Label]; !ok {
				lensPowers[lens.Label] = 1
			}

			lensPowers[lens.Label] *= (b + 1) * (l + 1) * lens.FocalLength
		}
	}

	for _, lp := range lensPowers {
		// fmt.Printf("%s: %-4d\n", ln, lp)
		power += lp
	}

	return power
}

func (b *Box) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Box %d:", b.ID))
	for _, lens := range b.Lenses {
		sb.WriteString(fmt.Sprintf(" [%s %d]", lens.Label, lens.FocalLength))
	}

	return sb.String()
}
