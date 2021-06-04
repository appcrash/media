package server

import (
	"fmt"
	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/wernerd/GoRTP/src/net/rtp"
	"net"
	"runtime/debug"
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

// receive rtcp packet
func (session *MediaSession) receiveCtrlLoop() {
	ctrlReceiver := session.rtpSession.CreateCtrlEventChan()
	for {
		select {
		case eventArray := <- ctrlReceiver:
			for _,event := range eventArray {
				if event.EventType == rtp.RtcpBye {
					// peer send bye, notify data send/receive loop to stop
					session.sndCtrlC <- "stop"
					session.rcvCtrlC <- "stop"
					// and also terminate myself
					return
				}
			}
		}
	}
}

func (session *MediaSession) receivePacketLoop() {
	// Create and store the data receive channel.
	defer func() {
		if r := recover(); r != nil {
			fmt.Errorf("receivePacketLoop panic(recovered)")
			debug.PrintStack()
		}
	}()

	rtpSession := session.rtpSession
	dataReceiver := rtpSession.CreateDataReceiveChan()

outLoop:
	for {
		select {
		case rp := <-dataReceiver:
			data := rp.Payload()
			// send received data to all sinkers, then free the packet
			//for _,s := range session.sinkerList {
			//	if shouldContinue := s.HandleData(session,data); !shouldContinue {
			//		break
			//	}
			//}
			session.sink.HandleData(session,data)
			rp.FreePacket()
		case cmd := <-session.rcvCtrlC:
			if cmd == "stop" {
				fmt.Println("stop local receive")
			}
			break outLoop
		}
	}

	// TODO: clean up here
}

func (session *MediaSession) sendPacketLoop() {
	var ts uint32 = 0
	ticker := time.NewTicker(20 * time.Millisecond)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("sendPacketLoop panic %v",r)
			debug.PrintStack()
		}
	}()

outLoop:
	for {
		select {
			case <- ticker.C:
				data,tsAdv := session.source.PullData(session)
				if data != nil {
					if session.rtpSession == nil {
						fmt.Println("########################")
					}
					session.rtpSession.SsrcStreamClose()
					packet := session.rtpSession.NewDataPacket(ts)
					packet.SetPayload(data)
					session.rtpSession.WriteData(packet)
					packet.FreePacket()
					ts += tsAdv
				}
			case cmd := <-session.sndCtrlC:
				if cmd == "stop" {
					break outLoop
				}
		}
	}

	// TODO: clean up here

}

func (session *MediaSession) StartSession() {
	session.rtpSession.StartSession()
	go session.receiveCtrlLoop()
	go session.receivePacketLoop()
	go session.sendPacketLoop()
}

func (session *MediaSession) AddRemote(ip string, port int) {
	ipaddr := net.ParseIP(ip)

	//println("peer ip port ",ipaddr,port)
	//println("session is ", session)
	session.rtpSession.AddRemote(&rtp.Address{
		IPAddr:   ipaddr,
		DataPort: port,
		CtrlPort: 1 + port,
		Zone:     "",
	})
}

