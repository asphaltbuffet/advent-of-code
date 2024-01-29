package exercises

import (
	"fmt"
	"slices"
	"strings"
)

type Conjunction struct {
	ID           string
	destinations []string
	memory       Memory
}

type Memory map[string]Pulse

func (m Memory) String() string {
	items := make([]string, 0, len(m))

	for k, v := range m {
		items = append(items, fmt.Sprintf("[%s: %s]", k, v))
	}

	return "[" + strings.Join(items, " ") + "]"
}

func (c *Conjunction) Receive(msg Message) {
	var sendPulse Pulse

	c.memory[msg.Src] = msg.Pulse

	if c.AllHigh() {
		sendPulse = Low
	} else {
		sendPulse = High
	}

	for _, dst := range c.destinations {
		msgQueue.Send(c.ID, dst, sendPulse)
	}
}

func (c *Conjunction) GetDests() []string {
	return c.destinations
}

func (c *Conjunction) GetInputs() []string {
	inputs := make([]string, 0, len(c.memory))

	for id := range c.memory {
		inputs = append(inputs, id)
	}

	slices.Sort(inputs) // Sort for deterministic order

	return inputs
}

func (c *Conjunction) SetMemory(id string, p Pulse) {
	c.memory[id] = p
}

func (c *Conjunction) GetMemory() Memory {
	return c.memory
}

func (c *Conjunction) AllHigh() bool {
	for _, p := range c.memory {
		if p == Low {
			return false
		}
	}

	return true
}

func (c *Conjunction) GetID() string {
	return c.ID
}
