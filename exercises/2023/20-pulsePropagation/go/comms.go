package exercises

import (
	"fmt"
	"slices"
	"strings"

	util "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

//go:generate stringer -type=Pulse
type Pulse int

const (
	Low Pulse = iota
	High
)

type ModConfig map[string]Module

type Module interface {
	GetID() string
	Receive(Message)
	GetDests() []string
}

type Conjuncter interface {
	Module

	GetMemory() Memory
	SetMemory(string, Pulse)
	GetInputs() []string
	AllHigh() bool
}

type FlipFlopper interface {
	Module

	GetState() State
	SetState(State)
}

var msgQueue MessageQueue

func loadInput(input string) (ModConfig, error) {
	lines := strings.Split(input, "\n")
	cfg := make(ModConfig, len(lines))

	for _, line := range lines {
		if err := cfg.addModule(line); err != nil {
			return nil, err
		}
	}

	cfg.InitModules()

	// cfg.DebugPrint()

	return cfg, nil
}

func (cfg ModConfig) DebugPrint() {
	fmt.Println("╭ Configuration:")

	for _, m := range cfg {
		fmt.Printf("├╴ID: %q\n", m.GetID())
		fmt.Printf("│  ⤷ Type: %T\n", m)
		fmt.Printf("│  ⤷ Dest: %v\n", m.GetDests())

		switch v := m.(type) {
		case Conjuncter:
			fmt.Printf("│  ⤷ Memory: %s\n", v.GetMemory())

		case FlipFlopper:
			fmt.Printf("│  ⤷ State: %s\n", v.GetState())
		}
	}

	fmt.Println("╰─┄┄┈")
}

// InitModules populates the memory of all conjunction modules
// and sets all flipflops to off.
func (cfg ModConfig) InitModules() {
	dests := make(map[string][]string, len(cfg))

	for _, m := range cfg {
		for _, d := range m.GetDests() {
			dests[d] = append(dests[d], m.GetID())
		}
	}

	for _, m := range cfg {
		switch v := m.(type) {
		case Conjuncter:
			for _, d := range dests[v.GetID()] {
				v.SetMemory(d, Low)
			}

		case FlipFlopper:
			// not necessary for initial state, but also works as a reset
			v.SetState(Off)

		default:
			continue
		}
	}
}

func (cfg ModConfig) addModule(s string) error {
	m, c, ok := strings.Cut(s, " -> ")
	if !ok {
		return fmt.Errorf("invalid input: %s", s)
	}

	modType, modName, conns := m[0], m[1:], strings.Split(c, ", ")

	switch modType {
	case 'b':
		cfg[BroadcasterID] = &Broadcast{destinations: conns}

	case '%':
		cfg[modName] = &FlipFlop{ID: modName, State: Off, destinations: conns}

	case '&':
		cfg[modName] = &Conjunction{ID: modName, destinations: conns, memory: map[string]Pulse{}}

	default:
		return fmt.Errorf("invalid module type: %c", modType)
	}

	return nil
}

func (cfg ModConfig) ProcessQueue() (int, int) {
	counts := map[Pulse]int{Low: 0, High: 0}

	for len(msgQueue) > 0 {
		var msg Message

		msg, msgQueue = msgQueue[0], msgQueue[1:]

		counts[msg.Pulse]++

		// drop it on the floor if module isn't a receiver
		if m, ok := cfg[msg.Dst]; ok {
			m.Receive(msg)
		}
	}

	return counts[Low], counts[High]
}

func (cfg ModConfig) ProcessUntilHigh(cj Conjuncter) int {
	targets := cj.GetInputs()
	cycleCounts := map[string]int{}
	lcm := 1

	// build up the cycle counts for each target and end when we've seen them all
	for i := 1; len(targets) != len(cycleCounts); i++ {
		// send a pulse from the button
		msgQueue.Send("button", "broadcaster", Low)

		for len(msgQueue) > 0 {
			msg := msgQueue[0]
			msgQueue = msgQueue[1:]

			m, ok := cfg[msg.Dst]
			if !ok {
				// drop it on the floor if module isn't a receiver
				continue
			}

			if msg.Dst == cj.GetID() && msg.Pulse == High {
				if _, seen := cycleCounts[msg.Src]; !seen {
					cycleCounts[msg.Src] = i
					lcm = util.LCM(lcm, i)
				}
			}

			m.Receive(msg)
		}
	}

	return lcm
}

// find the module that sends a signal to the target
func (cfg ModConfig) GetFeed(id string) (*Conjunction, error) {
	// we have to loop through everything since the receiver may not be registered
	for _, m := range cfg {
		if slices.Contains(m.GetDests(), id) {
			cj, ok := m.(*Conjunction)
			if !ok {
				return nil, fmt.Errorf("only works if the source is a conjunction")
			}

			return cj, nil
		}
	}

	return nil, fmt.Errorf("no feed found for %q", id)
}
