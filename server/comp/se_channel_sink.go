package comp

import "github.com/appcrash/media/server/event"

// ChanSink forward input message to go channel, once channel is linked, the channel would receive what node get,
// and be closed by node when exiting graph
type ChanSink struct {
	SessionNode

	C chan []byte
}

func (cs *ChanSink) Accept() []MessageType {
	return []MessageType{MtRawByte}
}

func (cs *ChanSink) Init() error {
	cs.SetHandler(MtRawByte, cs.handleRawByte)
	cs.SetHandler(MtChannelLink, cs.handleChannelLink)
	return nil
}

func (cs *ChanSink) handleRawByte(evt *event.Event) {
	if msg, ok := ToMessage[*RawByteMessage](evt); ok {
		select {
		case cs.C <- msg.Data:
		default:
		}
	}
}

func (cs *ChanSink) handleChannelLink(evt *event.Event) {
	if msg, ok := ToMessage[*ChannelLinkMessage](evt); ok {
		defer func() { msg.C <- nil }()

		if msg.LinkChannel != nil {
			cs.C = msg.LinkChannel
		}
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
