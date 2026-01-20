package exercises

import (
	"fmt"
	"slices"
	"strings"
)

type Devices map[string][]string

func GetDevices(s string) *Devices {
	lines := strings.Split(s, "\n")
	devices := make(Devices, len(lines))

	for _, l := range lines {
		tok := strings.Fields(l)
		devices[strings.Trim(tok[0], ":")] = tok[1:]
	}
	// fmt.Printf("%#v\n", devices)

	return &devices
}

func (d Devices) Trace(memo map[string]int, cur string, keys ...string) int {
	state := fmt.Sprint(cur, keys)
	if v, ok := memo[state]; ok {
		return v
	}

	n := len(keys)
	idx := 0

	if slices.Contains(keys, cur) {
		// we got to end
		if n == 1 {
			memo[state] = 1
			return 1
		}

		// we skipped a necessary path
		if cur != keys[0] {
			memo[state] = 0
			return 0
		}

		// we hit intermediate point
		idx = 1
	}

	paths := 0
	for _, k := range d[cur] {
		paths += d.Trace(memo, k, keys[idx:]...)
	}

	memo[state] = paths
	return paths
}

func (d Devices) PathExists(a, b string) bool {
	if a == b {
		return true
	}

	for _, k := range d[a] {
		if d.PathExists(k, b) {
			return true
		}
	}

	return false
}
