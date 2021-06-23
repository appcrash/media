package server

import (
	"fmt"
	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/wernerd/GoRTP/src/net/rtp"
	"net"
	"os"
	"runtime/debug"
	"time"
)

type MediaSession struct {
	sessionId  string
	rtpSession *rtp.Session
	payloadType uint32

	sndCtrlC   chan string
	rcvCtrlC   chan string

	source []Source
	sink []Sink
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


func createSession(localPort int,payloadType uint32) *MediaSession{
	tpLocal, _ := rtp.NewTransportUDP(local, localPort,"")
	session := rtp.NewSession(tpLocal, tpLocal)
	strLocalIdx, _ := session.NewSsrcStreamOut(&rtp.Address{
		IPAddr:   local.IP,
		DataPort: localPort,
		CtrlPort: 1 + localPort,
		Zone : "",
	}, 0, 0)
	session.SsrcStreamOutForIndex(strLocalIdx).SetPayloadType(byte(payloadType))

	ms := MediaSession{
			sessionId:  uuid.New().String(),
			rtpSession: session,
			payloadType: payloadType,
			sndCtrlC:   make(chan string,2),  // use buffered version avoiding deadlock
			rcvCtrlC:   make(chan string,2),
	}
	sessionMap.Set(ms.sessionId,&ms)
	return &ms
}

func (session *MediaSession) GetSessionId() string {
	return session.sessionId
}

func (session *MediaSession) GetPayloadType() uint32 {
	return session.payloadType
}

func (session *MediaSession) GetSource() []Source {
	return session.source
}

func (session *MediaSession) GetSink() []Sink {
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
					fmt.Println("rtp peer says bye")
					session.sndCtrlC <- "stop"
					session.rcvCtrlC <- "stop"
					fmt.Println("sent stop cmd to send/recv loops")
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

			//fmt.Printf("data len is %v",len(data))
			// push received data to all sinkers, then free the packet
			for _,s := range session.sink {
				if shouldContinue := s.HandleData(session,data); !shouldContinue {
					break
				}
			}
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

	file, _ := os.Create("/home/yh/Downloads/out.amr")
	defer file.Close()
	file.Write([]byte("#!AMR\n"))

outLoop:
	for {
		select {
			case <- ticker.C:
				var payload,data []byte
				var tsDelta uint32

				// pull data from all sources
				for _,source := range session.source {
					data,tsDelta = source.PullData(session,payload,tsDelta)
				}
				if data != nil {
					if session.rtpSession == nil {
						fmt.Println("########################")
					}
					//session.rtpSession.SsrcStreamClose()
					//toc := data[1] & 0x7f
					//file.Write([]byte{toc})
					//file.Write(data[2:])

					packet := session.rtpSession.NewDataPacket(ts)
					packet.SetPayload(data)

					session.rtpSession.WriteData(packet)
					packet.FreePacket()
					ts += tsDelta
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
	session.rtpSession.RtcpSessionBandwidth = 5000
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

