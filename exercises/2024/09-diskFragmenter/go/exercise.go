package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 9.
type Exercise struct {
	common.BaseExercise
}

// expand turns the dense disk map into a block slice where each entry is a file
// ID or -1 for free space.
func expand(instr string) []int {
	s := strings.TrimSpace(instr)
	var blocks []int
	id := 0
	for i := 0; i < len(s); i++ {
		n := int(s[i] - '0')
		val := -1
		if i%2 == 0 {
			val = id
			id++
		}
		for j := 0; j < n; j++ {
			blocks = append(blocks, val)
		}
	}
	return blocks
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	blocks := expand(instr)
	l, r := 0, len(blocks)-1
	for l < r {
		switch {
		case blocks[l] != -1:
			l++
		case blocks[r] == -1:
			r--
		default:
			blocks[l], blocks[r] = blocks[r], -1
			l++
			r--
		}
	}

	checksum := 0
	for i, v := range blocks {
		if v != -1 {
			checksum += i * v
		}
	}
	return checksum, nil
}

// file is a contiguous run on disk: file ID, start block, and length.
type file struct {
	id, start, length int
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	s := strings.TrimSpace(instr)

	var files []file
	type span struct{ start, length int }
	var frees []span

	pos := 0
	id := 0
	for i := 0; i < len(s); i++ {
		n := int(s[i] - '0')
		if i%2 == 0 {
			files = append(files, file{id, pos, n})
			id++
		} else if n > 0 {
			frees = append(frees, span{pos, n})
		}
		pos += n
	}

	// Move whole files in order of decreasing ID into the leftmost fitting gap.
	for i := len(files) - 1; i >= 0; i-- {
		f := &files[i]
		for j := range frees {
			g := &frees[j]
			if g.start >= f.start {
				break // no gap to the left is large enough
			}
			if g.length >= f.length {
				f.start = g.start
				g.start += f.length
				g.length -= f.length
				break
			}
		}
	}

	checksum := 0
	for _, f := range files {
		for k := 0; k < f.length; k++ {
			checksum += (f.start + k) * f.id
		}
	}
	return checksum, nil
}
