package codec

import (
	"encoding/binary"
)

const (
	// The NAL unit type octet has the following format:
	//+---------------+
	//|0|1|2|3|4|5|6|7|
	//+-+-+-+-+-+-+-+-+
	//|F|NRI| Type    |
	//+---------------+

	nalTypeStapa uint8 = 24
	nalTypeFua   uint8 = 28
	nalTypeSps   uint8 = 7
	nalTypePps   uint8 = 8

	bitmaskNalType uint8 = 0x1f
	bitmaskRefIdc  uint8 = 0x60
	bitmaskFuStart uint8 = 0x80
	bitmaskFuEnd   uint8 = 0x40
)

var annexbNalStartCode = []byte{0x00, 0x00, 0x00, 0x01}

func analyzeH264(payload []byte) {
	logger.Infof("payload len is %v", len(payload))
	nals := extractNals(payload)
	for _, n := range nals {
		printNal(n)
	}
}

// extractNals splits payload by start code
func extractNals(payload []byte) (nals [][]byte) {
	// start code can be of 4 bytes: 0x00,0x00,0x00,0x01 (sps,pps,first slice)
	// or 3 bytes: 0x00,0x00,0x01
	zeros := 0
	prevStart := 0
	totalLen := len(payload)
	for i, b := range payload {
		switch b {
		case 0x00:
			zeros++
			continue
		case 0x01:
			if zeros == 2 || zeros == 3 {
				// found a start code
				if prevStart != 0 {
					nal := payload[prevStart : i-zeros]
					nals = append(nals, nal)
				}
				prevStart = i + 1
				if prevStart >= totalLen {
					return
				}
			}
		}
		zeros = 0
	}
	if totalLen > prevStart {
		nals = append(nals, payload[prevStart:])
	}
	return
}

func printNal(nal []byte) {

	nalType := nal[0] & bitmaskNalType
	nalRefIdc := nal[0] & bitmaskRefIdc
	switch nalType {
	case nalTypeSps:
		logger.Infof("type sps")
	case nalTypePps:
		logger.Infof("type pps")
	default:
		logger.Infof("type: %d", nalType)
	}
	logger.Infof("nal len is %d,refIdc is %d", len(nal), nalRefIdc)
}

func makeStapA(nals ...[]byte) (rtpPayload []byte) {
	// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//|                     RTP Header                                |
	//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//| STAP-A NAL HDR|        NALU 1 Size            |  NALU 1 HDR   |
	//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//|             NALU 1 Data                                       |
	//:                                                               :
	//+               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//|               |     NALU 2 Size               | NALU 2 HDR    |
	//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//|                        NALU 2 Data                            |
	//:                                                               :
	//|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//|                               :...OPTIONAL RTP padding        |
	//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	// The value of NRI MUST be the maximum of all the NAL units carried
	// in the aggregation packet.
	rtpPayload = []byte{0x00} // header placeholder, set it later
	size := make([]byte, 2)
	var maxNri uint8
	for _, nal := range nals {
		binary.BigEndian.PutUint16(size, uint16(len(nal)))
		rtpPayload = append(rtpPayload, size...)
		rtpPayload = append(rtpPayload, nal...)
		nri := nal[0] & bitmaskRefIdc
		if maxNri < nri {
			maxNri = nri
		}
	}
	rtpPayload[0] = maxNri | nalTypeStapa
	return
}

func makeFuA(mtu int, nal []byte) (rtpPayload [][]byte) {
	// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//| FU indicator |  FU header     |                               |
	//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+                               |
	//|                                                               |
	//|                    FU payload                                 |
	//|                                                               |
	//|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//|                               :...OPTIONAL RTP padding        |
	//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	// The FU header has the following format:
	//+---------------+
	//|0|1|2|3|4|5|6|7|
	//+-+-+-+-+-+-+-+-+
	//|S|E|R| Type    |
	//+---------------+
	nri := nal[0] & bitmaskRefIdc
	nalType := nal[0] & bitmaskNalType
	indicator := nri | nalTypeFua
	isFirstFragment := true
	startPtr := 1 // skip the nal header
	remainingPayloadSize := len(nal) - startPtr
	maxPayloadSize := mtu - 2 // 2 = fu indicator + fu header

	for remainingPayloadSize > 0 {
		fragmentPayloadSize := remainingPayloadSize
		if fragmentPayloadSize > maxPayloadSize {
			fragmentPayloadSize = maxPayloadSize
		}
		payload := make([]byte, fragmentPayloadSize+2)
		payload[0] = indicator
		header := nalType
		if isFirstFragment {
			header |= 1 << 7 // set start bit
			isFirstFragment = false
		} else if fragmentPayloadSize == remainingPayloadSize {
			header |= 1 << 6 // set end bit
		}
		payload[1] = header
		copy(payload[2:], nal[startPtr:startPtr+fragmentPayloadSize])
		rtpPayload = append(rtpPayload, payload)
		startPtr += fragmentPayloadSize
	}
	return
}
