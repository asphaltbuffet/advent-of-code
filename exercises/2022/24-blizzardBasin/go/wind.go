package exercises

type wind struct {
	start     point
	direction point
	totalRows int
	totalCols int
	char      string
}

func (b wind) extrapolatePosition(steps int) point {
	p := point{
		x: (b.start.x + b.direction.x*steps) % b.totalCols,
		y: (b.start.y + b.direction.y*steps) % b.totalRows,
	}

	p.x += b.totalCols
	p.x %= b.totalCols

	p.y += b.totalRows
	p.y %= b.totalRows

	return p
}
