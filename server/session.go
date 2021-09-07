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
	"sync/atomic"
	"time"
)

var sessionMap = sync.Map{}

type sessionFinalizer func()

type MediaSession struct {
	sessionId         string
	rtpSession        *rtp.Session
	rtpPort           uint16
	payloadTypeNumber uint32
	payloadCodec      rpc.CodecType
	graphDesc         string
	stopCounter       uint32

	sndCtrlC     chan string
	rcvCtrlC     chan string
	rcvRtcpCtrlC chan string
	doneC        chan string // notify this channel when loop is done

	source    []Source
	sink      []Sink
	finalizer sessionFinalizer
	graph     *event.EventGraph
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
		rtpPort:           uint16(localPort),
		payloadTypeNumber: mediaParam.GetPayloadDynamicType(),
		payloadCodec:      mediaParam.GetPayloadCodecType(),
		graphDesc:         mediaParam.GetGraphDesc(),

		// use buffered version to avoid deadlock
		sndCtrlC:     make(chan string, 2),
		rcvCtrlC:     make(chan string, 2),
		rcvRtcpCtrlC: make(chan string, 2),
		doneC:        make(chan string, 3),
	}
	sessionMap.Store(ms.sessionId, &ms)
	return &ms
}

func (s *MediaSession) GetSessionId() string {
	return s.sessionId
}

func (s *MediaSession) GetPayloadType() uint32 {
	return s.payloadTypeNumber
}

func (s *MediaSession) GetCodecType() rpc.CodecType {
	return s.payloadCodec
}

func (s *MediaSession) GetGraphDescription() string {
	return s.graphDesc
}

func (s *MediaSession) GetSource() []Source {
	return s.source
}

func (s *MediaSession) GetSink() []Sink {
	return s.sink
}

func (s *MediaSession) GetEventGraph() *event.EventGraph {
	return s.graph
}

// receive rtcp packet
func (s *MediaSession) receiveCtrlLoop() {
	rtcpReceiver := s.rtpSession.CreateCtrlEventChan()
	ctrlC := s.rcvRtcpCtrlC

	defer func() {
		s.doneC <- "done"
	}()

	for {
		select {
		case eventArray, more := <-rtcpReceiver:
			if !more {
				// RTP stack closed rtcp channel, just return
				return
			}
			for _, evt := range eventArray {
				if evt.EventType == rtp.RtcpBye {
					// peer send bye, notify data send/receive loop to stop
					fmt.Println("rtp peer says bye")
					s.Stop()
					//s.sndCtrlC <- "stop"
					//s.rcvCtrlC <- "stop"
					//fmt.Println("sent stop cmd to send/recv loops")
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

func (s *MediaSession) receivePacketLoop() {
	// Create and store the data receive channel.
	defer func() {
		if r := recover(); r != nil {
			fmt.Errorf("receivePacketLoop panic(recovered)\n")
			debug.PrintStack()
		}
	}()

	defer func() {
		s.doneC <- "done"
	}()

	rtpSession := s.rtpSession
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
			for _, sk := range s.sink {
				data, shouldContinue = sk.HandleData(s, data)
				if !shouldContinue {
					break
				}
			}
			rp.FreePacket()
		case cmd := <-s.rcvCtrlC:
			if cmd == "stop" {
				fmt.Println("stop local receive")
				break outLoop
			}
		}
	}

}

func (s *MediaSession) sendPacketLoop() {
	var ts uint32 = 0
	timeStep := codec.GetCodecTimeStep(s.payloadCodec)
	ticker := time.NewTicker(time.Duration(timeStep) * time.Millisecond)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("sendPacketLoop panic %v", r)
			debug.PrintStack()
		}
	}()

	defer func() {
		s.doneC <- "done"
	}()

outLoop:
	for {
		select {
		case <-ticker.C:
			var data []byte
			var tsDelta uint32

			// pull data from all sources
			for _, source := range s.source {
				data, tsDelta = source.PullData(s, data, tsDelta)
			}
			if data != nil {
				if s.rtpSession == nil {
					break outLoop
				}
				packet := s.rtpSession.NewDataPacket(ts)
				packet.SetPayload(data)
				s.rtpSession.WriteData(packet)
				packet.FreePacket()
				ts += tsDelta
			}
		case cmd := <-s.sndCtrlC:
			if cmd == "stop" {
				fmt.Println("stop local send")
				break outLoop
			}
		}
	}

	ticker.Stop()
}

func (s *MediaSession) Start() (err error) {
	if err = s.rtpSession.StartSession(); err != nil {
		return
	}
	go s.receiveCtrlLoop()
	go s.receivePacketLoop()
	go s.sendPacketLoop()
	return
}

func (s *MediaSession) Stop() {
	if atomic.AddUint32(&s.stopCounter, 1) != 1 {
		// somebody has already called Stop()
		return
	}
	nbStopped := 0
	nbDone := 0
stop:
	for nbStopped < 3 {
		select {
		case s.sndCtrlC <- "stop":
			s.sndCtrlC = nil
			nbStopped++
		case s.rcvCtrlC <- "stop":
			s.rcvCtrlC = nil
			nbStopped++
		case s.rcvRtcpCtrlC <- "stop":
			s.rcvRtcpCtrlC = nil
			nbStopped++
		case <-time.After(2 * time.Second):
			// TODO: how to avoid memory leak
			fmt.Errorf("s(%v) stops timeout\n", s.sessionId)
			break stop
		}
	}

	// for debug purpose, check all loops are finished normally
done:
	for nbDone < 3 {
		select {
		case <-s.doneC:
			nbDone++
		case <-time.After(2 * time.Second):
			break done
		}
	}
	if nbDone != 3 {
		fmt.Errorf("s(%v) loops don't stop normally, finished number:%v\n", s.sessionId, nbDone)
	}
	s.rtpSession.CloseSession()
	if s.finalizer != nil {
		s.finalizer()
	}
	sessionMap.Delete(s.GetSessionId())
}

// AddNode add an event node to server-wide event graph
func (s *MediaSession) AddNode(node event.Node) {
	s.graph.AddNode(node)
}

// AddRemote add rtp peer to communicate
func (s *MediaSession) AddRemote(ip string, port int) (err error) {
	ipaddr := net.ParseIP(ip)

	_, err = s.rtpSession.AddRemote(&rtp.Address{
		IPAddr:   ipaddr,
		DataPort: port,
		CtrlPort: 1 + port,
		Zone:     "",
	})
	return
}
