package exercises

import (
	"strconv"
	"strings"
	"sync"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 7.
type Exercise struct {
	common.BaseExercise
}

// parseProgram parses a comma-separated Intcode program.
func parseProgram(s string) []int {
	s = strings.TrimSpace(s)
	parts := strings.Split(s, ",")
	mem := make([]int, len(parts))
	for i, p := range parts {
		v, _ := strconv.Atoi(strings.TrimSpace(p))
		mem[i] = v
	}
	return mem
}

// runIntcode runs an Intcode program with the given input queue, returns outputs.
//
//nolint:gocognit // Intcode VM opcode dispatch is inherently complex
func runIntcode(prog []int, inputs []int) []int {
	mem := make([]int, len(prog))
	copy(mem, prog)

	ip := 0
	inIdx := 0
	var outputs []int

	param := func(offset, mode int) int {
		v := mem[ip+offset]
		if mode == 0 {
			return mem[v]
		}
		return v
	}

	for {
		instr := mem[ip]
		op := instr % 100
		m1 := (instr / 100) % 10
		m2 := (instr / 1000) % 10

		switch op {
		case 1: // add
			mem[mem[ip+3]] = param(1, m1) + param(2, m2)
			ip += 4
		case 2: // mul
			mem[mem[ip+3]] = param(1, m1) * param(2, m2)
			ip += 4
		case 3: // input
			mem[mem[ip+1]] = inputs[inIdx]
			inIdx++
			ip += 2
		case 4: // output
			outputs = append(outputs, param(1, m1))
			ip += 2
		case 5: // jump-if-true
			if param(1, m1) != 0 {
				ip = param(2, m2)
			} else {
				ip += 3
			}
		case 6: // jump-if-false
			if param(1, m1) == 0 {
				ip = param(2, m2)
			} else {
				ip += 3
			}
		case 7: // less than
			if param(1, m1) < param(2, m2) {
				mem[mem[ip+3]] = 1
			} else {
				mem[mem[ip+3]] = 0
			}
			ip += 4
		case 8: // equals
			if param(1, m1) == param(2, m2) {
				mem[mem[ip+3]] = 1
			} else {
				mem[mem[ip+3]] = 0
			}
			ip += 4
		case 99:
			return outputs
		}
	}
}

// permutations returns all permutations of the given slice.
func permutations(nums []int) [][]int {
	if len(nums) == 0 {
		return [][]int{{}}
	}
	var result [][]int
	for i, v := range nums {
		rest := make([]int, 0, len(nums)-1)
		rest = append(rest, nums[:i]...)
		rest = append(rest, nums[i+1:]...)
		for _, p := range permutations(rest) {
			result = append(result, append([]int{v}, p...))
		}
	}
	return result
}

// runAmplifiers runs 5 amplifiers in series with the given phase settings.
func runAmplifiers(prog []int, phases []int) int {
	signal := 0
	for _, phase := range phases {
		outs := runIntcode(prog, []int{phase, signal})
		signal = outs[len(outs)-1]
	}
	return signal
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	prog := parseProgram(instr)
	phases := []int{0, 1, 2, 3, 4}
	maxSignal := 0
	first := true
	for _, perm := range permutations(phases) {
		sig := runAmplifiers(prog, perm)
		if first || sig > maxSignal {
			maxSignal = sig
			first = false
		}
	}
	return maxSignal, nil
}

// runIntcodeAsync runs an Intcode program reading from in and writing to out, then signals done.
//
//nolint:gocognit // Intcode VM opcode dispatch is inherently complex
func runIntcodeAsync(prog []int, in <-chan int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	mem := make([]int, len(prog))
	copy(mem, prog)
	ip := 0

	param := func(offset, mode int) int {
		v := mem[ip+offset]
		if mode == 0 {
			return mem[v]
		}
		return v
	}

	for {
		instr := mem[ip]
		op := instr % 100
		m1 := (instr / 100) % 10
		m2 := (instr / 1000) % 10

		switch op {
		case 1:
			mem[mem[ip+3]] = param(1, m1) + param(2, m2)
			ip += 4
		case 2:
			mem[mem[ip+3]] = param(1, m1) * param(2, m2)
			ip += 4
		case 3:
			mem[mem[ip+1]] = <-in
			ip += 2
		case 4:
			out <- param(1, m1)
			ip += 2
		case 5:
			if param(1, m1) != 0 {
				ip = param(2, m2)
			} else {
				ip += 3
			}
		case 6:
			if param(1, m1) == 0 {
				ip = param(2, m2)
			} else {
				ip += 3
			}
		case 7:
			if param(1, m1) < param(2, m2) {
				mem[mem[ip+3]] = 1
			} else {
				mem[mem[ip+3]] = 0
			}
			ip += 4
		case 8:
			if param(1, m1) == param(2, m2) {
				mem[mem[ip+3]] = 1
			} else {
				mem[mem[ip+3]] = 0
			}
			ip += 4
		case 99:
			return
		}
	}
}

// runFeedbackLoop runs 5 amplifiers in a feedback loop with the given phase settings.
// Returns the last value amp E outputs.
func runFeedbackLoop(prog []int, phases []int) int {
	const n = 5
	chans := make([]chan int, n)
	for i := range chans {
		chans[i] = make(chan int, 100)
	}

	// Prime each amp's input channel with its phase.
	for i, phase := range phases {
		chans[i] <- phase
	}
	// Send initial signal to amp A.
	chans[0] <- 0

	var wg sync.WaitGroup
	wg.Add(n)
	for i := range n {
		go runIntcodeAsync(prog, chans[i], chans[(i+1)%n], &wg)
	}

	// Collect the last value written to chans[0] (amp E's output = amp A's input).
	// We do this by draining chans[0] after all amps halt.
	wg.Wait()

	// After all goroutines finish, chans[0] holds the last output from amp E.
	last := 0
	for len(chans[0]) > 0 {
		last = <-chans[0]
	}
	return last
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	prog := parseProgram(instr)
	phases := []int{5, 6, 7, 8, 9}
	maxSignal := 0
	first := true
	for _, perm := range permutations(phases) {
		sig := runFeedbackLoop(prog, perm)
		if first || sig > maxSignal {
			maxSignal = sig
			first = false
		}
	}
	return maxSignal, nil
}
