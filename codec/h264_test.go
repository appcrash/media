package codec_test

import (
	"bytes"
	"github.com/appcrash/media/codec"
	"testing"
)

func makeDummyH264Nal(nalType uint8, nri uint8, dataCount int) []byte {
	data := make([]byte, dataCount)
	for i := 1; i < dataCount; i++ {
		data[i] = 0x01
	}
	data[0] = (nalType & codec.BitmaskNalType) | ((nri << 5) & codec.BitmaskRefIdc)
	return data
}

func joinNals(nals [][]byte) []byte {
	return bytes.Join(nals, []byte{0x00, 0x00, 0x01})
}

func TestPacketListFromH264(t *testing.T) {
	nal500 := makeDummyH264Nal(1, 0x11, 500)
	nal495 := makeDummyH264Nal(1, 0x11, 495) // 5 = 2 * nalSize(2 bytes) + stapA header(1 byte)
	nal496 := makeDummyH264Nal(1, 0x11, 496)
	nal1000 := makeDummyH264Nal(1, 0x11, 1000)
	nal1001 := makeDummyH264Nal(1, 0x11, 1001)

	pl := codec.PacketListFromH264(joinNals([][]byte{nal500, nal495}), 100, 1000)
	if pl.Len() != 1 && len(pl.Payload) != 1000 {
		t.Fatal("should be packed as a stapA")
	}
	pl = codec.PacketListFromH264(joinNals([][]byte{nal500, nal496}), 100, 1000)
	if pl.Len() != 2 && len(pl.Payload) != 500 && len(pl.Next.Payload) != 496 {
		t.Fatal("should be packed as 2 nals")
	}
	pl = codec.PacketListFromH264(nal1000, 100, 1000)
	if pl.Len() != 1 && len(pl.Payload) != 1000 {
		t.Fatal("should be packed as a single nal")
	}
	pl = codec.PacketListFromH264(nal1001, 100, 1000)
	ppl := pl.Next
	if pl.Len() != 2 &&
		len(pl.Payload) != 1000 && (pl.Payload[0]&codec.BitmaskNalType) != codec.NalTypeFua &&
		(pl.Payload[1]&codec.BitmaskFuStart) != 0x00 && len(ppl.Payload) != 4 &&
		(ppl.Payload[0]&codec.BitmaskNalType) != codec.NalTypeFua && (ppl.Payload[1]&codec.BitmaskFuEnd) != 00 {
		t.Fatal("should be packed as a 2 nals with proper indicator and header")
	}

	nals := joinNals([][]byte{nal500, nal496, nal1001, nal1000})
	nals = append([]byte{0x00, 0x00, 0x00, 0x01}, nals...) // with annexB start code
	pl = codec.PacketListFromH264(nals, 100, 1000)
	if pl.Len() != 5 {
		t.Fatal("should be packed as single, aggregation and fu nals")
	}
}
