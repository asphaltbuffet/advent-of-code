package exercises

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 24.
type Exercise struct {
	common.BaseExercise
}

type group struct {
	army       int // 0 = immune system, 1 = infection
	units      int
	hp         int
	damage     int
	attackType string
	initiative int
	weak       map[string]bool
	immune     map[string]bool
}

func (g *group) effectivePower() int { return g.units * g.damage }

// damageTo returns the damage this group would deal to the defender, accounting
// for the defender's weaknesses and immunities.
func (g *group) damageTo(d *group) int {
	if d.immune[g.attackType] {
		return 0
	}
	dmg := g.effectivePower()
	if d.weak[g.attackType] {
		dmg *= 2
	}
	return dmg
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	groups := parse(instr)
	_, units := fight(groups)
	return units, nil
}

// Two returns the answer to the smallest immune-system boost that wins.
func (e Exercise) Two(instr string) (any, error) {
	for boost := 0; ; boost++ {
		groups := parse(instr)
		for _, g := range groups {
			if g.army == 0 {
				g.damage += boost
			}
		}
		winner, units := fight(groups)
		if winner == 0 {
			return units, nil
		}
	}
}

// fight runs the battle to completion and returns the winning army (0 immune, 1
// infection, -1 stalemate) and its remaining units. A round that kills no units is
// a stalemate; the immune system has not won, so it is reported as infection.
func fight(groups []*group) (int, int) {
	for {
		byArmy := [2]int{}
		for _, g := range groups {
			if g.units > 0 {
				byArmy[g.army] += g.units
			}
		}
		if byArmy[0] == 0 || byArmy[1] == 0 {
			if byArmy[0] > 0 {
				return 0, byArmy[0]
			}
			return 1, byArmy[1]
		}

		if killed := round(groups); killed == 0 {
			return 1, byArmy[1] // stalemate: immune system did not win
		}
	}
}

// round performs one target-selection + attack phase and returns the total units
// killed (0 signals a stalemate).
func round(groups []*group) int {
	targets := selectTargets(groups)

	attackers := livingGroups(groups)
	sort.Slice(attackers, func(i, j int) bool {
		return attackers[i].initiative > attackers[j].initiative
	})

	killed := 0
	for _, atk := range attackers {
		if atk.units <= 0 {
			continue
		}
		def, ok := targets[atk]
		if !ok {
			continue
		}
		dead := min(atk.damageTo(def)/def.hp, def.units)
		def.units -= dead
		killed += dead
	}
	return killed
}

func livingGroups(groups []*group) []*group {
	out := make([]*group, 0, len(groups))
	for _, g := range groups {
		if g.units > 0 {
			out = append(out, g)
		}
	}
	return out
}

func selectTargets(groups []*group) map[*group]*group {
	selectors := livingGroups(groups)
	sort.Slice(selectors, func(i, j int) bool {
		a, b := selectors[i], selectors[j]
		if a.effectivePower() != b.effectivePower() {
			return a.effectivePower() > b.effectivePower()
		}
		return a.initiative > b.initiative
	})

	targets := map[*group]*group{}
	chosen := map[*group]bool{}
	for _, atk := range selectors {
		best := chooseBestTarget(atk, livingGroups(groups), chosen)
		if best != nil {
			targets[atk] = best
			chosen[best] = true
		}
	}
	return targets
}

func chooseBestTarget(atk *group, defenders []*group, chosen map[*group]bool) *group {
	var best *group
	bestDmg := 0
	for _, def := range defenders {
		if def.army == atk.army || chosen[def] {
			continue
		}
		dmg := atk.damageTo(def)
		if dmg == 0 {
			continue
		}
		if best == nil || dmg > bestDmg ||
			(dmg == bestDmg && def.effectivePower() > best.effectivePower()) ||
			(dmg == bestDmg && def.effectivePower() == best.effectivePower() && def.initiative > best.initiative) {
			best, bestDmg = def, dmg
		}
	}
	return best
}

var lineRe = regexp.MustCompile(
	`(\d+) units each with (\d+) hit points (?:\((.*?)\) )?with an attack that does (\d+) (\w+) damage at initiative (\d+)`,
)

// parse reads both armies into a flat slice of groups.
func parse(instr string) []*group {
	var groups []*group
	army := 0
	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		switch {
		case line == "":
			continue
		case strings.HasPrefix(line, "Immune System"):
			army = 0
			continue
		case strings.HasPrefix(line, "Infection"):
			army = 1
			continue
		}
		m := lineRe.FindStringSubmatch(line)
		g := &group{
			army:       army,
			units:      atoi(m[1]),
			hp:         atoi(m[2]),
			damage:     atoi(m[4]),
			attackType: m[5],
			initiative: atoi(m[6]),
			weak:       map[string]bool{},
			immune:     map[string]bool{},
		}
		for clause := range strings.SplitSeq(m[3], "; ") {
			clause = strings.TrimSpace(clause)
			if weakTypes, okWeak := strings.CutPrefix(clause, "weak to "); okWeak {
				for t := range strings.SplitSeq(weakTypes, ", ") {
					g.weak[t] = true
				}
			} else if immuneTypes, okImmune := strings.CutPrefix(clause, "immune to "); okImmune {
				for t := range strings.SplitSeq(immuneTypes, ", ") {
					g.immune[t] = true
				}
			}
		}
		groups = append(groups, g)
	}
	return groups
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
