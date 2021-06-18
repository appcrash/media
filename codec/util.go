package codec

import "math"

func BitrateToFrameSize(bitrate float64, frameIntervalMs float64) int32 {
	frameSize := bitrate / 8.0 / (1000.0 / frameIntervalMs)
	return int32(math.Floor(frameSize))
}
