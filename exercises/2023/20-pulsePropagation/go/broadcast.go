package exercises

type Broadcast struct {
	destinations []string
}

const BroadcasterID = "broadcaster"

func (b *Broadcast) Receive(msg Message) {
	for _, dst := range b.destinations {
		msgQueue.Send(BroadcasterID, dst, msg.Pulse)
	}
}

func (b *Broadcast) GetDests() []string {
	return b.destinations
}

func (b *Broadcast) GetID() string {
	return BroadcasterID
}
