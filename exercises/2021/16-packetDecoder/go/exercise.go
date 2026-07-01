package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 16.
type Exercise struct {
	common.BaseExercise
}

type packet struct {
	version  int
	typeID   int
	value    int
	children []packet
}

// hexToBits expands a hex transmission into its bit string (MSB first).
func hexToBits(hex string) (string, error) {
	var sb strings.Builder
	for _, c := range strings.TrimSpace(hex) {
		var v int
		switch {
		case c >= '0' && c <= '9':
			v = int(c - '0')
		case c >= 'A' && c <= 'F':
			v = int(c-'A') + 10
		case c >= 'a' && c <= 'f':
			v = int(c-'a') + 10
		default:
			return "", fmt.Errorf("invalid hex digit %q", c)
		}
		fmt.Fprintf(&sb, "%04b", v)
	}
	return sb.String(), nil
}

// bitsToInt reads the substring bits[a:b] as a big-endian integer.
func bitsToInt(bits string) int {
	n := 0
	for _, c := range bits {
		n = n<<1 | int(c-'0')
	}
	return n
}

// parsePacket reads one packet starting at pos and returns it with the position
// immediately after it.
func parsePacket(bits string, pos int) (packet, int) {
	var p packet
	p.version = bitsToInt(bits[pos : pos+3])
	p.typeID = bitsToInt(bits[pos+3 : pos+6])
	pos += 6

	if p.typeID == 4 {
		// Literal: 5-bit groups, first bit signals another group follows.
		var val int
		for {
			group := bits[pos : pos+5]
			pos += 5
			val = val<<4 | bitsToInt(group[1:])
			if group[0] == '0' {
				break
			}
		}
		p.value = val
		return p, pos
	}

	// Operator: length type id decides how sub-packets are delimited.
	lengthTypeID := bits[pos]
	pos++
	if lengthTypeID == '0' {
		totalLen := bitsToInt(bits[pos : pos+15])
		pos += 15
		end := pos + totalLen
		for pos < end {
			var child packet
			child, pos = parsePacket(bits, pos)
			p.children = append(p.children, child)
		}
	} else {
		count := bitsToInt(bits[pos : pos+11])
		pos += 11
		for i := 0; i < count; i++ {
			var child packet
			child, pos = parsePacket(bits, pos)
			p.children = append(p.children, child)
		}
	}

	return p, pos
}

func (p packet) versionSum() int {
	sum := p.version
	for _, c := range p.children {
		sum += c.versionSum()
	}
	return sum
}

func (p packet) eval() int {
	switch p.typeID {
	case 0:
		sum := 0
		for _, c := range p.children {
			sum += c.eval()
		}
		return sum
	case 1:
		prod := 1
		for _, c := range p.children {
			prod *= c.eval()
		}
		return prod
	case 2:
		m := p.children[0].eval()
		for _, c := range p.children[1:] {
			if v := c.eval(); v < m {
				m = v
			}
		}
		return m
	case 3:
		m := p.children[0].eval()
		for _, c := range p.children[1:] {
			if v := c.eval(); v > m {
				m = v
			}
		}
		return m
	case 4:
		return p.value
	case 5:
		return boolInt(p.children[0].eval() > p.children[1].eval())
	case 6:
		return boolInt(p.children[0].eval() < p.children[1].eval())
	case 7:
		return boolInt(p.children[0].eval() == p.children[1].eval())
	}
	return 0
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func decode(instr string) (packet, error) {
	bits, err := hexToBits(instr)
	if err != nil {
		return packet{}, fmt.Errorf("expanding hex: %w", err)
	}
	p, _ := parsePacket(bits, 0)
	return p, nil
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	p, err := decode(instr)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%d", p.versionSum()), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	p, err := decode(instr)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%d", p.eval()), nil
}
