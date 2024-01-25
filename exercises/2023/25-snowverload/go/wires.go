package exercises

import "strings"

type Wire struct {
	Src  string
	Dest string
}

func removeWire(components map[string][]string, wire Wire) {
	trimmedSrc := make([]string, 0, len(components[wire.Src]))
	for _, s := range components[wire.Src] {
		if s != wire.Dest {
			trimmedSrc = append(trimmedSrc, s)
		}
	}

	components[wire.Src] = trimmedSrc

	trimmedDest := make([]string, 0, len(components[wire.Dest]))
	for _, d := range components[wire.Dest] {
		if d != wire.Src {
			trimmedDest = append(trimmedDest, d)
		}
	}

	components[wire.Dest] = trimmedDest
}

func removeMax(wires map[Wire]int) Wire {
	var max int
	var maxWire Wire

	for wire, visits := range wires {
		if visits > max {
			max = visits
			maxWire = wire
		}
	}

	delete(wires, maxWire)

	return maxWire
}

func countVisits(components map[string][]string) map[Wire]int {
	encountered := map[Wire]int{}
	i := 0

	for c := range components {
		walk(components, c, encountered)
		i++

		// we don't need to visit all wires, just enough to find the right max. This is hacky.
		if i > 25 {
			break
		}
	}

	return encountered
}

func countComponents(components map[string][]string, start string) int {
	visited := make(map[string]bool, len(components))
	path := []string{start}

	for len(path) > 0 {
		cur := path[0]
		path = path[1:]

		for _, next := range components[cur] {
			if _, seen := visited[next]; seen {
				continue
			}

			path = append(path, next)
			visited[next] = true
		}
	}

	return len(visited)
}

func walk(components map[string][]string, start string, encountered map[Wire]int) {
	visited := map[string]bool{}
	path := []string{start} // TODO: is this faster with a list.List?

	for len(path) > 0 {
		cur := path[0]
		path = path[1:]

		for _, next := range components[cur] {
			if _, seen := visited[next]; seen {
				continue
			}

			path = append(path, next)
			visited[next] = true

			// sort so we don't count the same connection different ways
			var w Wire
			if cur < next {
				w = Wire{cur, next}
			} else {
				w = Wire{next, cur}
			}

			encountered[w]++
		}
	}
}

func parseComponents(in string) map[string][]string {
	lines := strings.Split(in, "\n")
	components := map[string][]string{}

	for _, line := range lines {
		src, rawDests, _ := strings.Cut(line, ": ")

		// add component if it doesn't exist yet
		if _, seen := components[src]; !seen {
			components[src] = []string{}
		}

		dests := strings.Fields(rawDests)
		for _, d := range dests {
			components[src] = append(components[src], d)

			// add component if it doesn't exist yet - handles if only on the right side
			if _, seen := components[d]; !seen {
				components[d] = []string{}
			}

			components[d] = append(components[d], src)
		}
	}

	return components
}
