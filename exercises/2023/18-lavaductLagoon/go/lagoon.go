package exercises

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
	"strings"

	util "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

type Point util.Point2D

type Lagoon struct {
	Path  []Point
	Steps Steps
	Min   Point
	Max   Point
}

func (lg *Lagoon) Height() int {
	return lg.Max.Y - lg.Min.Y + 1
}

func (lg *Lagoon) Width() int {
	return lg.Max.X - lg.Min.X + 1
}

type Step struct {
	Value     string
	Direction Direction
	Distance  int
}

type Steps []*Step

// Parse the input into a list of steps.
func parseInput(input string) (Steps, Steps, error) {
	if input == "" {
		return nil, nil, fmt.Errorf("input is empty")
	}

	lines := strings.Split(input, "\n")
	steps := make(Steps, len(lines))
	hexSteps := make(Steps, len(lines))

	for i, line := range lines {
		var err error
		steps[i], hexSteps[i], err = parseLine(line)
		if err != nil {
			return nil, nil, err
		}
	}

	return steps, hexSteps, nil
}

func (lg *Lagoon) DebugPrint() {
	// print map specs
	fmt.Printf("Map dimensions: %dx%d\n", lg.Width(), lg.Height())
	fmt.Printf("Map bounds: (%d, %d) - (%d, %d)\n", lg.Min.X, lg.Min.Y, lg.Max.X, lg.Max.Y)
	fmt.Printf("Vertices:  %v\n", lg.Path)
}

type Direction rune

const (
	Right Direction = 'R'
	Left  Direction = 'L'
	Up    Direction = 'U'
	Down  Direction = 'D'
)

var Move = map[Direction]Point{
	Right: {1, 0},
	Left:  {-1, 0},
	Up:    {0, -1},
	Down:  {0, 1},
}

var HexMove = map[string]Direction{
	"0": Right, // right
	"1": Down,  // down
	"2": Left,  // left
	"3": Up,    // up
}

var MoveString = map[Point]string{
	{1, 0}:  "R",
	{-1, 0}: "L",
	{0, -1}: "U",
	{0, 1}:  "D",
}

func parseLine(line string) (*Step, *Step, error) {
	tokens := strings.Fields(line)
	if len(tokens) != 3 {
		return nil, nil, fmt.Errorf("invalid line: %s", line)
	}

	dir := Direction(tokens[0][0])
	dist, err := strconv.Atoi(tokens[1])
	if err != nil {
		return nil, nil, err
	}

	cToken := tokens[2]
	d, err := strconv.ParseInt(cToken[2:len(cToken)-2], 16, 0)
	if err != nil {
		return nil, nil, err
	}

	return &Step{
			Value:     line,
			Direction: dir,
			Distance:  dist,
		},
		&Step{
			Value:     cToken[1 : len(cToken)-1],
			Direction: HexMove[cToken[len(cToken)-2:len(cToken)-1]],
			Distance:  int(d),
		},
		nil
}

func (ss Steps) GetBoundaryPoints() ([]Point, error) {
	points := make([]Point, 0, len(ss))
	var lastDir Direction
	curPoint := Point{0, 0}

	points = append(points, curPoint)

	for _, s := range ss {
		switch {
		case lastDir == Up && s.Direction == Right, lastDir == Right && s.Direction == Up:
			points = append(points, curPoint)

		case lastDir == Up && s.Direction == Left, lastDir == Left && s.Direction == Up:
			points = append(points, curPoint.Add(Point{0, 1}))

		case lastDir == Down && s.Direction == Right, lastDir == Right && s.Direction == Down:
			points = append(points, curPoint.Add(Point{1, 0}))

		case lastDir == Down && s.Direction == Left, lastDir == Left && s.Direction == Down:
			points = append(points, curPoint.Add(Point{1, 1}))
		}

		switch s.Direction {
		case Up:
			curPoint.Y -= s.Distance

		case Down:
			curPoint.Y += s.Distance

		case Left:
			curPoint.X -= s.Distance

		case Right:
			curPoint.X += s.Distance

		default:
			return nil, fmt.Errorf("invalid direction: %c", s.Direction)
		}

		lastDir = s.Direction
	}

	if curPoint.X != 0 || curPoint.Y != 0 {
		return nil, fmt.Errorf("invalid path: %v", points)
	}

	// end with 0,0 to close the polygon
	points = append(points, curPoint)

	return points, nil
}

func (p Point) Move(dir Direction, dist int) Point {
	return p.Add(Move[dir].Mul(dist))
}

func (p Point) Mul(scalar int) Point {
	return Point{p.X * scalar, p.Y * scalar}
}

func (p Point) Add(p2 Point) Point {
	return Point{p.X + p2.X, p.Y + p2.Y}
}

func shoelace(points []Point) int {
	n := len(points)

	if n < 3 {
		return 0
	}

	x := make([]int, n)
	y := make([]int, n)

	for i, p := range points {
		x[i] = p.X
		y[i] = p.Y
	}

	var area int

	for i := 0; i < n-1; i++ {
		area += (y[i] + y[i+1]) * (x[i] - x[i+1])
	}

	return util.AbsInt(area / 2) //nolint:gomnd // ignore dividing by 2
}

//nolint:gomnd // ignore image sizes
func (lg *Lagoon) plotPoints(points []Point, filename string) error {
	const padding int = 5

	// Create a new image with a white background
	img := image.NewRGBA(image.Rect(lg.Min.X-padding, lg.Min.Y-padding, lg.Max.X+padding, lg.Max.Y+padding))
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)

	// Draw the polygon
	for i, j := len(points)-1, 0; j < len(points); i, j = (i+1)%len(points), j+1 {
		drawLine(img, points[i], points[j], color.Black)
	}

	// Save the image to a file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func drawLine(img *image.RGBA, p0, p1 Point, c color.Color) {
	x0, y0 := p0.X, p0.Y
	x1, y1 := p1.X, p1.Y

	// Bresenham's line algorithm: https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
	dx := util.AbsInt(x1 - x0)
	dy := util.AbsInt(y1 - y0)

	var sx, sy int

	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}

	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	for {
		img.Set(x0, y0, c)

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err

		if e2 > -dy {
			err -= dy
			x0 += sx
		}

		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}
