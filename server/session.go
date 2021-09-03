package server

import (
	"fmt"
	"github.com/appcrash/GoRTP/rtp"
	"github.com/appcrash/media/codec"
	"github.com/appcrash/media/server/event"
	"github.com/appcrash/media/server/rpc"
	"github.com/google/uuid"
	"net"
	"runtime/debug"
	"sync"
	"time"
)

var sessionMap = sync.Map{}

type MediaSession struct {
	sessionId         string
	rtpSession        *rtp.Session
	payloadTypeNumber uint32
	payloadCodec      rpc.CodecType
	graphDesc         string

	sndCtrlC     chan string
	rcvCtrlC     chan string
	rcvRtcpCtrlC chan string

	source []Source
	sink   []Sink
	graph  *event.EventGraph
}

func profileOfCodec(c rpc.CodecType) (profile string) {
	switch c {
	case rpc.CodecType_PCM_ALAW:
		profile = "PCMA"
	case rpc.CodecType_AMRNB:
		profile = "AMR"
	case rpc.CodecType_AMRWB:
		profile = "AMR-WB"
	}
	return
}

func createSession(ipAddr *net.IPAddr, localPort int, mediaParam *rpc.MediaParam) *MediaSession {
	tpLocal, _ := rtp.NewTransportUDP(ipAddr, localPort, "")
	session := rtp.NewSession(tpLocal, tpLocal)
	strLocalIdx, _ := session.NewSsrcStreamOut(&rtp.Address{
		IPAddr:   ipAddr.IP,
		DataPort: localPort,
		CtrlPort: 1 + localPort,
		Zone:     "",
	}, 0, 0)

	if profile := profileOfCodec(mediaParam.GetPayloadCodecType()); profile != "" {
		session.SsrcStreamOutForIndex(strLocalIdx).SetProfile(profile, byte(mediaParam.GetPayloadDynamicType()))
	} else {
		session.CloseSession()
		return nil
	}

	ms := MediaSession{
		sessionId:         uuid.New().String(),
		rtpSession:        session,
		payloadTypeNumber: mediaParam.GetPayloadDynamicType(),
		payloadCodec:      mediaParam.GetPayloadCodecType(),
		graphDesc:         mediaParam.GetGraphDesc(),

		// use buffered version to avoid deadlock
		sndCtrlC:     make(chan string, 2),
		rcvCtrlC:     make(chan string, 2),
		rcvRtcpCtrlC: make(chan string, 2),
	}
	sessionMap.Store(ms.sessionId, &ms)
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

func (session *MediaSession) GetGraphDescription() string {
	return session.graphDesc
}

func (session *MediaSession) GetSource() []Source {
	return session.source
}

func (session *MediaSession) GetSink() []Sink {
	return session.sink
}

func (session *MediaSession) GetEventGraph() *event.EventGraph {
	return session.graph
}

// receive rtcp packet
func (session *MediaSession) receiveCtrlLoop() {
	rtcpReceiver := session.rtpSession.CreateCtrlEventChan()
	ctrlC := session.rcvRtcpCtrlC

	for {
		select {
		case eventArray, more := <-rtcpReceiver:
			if !more {
				// RTP stack closed rtcp channel, just return
				return
			}
			for _, event := range eventArray {
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
		case rp, more := <-dataReceiver:
			var shouldContinue bool
			if !more {
				// RTP stack closed this channel, so stop receiving anymore
				return
			}
			data := rp.Payload()

			// push received data to all sinks, then free the packet
			for _, s := range session.sink {
				data, shouldContinue = s.HandleData(session, data)
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
					break outLoop
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

	ticker.Stop()
}

func (session *MediaSession) Start() {
	session.rtpSession.StartSession()
	go session.receiveCtrlLoop()
	go session.receivePacketLoop()
	go session.sendPacketLoop()
}

func (session *MediaSession) Stop() {
	nbStopped := 0
	for nbStopped < 3 {
		select {
		case session.sndCtrlC <- "stop":
			session.sndCtrlC = nil
			nbStopped++
		case session.rcvCtrlC <- "stop":
			session.rcvCtrlC = nil
			nbStopped++
		case session.rcvRtcpCtrlC <- "stop":
			session.rcvRtcpCtrlC = nil
			nbStopped++
		case <-time.After(5 * time.Second):
			// TODO: how to avoid memory leak
			fmt.Errorf("session(%v) stops timeout", session.sessionId)
		}
	}
	session.rtpSession.CloseSession()
	sessionMap.Delete(session.GetSessionId())
}

// AddNode add an event node to server-wide event graph
func (session *MediaSession) AddNode(node event.Node) {
	session.graph.AddNode(node)
}

// AddRemote add rtp peer to communicate
func (session *MediaSession) AddRemote(ip string, port int) {
	ipaddr := net.ParseIP(ip)

	session.rtpSession.AddRemote(&rtp.Address{
		IPAddr:   ipaddr,
		DataPort: port,
		CtrlPort: 1 + port,
		Zone:     "",
	})
}
