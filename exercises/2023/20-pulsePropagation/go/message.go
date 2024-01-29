package exercises

type MessageQueue []Message

type Message struct {
	Src   string
	Dst   string
	Pulse Pulse
}

func (mq *MessageQueue) Send(src, dst string, p Pulse) {
	*mq = append(*mq, Message{Src: src, Dst: dst, Pulse: p})
}
