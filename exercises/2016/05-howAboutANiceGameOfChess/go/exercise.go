package exercises

import (
	"crypto/md5"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 5.
type Exercise struct {
	common.BaseExercise
}

const hexDigits = "0123456789abcdef"

// interesting reports whether h has five leading zero hex digits, returning the
// 6th and 7th hex nibbles (used as char / position depending on the part).
func interesting(h [16]byte) (byte, byte, bool) {
	if h[0] == 0 && h[1] == 0 && h[2]>>4 == 0 {
		return h[2] & 0x0f, h[3] >> 4, true
	}
	return 0, 0, false
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	door := strings.TrimSpace(instr)
	var pw strings.Builder
	for i := 0; pw.Len() < 8; i++ {
		h := md5.Sum([]byte(door + strconv.Itoa(i)))
		if sixth, _, ok := interesting(h); ok {
			pw.WriteByte(hexDigits[sixth])
		}
	}
	return pw.String(), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	door := strings.TrimSpace(instr)
	pw := []byte("________")
	filled := 0
	for i := 0; filled < 8; i++ {
		h := md5.Sum([]byte(door + strconv.Itoa(i)))
		sixth, seventh, ok := interesting(h)
		if !ok || sixth > 7 || pw[sixth] != '_' {
			continue
		}
		pw[sixth] = hexDigits[seventh]
		filled++
	}
	return string(pw), nil
}
