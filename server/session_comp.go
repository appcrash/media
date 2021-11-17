package server

import (
	"errors"
	"fmt"
	"github.com/appcrash/GoRTP/rtp"
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/prom"
	"github.com/appcrash/media/server/rpc"
	"github.com/google/uuid"
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
	case rpc.CodecType_TELEPHONE_EVENT:
		profile = "TELEPHONE-EVENT"
	}
	return
}

func newSession(srv *MediaServer, mediaParam *rpc.CreateParam) (*MediaSession, error) {
	var localPort uint16
	if localPort = srv.getNextAvailableRtpPort(); localPort == 0 {
		return nil, errors.New("server runs out of port resource")
	}
	//instanceId := mediaParam.InstanceId
	//if !channel.GetSystemChannel().HasInstance(instanceId) {
	//	return nil, fmt.Errorf("the instance %v not registered, cannot create session", instanceId)
	//}
	sid := uuid.New().String()
	sid = strings.Replace(sid, "-", "", -1) // ID in nmd language doesn't contains '-'
	gd := mediaParam.GetGraphDesc()
	composer := comp.NewSessionComposer(sid)
	if err := composer.ParseGraphDescription(gd); err != nil {
		logger.Errorln(err)
		return nil, errors.New("composer parse graph description failed")
	}
	now := time.Now()
	s := MediaSession{
		server:    srv,
		sessionId: sid,
		localIp:   srv.rtpServerIpAddr,
		localPort: int(localPort),
		rtpPort:   localPort,
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
		return nil, errors.New("create s without any codec info")
	}
	for _, ci := range codecInfos {
		switch ci.PayloadType {
		case rpc.CodecType_PCM_ALAW, rpc.CodecType_AMRNB, rpc.CodecType_AMRWB:
			if s.audioPayloadNumber != 0 {
				return nil, fmt.Errorf("create session with more than one audio type:"+
					" previous number:%v, this number:%v", s.audioPayloadNumber, ci.PayloadNumber)
			}
			s.audioPayloadNumber = uint8(ci.PayloadNumber)
			s.audioPayloadCodec = ci.PayloadType
			s.audioCodecParam = ci.CodecParam
		case rpc.CodecType_TELEPHONE_EVENT:
			s.telephoneEventPayloadNumber = uint8(ci.PayloadNumber)
		}
	}
	if s.audioPayloadNumber == 0 {
		return nil, errors.New("create session without any audio codec info")
	}

	return &s, nil
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
	if err := s.composer.ComposeNodes(s.server.graph); err != nil {
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
	if profile := profileOfCodec(s.audioPayloadCodec); profile != "" {
		s.rtpSession.SsrcStreamOutForIndex(strLocalIdx).SetProfile(profile, byte(s.audioPayloadNumber))
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
	prom.StartedSession.Dec()
	s.server.removeFromSessionMap(s)
}
