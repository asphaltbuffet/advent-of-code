package exercises

type point struct {
	x int
	y int
}

func (p point) add(p2 point) point {
	return point{
		x: p.x + p2.x,
		y: p.y + p2.y,
	}
}

func (p *point) normalize(maxX, maxY int) {
	p.x = (p.x + maxX) % maxX
	p.y = (p.y + maxY) % maxY
}
