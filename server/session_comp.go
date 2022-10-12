package server

import (
	"errors"
	"fmt"
	"github.com/appcrash/GoRTP/rtp"
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/prom"
	"github.com/appcrash/media/server/rpc"
	"github.com/google/uuid"
	"net"
	"strings"
	"time"
)

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
	}
	return
}

func newSession(srv *MediaServer, mediaParam *rpc.CreateParam) (s *MediaSession, err error) {
	var localPort, remotePort uint16
	var remoteIp *net.IPAddr

	defer func() {
		// if create session failed, avoid port leaking
		if err != nil && localPort > 0 {
			srv.reclaimRtpPort(localPort)
		}
	}()

	if localPort = srv.getNextAvailableRtpPort(); localPort == 0 {
		err = errors.New("server runs out of port resource")
		return
	}
	//instanceId := mediaParam.InstanceId
	//if !channel.GetSystemChannel().HasInstance(instanceId) {
	//	return nil, fmt.Errorf("the instance %v not registered, cannot create session", instanceId)
	//}
	if remoteIp, err = net.ResolveIPAddr("ip", mediaParam.GetPeerIp()); err != nil {
		return nil, fmt.Errorf("invalid peer ip address: %v", mediaParam.GetPeerIp())
	}
	if mediaParam.GetPeerPort()&0xffff0000 != 0 {
		// not a uint16 port number
		return nil, fmt.Errorf("invalid peer port: %v", mediaParam.GetPeerPort())
	}
	remotePort = uint16(mediaParam.GetPeerPort())
	sid := uuid.New().String()
	sid = strings.Replace(sid, "-", "", -1) // ID in nmd language doesn't contains '-'
	gd := mediaParam.GetGraphDesc()
	composer := comp.NewSessionComposer(sid)
	if err = composer.ParseGraphDescription(gd); err != nil {
		logger.Errorf("parse graph error: %v", err)
		return nil, errors.New("composer parse graph description failed")
	}
	now := time.Now()
	s = &MediaSession{
		server:     srv,
		sessionId:  sid,
		localIp:    srv.rtpServerIpAddr,
		localPort:  localPort,
		remoteIp:   remoteIp,
		remotePort: remotePort,
		//instanceId: instanceId,

		createTimestamp:      now,
		activeCheckTimestamp: now,
		activeEchoTimestamp:  now,

		// use buffered version to avoid deadlock
		doneC:  make(chan string, 3),
		status: sessionStatusCreated,

		composer: composer,
	}
	codecInfos := mediaParam.GetCodecs()
	if len(codecInfos) == 0 {
		err = errors.New("create session without any codec info")
		return
	}
	for _, ci := range codecInfos {
		switch ci.PayloadType {
		case rpc.CodecType_PCM_ALAW, rpc.CodecType_AMRNB, rpc.CodecType_AMRWB, rpc.CodecType_H264:
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

	return
}

func (s *MediaSession) setupGraph() error {
	if err := s.composer.ComposeNodes(s.server.graph); err != nil {
		return err
	}
	// search rtp packet provider and consumer, this is the edge of rtp stack and graph
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
		return fmt.Errorf("session(%v) missing rtp provider ch:(%v) or consumer ch:(%v)",
			s.GetSessionId(), s.pullC, s.handleC)
	}
	return nil
}

// activate carry out actual work, such as listen on udp port, create rtp stream, create event node instances and
// add them to graph
func (s *MediaSession) activate() (err error) {
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
	if s.localPort != 0 {
		s.server.reclaimRtpPort(s.localPort)
		s.localPort = 0
	}
	prom.StartedSession.Dec()
	s.server.removeFromSessionMap(s)
}
