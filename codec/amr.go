package codec

// shamelessly copied from ffmpeg(amr.c)  :)

var amrnbPackedSize = [16]int{
	13, 14, 16, 18, 20, 21, 27, 32, 6, 1, 1, 1, 1, 1, 1, 1,
}

var amrwbPackedSize = [16]int{
	18, 24, 33, 37, 41, 47, 51, 59, 61, 6, 1, 1, 1, 1, 1, 1,
}

// AmrSplitToFrames transform payload read from amr file into frames
func AmrSplitToFrames(payload []byte, isAmrwb bool) (frames [][]byte) {
	plen := len(payload)
	var i int
	packedSize := amrnbPackedSize
	if isAmrwb {
		packedSize = amrwbPackedSize
	}
	for i < plen {
		toc := payload[i]
		mode := (toc >> 3) & 0x0f
		frameLen := packedSize[mode]
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

func AmrFrameToRtpPayload(frame [][]byte) (rtpPayload [][]byte) {
	for _, p := range frame {
		toc := []byte{p[0] & 0x7f} // set F bit zero, so this is the last frame in the frame(only one)
		//if i == 0 {
		//	header = []byte{0x70}   // CMR mode 7 12.2kbit/s
		//} else {
		//	header = []byte{0xf0}   // CMR is 15, no mode change required
		//}

		rp := append([]byte{0xf0}, toc...)
		rp = append(rp, p[1:]...)
		rtpPayload = append(rtpPayload, rp)
	}
	return
}

// AmrRtpPayloadToFrame skip 1 byte in the header, count the frames in the index, then extract each frame data
func AmrRtpPayloadToFrame(payload []byte, isAmrwb bool) (frames [][]byte) {
	var nbFrame = 1
	pl := len(payload)
	packedSize := amrnbPackedSize
	if isAmrwb {
		packedSize = amrwbPackedSize
	}
	for nbFrame < pl-1 {
		if (payload[nbFrame] & 0x80) != 0 {
			// F bit is set, go to next
			nbFrame++
			continue
		}
		break
	}
	if nbFrame+1 >= pl {
		// invalid frame
		return
	}
	speechData := 1 + nbFrame
	for i := 0; i < nbFrame; i++ {
		toc := payload[1+i]
		mode := (toc >> 3) & 0x0f
		frameEnd := speechData + packedSize[mode] - 1 // packedSize counts the toc
		if frameEnd > pl {
			return
		}
		f := []byte{toc & 0x7c} // keep Q	bit
		if frameEnd > speechData {
			f = append(f, payload[speechData:frameEnd]...)
		}
		speechData = frameEnd
		frames = append(frames, f)
	}
	return
}
