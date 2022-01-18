package codec

import "encoding/binary"

const (
	nalTypeStapa = 24
	nalTypeFua   = 28
	nalTypeSps   = 7
	nalTypePps   = 8

	bitmaskNalType = 0x1f
	bitmaskRefIdc  = 0x60
	bitmaskFuStart = 0x80
	bitmaskFuEnd   = 0x40

	nalHeaderStapa = 0x78 // (nri: 0x11,type: 24)
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
	// The NAL unit type octet has the following format:
	//+---------------+
	//|0|1|2|3|4|5|6|7|
	//+-+-+-+-+-+-+-+-+
	//|F|NRI| Type    |
	//+---------------+
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
	rtpPayload = append(rtpPayload, nalHeaderStapa)
	header := make([]byte, 2)
	for _, nal := range nals {
		binary.BigEndian.PutUint16(header, uint16(len(nal)))
		rtpPayload = append(rtpPayload, header...)
		rtpPayload = append(rtpPayload, nal...)
	}
	return
}

func makeFuA(mtu int, nal []byte) (rtpPayload [][]byte) {

}
