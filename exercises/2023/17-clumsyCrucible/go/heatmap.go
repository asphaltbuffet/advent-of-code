//nolint:gomnd // ignore magic numbers
package exercises

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strconv"

	"github.com/kettek/apng"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var p = color.Palette{
	color.RGBA{0x80, 0x80, 0x80, 0xff}, // 0 (unknown - grey)
	color.RGBA{0xff, 0x33, 0x33, 0xff}, // 9 (hot - red)
	color.RGBA{0xfd, 0x72, 0x43, 0xff}, // 8
	color.RGBA{0xfa, 0xb9, 0x54, 0xff}, // 7
	color.RGBA{0xf9, 0xe8, 0x60, 0xff}, // 6
	color.RGBA{0xf8, 0xfc, 0x65, 0xff}, // 5 (yellow)
	color.RGBA{0xe4, 0xf2, 0x78, 0xff}, // 4
	color.RGBA{0xb6, 0xdb, 0xa2, 0xff}, // 3
	color.RGBA{0x77, 0xbb, 0xda, 0xff}, // 2
	color.RGBA{0x4f, 0xa7, 0xff, 0xff}, // 1 (cold - blue)
}

const squareSize = 25

func (c *City) GenerateImage(path []Pather, outfile string) error {
	img := generateBackgroundImage(c, path)

	f, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = png.Encode(f, img); err != nil {
		return err
	}

	return nil
}

func generateBackgroundImage(c *City, path []Pather) *image.RGBA {
	// add a row to the bottom for debug info
	img := image.NewRGBA(image.Rect(0, 0, c.Width*squareSize, (c.Height+1)*squareSize))

	for y := 0; y < c.Height+1; y++ {
		for x := 0; x < c.Width; x++ {
			if b, ok := c.Blocks[Point{x, y}]; ok {
				drawSquare(img, Point{x, y}, p[b.HeatLoss])
			} else {
				drawSquare(img, Point{x, y}, color.White)
			}
		}
	}

	// mark the path
	for i := 0; i < len(path)-1; i++ {
		drawLine(img, path[i].(*Block).Position, path[i+1].(*Block).Position, color.RGBA{0, 255, 0, 255})
	}

	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			if b, ok := c.Blocks[Point{x, y}]; ok {
				drawLabels(img, Point{x, y}, b.HeatLoss)
			}
		}
	}

	return img
}

func drawSquare(img *image.RGBA, p Point, c color.Color) {
	maxX, maxY := p.X*squareSize+squareSize, p.Y*squareSize+squareSize

	for x := p.X * squareSize; x < maxX; x++ {
		for y := p.Y * squareSize; y < maxY; y++ {
			img.Set(x, y, c)
		}
	}
}

func drawLabels(img *image.RGBA, p Point, val int) {
	textX := p.X*squareSize + squareSize/2 - 3
	textY := p.Y*squareSize + squareSize/2 + 5

	addLabel(img, textX, textY, color.Black, strconv.Itoa(val))
}

func drawLine(img *image.RGBA, p1, p2 Point, c color.RGBA) {
	// draw a line between the two points
	switch {
	case p1.X == p2.X:
		if p1.Y > p2.Y {
			p1, p2 = p2, p1
		}

		// vertical line
		x := p1.X*squareSize + squareSize/2

		for y := p1.Y*squareSize + squareSize/2; y < p2.Y*squareSize+squareSize/2; y++ {
			img.Set(x-2, y, c)
			img.Set(x-1, y, c)
			img.Set(x, y, c)
			img.Set(x+1, y, c)
			img.Set(x+2, y, c)
		}

	case p1.Y == p2.Y:
		if p1.X > p2.X {
			p1, p2 = p2, p1
		}

		// horizontal line
		y := p1.Y*squareSize + squareSize/2

		for x := p1.X*squareSize + squareSize/2; x < p2.X*squareSize+squareSize/2; x++ {
			img.Set(x, y-2, c)
			img.Set(x, y-1, c)
			img.Set(x, y, c)
			img.Set(x, y+1, c)
			img.Set(x, y+2, c)
		}
	}
}

func addLabel(img *image.RGBA, x, y int, c color.Color, label string) {
	pt := fixed.P(x, y)

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(c),
		Face: basicfont.Face7x13,
		Dot:  pt,
	}

	d.DrawString(label)
}

func (c *City) GenerateFrame(path []Pather, s string) (*image.RGBA, error) {
	// img := generateBackgroundImage(c, path)
	img := image.NewRGBA(image.Rect(0, 0, c.Width*squareSize, (c.Height+1)*squareSize))

	// add debug info
	addLabel(img, 5, c.Height*squareSize+15, color.Black, s)

	// mark the path
	for i := 0; i < len(path)-1; i++ {
		drawLine(img, path[i].(*Block).Position, path[i+1].(*Block).Position, color.RGBA{0, 255, 0, 255})
	}

	for k, block := range c.Blocks {
		midX, midY := k.X*squareSize+squareSize/2, k.Y*squareSize+squareSize/2

		node, ok := nm[block]
		if ok && node.closed {
			// circle closed blocks in black
			drawCircle(img, midX, midY, squareSize/2, color.RGBA{0, 0, 0, 255})
		}

		if ok && node.open {
			// circle estimated blocks in white
			drawCircle(img, midX, midY, squareSize/2, color.RGBA{255, 255, 255, 255})
		}

		// circle current block in pink
		if block == path[0].(*Block) {
			drawCircle(img, midX, midY, squareSize/2+1, color.RGBA{255, 0, 255, 255})
		}
	}

	return img, nil
}

func drawCircle(img *image.RGBA, cX, cY, radius int, c color.RGBA) {
	// line thickness for circle drawn
	thickness := 1

	for x := cX - radius; x < cX+squareSize; x++ {
		for y := cY - radius; y < cY+squareSize; y++ {
			dist := math.Sqrt(float64((x-cX)*(x-cX) + (y-cY)*(y-cY)))
			if dist > float64(radius-thickness) && dist < float64(radius+thickness) {
				img.Set(x, y, c)
			}
		}
	}
}

func GenerateAPNG(frames []*image.RGBA, output string) error {
	if len(frames) == 0 {
		return fmt.Errorf("no frames to generate animation")
	}

	a := apng.APNG{
		Frames: make([]apng.Frame, len(frames)),
	}

	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	for i, img := range frames {
		img := img
		a.Frames[i].Image = img
		a.Frames[i].DelayNumerator = 1
		a.Frames[i].DelayDenominator = 1
	}

	return apng.Encode(out, a)
}
