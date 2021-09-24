package codec

import (
	"encoding/binary"
)

// shamelessly copied from ffmpeg(amr.c)  :)
// packedSize = frameSize + 1byte(toc), used in octet-align mode and storage format
var amrnbPackedSize = [16]int{
	13, 14, 16, 18, 20, 21, 27, 32, 6 /*SID*/, 1, 1, 1, 1, 1, 1, 1,
}
var amrwbPackedSize = [16]int{
	18, 24, 33, 37, 41, 47, 51, 59, 61, 6, 1, 1, 1, 1, 1, 1,
}

// bandwidth efficient mode bits for each mode
var amrnbFrameBit = [16]int{
	95, 103, 118, 134, 148, 159, 204, 244, 39 /*SID*/, 0, 0, 0, 0, 0, 0, 0,
}

// 3GPP TS 26.201
var amrwbFrameBit = [16]int{
	132, 177, 253, 285, 317, 365, 397, 461, 477, 40 /*SID*/, 0, 0, 0, 0, 0, 0,
}

// AmrSplitToFrames transform data read from amr file into frames (toc+data for each frame)
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

// AmrFrameToRtpPayload only support one frame in every rtp packet
func AmrFrameToRtpPayload(frames [][]byte, isAmrwb bool, isOctetAlignMode bool) (rtpPayload [][]byte) {
	if isOctetAlignMode {
		return amrOctetModeFrameToRtpPayload(frames)
	} else {
		return amrBandwidthEfficientModeFrameToRtpPayload(frames, isAmrwb)
	}
}

// AmrRtpPayloadToFrame skip the header, count the frames in the index, then extract each frame data
func AmrRtpPayloadToFrame(payload []byte, isAmrwb bool, isOctetAlignMode bool) (frames [][]byte) {
	if isOctetAlignMode {
		return amrOctetModeRtpPayloadToFrame(payload, isAmrwb)
	} else {
		return amrBandwidthEfficientModeRtpPayloadToFrame(payload, isAmrwb)
	}

}

// Octet align mode
//  0 1 2 3 4 5 6 7
// +-+-+-+-+-+-+-+-+
// |F|  FT   |Q|P|P|
// +-+-+-+-+-+-+-+-+
// F (1 bit): see definition in Section 4.3.2.
// FT (4 bits, unsigned integer): see definition in Section 4.3.2.
// Q (1 bit): see definition in Section 4.3.2.
// P bits: padding bits, MUST be set to zero, and MUST be ignored on
func amrOctetModeRtpPayloadToFrame(payload []byte, isAmrwb bool) (frames [][]byte) {
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
		f := []byte{toc & 0x7c} // keep Q bit
		if frameEnd > speechData {
			f = append(f, payload[speechData:frameEnd]...)
		}
		speechData = frameEnd
		frames = append(frames, f)
	}
	return
}

func amrOctetModeFrameToRtpPayload(frame [][]byte) (rtpPayload [][]byte) {
	for _, p := range frame {
		toc := []byte{p[0] & 0x7f} // set F bit zero, so this is the last frame in the frame(only one)
		rp := append([]byte{0xf0}, toc...)
		rp = append(rp, p[1:]...)
		rtpPayload = append(rtpPayload, rp)
	}
	return
}

// In bandwidth-efficient mode, a ToC entry takes the following format:
//  0 1 2 3 4 5
// +-+-+-+-+-+-+
// |F|  FT   |Q|
// +-+-+-+-+-+-+
// F (1 bit): If set to 1, indicates that this frame is followed by another speech frame in this payload; if set to 0,
// indicates that this frame is the last frame in this payload.
// FT (4 bits): Frame type index, indicating either the AMR or AMR-WB speech coding mode or comfort noise (SID) mode
// of the corresponding frame carried in this payload.
// Q (1 bit): Frame quality indicator. If set to 0, indicates the corresponding frame is severely damaged, and
// the receiver should set the RX_TYPE (see [6]) to either SPEECH_BAD or SID_BAD depending on the frame type (FT).
func amrBandwidthEfficientModeRtpPayloadToFrame(payload []byte, isAmrwb bool) (frames [][]byte) {
	pl := len(payload)
	if pl <= 2 {
		return
	}
	var frameTypes []uint16
	payloadBitLen := pl * 8
	frameStartBit := 4       // skip CMR
	var val0xff uint8 = 0xff // keep shift result is of size uint8
	for frameStartBit < payloadBitLen-8 /* read uint16 wide */ {
		// count the frames in toc
		startByte := frameStartBit / 8
		shift := 16 - 6 - (frameStartBit & 0x07)
		twoBytes := binary.BigEndian.Uint16(payload[startByte:])
		toc := (twoBytes >> shift) & 0x3f
		ft := (toc >> 1) & 0x0f
		frameTypes = append(frameTypes, ft)
		frameStartBit += 6
		if toc&0x20 == 0 {
			// F bit is 0
			break
		}
	}
	// frame's count and type collected
	packedSize := amrnbPackedSize
	frameBitSize := amrnbFrameBit
	if isAmrwb {
		packedSize = amrwbPackedSize
		frameBitSize = amrwbFrameBit
	}
	for _, ft := range frameTypes {
		size := packedSize[ft] - 1
		nbFrameBit := frameBitSize[ft]
		if size == 0 {
			continue
		}
		if frameStartBit+nbFrameBit > payloadBitLen {
			logger.Errorln("amr bandwidth efficient frame exceeds its payload length")
			return
		}
		frame := make([]byte, size)
		// extract bytes from bits flow ...
		startByte := frameStartBit / 8
		remainingBits := nbFrameBit & 0x07
		leftShift := frameStartBit & 0x07
		rightShift := 8 - leftShift
		if leftShift == 0 {
			// nice, already aligned
			copy(frame[:size], payload[startByte:startByte+size])
			if remainingBits != 0 {
				frame[size-1] &= val0xff << (8 - remainingBits)
			}
		} else {
			// not aligned, extract from every two bytes into one
			mask1 := uint8(1<<rightShift - 1)
			mask2 := uint8(1<<leftShift - 1)
			for i := 0; i < size-1; i++ {
				b1, b2 := payload[startByte], payload[startByte+1]
				frame[i] = ((b1 & mask1) << leftShift) | ((b2 >> rightShift) & mask2)
				startByte += 1
			}
			// be careful to handle last byte when not aligned
			if remainingBits != 0 {
				lastByte := payload[startByte]
				if remainingBits <= rightShift {
					// last byte does not exceed boundary
					frame[size-1] = (lastByte << leftShift) & (val0xff << (8 - remainingBits))
				} else {
					// last byte sits across two bytes
					twoBytes := binary.BigEndian.Uint16(payload[startByte:])
					shift := 16 - leftShift - remainingBits
					frame[size-1] = uint8((twoBytes >> shift) & (1<<remainingBits - 1))
				}
			} else {
				b1, b2 := payload[startByte], payload[startByte+1]
				frame[size-1] = ((b1 & mask1) << leftShift) | ((b2 >> rightShift) & mask2)
			}
		}
		frameStartBit += nbFrameBit
		frames = append(frames, frame)
	}
	return
}

func amrBandwidthEfficientModeFrameToRtpPayload(frames [][]byte, isAmrwb bool) (rtpPayload [][]byte) {
	var hasExtraByte bool
	frameBitSize := amrnbFrameBit
	if isAmrwb {
		frameBitSize = amrwbFrameBit
	}
	for _, frm := range frames {
		mode := (frm[0] >> 3) & 0x0f
		frmBitSize := frameBitSize[mode]
		nbFrameBit := frmBitSize + 4 + 6 /* extra CMR and TOC */
		if nbFrameBit == 0 {
			continue
		}
		nbByte := nbFrameBit / 8
		if extraBits := frmBitSize & 0x07; extraBits != 0 {
			nbByte++
			// frame = toc(1 byte) + octet-align data (n bytes with last byte padding maybe)
			// rtp_payload = cmr_toc(10 bits) + bit-data (nbFrameBit)
			// rp has 2 more bits than frame in header, if the last byte of frame has at least 2 padding bits,
			// then the resulting rp has the same length(in bytes) of the frame, otherwise rp needs one more
			// extra byte to hold the data
			if extraBits == 7 || extraBits == 0 /* no padding */ {
				hasExtraByte = true
			}
		}
		if (hasExtraByte && len(frm) > nbByte-1) ||
			(!hasExtraByte && len(frm) > nbByte) {
			logger.Errorf("amr bandwidth efficient mode: frame is too large to transform: "+
				"frame len:%v, expected payload len:%v, hasExtraByte:%v", len(frm), nbByte, hasExtraByte)
			continue
		}
		rp := make([]byte, nbByte)
		cmrAndToc := 0xf0 | ((mode >> 1) & 0x07) // CMR=15, F-bit=0, 3-bit of toc
		remainToc := ((mode & 0x01) << 7) | 0x40 // 1-bit of toc, Q-bit=1
		copy(rp[:2], []byte{cmrAndToc, remainToc})
		rightShift := 2 /* (4+6) mod 8 */
		leftShift := 6
		leftPart := rp[1] & 0xc0
		for i := 1; i < len(frm); i++ {
			b := frm[i]
			rp[i] = leftPart | (b >> rightShift)
			leftPart = (b & 0x03) << leftShift
		}
		if hasExtraByte {
			rp[nbByte-1] = leftPart & 0xc0
		}
		rtpPayload = append(rtpPayload, rp)
	}
	return
}
