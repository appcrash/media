package comp

import (
	"errors"
	"time"
)

// ChanSink forward input message to go channel, once channel is linked, the channel would receive what node gets,
// and be closed when node exiting graph
type ChanSink struct {
	SessionNode

	C chan []byte
}

func (n *ChanSink) handleRawByte(msg *RawByteMessage) {
	logger.Infof("%v raw byte", n)
	select {
	case n.C <- msg.Data:
	default:
	}
}

func (n *ChanSink) handleChannelLink(msg *ChannelLinkRequestMessage) {
	if n.C != nil {
		logger.Errorf("chan_sink %v already linked", n)
		msg.C <- "all ready linked"
		return
	}

	if msg.LinkChannel != nil {
		n.C = msg.LinkChannel
		msg.C <- nil
	} else {
		msg.C <- "linking channel is nil"
	}
}

func (n *ChanSink) LinkMe(c chan []byte) error {
	msg := &ChannelLinkRequestMessage{
		LinkChannel: c,
	}
	msg.C = make(chan interface{})
	n.DeliverToStream(msg)
	select {
	case resp := <-msg.C:
		if resp != nil {
			return errors.New("link failed")
		}
	case <-time.After(2 * time.Second):
		return errors.New("link timeout")
	}
	return nil
}

func (n *ChanSink) OnExit() {
	if n.C == nil {
		return
	}
	close(n.C)
	n.C = nil
}
