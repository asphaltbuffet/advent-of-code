package exercises

import (
	"errors"
	"fmt"
	"os"
)

var (
	root = Tile{Coord: Point{500, 0}, Type: Source}

	// ErrVoidPath is returned when a path will go outside rock boundaries.
	ErrVoidPath = errors.New("path is void")
)

// AddRocks adds the rocks to the graph between the given points on each row.
func (d *Day14) AddRocks(points [][]Point) error {
	d.MinX = points[0][0].X
	d.MaxX = points[0][0].X
	d.MaxY = points[0][0].Y

	for _, p := range points {
		for i := 1; i < len(p); i++ {
			// update the size of the graph
			if p[i].X > d.MaxX {
				d.MaxX = p[i].X
			}

			if p[i].Y > d.MaxY {
				d.MaxY = p[i].Y
			}

			if p[i].X < d.MinX {
				d.MinX = p[i].X
			}

			rocks, err := GetPointsBetween(p[i-1], p[i])
			if err != nil {
				return fmt.Errorf("getting points between %+v and %+v: %w", p[i-1], p[i], err)
			}

			for _, r := range rocks {
				d.Tiles[r] = Tile{Coord: r, Type: Rock}
			}
		}
	}

	return nil
}

// BuildGraph generates the possible paths of sand falling through environment.
func (d *Day14) BuildGraph(src Point) error {
	if _, found := d.Tiles[src]; found {
		// we've already been here, go back.
		return nil
	}

	if src.Y > d.MaxY { // may want to optimize by checking X as well
		return ErrVoidPath
	}

	// go down first
	down := Point{src.X, src.Y + 1}

	if err := d.BuildGraph(down); err != nil {
		switch err {
		case ErrVoidPath:
			return ErrVoidPath
		default:
			return fmt.Errorf("calculating path down: %w", err) // wrapping this may be real ugly with no additional information
		}
	}

	// go left
	left := Point{src.X - 1, src.Y + 1}

	if err := d.BuildGraph(left); err != nil {
		switch err {
		case ErrVoidPath:
			return ErrVoidPath
		default:
			return fmt.Errorf("calculating path left: %w", err)
		}
	}

	// go right
	right := Point{src.X + 1, src.Y + 1}

	if err := d.BuildGraph(right); err != nil {
		switch err {
		case ErrVoidPath:
			return ErrVoidPath
		default:
			return fmt.Errorf("calculating path down: %w", err)
		}
	}
	// all paths are down are valid/full; add this tile to the graph.
	d.Tiles[src] = Tile{Coord: src, Type: Sand}

	return nil
}

// RenderGraph renders the graph in Graphviz format to a file.
func (d *Day14) RenderGraph() {
	f, err := os.Create("graph.txt")
	if err != nil {
		panic(err)
	}

	defer func() {
		if defErr := f.Close(); defErr != nil {
			panic(defErr)
		}
	}()

	for j := 0; j <= d.MaxY; j++ {
		for i := d.MinX; i <= d.MaxX; i++ {
			t := d.Tiles[Point{i, j}]
			switch t.Type {
			case Unknown:
				_, err = f.WriteString(string(Air) + " ")
				if err != nil {
					panic(err)
				}
			default:
				_, err = f.WriteString(string(t.Type) + " ")
				if err != nil {
					panic(err)
				}
			}
		}

		// fmt.Println()
		_, err = f.WriteString("\n")
		if err != nil {
			panic(err)
		}
	}
}

// BuildGraphWithFloor generates the possible paths of sand falling through environment with a floor.
func (d *Day14) BuildGraphWithFloor(src Point) error {
	if _, found := d.Tiles[src]; found {
		// we've already been here, go back.
		return nil
	}

	if src.Y == d.MaxY {
		// we've reached the bottom, add this tile to the graph.
		d.Tiles[src] = Tile{Coord: src, Type: Rock}

		if src.X > d.MaxX {
			d.MaxX = src.X
		}

		if src.X < d.MinX {
			d.MinX = src.X
		}

		return nil
	}

	// go down first
	down := Point{src.X, src.Y + 1}

	if err := d.BuildGraphWithFloor(down); err != nil {
		switch err {
		case ErrVoidPath:
			return ErrVoidPath
		default:
			return fmt.Errorf("calculating path down: %w", err) // wrapping this may be real ugly with no additional information
		}
	}

	// go left
	left := Point{src.X - 1, src.Y + 1}

	if err := d.BuildGraphWithFloor(left); err != nil {
		switch err {
		case ErrVoidPath:
			return ErrVoidPath
		default:
			return fmt.Errorf("calculating path left: %w", err)
		}
	}

	// go right
	right := Point{src.X + 1, src.Y + 1}

	if err := d.BuildGraphWithFloor(right); err != nil {
		switch err {
		case ErrVoidPath:
			return ErrVoidPath
		default:
			return fmt.Errorf("calculating path down: %w", err)
		}
	}
	// all paths are down are valid/full; add this tile to the graph.
	d.Tiles[src] = Tile{Coord: src, Type: Sand}

	return nil
}
