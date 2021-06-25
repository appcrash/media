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
		//frameLen := amrnbPackedSize[mode] // TODO: support amr-wb
		frameLen := amrwbPackedSize[mode] // TODO: support amr-wb
		frameEnd := i + frameLen
		if frameEnd > plen {
			break
		}
		frames = append(frames, payload[i:frameEnd])
		i = frameEnd
	}
	return
}


//  0 1 2 3 4 5 6 7
// +-+-+-+-+-+-+-+-+
// |F|  FT   |Q|P|P|
// +-+-+-+-+-+-+-+-+
// F (1 bit): see definition in Section 4.3.2.
// FT (4 bits, unsigned integer): see definition in Section 4.3.2.
// Q (1 bit): see definition in Section 4.3.2.
// P bits: padding bits, MUST be set to zero, and MUST be ignored on
func AmrFrameToRtpPayload(payload [][]byte) (rtpPayload [][]byte) {
	for _,p := range payload {
		toc := []byte{p[0] & 0x7f} // set F bit zero, so this is the last frame in the payload(only one)
		//if i == 0 {
		//	header = []byte{0x70}   // CMR mode 7 12.2kbit/s
		//} else {
		//	header = []byte{0xf0}   // CMR is 15, no mode change required
		//}

		rp := append([]byte{0xf0},toc...)
		rp = append(rp,p[1:]...)
		rtpPayload = append(rtpPayload,rp)
	}
	return
}
