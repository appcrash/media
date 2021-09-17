package utils

import (
	"github.com/appcrash/media/codec"
	"github.com/appcrash/media/server/rpc"
)

// transform payload of rtp packet into proper frame that can be used by event graph
// i.e. transcoding, recording
func TransformPayloadToFrame(codecType rpc.CodecType, payload []byte) (frame [][]byte) {
	switch codecType {
	case rpc.CodecType_AMRNB:
		frame = codec.AmrRtpPayloadToFrame(payload, false)
	case rpc.CodecType_AMRWB:
		frame = codec.AmrRtpPayloadToFrame(payload, true)
	case rpc.CodecType_PCM_ALAW:
		fallthrough
	default:
		frame = [][]byte{payload}
	}
	return
}

// transform audio frame into the format of proper rtp payload
func TransformFrameToPayload(codecType rpc.CodecType, frame []byte) (payload []byte) {
	switch codecType {
	case rpc.CodecType_AMRNB:
		fallthrough
	case rpc.CodecType_AMRWB:
		payload = codec.AmrFrameToRtpPayload([][]byte{frame})[0]
	case rpc.CodecType_PCM_ALAW:
		fallthrough
	default:
		payload = frame
	}
	return
}
