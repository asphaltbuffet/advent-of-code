package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 14.
type Exercise struct {
	common.BaseExercise
}

// op is one program line: either a new mask, or a memory write.
type op struct {
	mask     string // set when this line is "mask = ..."
	addr     uint64
	value    uint64
	isMemory bool
}

func parse(instr string) ([]op, error) {
	var ops []op
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "mask") {
			ops = append(ops, op{mask: strings.TrimPrefix(line, "mask = ")})
			continue
		}
		// mem[ADDR] = VALUE
		inner := line[strings.IndexByte(line, '[')+1 : strings.IndexByte(line, ']')]
		addr, err := strconv.ParseUint(inner, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing address in %q: %w", line, err)
		}
		val, err := strconv.ParseUint(strings.TrimSpace(line[strings.IndexByte(line, '=')+1:]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing value in %q: %w", line, err)
		}
		ops = append(ops, op{addr: addr, value: val, isMemory: true})
	}
	return ops, nil
}

// One applies the value mask (0/1 overwrite bits, X passes through) and sums the
// resulting memory.
func (e Exercise) One(instr string) (any, error) {
	ops, err := parse(instr)
	if err != nil {
		return nil, err
	}

	var orMask, andMask uint64
	mem := map[uint64]uint64{}

	for _, o := range ops {
		if !o.isMemory {
			// Build OR (force 1s) and AND (force 0s) masks from the string.
			orMask, andMask = 0, 0
			for _, c := range o.mask {
				orMask <<= 1
				andMask <<= 1
				switch c {
				case '1':
					orMask |= 1
					andMask |= 1
				case '0':
					// leave both bits 0: OR keeps value bit, AND clears it
				default: // 'X'
					andMask |= 1 // pass through
				}
			}
			continue
		}
		mem[o.addr] = (o.value | orMask) & andMask
	}

	var sum uint64
	for _, v := range mem {
		sum += v
	}
	return fmt.Sprintf("%d", sum), nil
}

// addresses expands a base address under a v2 mask into every concrete address:
// mask '1' forces a bit on, '0' leaves it, 'X' is floating (both values).
func addresses(mask string, base uint64) []uint64 {
	addrs := []uint64{0}
	n := len(mask)
	for i, c := range mask {
		bit := uint64(1) << (n - 1 - i)
		switch c {
		case '0':
			// keep the base's bit
			keep := base & bit
			for j := range addrs {
				addrs[j] |= keep
			}
		case '1':
			for j := range addrs {
				addrs[j] |= bit
			}
		default: // 'X' floating: fork every address into 0 and 1
			doubled := make([]uint64, 0, len(addrs)*2)
			for _, a := range addrs {
				doubled = append(doubled, a, a|bit)
			}
			addrs = doubled
		}
	}
	return addrs
}

// Two decodes each write's address through the mask (floating bits fan out to
// many addresses) and sums the resulting memory.
func (e Exercise) Two(instr string) (any, error) {
	ops, err := parse(instr)
	if err != nil {
		return nil, err
	}

	var mask string
	mem := map[uint64]uint64{}
	for _, o := range ops {
		if !o.isMemory {
			mask = o.mask
			continue
		}
		for _, a := range addresses(mask, o.addr) {
			mem[a] = o.value
		}
	}

	var sum uint64
	for _, v := range mem {
		sum += v
	}
	return fmt.Sprintf("%d", sum), nil
}
