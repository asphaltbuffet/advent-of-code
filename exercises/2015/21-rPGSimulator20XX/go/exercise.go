package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 21.
type Exercise struct {
	common.BaseExercise
}

const playerHP = 100

// item is a shop purchase contributing cost, damage, and armor.
type item struct{ cost, damage, armor int }

// The shop tables from the problem statement (not the input).
var weapons = []item{{8, 4, 0}, {10, 5, 0}, {25, 6, 0}, {40, 7, 0}, {74, 8, 0}}

var armors = []item{{0, 0, 0}, {13, 0, 1}, {31, 0, 2}, {53, 0, 3}, {75, 0, 4}, {102, 0, 5}}

var rings = []item{{25, 1, 0}, {50, 2, 0}, {100, 3, 0}, {20, 0, 1}, {40, 0, 2}, {80, 0, 3}}

// parseBoss reads the boss's HP, damage, and armor from the input.
func parseBoss(instr string) (hp, dmg, arm int) {
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		_, val, _ := strings.Cut(line, ": ")
		n, _ := strconv.Atoi(strings.TrimSpace(val))
		switch {
		case strings.HasPrefix(line, "Hit Points"):
			hp = n
		case strings.HasPrefix(line, "Damage"):
			dmg = n
		case strings.HasPrefix(line, "Armor"):
			arm = n
		}
	}
	return
}

// playerWins reports whether the player (100 HP, attacking first) defeats the
// boss given both sides' damage and armor. Each side needs ceil(targetHP /
// effectiveDamage) hits; the player wins ties because they strike first.
func playerWins(pDmg, pArm, bHP, bDmg, bArm int) bool {
	turnsToKillBoss := ceilDiv(bHP, max(1, pDmg-bArm))
	turnsToKillPlayer := ceilDiv(playerHP, max(1, bDmg-pArm))
	return turnsToKillBoss <= turnsToKillPlayer
}

func ceilDiv(a, b int) int { return (a + b - 1) / b }

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// loadouts yields every legal equipment combination: exactly one weapon, zero
// or one armor, and zero/one/two distinct rings, as (cost, damage, armor).
func loadouts() [][3]int {
	var out [][3]int

	// Ring choices: none, one, or two distinct rings (by index).
	var ringSets [][]item
	ringSets = append(ringSets, nil)
	for i := range rings {
		ringSets = append(ringSets, []item{rings[i]})
		for j := i + 1; j < len(rings); j++ {
			ringSets = append(ringSets, []item{rings[i], rings[j]})
		}
	}

	for _, w := range weapons {
		for _, a := range armors {
			for _, rs := range ringSets {
				cost, dmg, arm := w.cost+a.cost, w.damage+a.damage, w.armor+a.armor
				for _, r := range rs {
					cost += r.cost
					dmg += r.damage
					arm += r.armor
				}
				out = append(out, [3]int{cost, dmg, arm})
			}
		}
	}

	return out
}

// One returns the least gold spent that still wins the fight.
func (e Exercise) One(instr string) (any, error) {
	bHP, bDmg, bArm := parseBoss(instr)

	best := -1
	for _, l := range loadouts() {
		cost, dmg, arm := l[0], l[1], l[2]
		if playerWins(dmg, arm, bHP, bDmg, bArm) && (best < 0 || cost < best) {
			best = cost
		}
	}

	return best, nil
}

// Two returns the most gold spent that still loses the fight.
func (e Exercise) Two(instr string) (any, error) {
	bHP, bDmg, bArm := parseBoss(instr)

	worst := 0
	for _, l := range loadouts() {
		cost, dmg, arm := l[0], l[1], l[2]
		if !playerWins(dmg, arm, bHP, bDmg, bArm) && cost > worst {
			worst = cost
		}
	}

	return worst, nil
}
