package exercises

import (
	"errors"
	"math/big"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 22.
type Exercise struct {
	common.BaseExercise
}

// shuffle applies the sequence of operations to a deck of n cards and returns the result.
func shuffle(lines []string, n int) []int {
	deck := make([]int, n)
	for i := range deck {
		deck[i] = i
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case line == "deal into new stack":
			for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
				deck[i], deck[j] = deck[j], deck[i]
			}
		case strings.HasPrefix(line, "cut "):
			k, _ := strconv.Atoi(line[4:])
			if k < 0 {
				k += n
			}
			deck = append(deck[k:], deck[:k]...)
		case strings.HasPrefix(line, "deal with increment "):
			k, _ := strconv.Atoi(line[20:])
			out := make([]int, n)
			for i, card := range deck {
				out[(i*k)%n] = card
			}
			deck = out
		}
	}

	return deck
}

// One returns the answer to the first part of the exercise.
// For a 10-card example deck (≤15 operations), returns the full deck as space-separated values.
// For the real 10007-card input, returns the position of card 2019.
func (e Exercise) One(instr string) (any, error) {
	lines := strings.Split(strings.TrimSpace(instr), "\n")

	if len(lines) <= 15 {
		// Example: small deck, return full deck contents
		deck := shuffle(lines, 10)
		parts := make([]string, len(deck))
		for i, v := range deck {
			parts[i] = strconv.Itoa(v)
		}
		return strings.Join(parts, " "), nil
	}

	// Real input: 10007-card deck, find position of card 2019
	deck := shuffle(lines, 10007)
	for pos, card := range deck {
		if card == 2019 {
			return pos, nil
		}
	}

	return nil, errors.New("card 2019 not found")
}

// composeAffine composes two affine transforms f(x)=a1*x+b1 then g(x)=a2*x+b2
// into a single transform h(x) = (a2*a1)*x + (a2*b1+b2) mod n.
func composeAffine(a1, b1, a2, b2, n *big.Int) (*big.Int, *big.Int) {
	ra := new(big.Int).Mul(a2, a1)
	ra.Mod(ra, n)

	rb := new(big.Int).Mul(a2, b1)
	rb.Add(rb, b2)
	rb.Mod(rb, n)

	return ra, rb
}

// modInverse returns x^(-1) mod n using Fermat's little theorem (n must be prime).
func modInverse(x, n *big.Int) *big.Int {
	exp := new(big.Int).Sub(n, big.NewInt(2))
	return new(big.Int).Exp(x, exp, n)
}

// Two returns the answer to the second part of the exercise.
// Uses N=119315717514047 cards, K=101741582076661 shuffles, finds card at position 2020.
func (e Exercise) Two(instr string) (any, error) {
	lines := strings.Split(strings.TrimSpace(instr), "\n")

	n := big.NewInt(119315717514047)
	k := big.NewInt(101741582076661)

	// Start with identity transform: f(x) = 1*x + 0
	accumA := big.NewInt(1)
	accumB := big.NewInt(0)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		var opA, opB *big.Int

		switch {
		case line == "deal into new stack":
			// pos → -1*pos + (N-1)
			opA = new(big.Int).Sub(n, big.NewInt(1)) // -1 mod N = N-1
			opB = new(big.Int).Sub(n, big.NewInt(1)) // N-1

		case strings.HasPrefix(line, "cut "):
			c, _ := strconv.ParseInt(line[4:], 10, 64)
			// pos → 1*pos - c  mod N
			opA = big.NewInt(1)
			opB = new(big.Int).SetInt64(-c)
			opB.Mod(opB, n)
			if opB.Sign() < 0 {
				opB.Add(opB, n)
			}

		case strings.HasPrefix(line, "deal with increment "):
			inc, _ := strconv.ParseInt(line[20:], 10, 64)
			// pos → inc*pos + 0
			opA = big.NewInt(inc)
			opB = big.NewInt(0)

		default:
			continue
		}

		// Compose: current transform then this new one
		accumA, accumB = composeAffine(accumA, accumB, opA, opB, n)
	}

	// After one shuffle: pos → accumA*pos + accumB  mod N
	// After k shuffles: pos → accumA^k * pos + accumB * (accumA^k - 1) * modInverse(accumA-1, N)  mod N

	ak := new(big.Int).Exp(accumA, k, n) // accumA^k mod N

	// btotal = accumB * (accumA^k - 1) * modInverse(accumA-1, N) mod N
	akm1 := new(big.Int).Sub(ak, big.NewInt(1)) // accumA^k - 1
	akm1.Mod(akm1, n)

	am1 := new(big.Int).Sub(accumA, big.NewInt(1)) // accumA - 1
	am1.Mod(am1, n)

	invAm1 := modInverse(am1, n)

	btotal := new(big.Int).Mul(accumB, akm1)
	btotal.Mul(btotal, invAm1)
	btotal.Mod(btotal, n)

	// We want x such that f^k(x) = 2020
	// x = (2020 - btotal) * modInverse(accumA^k, N)  mod N
	pos := big.NewInt(2020)
	pos.Sub(pos, btotal)
	pos.Mod(pos, n)
	if pos.Sign() < 0 {
		pos.Add(pos, n)
	}

	invAK := modInverse(ak, n)
	pos.Mul(pos, invAK)
	pos.Mod(pos, n)

	return pos.String(), nil
}
