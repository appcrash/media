package codec

// PcmaSplitToFrames bitrate(8k x 8bit) pcm alaw payload, 20ms per frame
func PcmaSplitToFrames(payload []byte, intervalMs int) (frames [][]byte) {
	plen := len(payload)
	frameLenInBytes := 8000 / 1000 * intervalMs
	var i int
	for i < plen {
		frameEnd := i + frameLenInBytes
		if frameEnd > plen {
			if i < plen-1 {
				// remaining little frame
				frames = append(frames, payload[i:plen])
			}
			break
		}
		frames = append(frames, payload[i:frameEnd])
		i = frameEnd
	}
	return
}
