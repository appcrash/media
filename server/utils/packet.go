package utils

import "github.com/appcrash/GoRTP/rtp"

// RtpPacketList is either received RTP data packet or generated packets by codecs that can be readily put to
// stack for transmission. audio data is usually one packet at a time as no pts is required, but video codecs can
// build multiple packets of the same pts. those packets can be linked and send to rtp stack as a whole.
type RtpPacketList struct {
	Payload     []byte // rtp payload
	RawBuffer   []byte // rtp payload + rtp header
	PayloadType uint8
	Pts         uint32 // presentation timestamp
	PrevPts     uint32 // previous packet's pts
	Marker      bool   // should mark-bit in rtp header be set?
	Ssrc        uint32
	Csrc        []uint32

	next *RtpPacketList // more RtpPacketList, if any
}

func NewPacketListFromRtpPacket(packet *rtp.DataPacket) *RtpPacketList {
	if packet.InUse() <= 0 || packet.Buffer() == nil {
		return nil
	}
	return &RtpPacketList{
		Payload:     packet.Payload(),
		RawBuffer:   packet.Buffer()[:packet.InUse()],
		PayloadType: packet.PayloadType(),
		Pts:         packet.Timestamp(),
		Marker:      packet.Marker(),
		Ssrc:        packet.Ssrc(),
		Csrc:        packet.CsrcList(),
	}
}

func (pl *RtpPacketList) Iterate(f func(p *RtpPacketList)) {
	ppl := pl
	for ppl != nil {
		f(ppl)
		ppl = ppl.next
	}
}

func (pl RtpPacketList) CloneSingle() *RtpPacketList {
	return &RtpPacketList{
		Payload:     pl.Payload,
		RawBuffer:   pl.RawBuffer,
		PayloadType: pl.PayloadType,
		Pts:         pl.Pts,
		Marker:      pl.Marker,
		Ssrc:        pl.Ssrc,
		Csrc:        pl.Csrc,
	}
}

func (pl *RtpPacketList) Clone() *RtpPacketList {
	var cloned, current *RtpPacketList
	pl.Iterate(func(packet *RtpPacketList) {
		newPacket := packet.CloneSingle()
		if cloned == nil {
			cloned = newPacket
		} else {
			current.next = newPacket
		}

		current = newPacket
	})
	return cloned
}

func (pl *RtpPacketList) Next() *RtpPacketList {
	return pl.next
}

func (pl *RtpPacketList) SetNext(npl *RtpPacketList) {
	pl.next = npl
}

func (pl *RtpPacketList) GetLast() *RtpPacketList {
	ppl := pl
	for ppl.next != nil {
		ppl = ppl.next
	}
	return ppl
}

func (pl *RtpPacketList) Len() (length int) {
	pl.Iterate(func(ppl *RtpPacketList) {
		length++
	})
	return
}
