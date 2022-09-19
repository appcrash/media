package comp

import (
	"context"
	"github.com/appcrash/media/server/utils"
)

// RtpSrc is an entry node of graph as well as a rtp packet consumer, it connects rtp stack and graph
type RtpSrc struct {
	SessionNode
	InitiatorNode

	cancelF context.CancelFunc
	handleC chan *utils.RtpPacketList
}

func (n *RtpSrc) Init() error {
	n.handleC = make(chan *utils.RtpPacketList)
	ctx, cancel := context.WithCancel(context.Background())
	n.cancelF = cancel
	go n.loop(ctx)
	return nil
}

func (n *RtpSrc) Offer() []MessageType {
	return []MessageType{MtRtpPacket}
}

func (n *RtpSrc) HandlePacketChannel() chan<- *utils.RtpPacketList {
	return n.handleC
}

func (n *RtpSrc) OnExit() {
	n.cancelF()
}

func (n *RtpSrc) loop(ctx context.Context) {
	done := ctx.Done()

	// use the first link point only
	lp := n.GetLinkPoint(0)
	if lp == nil {
		logger.Errorf("rtp_src %v has no output node", n)
		return
	}
	for {
		select {
		case packet, more := <-n.handleC:
			if !more {
				return
			}
			msg := &RtpPacketMessage{Packet: packet}
			lp.SendMessage(msg)
		case <-done:
			return
		}
	}
}
