package comp

import "github.com/appcrash/media/server/utils"

// RtpSink accept rtp packet message and forward packet list to rtp stack
type RtpSink struct {
	SessionNode

	pullC chan *utils.RtpPacketList
}

func (n *RtpSink) Init() error {
	n.pullC = make(chan *utils.RtpPacketList)
	return nil
}

func (n *RtpSink) OnExit() {
	close(n.pullC)
	n.pullC = nil
}

func (n *RtpSink) handleRtpPacket(msg *RtpPacketMessage) {
	select {
	case n.pullC <- msg.Packet:
	default:
	}
}

func (n *RtpSink) PullPacketChannel() <-chan *utils.RtpPacketList {
	return n.pullC
}
