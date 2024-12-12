package alphametics

import (
	"errors"
	"fmt"
	"math"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

func Solve(puzzle string) (map[string]int, error) {
	input := strings.Trim(puzzle, "+=")
	letters := uniqueLetters(input)

	if letters.Cardinality() > 10 {
		return nil, errors.New("too many letters")
	}

	letterVal := letters.ToSlice()
	fmt.Println("Letters: ", letterVal)

	c := getCombinations(string(letterVal)).ToSlice()

	for _, p := range c {
		cipher := map[string]int{}

		for i, l := range p {
			if l != ' ' {
				cipher[string(l)] = i
			}
		}

		if calc(cipher, input) {
			return cipher, nil
		}
	}

	return nil, errors.New("no solution found")
}

func calc(c map[string]int, s string) bool {
	// split into words and separate sum
	words := strings.Fields(s)

	// toss out any combination that a word starts with 0
	for _, w := range words {
		if c[w[:1]] == 0 {
			return false
		}
	}

	sum := 0

	for _, w := range words[:len(words)-1] {
		sum += cipherToInt(c, w)
	}

	last := cipherToInt(c, words[len(words)-1])

	return sum == last
}

func cipherToInt(c map[string]int, s string) int {
	sub := 0
	pos := len(s) - 1

	for _, l := range s {
		sub += c[string(l)] * int(math.Pow10(pos))
		pos--
	}

	return sub
}

func uniqueLetters(s string) mapset.Set[rune] {
	all := mapset.NewSet[rune]()

	for _, c := range strings.ToUpper(s) {
		if c >= 'A' && c <= 'Z' {
			all.Add(c)
		}
	}

	return all
}

func getCombinations(s string) mapset.Set[string] {
	all := mapset.NewSet[string]()

	Perm([]rune(fmt.Sprintf("%10s", s)), func(a []rune) {
		all.Add(string(a))
	})

	return all
}

func Perm(a []rune, f func([]rune)) {
	perm(a, f, 0)
}

func perm(a []rune, f func([]rune), i int) {
	if i > len(a) {
		f(a)
		return
	}

	perm(a, f, i+1)

	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		perm(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
}
