package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Vis plots the probe trajectories as an SVG. Every launch velocity that lands in
// the target is drawn as a faint arc, so the envelope of all 1566 winning shots
// (part two) is visible; the single arc reaching the greatest height (part one)
// is highlighted, and the target area is a labeled box. y is flipped so "up" is
// up. The highlight and target read by position and brightness, so the plot works
// in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	t, err := parse(instr)
	if err != nil {
		return err
	}

	type shot struct {
		path [][2]int
		peak int
	}
	var winners []shot
	best := shot{peak: -1 << 30}

	for vx := 1; vx <= t.xMax; vx++ {
		for vy := t.yMin; vy <= -t.yMin; vy++ {
			path, peak, hit := trace(vx, vy, t)
			if !hit {
				continue
			}
			s := shot{path, peak}
			winners = append(winners, s)
			if peak > best.peak {
				best = s
			}
		}
	}

	// World bounds: from launch (0,0) out past the target, up to the best peak.
	minX, maxX := 0, t.xMax
	minY, maxY := t.yMin, best.peak
	worldW := float64(maxX - minX)
	worldH := float64(maxY - minY)

	const (
		W     = 1000
		H     = 620
		mLeft = 50
		mTop  = 40
		mBot  = 40
		mRt   = 20
	)
	plotW := float64(W - mLeft - mRt)
	plotH := float64(H - mTop - mBot)
	// Independent x/y scales: the peak height dwarfs the target distance, so a
	// uniform scale would squeeze everything into a sliver. Non-uniform scaling
	// fills the canvas; the shape of the arcs is still faithful, only stretched.
	scaleX := plotW / worldW
	scaleY := plotH / worldH
	sx := func(x int) float64 { return mLeft + (float64(x)-float64(minX))*scaleX }
	sy := func(y int) float64 { return mTop + (float64(maxY)-float64(y))*scaleY }

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="26" fill="#e8ecf4" font-size="15" text-anchor="middle">Probe trajectories that hit the target (%d of them); highest arc highlighted</text>`, W/2, len(winners))

	// Target box.
	tx0, ty0 := sx(t.xMin), sy(t.yMax)
	tw := float64(t.xMax-t.xMin+1) * scaleX
	th := float64(t.yMax-t.yMin+1) * scaleY
	fmt.Fprintf(&sb, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#009E73" fill-opacity="0.25" stroke="#009E73" stroke-width="1.5"/>`, tx0, ty0, tw, th)
	fmt.Fprintf(&sb, `<text x="%.1f" y="%.1f" fill="#3fe0b0" font-size="12" text-anchor="middle">target</text>`, tx0+tw/2, ty0+th+16)

	// Zero line (launch height) for reference.
	fmt.Fprintf(&sb, `<line x1="%.1f" y1="%.1f" x2="%d" y2="%.1f" stroke="#2a2f3a" stroke-dasharray="3 3"/>`, sx(minX), sy(0), W-mRt, sy(0))

	// Faint winning arcs.
	fmt.Fprint(&sb, `<g fill="none" stroke="#56B4E9" stroke-opacity="0.12" stroke-width="1">`)
	for _, s := range winners {
		fmt.Fprintf(&sb, `<polyline points="%s"/>`, polyPoints(s.path, sx, sy))
	}
	fmt.Fprint(&sb, `</g>`)

	// Highlighted highest arc.
	fmt.Fprintf(&sb, `<polyline points="%s" fill="none" stroke="#E69F00" stroke-width="2.5"/>`, polyPoints(best.path, sx, sy))
	// Peak marker.
	px, py := sx(best.path[peakIndex(best.path)][0]), sy(best.peak)
	fmt.Fprintf(&sb, `<circle cx="%.1f" cy="%.1f" r="4" fill="#E69F00"/>`, px, py)
	fmt.Fprintf(&sb, `<text x="%.1f" y="%.1f" fill="#E69F00" font-size="13" text-anchor="middle">peak y = %d</text>`, px, py-10, best.peak)

	// Launch point.
	fmt.Fprintf(&sb, `<circle cx="%.1f" cy="%.1f" r="4" fill="#D55E00"/>`, sx(0), sy(0))
	fmt.Fprintf(&sb, `<text x="%.1f" y="%.1f" fill="#D55E00" font-size="12">launch</text>`, sx(0)+8, sy(0)-6)

	sb.WriteString(`</svg>`)

	return os.WriteFile(filepath.Join(outdir, "trick-shot.svg"), []byte(sb.String()), 0o600)
}

// trace simulates a shot and returns its full path, peak height, and hit status.
func trace(vx, vy int, t target) (path [][2]int, peak int, hit bool) {
	x, y := 0, 0
	path = append(path, [2]int{0, 0})
	peak = 0
	for {
		x += vx
		y += vy
		if vx > 0 {
			vx--
		} else if vx < 0 {
			vx++
		}
		vy--
		if y > peak {
			peak = y
		}
		path = append(path, [2]int{x, y})

		if x >= t.xMin && x <= t.xMax && y >= t.yMin && y <= t.yMax {
			return path, peak, true
		}
		if x > t.xMax || y < t.yMin {
			return path, peak, false
		}
	}
}

func polyPoints(path [][2]int, sx, sy func(int) float64) string {
	var b strings.Builder
	for _, p := range path {
		fmt.Fprintf(&b, "%.1f,%.1f ", sx(p[0]), sy(p[1]))
	}
	return strings.TrimSpace(b.String())
}

func peakIndex(path [][2]int) int {
	best, bi := path[0][1], 0
	for i, p := range path {
		if p[1] > best {
			best, bi = p[1], i
		}
	}
	return bi
}
