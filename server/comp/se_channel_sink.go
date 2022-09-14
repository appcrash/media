package comp

// ChanSink forward input message to go channel, once channel is linked, the channel would receive what node gets,
// and be closed when node exiting graph
type ChanSink struct {
	SessionNode

	C chan []byte
}

func (cs *ChanSink) handleRawByte(msg *RawByteMessage) {
	select {
	case cs.C <- msg.Data:
	default:
	}
}

func (cs *ChanSink) handleChannelLink(msg *ChannelLinkMessage) {
	defer func() { msg.C <- nil }()

	if msg.LinkChannel != nil {
		cs.C = msg.LinkChannel
	}
}

func (cs *ChanSink) ChannelLink(c chan []byte) {
	msg := &ChannelLinkMessage{
		LinkChannel: c,
	}
	msg.C = make(chan interface{})
	cs.DeliverToStream(msg)
}

func (cs *ChanSink) OnExit() {
	close(cs.C)
	cs.C = nil
}
