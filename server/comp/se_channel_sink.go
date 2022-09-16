package comp

// ChanSink forward input message to go channel, once channel is linked, the channel would receive what node gets,
// and be closed when node exiting graph
type ChanSink struct {
	SessionNode

	C chan []byte
}

func (n *ChanSink) handleRawByte(msg *RawByteMessage) {
	select {
	case n.C <- msg.Data:
	default:
	}
}

func (n *ChanSink) handleChannelLink(msg *ChannelLinkRequestMessage) {
	defer func() { msg.C <- nil }()

	if msg.LinkChannel != nil {
		n.C = msg.LinkChannel
	}
}

func (n *ChanSink) LinkMe(c chan []byte) {
	msg := &ChannelLinkRequestMessage{
		LinkChannel: c,
	}
	msg.C = make(chan interface{})
	n.DeliverToStream(msg)
}

func (n *ChanSink) OnExit() {
	if n.C == nil {
		return
	}
	close(n.C)
	n.C = nil
}
