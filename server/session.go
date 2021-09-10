package server

import (
	"errors"
	"github.com/appcrash/GoRTP/rtp"
	"github.com/appcrash/media/codec"
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/event"
	"github.com/appcrash/media/server/rpc"
	"github.com/google/uuid"
	"net"
	"runtime/debug"
	"sync/atomic"
	"time"
)

type MediaSession struct {
	server            *MediaServer
	isStarted         bool
	sessionId         string
	localIp           *net.IPAddr
	localPort         int
	rtpPort           uint16
	rtpSession        *rtp.Session
	payloadTypeNumber uint32
	payloadCodec      rpc.CodecType
	stopCounter       uint32

	sndCtrlC     chan string
	rcvCtrlC     chan string
	rcvRtcpCtrlC chan string
	doneC        chan string // notify this channel when loop is done

	source   []Source
	sink     []Sink
	composer *comp.Composer
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

func (s *MediaSession) GetEventGraph() *event.Graph {
	return s.server.graph
}

func (s *MediaSession) GetController() comp.Controller {
	return s.composer.GetController()
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

func newSession(srv *MediaServer, mediaParam *rpc.CreateParam) (*MediaSession, error) {
	var localPort uint16
	if localPort = srv.getNextAvailableRtpPort(); localPort == 0 {
		return nil, errors.New("server runs out of port resource")
	}
	sid := uuid.New().String()
	gd := mediaParam.GetGraphDesc()
	composer := comp.NewSessionComposer(sid)
	if err := composer.ParseGraphDescription(gd); err != nil {
		logger.Errorln(err)
		return nil, errors.New("composer parse graph description failed")
	}

	session := MediaSession{
		server:            srv,
		sessionId:         sid,
		localIp:           srv.rtpServerIpAddr,
		localPort:         int(localPort),
		rtpPort:           localPort,
		payloadTypeNumber: mediaParam.GetPayloadDynamicType(),
		payloadCodec:      mediaParam.GetPayloadCodecType(),

		// use buffered version to avoid deadlock
		sndCtrlC:     make(chan string, 2),
		rcvCtrlC:     make(chan string, 2),
		rcvRtcpCtrlC: make(chan string, 2),
		doneC:        make(chan string, 3),

		composer: composer,
	}
	return &session, nil
}

func (s *MediaSession) setupGraph() error {
	// search any source or sink is interested in composer
	var ca []comp.ComposerAware
	for _, src := range s.source {
		if cs, ok := src.(comp.ComposerAware); ok {
			ca = append(ca, cs)
		}
	}
	for _, sink := range s.sink {
		if cs, ok := sink.(comp.ComposerAware); ok {
			ca = append(ca, cs)
		}
	}
	// call pre- and post- setup
	for _, cai := range ca {
		if err := cai.PreSetup(s.composer); err != nil {
			return err
		}
	}
	if err := s.composer.PrepareNodes(s.server.graph); err != nil {
		return err
	}
	for _, cai := range ca {
		if err := cai.PostSetup(s.composer); err != nil {
			return err
		}
	}
	return nil
}

// activate carry out actual work, such as listen on udp port, create rtp stream, create event node instances and
// add them to graph
func (s *MediaSession) activate() (err error) {
	var tpLocal *rtp.TransportUDP
	if err = s.setupGraph(); err != nil {
		return
	}
	if tpLocal, err = rtp.NewTransportUDP(s.localIp, s.localPort, ""); err != nil {
		return
	}
	s.rtpSession = rtp.NewSession(tpLocal, tpLocal)
	strLocalIdx, errStr := s.rtpSession.NewSsrcStreamOut(&rtp.Address{
		IPAddr:   s.localIp.IP,
		DataPort: s.localPort,
		CtrlPort: 1 + s.localPort,
		Zone:     "",
	}, 0, 0)
	if errStr != "" {
		return errors.New(string(errStr))
	}
	if profile := profileOfCodec(s.payloadCodec); profile != "" {
		s.rtpSession.SsrcStreamOutForIndex(strLocalIdx).SetProfile(profile, byte(s.payloadTypeNumber))
	} else {
		return errors.New("unsupported rtp payload profile")
	}
	return nil
}

// release all resources this session occupied
func (s *MediaSession) finalize() {
	if s.composer != nil {
		s.composer.ExitGraph()
	}
	if s.rtpSession != nil {
		s.rtpSession.CloseSession()
	}
	if s.rtpPort != 0 {
		s.server.reclaimRtpPort(s.rtpPort)
	}
	sessionMap.Delete(s.sessionId)
}

// receive rtcp packet
func (s *MediaSession) receiveCtrlLoop() {
	rtcpReceiver := s.rtpSession.CreateCtrlEventChan()
	ctrlC := s.rcvRtcpCtrlC

	defer func() {
		logger.Debugf("session:%v stop ctrl recv", s.GetSessionId())
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
					logger.Debugln("rtp peer says bye")
					s.Stop()
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
			logger.Fatalln("receivePacketLoop panic(recovered)")
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
				logger.Debugf("session:%v stop local receive", s.GetSessionId())
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
			logger.Fatalln("sendPacketLoop panic %v", r)
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
				_, _ = s.rtpSession.WriteData(packet)
				packet.FreePacket()
				ts += tsDelta
			}
		case cmd := <-s.sndCtrlC:
			if cmd == "stop" {
				logger.Debugf("session:%v stop local send", s.GetSessionId())
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
			logger.Errorf("s(%v) stops timeout", s.sessionId)
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
		logger.Errorf("s(%v) loops don't stop normally, finished number:%v", s.sessionId, nbDone)
	}

	s.finalize()
}

// AddNode add an event node to server-wide event graph
func (s *MediaSession) AddNode(node event.Node) {
	s.server.graph.AddNode(node)
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

// statically or dynamically add a channel to event graph
func (s *MediaSession) RegisterChannel(name string, ch chan<- *event.Event) {
	s.composer.RegisterChannel(name, ch)
}
