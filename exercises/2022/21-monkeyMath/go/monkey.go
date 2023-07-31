package exercises

import (
	"fmt"
	"strconv"
	"strings"
)

func parse(instr string) map[string]string {
	out := make(map[string]string, strings.Count(instr, "\n")+1)

	for _, line := range strings.Split(instr, "\n") {
		name, action, _ := strings.Cut(line, ": ")
		out[name] = action
	}

	return out
}

func calc(name string, raw map[string]string, done map[string]int) (int, error) {
	if val, ok := done[name]; ok {
		return val, nil
	}

	action := strings.Split(raw[name], " ")

	switch len(action) {
	case 1:
		if val, err := strconv.Atoi(action[0]); err != nil {
			return 0, err
		} else {
			done[name] = val
		}
	case 3:
		left, err := calc(action[0], raw, done)
		if err != nil {
			return 0, fmt.Errorf("caculating left of %q: %w", name, err)
		}

		right, err := calc(action[2], raw, done)
		if err != nil {
			return 0, fmt.Errorf("caculating right of %q: %w", name, err)
		}

		switch action[1] {
		case "+":
			done[name] = left + right
		case "-":
			done[name] = left - right
		case "*":
			done[name] = left * right
		case "/":
			done[name] = left / right
		default:
			return 0, fmt.Errorf("nnknown operator: %s", action[1])
		}
	}

	return done[name], nil
}
