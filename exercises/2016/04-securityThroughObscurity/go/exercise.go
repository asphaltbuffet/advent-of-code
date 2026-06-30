package exercises

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 4.
type Exercise struct {
	common.BaseExercise
}

var roomRe = regexp.MustCompile(`([a-z-]+)-(\d+)\[([a-z]+)\]`)

type room struct {
	name     string // includes dashes
	sector   int
	checksum string
}

func parseRooms(instr string) []room {
	var rooms []room
	for _, m := range roomRe.FindAllStringSubmatch(instr, -1) {
		sector, _ := strconv.Atoi(m[2])
		rooms = append(rooms, room{name: m[1], sector: sector, checksum: m[3]})
	}
	return rooms
}

// real reports whether the room's checksum equals the five most common letters
// in its name (ties broken alphabetically).
func (r room) real() bool {
	counts := map[rune]int{}
	for _, ch := range r.name {
		if ch != '-' {
			counts[ch]++
		}
	}
	letters := make([]rune, 0, len(counts))
	for ch := range counts {
		letters = append(letters, ch)
	}
	sort.Slice(letters, func(i, j int) bool {
		if counts[letters[i]] != counts[letters[j]] {
			return counts[letters[i]] > counts[letters[j]]
		}
		return letters[i] < letters[j]
	})
	top := letters
	if len(top) > 5 {
		top = top[:5]
	}
	return string(top) == r.checksum
}

// decrypt shifts each letter forward by the sector ID; dashes become spaces.
func (r room) decrypt() string {
	shift := rune(r.sector % 26)
	var b strings.Builder
	for _, ch := range r.name {
		if ch == '-' {
			b.WriteByte(' ')
		} else {
			b.WriteRune('a' + (ch-'a'+shift)%26)
		}
	}
	return b.String()
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	sum := 0
	for _, r := range parseRooms(instr) {
		if r.real() {
			sum += r.sector
		}
	}
	return sum, nil
}

// Two returns the answer to the second part of the exercise: the sector ID of
// the room whose decrypted name mentions north pole object storage.
func (e Exercise) Two(instr string) (any, error) {
	for _, r := range parseRooms(instr) {
		if r.real() && strings.Contains(r.decrypt(), "northpole") {
			return r.sector, nil
		}
	}
	return 0, nil
}
