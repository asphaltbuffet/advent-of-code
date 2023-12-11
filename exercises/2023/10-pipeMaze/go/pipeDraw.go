package exercises

import (
	"image"
	"image/draw"
	"image/png"
	"os"
	"strings"
)

// loadImage loads an image from a file
func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func drawPipes(s string) {
	m, start, _ := parseInput(s)

	path, onpath, _ := findPath(m, start)
	path = path
	charHeight, charWidth := 50, 50
	lines := strings.Split(s, "\n")
	imgWidth := len(lines[0]) * charWidth
	imgHeight := len(lines) * charWidth

	// Create a new RGBA image with the calculated dimensions
	rgba := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	for y := 0; y < len(lines); y++ {
		for x := 0; x < len(lines[0]); x++ {
			tmp := (Point{x, y})
			var charImg image.Image
			var err error

			switch {
			case onpath[tmp]:
				charImg, err = loadImage(pipeFile(m[tmp]) + ".png") // Load the image for the character
			// case windingNumber(tmp, path):
			// 	charImg, err = loadImage("nest.png") // Load the image for the character
			default:
				charImg, err = loadImage("blank.png") // Load the image for the character
			}

			if err != nil {
				panic(err)
			}

			// Set the offset for where to draw the character image
			offset := image.Pt(x*charWidth, y*charHeight)

			// Draw the character image onto the final image
			draw.Draw(rgba, charImg.Bounds().Add(offset), charImg, image.Point{}, draw.Over)
		}
	}

	// Save the final image to a file
	outputFile, err := os.Create("finalImage.png")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, rgba)
	if err != nil {
		panic(err)
	}
}

func pipeFile(p *Pipe) string {
	switch p.Shape {
	case StartPipe:
		return "start"
	case HorizontalPipe:
		return "WE"
	case VerticalPipe:
		return "NS"
	case NECornerPipe:
		return "NE"
	case NWCornerPipe:
		return "NW"
	case SECornerPipe:
		return "SE"
	case SWCornerPipe:
		return "SW"
	default:
		return "blank"
	}
}
