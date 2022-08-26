package comp

import (
	"context"
	"errors"
	"time"
)

type ChanSrc struct {
	SessionNode

	context context.Context
	cancelF context.CancelFunc
	C       chan []byte
}

func (n *ChanSrc) Offer() []MessageType {
	return []MessageType{MtRawByte}
}

func (n *ChanSrc) Init() error {
	n.context, n.cancelF = context.WithCancel(context.Background())
	return nil
}

func (n *ChanSrc) OnExit() {
	n.cancelF()
}

func (n *ChanSrc) handleChannelLink(msg *ChannelLinkRequestMessage) {
	var resp interface{}

	defer func() { msg.C <- resp }()
	if n.C != nil {
		logger.Errorf("chan_src %v already linked", n)
		resp = "all ready linked"
		return
	}

	if msg.LinkChannel != nil {
		n.C = msg.LinkChannel
		go n.loop()
	} else {
		resp = "linking channel is nil"
	}
}

func (n *ChanSrc) LinkMe(c chan []byte) error {
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

func (n *ChanSrc) loop() {
	lp := n.GetLinkPoint(0)
	if lp == nil {
		logger.Warnf("chan_src %v has no output node", n)
		return
	}
	done := n.context.Done()
	for {
		select {
		case rawBytes, more := <-n.C:
			if !more {
				return
			}
			msg := &RawByteMessage{
				Data: rawBytes,
			}
			lp.SendMessage(msg)
		case <-done:
			return
		}
	}
}
