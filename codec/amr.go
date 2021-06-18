package codec

// shamelessly copied from ffmpeg(amr.c)  :)

var amrnbPackedSize = [16]int{
	13, 14, 16, 18, 20, 21, 27, 32, 6, 1, 1, 1, 1, 1, 1, 1,
}

var amrwbPackedSize = [16]int{
	18, 24, 33, 37, 41, 47, 51, 59, 61, 6, 1, 1, 1, 1, 1, 1,
}

func AmrSplitToFrames(payload []byte) (frames [][]byte) {
	plen := len(payload)
	var i int
	for i < plen {
		toc := payload[i]
		mode := (toc >> 3) & 0x0f
		frameLen := amrnbPackedSize[mode] // TODO: support amr-wb
		frameEnd := i + frameLen
		if frameEnd > plen {
			break
		}
		frames = append(frames, payload[i:frameEnd])
		i = frameEnd
	}
	return
}
