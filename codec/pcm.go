package codec

// PcmaSplitToFrames bitrate(8k x 8bit) pcm alaw payload, 20ms per frame
func PcmaSplitToFrames(payload []byte) (frames [][]byte) {
	plen := len(payload)
	//frameLen := int(BitrateToFrameSize(8000.0*8, 20.0))
	frameLen := 160
	var i int
	for i < plen {
		frameEnd := i + frameLen
		if frameEnd > plen {
			if i < plen - 1 {
				// remaining little frame
				frames = append(frames,payload[i:plen])
			}
			break
		}
		frames = append(frames,payload[i:frameEnd])
		i = frameEnd
	}
	return
}
