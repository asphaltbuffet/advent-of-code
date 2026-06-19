package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 22.
type Exercise struct {
	common.BaseExercise
}

// state is the full fight state. Timers count remaining ticks for each effect.
type state struct {
	php, mana, bhp, bdmg     int
	shield, poison, recharge int
	spent                    int
}

// applyEffects ticks the active effects once (poison damage, recharge mana) and
// decrements each timer. Shield's armor is read directly from its timer.
func (s *state) applyEffects() {
	if s.poison > 0 {
		s.bhp -= 3
		s.poison--
	}
	if s.recharge > 0 {
		s.mana += 101
		s.recharge--
	}
	if s.shield > 0 {
		s.shield--
	}
}

// search runs a branch-and-bound DFS for the least mana spent to win. hard adds
// the part-two rule (the player loses 1 HP at the start of each of their turns).
func search(s state, hard bool, best *int) {
	if s.spent >= *best {
		return // can't beat the current best from here
	}

	// --- Player turn ---
	if hard {
		s.php--
		if s.php <= 0 {
			return // died before acting
		}
	}
	s.applyEffects()
	if s.bhp <= 0 {
		if s.spent < *best {
			*best = s.spent
		}
		return
	}

	// Try each castable spell (instant spells resolve now; effect spells may
	// only be cast when not already active).
	type spell struct {
		cost int
		cast func(*state)
	}
	spells := []spell{
		{53, func(n *state) { n.bhp -= 4 }},                       // Magic Missile
		{73, func(n *state) { n.bhp -= 2; n.php += 2 }},           // Drain
		{113, func(n *state) { n.shield = 6 }},                    // Shield
		{173, func(n *state) { n.poison = 6 }},                    // Poison
		{229, func(n *state) { n.recharge = 5 }},                  // Recharge
	}

	for _, sp := range spells {
		if sp.cost > s.mana {
			continue
		}
		// Effect spells can't be cast while still active.
		switch sp.cost {
		case 113:
			if s.shield > 0 {
				continue
			}
		case 173:
			if s.poison > 0 {
				continue
			}
		case 229:
			if s.recharge > 0 {
				continue
			}
		}

		next := s
		next.mana -= sp.cost
		next.spent += sp.cost
		sp.cast(&next)

		if next.bhp <= 0 {
			if next.spent < *best {
				*best = next.spent
			}
			continue
		}

		// --- Boss turn ---
		next.applyEffects()
		if next.bhp <= 0 {
			if next.spent < *best {
				*best = next.spent
			}
			continue
		}
		armor := 0
		if next.shield > 0 {
			armor = 7
		}
		hit := next.bdmg - armor
		if hit < 1 {
			hit = 1
		}
		next.php -= hit
		if next.php <= 0 {
			continue // player died
		}

		search(next, hard, best)
	}
}

// parseBoss reads the boss's hit points and damage.
func parseBoss(instr string) (hp, dmg int) {
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		_, val, _ := strings.Cut(line, ": ")
		n, _ := strconv.Atoi(strings.TrimSpace(val))
		if strings.HasPrefix(line, "Hit Points") {
			hp = n
		} else if strings.HasPrefix(line, "Damage") {
			dmg = n
		}
	}
	return
}

func solve(instr string, hard bool) int {
	bhp, bdmg := parseBoss(instr)
	best := 1 << 30
	search(state{php: 50, mana: 500, bhp: bhp, bdmg: bdmg}, hard, &best)
	return best
}

// One returns the least mana spent to win the fight.
func (e Exercise) One(instr string) (any, error) {
	return solve(instr, false), nil
}

// Two returns the least mana spent to win on hard mode.
func (e Exercise) Two(instr string) (any, error) {
	return solve(instr, true), nil
}
