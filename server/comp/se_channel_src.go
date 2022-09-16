package comp

import "context"

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

func (n *ChanSrc) handleChannelLink(msg *ChannelLinkMessage) {
	defer func() { msg.C <- nil }()

	if n.C != nil {
		logger.Errorf("chan_src %v already linked", n)
		return
	}

	if msg.LinkChannel != nil {
		n.C = msg.LinkChannel
		go n.loop()
	}
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
