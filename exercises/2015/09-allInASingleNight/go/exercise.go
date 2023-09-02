package exercises

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
	"github.com/asphaltbuffet/advent-of-code/pkg/permutation"
)

// Exercise for Advent of Code 2015 day 9.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	distances, _ := parse(instr)
	paths := calcPaths(distances)

	min := math.MaxInt32
	for _, v := range paths {
		if v < min {
			min = v
		}
	}

	return min, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	distances, _ := parse(instr)
	paths := calcPaths(distances)

	max := 0
	for _, v := range paths {
		if v > max {
			max = v
		}
	}

	return max, nil
}

func getUniquePaths(pp []string) [][]string {
	seen := map[string]bool{}
	result := [][]string{}

	for perm := range permutation.NewPermutationChan(pp) {
		hash := strings.Join(perm, "-")

		// fmt.Printf("checking %s\n", hash)

		if !seen[hash] {
			result = append(result, perm)

			revHash := strings.Join(reverse(perm), "-")
			seen[hash] = true
			seen[revHash] = true
		}
	}

	return result
}

func reverse[T any](s []T) []T {
	reversed := s
	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}

	return reversed
}

func calcPaths(distances map[string]map[string]int) map[string]int {
	routes := map[string]int{}

	var keys []string

	for k := range distances {
		keys = append(keys, k)
	}

	paths := getUniquePaths(keys)
	// fmt.Printf("%d unique paths found\n", len(paths))

	// validate the path is possible
	for _, path := range paths {
		sum := 0
		valid := true

		for i := 0; i < len(path)-1; i++ {
			if d, ok := distances[path[i]][path[i+1]]; ok {
				sum += d
			} else {
				// fmt.Printf("invalid path: %s\n", path)
				valid = false
			}
		}

		if valid {
			// fmt.Printf("valid path: %s = %d\n", path, sum)
			routes[strings.Join(path, "-")] = sum
		}
	}

	// fmt.Printf("%d valid paths found\n", len(routes))

	return routes
}

func parse(instr string) (map[string]map[string]int, error) {
	distances := map[string]map[string]int{}
	re := regexp.MustCompile(`^(.+) to (.+) = (\d+)$`)

	for _, line := range strings.Split(instr, "\n") {
		var (
			c1, c2 string
			dist   int
		)

		matches := re.FindStringSubmatch(line)

		if len(matches) == 4 {
			c1 = matches[1]
			c2 = matches[2]
			dist, _ = strconv.Atoi(matches[3])
		} else {
			return nil, fmt.Errorf("invalid input: %s", line)
		}

		if _, ok := distances[c1]; !ok {
			distances[c1] = map[string]int{}
		}

		if _, ok := distances[c2]; !ok {
			distances[c2] = map[string]int{}
		}

		distances[c1][c2] = dist
		distances[c2][c1] = dist
	}

	return distances, nil
}
