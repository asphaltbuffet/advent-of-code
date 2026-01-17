package exercises

import (
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

func (d Devices) Trace(key string) int {
	if key == "out" {
		return 1
	}

	paths := 0
	for _, k := range d[key] {
		paths += d.Trace(k)
	}

	return paths
}
