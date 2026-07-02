package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 25.
type Exercise struct {
	common.BaseExercise
}

const (
	modulus = 20201227
	subject = 7
)

func parse(instr string) (int, int, error) {
	lines := strings.Fields(strings.TrimSpace(instr))
	if len(lines) != 2 {
		return 0, 0, fmt.Errorf("expected two public keys")
	}
	a, err := strconv.Atoi(lines[0])
	if err != nil {
		return 0, 0, err
	}
	b, err := strconv.Atoi(lines[1])
	if err != nil {
		return 0, 0, err
	}
	return a, b, nil
}

// findLoop brute-forces the loop size: how many times 7 is transformed to reach
// the given public key.
func findLoop(pub int) int {
	value := 1
	loop := 0
	for value != pub {
		value = value * subject % modulus
		loop++
	}
	return loop
}

// modPow computes base^exp mod modulus by fast exponentiation.
func modPow(base, exp int) int {
	result := 1
	base %= modulus
	for exp > 0 {
		if exp&1 == 1 {
			result = result * base % modulus
		}
		base = base * base % modulus
		exp >>= 1
	}
	return result
}

// One recovers one card's loop size and applies it to the other public key to
// derive the shared encryption key.
func (e Exercise) One(instr string) (any, error) {
	card, door, err := parse(instr)
	if err != nil {
		return nil, err
	}

	cardLoop := findLoop(card)
	return fmt.Sprintf("%d", modPow(door, cardLoop)), nil
}

// Two is the free star for completing all other days; there is no puzzle here.
func (e Exercise) Two(_ string) (any, error) {
	return "Merry Christmas!", nil
}
