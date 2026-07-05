package exercises

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 20.
type Exercise struct {
	common.BaseExercise
}

type v3 struct{ x, y, z int }

func (a v3) add(b v3) v3 { return v3{a.x + b.x, a.y + b.y, a.z + b.z} }
func (a v3) manhattan() int {
	return abs(a.x) + abs(a.y) + abs(a.z)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type particle struct {
	pos, vel, acc v3
}

var signedInt = regexp.MustCompile(`-?\d+`)

// parseParticles reads each line's nine integers into a particle.
func parseParticles(instr string) []particle {
	var ps []particle
	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		nums := signedInt.FindAllString(line, -1)
		n := make([]int, 9)
		for i := range 9 {
			n[i], _ = strconv.Atoi(nums[i])
		}
		ps = append(ps, particle{
			pos: v3{n[0], n[1], n[2]},
			vel: v3{n[3], n[4], n[5]},
			acc: v3{n[6], n[7], n[8]},
		})
	}
	return ps
}

// One returns the index of the particle that stays closest to the origin in the
// long run: smallest acceleration magnitude, tie-broken by velocity then
// position.
func (e Exercise) One(instr string) (any, error) {
	ps := parseParticles(instr)
	best := 0
	for i := 1; i < len(ps); i++ {
		if less(ps[i], ps[best]) {
			best = i
		}
	}
	return best, nil
}

// less reports whether particle a will end up closer to the origin than b.
func less(a, b particle) bool {
	if aa, ba := a.acc.manhattan(), b.acc.manhattan(); aa != ba {
		return aa < ba
	}
	if av, bv := a.vel.manhattan(), b.vel.manhattan(); av != bv {
		return av < bv
	}
	return a.pos.manhattan() < b.pos.manhattan()
}

// Two simulates the swarm, removing particles that share a position each tick,
// and returns how many survive. Several hundred ticks suffice: once
// accelerations pull particles apart they never meet again.
func (e Exercise) Two(instr string) (any, error) {
	ps := parseParticles(instr)
	alive := make([]bool, len(ps))
	for i := range alive {
		alive[i] = true
	}

	for range 1000 {
		for i := range ps {
			if !alive[i] {
				continue
			}
			ps[i].vel = ps[i].vel.add(ps[i].acc)
			ps[i].pos = ps[i].pos.add(ps[i].vel)
		}
		removeCollisions(ps, alive)
	}

	count := 0
	for _, a := range alive {
		if a {
			count++
		}
	}
	return count, nil
}

func removeCollisions(ps []particle, alive []bool) {
	at := map[v3][]int{}
	for i := range ps {
		if alive[i] {
			at[ps[i].pos] = append(at[ps[i].pos], i)
		}
	}
	for _, group := range at {
		if len(group) > 1 {
			for _, i := range group {
				alive[i] = false
			}
		}
	}
}
