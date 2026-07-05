package exercises

import (
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"
	"path/filepath"
	"sort"
)

func sortFloats(s []float64) { sort.Float64s(s) }

// Vis animates the swarm as a GIF: particles are projected onto the x/y plane
// and collisions flash red before the particles vanish. The frame window is
// scaled to the 85th-percentile extent so outliers don't shrink the dense swarm.
func (e Exercise) Vis(instr, outdir string) error {
	const frames = 90
	const size = 600

	base := parseParticles(instr)
	scale := computeScale(base, frames, size)

	pal := color.Palette{
		color.RGBA{0x0d, 0x0f, 0x18, 0xff}, // 0 background
		color.RGBA{0x2f, 0x8a, 0xff, 0xff}, // 1 live particle
		color.RGBA{0xff, 0xc8, 0x4a, 0xff}, // 2 dense/overlapping
		color.RGBA{0xff, 0x44, 0x55, 0xff}, // 3 collision flash
	}

	g := renderParticleGIF(base, frames, size, scale, pal)

	f, err := os.Create(filepath.Join(outdir, "particle-swarm.gif"))
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, g)
}

func computeScale(base []particle, frames, size int) float64 {
	sim := make([]particle, len(base))
	copy(sim, base)
	alive := make([]bool, len(base))
	for i := range alive {
		alive[i] = true
	}
	var extents []float64
	for range frames {
		for i := range sim {
			if !alive[i] {
				continue
			}
			sim[i].vel = sim[i].vel.add(sim[i].acc)
			sim[i].pos = sim[i].pos.add(sim[i].vel)
			extents = append(extents,
				math.Abs(float64(sim[i].pos.x)), math.Abs(float64(sim[i].pos.y)))
		}
		removeCollisions(sim, alive)
	}
	sortFloats(extents)
	maxExt := extents[len(extents)*85/100]
	if maxExt < 1 {
		maxExt = 1
	}
	return (float64(size)/2 - 12) / maxExt
}

func renderParticleGIF(base []particle, frames, size int, scale float64, pal color.Palette) *gif.GIF {
	cx, cy := size/2, size/2
	sim := make([]particle, len(base))
	copy(sim, base)
	alive := make([]bool, len(base))
	for i := range alive {
		alive[i] = true
	}
	project := func(p v3) (int, int) {
		return cx + int(float64(p.x)*scale), cy - int(float64(p.y)*scale)
	}

	g := &gif.GIF{}
	for range frames {
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
		img := drawParticleFrame(at, size, project, pal)
		g.Image = append(g.Image, img)
		g.Delay = append(g.Delay, 8)
		removeCollisions(sim, alive)
	}
	return g
}

func drawParticleFrame(at map[v3][]int, size int, project func(v3) (int, int), pal color.Palette) *image.Paletted {
	img := image.NewPaletted(image.Rect(0, 0, size, size), pal)
	for pos, group := range at {
		x, y := project(pos)
		idx := uint8(1)
		if len(group) > 1 {
			idx = 3
		}
		plotParticle(img, x, y, size, idx)
	}
	return img
}

func plotParticle(img *image.Paletted, x, y, size int, idx uint8) {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			px, py := x+dx, y+dy
			if px >= 0 && px < size && py >= 0 && py < size {
				img.SetColorIndex(px, py, idx)
			}
		}
	}
}
