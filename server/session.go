package server

import (
	"fmt"
	"github.com/appcrash/media/codec"
	"github.com/appcrash/media/server/rpc"
	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/wernerd/GoRTP/src/net/rtp"
	"net"
	"runtime/debug"
	"time"
)

type MediaSession struct {
	sessionId         string
	rtpSession        *rtp.Session
	payloadTypeNumber uint32
	payloadCodec      rpc.CodecType
	recordFile        string

	sndCtrlC     chan string
	rcvCtrlC     chan string
	rcvRtcpCtrlC chan string

	source []Source
	sink   []Sink
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

func createSession(localPort int, mediaParam *rpc.MediaParam) *MediaSession {
	tpLocal, _ := rtp.NewTransportUDP(local, localPort, "")
	session := rtp.NewSession(tpLocal, tpLocal)
	strLocalIdx, _ := session.NewSsrcStreamOut(&rtp.Address{
		IPAddr:   local.IP,
		DataPort: localPort,
		CtrlPort: 1 + localPort,
		Zone:     "",
	}, 0, 0)
	session.SsrcStreamOutForIndex(strLocalIdx).SetPayloadType(byte(mediaParam.GetPayloadRtpType()))

	ms := MediaSession{
		sessionId:   uuid.New().String(),
		rtpSession:  session,
		payloadTypeNumber: mediaParam.GetPayloadRtpType(),
		payloadCodec: mediaParam.GetRecordType(),
		recordFile: mediaParam.GetRecordFile(),

		// use buffered version avoiding deadlock
		sndCtrlC:     make(chan string, 2),
		rcvCtrlC:     make(chan string, 2),
		rcvRtcpCtrlC: make(chan string, 2),
	}
	sessionMap.Set(ms.sessionId, &ms)
	return &ms
}

func (session *MediaSession) GetSessionId() string {
	return session.sessionId
}

func (session *MediaSession) GetPayloadType() uint32 {
	return session.payloadTypeNumber
}

func (session *MediaSession) GetCodecType() rpc.CodecType {
	return session.payloadCodec
}

func (session *MediaSession) GetRecordFile() string {
	return session.recordFile
}

func (session *MediaSession) GetSource() []Source {
	return session.source
}

func (session *MediaSession) GetSink() []Sink {
	return session.sink
}

// receive rtcp packet
func (session *MediaSession) receiveCtrlLoop() {
	rtcpReceiver := session.rtpSession.CreateCtrlEventChan()
	ctrlC := session.rcvRtcpCtrlC

	for {
		select {
		case eventArray := <-rtcpReceiver:
			for _, event := range eventArray {
				//if packetType, exist := codec.RtcpPacketTypeMap[event.EventType]; exist {
				//	fmt.Printf("rctp: [%v] --- reason %v \n", packetType, event.Reason)
				//}

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
		case msg := <-ctrlC:
			if msg == "stop" {
				return
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
			var shouldContinue bool
			data := rp.Payload()

			//fmt.Printf("data len is %v",len(data))
			// push received data to all sinkers, then free the packet
			for _, s := range session.sink {
				data,shouldContinue = s.HandleData(session, data)
				if !shouldContinue {
					break
				}
			}
			rp.FreePacket()
		case cmd := <-session.rcvCtrlC:
			if cmd == "stop" {
				fmt.Println("stop local receive")
				break outLoop
			}
		}
	}

	// TODO: clean up here
}

func (session *MediaSession) sendPacketLoop() {
	var ts uint32 = 0
	timeStep := codec.GetCodecTimeStep(session.payloadCodec)
	ticker := time.NewTicker(time.Duration(timeStep) * time.Millisecond)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("sendPacketLoop panic %v", r)
			debug.PrintStack()
		}
	}()

outLoop:
	for {
		select {
		case <-ticker.C:
			var data []byte
			var tsDelta uint32

			// pull data from all sources
			for _, source := range session.source {
				data, tsDelta = source.PullData(session, data, tsDelta)
			}
			if data != nil {
				if session.rtpSession == nil {
					fmt.Println("########################")
				}

				packet := session.rtpSession.NewDataPacket(ts)
				packet.SetPayload(data)

				session.rtpSession.WriteData(packet)
				packet.FreePacket()
				ts += tsDelta
			}
		case cmd := <-session.sndCtrlC:
			if cmd == "stop" {
				fmt.Println("stop local send")
				break outLoop
			}
		}
	}

	// TODO: clean up here

}

func (session *MediaSession) Start() {
	session.rtpSession.RtcpSessionBandwidth = 5000
	session.rtpSession.StartSession()
	go session.receiveCtrlLoop()
	go session.receivePacketLoop()
	go session.sendPacketLoop()
}

func (session *MediaSession) Stop() {
	// nonblock send message to loops
	select {
	case session.sndCtrlC <- "stop":
	default:
		fmt.Errorf("send stop to session send loop failed")
	}
	select {
	case session.rcvCtrlC <- "stop":
	default:
		fmt.Errorf("send stop to session receive loop failed")
	}
	select {
	case session.rcvRtcpCtrlC <- "stop":
	default:
		fmt.Errorf("send stop to session RTCP receive loop failed")
	}
	session.rtpSession.CloseSession()
	sessionMap.Remove(session.GetSessionId())
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
