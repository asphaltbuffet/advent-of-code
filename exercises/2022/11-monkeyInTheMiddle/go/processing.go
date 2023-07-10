package exercises

import (
	"fmt"
)

// ProcessRound processes a round of monkeys.
func (d *Day11) ProcessRound() error {
	for _, m := range d.Monkeys {
		m := m

		for _, i := range m.Items {
			item := 0

			switch m.Operator {
			case "+":
				item = i + m.Scalar
			case "*":
				item = i * m.Scalar
			case "^":
				item = i * i
			default:
				return fmt.Errorf("unknown operator: %s", m.Operator)
			}

			item /= 3

			if item%m.Divisor == 0 {
				d.Monkeys[m.TargetOne].Items = append(d.Monkeys[m.TargetOne].Items, item)
			} else {
				d.Monkeys[m.TargetTwo].Items = append(d.Monkeys[m.TargetTwo].Items, item)
			}

			m.Count++
		}

		m.Items = nil
	}

	// for j, m := range d.Monkeys {
	// 	log.Printf("monkey %d items: %v", j, m.Items)
	// }

	return nil
}

// ProcessRoundPart2 processes a round of monkeys.
func (d *Day11) ProcessRoundPart2() error {
	for _, m := range d.Monkeys {
		m := m

		for _, i := range m.Items {
			item := 0

			switch m.Operator {
			case "+":
				item = i + m.Scalar
			case "*":
				item = i * m.Scalar
			case "^":
				item = i * i
			default:
				return fmt.Errorf("unknown operator: %s", m.Operator)
			}

			item %= d.Product

			if item%m.Divisor == 0 {
				d.Monkeys[m.TargetOne].Items = append(d.Monkeys[m.TargetOne].Items, item)
			} else {
				d.Monkeys[m.TargetTwo].Items = append(d.Monkeys[m.TargetTwo].Items, item)
			}

			m.Count++
		}

		m.Items = nil
	}

	// // debug print monkey items
	// for j, m := range d.Monkeys {
	// 	log.Printf("monkey %d items: %v", j, m.Items)
	// }

	return nil
}
