package codec

import "math"
import "github.com/appcrash/media/server/rpc"

func BitrateToFrameSize(bitrate float64, frameIntervalMs float64) int32 {
	frameSize := bitrate / 8.0 / (1000.0 / frameIntervalMs)
	return int32(math.Floor(frameSize))
}

// send time step in milliseconds
func GetCodecTimeStep(codec rpc.CodecType) int {
	switch codec {
	case rpc.CodecType_PCM_ALAW:
		fallthrough
	case rpc.CodecType_AMRNB:
		fallthrough
	case rpc.CodecType_AMRWB:
		return 20
	case rpc.CodecType_H264:
		return 1000 / 25 // TODO: use dynamic framerate
	}

	// unknown codec
	return 20
}

var RtcpPacketTypeMap = map[int]string{
	200: "SENDER REPORT",
	201: "RECEIVER REPORT",
	202: "SOURCE DESCRIPTION",
	203: "GOODBYE",
}
