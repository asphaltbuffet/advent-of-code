package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// sortFloats sorts a slice of float64 ascending, in place.
func sortFloats(s []float64) { sort.Float64s(s) }

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
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		nums := signedInt.FindAllString(line, -1)
		n := make([]int, 9)
		for i := 0; i < 9; i++ {
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

	for tick := 0; tick < 1000; tick++ {
		for i := range ps {
			if !alive[i] {
				continue
			}
			ps[i].vel = ps[i].vel.add(ps[i].acc)
			ps[i].pos = ps[i].pos.add(ps[i].vel)
		}
		// Remove any particles now sharing a position.
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

	count := 0
	for _, a := range alive {
		if a {
			count++
		}
	}
	return count, nil
}

// Vis animates the swarm as a GIF: particles are projected onto the x/y plane
// with a signed-log scale (so the dense centre and the flung-out survivors both
// stay on canvas), and collisions flash as they resolve.
func (e Exercise) Vis(instr, outdir string) error {
	const frames = 90
	const size = 600
	const cx, cy = size / 2, size / 2

	base := parseParticles(instr)

	// First pass: simulate to find the median projected extent, so the window
	// frames the dense central swarm (where collisions happen) rather than the
	// few particles that fly far out. A linear scale reads as a real cloud.
	sim := make([]particle, len(base))
	copy(sim, base)
	aliveScan := make([]bool, len(base))
	for i := range aliveScan {
		aliveScan[i] = true
	}
	var extents []float64
	for f := 0; f < frames; f++ {
		for i := range sim {
			if !aliveScan[i] {
				continue
			}
			sim[i].vel = sim[i].vel.add(sim[i].acc)
			sim[i].pos = sim[i].pos.add(sim[i].vel)
			extents = append(extents,
				math.Abs(float64(sim[i].pos.x)), math.Abs(float64(sim[i].pos.y)))
		}
		at := map[v3][]int{}
		for i := range sim {
			if aliveScan[i] {
				at[sim[i].pos] = append(at[sim[i].pos], i)
			}
		}
		for _, g := range at {
			if len(g) > 1 {
				for _, i := range g {
					aliveScan[i] = false
				}
			}
		}
	}
	// Frame to roughly the 85th-percentile extent so outliers don't shrink the
	// swarm to a dot; particles beyond the window simply clip.
	sortFloats(extents)
	maxExt := extents[len(extents)*85/100]
	if maxExt < 1 {
		maxExt = 1
	}
	scale := (size/2 - 12) / maxExt

	// Palette: background, live particle, collision flash.
	pal := color.Palette{
		color.RGBA{0x0d, 0x0f, 0x18, 0xff}, // 0 background
		color.RGBA{0x2f, 0x8a, 0xff, 0xff}, // 1 live particle
		color.RGBA{0xff, 0xc8, 0x4a, 0xff}, // 2 dense/overlapping
		color.RGBA{0xff, 0x44, 0x55, 0xff}, // 3 collision flash
	}

	g := &gif.GIF{}

	// Second pass: render each frame.
	copy(sim, base)
	alive := make([]bool, len(base))
	for i := range alive {
		alive[i] = true
	}
	project := func(p v3) (int, int) {
		return cx + int(float64(p.x)*scale), cy - int(float64(p.y)*scale)
	}
	plot := func(img *image.Paletted, x, y int, idx uint8) {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				px, py := x+dx, y+dy
				if px >= 0 && px < size && py >= 0 && py < size {
					img.SetColorIndex(px, py, idx)
				}
			}
		}
	}

	for f := 0; f < frames; f++ {
		for i := range sim {
			if !alive[i] {
				continue
			}
			sim[i].vel = sim[i].vel.add(sim[i].acc)
			sim[i].pos = sim[i].pos.add(sim[i].vel)
		}
		at := map[v3][]int{}
		for i := range sim {
			if alive[i] {
				at[sim[i].pos] = append(at[sim[i].pos], i)
			}
		}

		img := image.NewPaletted(image.Rect(0, 0, size, size), pal)
		// Draw survivors; flag colliding cells.
		for pos, group := range at {
			x, y := project(pos)
			idx := uint8(1)
			if len(group) > 1 {
				idx = 3 // collision this frame
			}
			plot(img, x, y, idx)
		}
		g.Image = append(g.Image, img)
		g.Delay = append(g.Delay, 8) // ~80ms/frame

		// Resolve collisions for the next frame.
		for _, group := range at {
			if len(group) > 1 {
				for _, i := range group {
					alive[i] = false
				}
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "particle-swarm.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, g)
}
