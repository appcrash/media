package server

import (
	"errors"
	"fmt"
	"github.com/appcrash/GoRTP/rtp"
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/event"
	"github.com/appcrash/media/server/rpc"
	"net"
	"strconv"
	"sync/atomic"
)

type SessionIdType uint32

var sessionIdCounter uint32 // fetch new id from here, use atomic increment

func (id SessionIdType) String() string {
	// math.MaxUint32 == 4294967295, max 10 zero ...
	return fmt.Sprintf("%010d", id)
}

func SessionIdFromString(s string) (SessionIdType, error) {
	id, err := strconv.ParseUint(s, 10, 32)
	return SessionIdType(id), err
}

func profileOfCodec(c rpc.CodecType) (profile string) {
	switch c {
	case rpc.CodecType_PCM_ALAW:
		profile = "PCMA"
	case rpc.CodecType_AMRNB:
		profile = "AMR"
	case rpc.CodecType_AMRWB:
		profile = "AMR-WB"
	case rpc.CodecType_H264:
		profile = "H264"
	case rpc.CodecType_TELEPHONE_EVENT_8K, rpc.CodecType_TELEPHONE_EVENT_16K:
		profile = "TELEPHONE-EVENT"
	case rpc.CodecType_EVS:
		profile = "EVS"
	}
	return
}

func NewRtpMediaSession(localIp, remoteIp *net.IPAddr, localPort, remotePort uint16,
	codecInfos []*rpc.CodecInfo, gd string, graph *event.Graph) (s *RtpMediaSession, err error) {
	sid := SessionIdType(atomic.AddUint32(&sessionIdCounter, 1))

	composer := comp.NewSessionComposer(sid.String(), "")
	if err = composer.ParseGraphDescription(gd); err != nil {
		logger.Errorf("parse graph error: %v", err)
		return nil, errors.New("composer parse graph description failed")
	}
	s = &RtpMediaSession{
		sessionId:  sid,
		localIp:    localIp,
		localPort:  localPort,
		remoteIp:   remoteIp,
		remotePort: remotePort,
		instanceId: "",

		// use buffered version to avoid deadlock
		doneC:  make(chan string, 3),
		status: sessionStatusCreated,

		composer: composer,
		graph:    graph,
	}

	for _, ci := range codecInfos {
		switch ci.PayloadType {
		case rpc.CodecType_PCM_ALAW, rpc.CodecType_AMRNB, rpc.CodecType_AMRWB, rpc.CodecType_H264, rpc.CodecType_EVS:
			if s.avPayloadNumber != 0 {
				err = fmt.Errorf("create session with more than one audio/video type:"+
					" previous number:%v, this number:%v", s.avPayloadNumber, ci.PayloadNumber)
				return
			}
			s.avPayloadNumber = uint8(ci.PayloadNumber)
			s.avPayloadCodec = ci.PayloadType
			s.avCodecParam = ci.CodecParam
		case rpc.CodecType_TELEPHONE_EVENT_8K, rpc.CodecType_TELEPHONE_EVENT_16K:
			s.telephoneEventPayloadNumber = uint8(ci.PayloadNumber)
			s.telephoneEventPayloadCodec = ci.PayloadType
			s.telephoneEventCodecParam = ci.CodecParam
		}
	}
	if s.avPayloadNumber == 0 {
		err = errors.New("create session without any audio/video codec info")
	}

	// everything is checked, setup the watchdog
	s.watchdog = newWatchDog(s)
	return
}

func (s *RtpMediaSession) setupGraph() error {
	if err := s.composer.ComposeNodes(s.graph); err != nil {
		return err
	}
	// search rtp packet provider and consumer, this is the edge between rtp stack and graph
	s.composer.IterateNode(func(name string, node comp.SessionAware) {
		if s.pullC != nil && s.handleC != nil {
			return
		}
		if provider := comp.NodeTo[RtpPacketProvider](node); provider != nil {
			if s.pullC != nil {
				logger.Errorf("session(%v) has more than one rtp packet provider", s.GetSessionId())
			} else {
				s.pullC = provider.PullPacketChannel()
			}
		}

		if consumer := comp.NodeTo[RtpPacketConsumer](node); consumer != nil {
			if s.handleC != nil {
				logger.Errorf("session(%v) has more than one rtp packet consumer", s.GetSessionId())
			} else {
				s.handleC = consumer.HandlePacketChannel()
			}
		}
	})
	if s.pullC == nil || s.handleC == nil {
		return fmt.Errorf("session(%v) has invalid rtp provider(with channel:%v) or consumer(with channel:%v) ",
			s.GetSessionId(), s.pullC, s.handleC)
	}
	return nil
}

// activate carry out actual work, such as listen on udp port, create rtp stream, create event node instances and
// add them to graph
func (s *RtpMediaSession) activate() (err error) {
	var tpLocal *rtp.TransportUDP
	var localPort = int(s.localPort)
	if err = s.setupGraph(); err != nil {
		return
	}
	if tpLocal, err = rtp.NewTransportUDP(s.localIp, localPort, ""); err != nil {
		return
	}
	s.rtpSession = rtp.NewSession(tpLocal, tpLocal)
	strLocalIdx, errStr := s.rtpSession.NewSsrcStreamOut(&rtp.Address{
		IPAddr:   s.localIp.IP,
		DataPort: localPort,
		CtrlPort: 1 + localPort,
		Zone:     "",
	}, 0, 0)
	if errStr != "" {
		return errors.New(string(errStr))
	}
	if profile := profileOfCodec(s.avPayloadCodec); profile != "" {
		s.rtpSession.SsrcStreamOutForIndex(strLocalIdx).SetProfile(profile, byte(s.avPayloadNumber))
		s.rtpSessionLocalId = strLocalIdx //add by sean
	} else {
		return errors.New("unsupported rtp payload profile")
	}

	//s.watchdog.start()
	return nil
}

//author:sean. purpose:update rtp params.but does not use
func (s *RtpMediaSession) UpdateRtpParams() (err error) {
	if s.rtpSession == nil {
		logger.Errorln("Please initialize mediaSession before update payload")
	}

	if profile := profileOfCodec(s.avPayloadCodec); profile != "" {
		ssrcStream := s.rtpSession.SsrcStreamOutForIndex(s.rtpSessionLocalId)
		ssrcStream.SetProfile(profile, byte(s.avPayloadNumber))
		logger.Infof("update payload_number.current=%v", ssrcStream.PayloadTypeNumber())
	} else {
		return errors.New("unsupported rtp payload profile")
	}

	return nil
}

func (s *RtpMediaSession) onSystemEvent(se *rpc.SystemEvent) {
	switch se.Cmd {
	case rpc.SystemCommand_USER_EVENT:
		// TODO: use user event
	case rpc.SystemCommand_SESSION_INFO:
		s.watchdog.reportSessionInfo(se)
	}

}
