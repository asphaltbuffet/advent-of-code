package exercises

//go:generate stringer -type=State
type State int

const (
	Off State = iota
	On
)

type FlipFlop struct {
	ID           string
	destinations []string
	State        State
}

func (ff *FlipFlop) Receive(msg Message) {
	if msg.Pulse == High {
		return
	}

	var sendPulse Pulse

	switch ff.State {
	case On:
		ff.State = Off
		sendPulse = Low

	case Off:
		ff.State = On
		sendPulse = High
	}

	for _, dst := range ff.destinations {
		msgQueue.Send(ff.ID, dst, sendPulse)
	}
}

func (ff *FlipFlop) GetDests() []string {
	return ff.destinations
}

func (ff *FlipFlop) GetID() string {
	return ff.ID
}

func (ff *FlipFlop) GetState() State {
	return ff.State
}

func (ff *FlipFlop) SetState(s State) {
	ff.State = s
}
