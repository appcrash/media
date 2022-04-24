package utils

import "github.com/appcrash/GoRTP/rtp"

// PacketList is either received RTP data packet or generated packets by codecs that can be readily put to
// stack for transmission. audio data is usually one packet at a time as no pts is required, but video codecs can
// build multiple packets of the same pts. those packets can be linked and send to rtp stack as a whole.
type PacketList struct {
	Payload     []byte // rtp payload
	RawBuffer   []byte // rtp payload + rtp header
	PayloadType uint8
	Pts         uint32 // presentation timestamp
	Marker      bool   // should mark-bit in rtp header be set?
	Ssrc        uint32
	Csrc        []uint32

	next *PacketList // more PacketList, if any
}

func NewPacketListFromRtpPacket(packet *rtp.DataPacket) *PacketList {
	if packet.InUse() <= 0 || packet.Buffer() == nil {
		return nil
	}
	return &PacketList{
		Payload:     packet.Payload(),
		RawBuffer:   packet.Buffer()[:packet.InUse()],
		PayloadType: packet.PayloadType(),
		Pts:         packet.Timestamp(),
		Marker:      packet.Marker(),
		Ssrc:        packet.Ssrc(),
		Csrc:        packet.CsrcList(),
	}
}

func (pl *PacketList) HasMore() bool {
	return pl.next != nil
}

func (pl *PacketList) Iterate(f func(p *PacketList)) {
	ppl := pl
	for ppl != nil {
		f(ppl)
		ppl = ppl.next
	}
}

func (pl *PacketList) Next() *PacketList {
	return pl.next
}

func (pl *PacketList) SetNext(npl *PacketList) {
	pl.next = npl
}

func (pl *PacketList) Len() (length int) {
	pl.Iterate(func(ppl *PacketList) {
		length++
	})
	return
}
