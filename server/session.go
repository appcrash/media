package server

import (
	"fmt"
	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/wernerd/GoRTP/src/net/rtp"
	"net"
	"time"
)

type MediaSession struct {
	sessionId  string
	rtpSession *rtp.Session
	sndCtrlC   chan string
	rcvCtrlC   chan string

	source Source
	sink Sink
}

var sessionMap = cmap.New()

const PCMAPayLoadLength = 160

func WavPayload(wavData []byte) (rtpPayload []byte) {
	start := 0
	for i := 36; i < len(wavData); i++ {
		if string(wavData[i:i+4]) == "data" {
			start = i + 8
			break
		}
	}
	if start == 0 {
		fmt.Errorf("data ERROR")
		return rtpPayload
	}
	rtpPayload = wavData[start:]
	return rtpPayload
}


func createSession(localPort int) *MediaSession{
	tpLocal, _ := rtp.NewTransportUDP(local, localPort,"")
	session := rtp.NewSession(tpLocal, tpLocal)
	strLocalIdx, _ := session.NewSsrcStreamOut(&rtp.Address{
		IPAddr:   local.IP,
		DataPort: localPort,
		CtrlPort: 1 + localPort,
		Zone : "",
	}, 0, 0)
	session.SsrcStreamOutForIndex(strLocalIdx).SetPayloadType(8)

	ms := MediaSession{
			sessionId:  uuid.New().String(),
			rtpSession: session,
			sndCtrlC:   make(chan string,2),  // use buffered version avoiding deadlock
			rcvCtrlC:   make(chan string,2),
	}
	sessionMap.Set(ms.sessionId,&ms)
	return &ms
}

func (session *MediaSession) GetSessionId() string {
	return session.sessionId
}

func (session *MediaSession) GetSource() Source {
	return session.source
}

func (session *MediaSession) GetSink() Sink {
	return session.sink
}

func (session *MediaSession) receivePacketLoop() {
	// Create and store the data receive channel.
	rtpSession := session.rtpSession
	dataReceiver := rtpSession.CreateDataReceiveChan()
	var cnt int
	for {
		select {
		case rp := <-dataReceiver:
			if (cnt % 50) == 0 {
				println("Remote receiver got:", cnt, "packets")
			}
			data := rp.Payload()

			// send received data to all sinkers, then free the packet
			//for _,s := range session.sinkerList {
			//	if shouldContinue := s.HandleData(session,data); !shouldContinue {
			//		break
			//	}
			//}
			session.sink.HandleData(session,data)

			cnt++
			rp.FreePacket()
		case cmd := <-session.rcvCtrlC:
			if cmd == "stop" {
				println("stop local receive")
			}
			return
		}
	}
}

func (session *MediaSession) sendPacketLoop() {
	//sendBuffer := []byte{}
	//sendPtr := int(math.MaxInt32)
	//sendStep := int(0)
	var ts uint32 = 0

	ticker := time.NewTicker(20 * time.Millisecond)

	for {
		select {
			case <- ticker.C:
				// if send buffer is done, stop the ticker until send buffer activated again
				// otherwise forward the send point and feed the data to rtp
				//if sendPtr >= len(sendBuffer) {
				//	fmt.Println("sendPtr out of buffer range, stop ticker")
				//	ticker.Stop()
				//} else {
				//	nextPtr := sendPtr + sendStep
				//	if nextPtr > len(sendBuffer) {
				//		nextPtr = len(sendBuffer)
				//	}
				//	payload := sendBuffer[sendPtr : nextPtr]
				//	packet := session.rtpSession.NewDataPacket(ts)
				//	ts += uint32(sendStep)  // TODO: use real timestamp step
				//	packet.SetPayload(payload)
				//	session.rtpSession.WriteData(packet)
				//	packet.FreePacket()
				//	sendPtr = nextPtr
				//}
				data,tsAdv := session.source.PullData(session)
				if data != nil {
					packet := session.rtpSession.NewDataPacket(ts)
					packet.SetPayload(data)
					session.rtpSession.WriteData(packet)
					packet.FreePacket()
					ts += tsAdv
				}
			case cmd := <-session.sndCtrlC:
				//if cmd == "stop" {
				//	break
				//} else {
				//	if mediaData, err := ioutil.ReadFile(cmd); err == nil {
				//		fmt.Println("start playing file: %v",cmd)
				//		wavData := WavPayload(mediaData)
				//		sendBuffer = wavData
				//		sendPtr = 0
				//		sendStep = PCMAPayLoadLength
				//		ticker.Reset(20 * time.Millisecond)
				//	} else {
				//		fmt.Println("error: can not find file %v when change send buffer",cmd)
				//	}
				//}
				if cmd == "stop" {
					break
				}
		}
	}

}

func (session *MediaSession) StartSession() {
	session.rtpSession.StartSession()
	go session.receivePacketLoop()
	go session.sendPacketLoop()
}

func (session *MediaSession) AddRemote(ip string, port int) {
	ipaddr := net.ParseIP(ip)

	println("peer ip port ",ipaddr,port)
	println("session is ", session)
	session.rtpSession.AddRemote(&rtp.Address{
		IPAddr:   ipaddr,
		DataPort: port,
		CtrlPort: 1 + port,
		Zone:     "",
	})
}

