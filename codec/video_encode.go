package codec

//#cgo pkg-config: libavformat libavcodec libavutil libswresample libavfilter freetype2
//
//#include <stdio.h>
//#include <stdlib.h>
//#include <inttypes.h>
//#include <stdint.h>
//#include <string.h>
//#include <libavformat/avformat.h>
//#include "codec.h"
import "C"
import (
	"fmt"
	"unsafe"
)

type (
	VideoContext = C.struct_VideoContext
	AVPacket     = C.struct_AVPacket
)

func NewVideoContext() *VideoContext {
	return (*VideoContext)(C.video_init())
}

func (vc *VideoContext) Iterate() {
	C.video_iterate(vc)
}

func (vc *VideoContext) PacketNumber() int {
	return *(*int)(unsafe.Pointer(&vc.nb_packet))
}

func (vc *VideoContext) GetPacket(i int) *H264Packet {
	packets := *(*[256]*AVPacket)(unsafe.Pointer(&vc.packet_data[0]))
	pkt := packets[i]
	return &H264Packet{
		Payload: C.GoBytes(unsafe.Pointer(pkt.data), pkt.size),
		Pts:     *(*int64)(unsafe.Pointer(&pkt.pts)),
	}
}

func EncodeText() {
	vc := NewVideoContext()
	for i := 0; i < 100; i++ {
		for vc.Iterate(); vc.PacketNumber() == 0; vc.Iterate() {
		}
		fmt.Printf("i is %v\n", i)
		fmt.Printf("nb_packet is %v\n", vc.PacketNumber())
		pkt := vc.GetPacket(0)
		payload := pkt.Payload
		fmt.Printf("packet size is %v,pts is %v\n", len(payload), pkt.Pts)
		for _, nal := range extractNals(payload) {
			printNal(nal)
		}
		pl := PacketListFromH264(payload, 100, 1400)
		logger.Infoln("...................")
		for pl != nil {
			printNal(pl.Payload)
			pl = pl.Next
		}
	}

}
