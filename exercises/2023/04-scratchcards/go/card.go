package exercises

import (
	"fmt"
	"math"
	"strings"
)

type Card struct {
	ID        int
	Winning   []string
	Scratched []string
}

func New(s string) *Card {
	// Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
	// Card   1: 33 56 23 64 92 86 94  7 59 13 | 86 92 64 43 10 70 16 55 79 33 56  8  7 25 82 14 31 96 94 13 99 29 69 75 23
	id, data, _ := strings.Cut(s, ":")
	win, scratch, _ := strings.Cut(data, "|")

	var n int
	fmt.Sscanf(id, "Card %d", &n)

	return &Card{
		ID:        n,
		Winning:   strings.Fields(win),
		Scratched: strings.Fields(scratch),
	}
}

func (c *Card) Score() int {
	count := c.Count()
	return int(math.Pow(2, float64(count-1)))
}

func (c *Card) Count() int {
	var count int

	for _, w := range c.Winning {
		for _, s := range c.Scratched {
			if w == s {
				count++
				break
			}
		}
	}

	return count
}

func countTotalCards(lines []string) int {
	cardQtys := make(map[int]int, len(lines))
	cardRef := make(map[int]*Card, len(lines))

	for _, line := range lines {
		card := New(line)
		cardQtys[card.ID]++
		cardRef[card.ID] = card
	}

	// card IDs are 1-indexed
	for i := 1; i <= len(lines); i++ {
		// fmt.Printf("card %3d: increase next %2d  (+%d)\n", cardRef[i].ID, cardRef[i].Count(), cardQtys[cardRef[i].ID])

		// update the <count> cards after <id> by amount in cardQtys
		addWonCards(cardRef[i].ID, cardRef[i].Count(), cardQtys)
	}

	var total int
	for _, qty := range cardQtys {
		total += qty
	}

	return total
}

func addWonCards(id, count int, qty map[int]int) {
	for j := 1; j <= count; j++ {
		next := id + j

		// safety check to prevent a runaway loop... again
		if _, ok := qty[next]; !ok {
			panic("card not found")
		}

		// fmt.Printf("+ %d card %ds\n", mult[id], next)

		// increase qty of "next" cards by qty of active card
		qty[next] += qty[id]
	}
}
