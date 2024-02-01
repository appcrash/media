package codec

import (
	"encoding/binary"
)

// shamelessly copied from ffmpeg(amr.c)  :)
// packedSize = frameSize + 1byte(toc), used in octet-align mode and storage format
//sean:ffmpeg rtpdec_amr.c give the array for speech data.there is different in wb and ft=10 (5)
//but here is 1 (should be 6 ?) but it does not matter because ft=10 is not useful class
var amrnbPackedSize = [16]int{
	13, 14, 16, 18, 20, 21, 27, 32, 6 /*SID*/, 1, 1, 1, 1, 1, 1, 1,
}
var amrwbPackedSize = [16]int{
	18, 24, 33, 37, 41, 47, 51, 59, 61, 6, 1, 1, 1, 1, 1, 1,
}

// bandwidth efficient mode bits for each mode
//sean:rfc4867 page32:
//compare with octetAlign and bandwidth-efficient when ft=5 indicate array as follow for speech data
//or consult from rfc3867 page8 table1 for speech data
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
	}
	return amrBandwidthEfficientModeFrameToRtpPayloads(frames, isAmrwb)
}

// AmrRtpPayloadToFrame skip the header, count the frames in the index, then extract each frame data
func AmrRtpPayloadToFrame(payload []byte, isAmrwb bool, isOctetAlignMode bool) (frames [][]byte) {
	if isOctetAlignMode {
		return amrOctetModeRtpPayloadToFrame(payload, isAmrwb)
	} else {
		return amrBandwidthEfficientModeRtpPayloadToFrames(payload, isAmrwb)
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

// FromInstanceC bandwidth-efficient mode, a ToC entry takes the following format:
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
	//cmr:=payload[0:4]
	var val0xff uint8 = 0xff // keep shift result is of size uint8
	for frameStartBit < payloadBitLen-8 /* read uint16 wide */ {
		// count the frames in toc
		startByte := frameStartBit / 8
		//shift := 16 - 6 - (frameStartBit & 0x07)
		shift := 16 - 10 - (frameStartBit & 0x07)
		twoBytes := binary.BigEndian.Uint16(payload[startByte:])
		toc := (twoBytes >> shift) & 0x3f
		ft := (toc >> 1) & 0x0f
		//logger.Info("ft=",ft,";pl=",pl)
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
				startByte++
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


//edit by sean.only support single channel and single frame or mul-frames in one packet
//RFC4867+3GPP TS 26.201
func amrBandwidthEfficientModeRtpPayloadToFrames(payload []byte, isAmrwb bool) (frames [][]byte) {
	pl := len(payload)
	if pl <= 2 {
		logger.Errorln("the length of rtp payload is invalid.len=",pl)
		return
	}

	var tocs [] uint8

	start:=0
	data:=binary.BigEndian.Uint16(payload[start:start+2])

	leftBits:=0

	//handle first toc 4+6+6
	toc:=uint8((data>>6)&0x3f)
	tocs=append(tocs,toc)
	leftBits=6 //maybe speech data
	toc_f:=toc & 0x20
	if toc_f==1{
		toc=uint8(data)&0x3f
		tocs=append(tocs,toc)
		toc_f=toc & 0x20
		leftBits=0
	}

	//handle other tocs
	var last_data uint16
	u16Count:=1

	for toc_f==1{
		start+=2
		last_data=data
		if start>pl{
			logger.Errorln("the amr header of rtp payload is invalid.")
			return
		}
		data=binary.BigEndian.Uint16(payload[start:start+2])
		switch leftBits {
		case 0: //6+6+4
			toc=uint8(data>>10)&0x3f
			tocs=append(tocs,toc)
			if toc&0x20==1{
				toc=uint8(data>>4)&0x3f
				tocs=append(tocs,toc)
			}
			leftBits=4
			break
		case 4: //2+6+6+2
			high4Bits:=uint8(last_data)&0x0f
			low2Bits:=uint8(data>>14)&0x03
			toc=(high4Bits<<2)+low2Bits
			tocs=append(tocs,toc)
			if toc&0x20==1{
				toc=uint8(data>>8)&0x3f
				tocs=append(tocs,toc)
				if toc&0x20==1{
					toc=uint8(data>>2)&0x3f
					tocs=append(tocs,toc)
				}
			}//if
			leftBits=2
			break
		case 2: //4+6+6
			high2Bits:=uint8(last_data)&0x03
			low4Bits:=uint8(data>>12)&0x0f
			toc=(high2Bits<<4)+low4Bits
			tocs=append(tocs,toc)
			if toc&0x20==1{
				toc=uint8(data>>6)&0x3f
				tocs=append(tocs,toc)
				if toc&0x20==1{
					toc=uint8(data)&0x3f
					tocs=append(tocs,toc)
				}
			}
			leftBits=0
			break
		default:
			logger.Errorln("The data of payload is invalid.it cant parse toc field")
			break
		}//switch

		u16Count++
		if toc&0x20==0{
			break
		}
	}//for

	//logger.Info("-----------------tocs len=",len(tocs),";pl=",pl)

	//speech data
	frameBitSize := amrnbFrameBit
	if isAmrwb {
		frameBitSize = amrwbFrameBit
	}
	startBytes:=u16Count*2
	for _,_toc:=range tocs{
		ft:=(_toc>>1)&0x0f
		bitFrame:=frameBitSize[ft] //size of speech data
		frame:=[]byte{(_toc<<2)&0x7c} //keep the same as octec-align mode

		byteFrame:=bitFrame/8
		remainBit:=bitFrame%8 //equal &0x07

		if startBytes+byteFrame>pl{
			logger.Errorln("amr bandwidth efficient mode frame exceeds its payload length")
			return
		}
		//logger.Info("-------toc=",_toc,";ft=",ft)
		if leftBits==0{ //aligned byte, perfect case
			frame=append(frame,payload[startBytes:startBytes+byteFrame]...)
			if remainBit!=0{ //last bits must be in the begin of byte
				leftMoveBits:=8-remainBit
				//remainFilter:=uint8(math.Pow(2,float64(remainBit))-1)<<leftMoveBits
				remainFilter:=uint8(((1<<remainBit)-1)<<leftMoveBits)
				remainData:=payload[startBytes+byteFrame-1]//end of payload
				lastByte:=remainData&remainFilter
				frame=append(frame,lastByte)
				leftBits=8-remainBit //for next frame
			}
		}else{
			cur:=startBytes-1 //because the byte includes toc and speech data or different frames
			//logger.Info("-------cur=",cur,";totalBits=",bitFrame,";leftBits=",leftBits)
			leftShift:=leftBits
			rightShift:=8-leftShift
			mask1 := uint8(1<<rightShift - 1)
			mask2 := uint8(1<<leftShift - 1)
			for bitFrame>0{
				if bitFrame>=8{
					highData:=(payload[cur] & mask2)<<(rightShift)
					lowData:=(payload[cur+1]>>leftShift) & mask1
					_data:=highData | lowData
					frame=append(frame,_data)
				}else{
					if bitFrame!=remainBit{
						logger.Errorln("amr bandwidth efficient frame is invalid.")
						return
					}
					//logger.Info("-------cur=",cur,";remainBits=",bitFrame,";leftBits=",leftBits)
					if bitFrame<=leftBits{//last byte is in one byte
						usedBits:=8-leftBits
						_data:=(payload[cur]<<usedBits)&(0xff<<(8-bitFrame))
						frame=append(frame,_data)
						leftBits-=bitFrame
					}else{//last byte crosses tow bytes
						m1:=uint8(1<<bitFrame-1)
						highData:=(payload[cur-1] & mask2)<<(rightShift)
						lowData:=(payload[cur]>>leftBits)&(m1<<(8-leftBits-bitFrame))
						_data:=highData | lowData
						frame=append(frame,_data)
						leftBits=8-(bitFrame-leftBits) //for next frame
					}

					//if leftBits+bitFrame>8{//one byte,bitFrame must be in the last byte
					//	usedBits:=8-leftBits
					//	_data:=(payload[cur]<<usedBits)&(0xff<<(8-bitFrame))
					//	frame=append(frame,_data)
					//	leftBits-=remainBit //for next frame
					//}else{//last byte crosses two byte
					//	m1:=uint8(1<<bitFrame-1)
					//	highData:=(payload[cur-1] & mask2)<<(rightShift)
					//	lowData:=(payload[cur]>>leftBits)&(m1<<(8-leftBits-bitFrame))
					//	_data:=highData | lowData
					//	frame=append(frame,_data)
					//	leftBits=8-remainBit //for next frame
					//}
				}
				bitFrame-=8
				cur++ //update payload current
			}//for
		}//else
		//logger.Info("-----------------frame size=",len(frame))
		frames=append(frames,frame)
		startBytes+=byteFrame//update startByte
	}//for

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


//author:sean consult RFC4867+3GPP TS 26.201
//input frames:toc+speech_data with octetAlign arrays.output:cmr+toc+speech_data for bandwidth-efficient
func amrBandwidthEfficientModeFrameToRtpPayloads(frames [][]byte, isAmrwb bool) (rtpPayload [][]byte){
	frameBitSize := amrnbFrameBit
	if isAmrwb {
		frameBitSize = amrwbFrameBit
	}

	for _,frm:= range frames{
		mode := (frm[0] >> 3) & 0x0f
		frmBitSize := frameBitSize[mode]
		totalBit := frmBitSize + 4 + 6 /* extra CMR and TOC */
		if totalBit<=10{ //speech_data=nullptr //consider especial situation silence voice
			logger.Warnln("The data of frame is invalid this time and will ignore")
			continue
		}

		needBytes:=totalBit/8
		extraBits:=totalBit%8

		if extraBits>0{
			needBytes++ //for extra data
		}

		//4+6=8+2 left:2  right:6
		rp:=make([]byte,needBytes)
		cmrAndToc := 0xf0 | ((mode >> 1) & 0x07) // CMR=15, F-bit=0, 3-bit of toc
		//cmrAndToc := (mode<<4)  | ((mode >> 1) & 0x07) // CMR=ft, F-bit=0, 3-bit of toc
		remainToc := ((mode & 0x01) << 7) | 0x40 // 1-bit of toc, Q-bit=1
		copy(rp[:2], []byte{cmrAndToc, remainToc})
		rightShift := 2 /* (4+6) mod 8 */
		leftShift := 6
		leftPart := rp[1] & 0xc0 //

		//load speech_data
		for i:=1;i<len(frm);i++{
			d:=frm[i]
			rp[i]=leftPart | (d>>rightShift)
			leftPart=(d&0x03) << leftShift
		}

		//handle last byte of frm
		pad:=(len(frm)-1)*8-frmBitSize
		if pad<2{ //extra byte to hold leftPart but less than two bits
			rp[needBytes-1]=leftPart & 0xc0
		}

		rtpPayload = append(rtpPayload, rp)

	}//for frames


	return
}





