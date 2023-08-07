package exercises

var relative = map[string]point{
	">": {1, 0},
	"v": {0, 1},
	"<": {-1, 0},
	"^": {0, -1},
}

type point struct {
	x, y int
}

func (p point) add(p2 point) point {
	return point{p.x + p2.x, p.y + p2.y}
}
