package exercises

import (
	"image"

	"github.com/asphaltbuffet/advent-of-code/pkg/set"
)

type block struct {
	up             *adjBlock
	right          *adjBlock
	left           *adjBlock
	down           *adjBlock
	dims           image.Rectangle
	adjascentZones set.Set[image.Rectangle]
}

// add b as a neighbor to a in direction dir
// and vice-versa with the inverse direction
func (bl *block) addAdjascent(bl2 *block, fromDir, toDir direction) (added bool) {
	if bl.adjascentZones.Has(bl2.dims) {
		return false
	}

	inverseToDir := toDir.invert()
	inverseFromDir := fromDir.invert()

	if bl.hasAdjascent(fromDir) || bl2.hasAdjascent(inverseToDir) {
		return false
	}

	bl.setAdjascent(bl2, fromDir, toDir)
	bl2.setAdjascent(bl, inverseToDir, inverseFromDir)

	return true
}

func (bl *block) setAdjascent(bl2 *block, fromDir, toDir direction) {
	switch fromDir {
	case upD:
		bl.up = &adjBlock{bl2, toDir}
	case rightD:
		bl.right = &adjBlock{bl2, toDir}
	case leftD:
		bl.left = &adjBlock{bl2, toDir}
	case downD:
		bl.down = &adjBlock{bl2, toDir}
	}

	bl.adjascentZones.Add(bl2.dims)
}

func (bl *block) hasAdjascent(dir direction) bool {
	switch dir {
	case upD:
		return bl.up != nil
	case rightD:
		return bl.right != nil
	case leftD:
		return bl.left != nil
	case downD:
		return bl.down != nil
	}

	return false
}
