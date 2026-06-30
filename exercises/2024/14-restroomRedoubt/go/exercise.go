package exercises

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 14.
type Exercise struct {
	common.BaseExercise
}

type robot struct {
	x, y, vx, vy int
}

var numRe = regexp.MustCompile(`-?\d+`)

// parseRobots scans every integer and groups them in fours (px, py, vx, vy).
func parseRobots(instr string) []robot {
	nums := numRe.FindAllString(instr, -1)
	var robots []robot
	for i := 0; i+3 < len(nums); i += 4 {
		v := make([]int, 4)
		for j := 0; j < 4; j++ {
			v[j], _ = strconv.Atoi(nums[i+j])
		}
		robots = append(robots, robot{v[0], v[1], v[2], v[3]})
	}
	return robots
}

// gridSize returns the grid dimensions: 11x7 for the example (all robots fit),
// otherwise the real 101x103 lab.
func gridSize(robots []robot) (w, h int) {
	w, h = 11, 7
	for _, r := range robots {
		if r.x >= 11 || r.y >= 7 {
			return 101, 103
		}
	}
	return w, h
}

func mod(a, m int) int {
	a %= m
	if a < 0 {
		a += m
	}
	return a
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	robots := parseRobots(instr)
	w, h := gridSize(robots)

	var q [4]int
	for _, r := range robots {
		x := mod(r.x+r.vx*100, w)
		y := mod(r.y+r.vy*100, h)
		switch {
		case x == w/2 || y == h/2:
			// on a middle line: ignored
		case x < w/2 && y < h/2:
			q[0]++
		case x > w/2 && y < h/2:
			q[1]++
		case x < w/2 && y > h/2:
			q[2]++
		default:
			q[3]++
		}
	}
	return q[0] * q[1] * q[2] * q[3], nil
}

// findEasterEgg returns the first time the robots arrange into the Easter-egg
// picture, detected as the unique arrangement where every robot occupies a
// distinct cell (-1 if none within one full cycle).
func findEasterEgg(robots []robot, w, h int) int {
	for t := 0; t < w*h; t++ {
		seen := make(map[int]bool, len(robots))
		unique := true
		for _, r := range robots {
			x := mod(r.x+r.vx*t, w)
			y := mod(r.y+r.vy*t, h)
			key := y*w + x
			if seen[key] {
				unique = false
				break
			}
			seen[key] = true
		}
		if unique {
			return t
		}
	}
	return -1
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	robots := parseRobots(instr)
	w, h := gridSize(robots)
	return findEasterEgg(robots, w, h), nil
}

// Vis renders the robot positions at the Easter-egg time (Part Two answer).
func (e Exercise) Vis(instr string, outdir string) error {
	robots := parseRobots(instr)
	w, h := gridSize(robots)
	t := findEasterEgg(robots, w, h)
	if t < 0 {
		t = 0
	}

	const scale = 6
	const pad = 8
	img := image.NewRGBA(image.Rect(0, 0, w*scale+2*pad, h*scale+2*pad))

	bg := color.RGBA{0x06, 0x12, 0x0a, 0xff}    // deep green-black
	fg := color.RGBA{0x6c, 0xff, 0x8a, 0xff}    // bright phosphor green
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.SetRGBA(x, y, bg)
		}
	}

	for _, r := range robots {
		x := mod(r.x+r.vx*t, w)
		y := mod(r.y+r.vy*t, h)
		x0 := pad + x*scale
		y0 := pad + y*scale
		for dy := 0; dy < scale; dy++ {
			for dx := 0; dx < scale; dx++ {
				img.SetRGBA(x0+dx, y0+dy, fg)
			}
		}
	}

	f, err := os.Create(filepath.Join(outdir, "restroom-redoubt.png"))
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
