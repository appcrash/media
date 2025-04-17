package server

import (
	"context"
	"fmt"
	"github.com/appcrash/GoRTP/rtp"
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/event"
	"github.com/appcrash/media/server/prom"
	"github.com/appcrash/media/server/rpc"
	"github.com/appcrash/media/server/utils"
	"net"
	"sync"
	"time"
)

const (
	sessionStatusCreated = iota
	sessionStatusUpdated
	sessionStatusStarted
	sessionStatusStopped
)

type RtpMediaSession struct {
	sessionId             SessionIdType
	localIp, remoteIp     *net.IPAddr
	localPort, remotePort uint16
	rtpSession            *rtp.Session
	rtpSessionLocalId     uint32 // rtpSession id which update rtp params
	instanceId            string // which instance created this session

	avPayloadNumber uint8
	avPayloadCodec  rpc.CodecType
	avCodecParam    string

	telephoneEventPayloadNumber uint8
	telephoneEventPayloadCodec  rpc.CodecType
	telephoneEventCodecParam    string

	mutex sync.Mutex

	status     int
	cancelFunc context.CancelFunc
	doneC      chan string // notify this channel when loop is done

	pullC        <-chan *utils.RtpPacketList
	handleC      chan<- *utils.RtpPacketList
	interceptors []RtpPacketInterceptor
	composer     *comp.Composer
	watchdog     *WatchDog
	graph        *event.Graph
}

func (s *RtpMediaSession) GetSessionId() SessionIdType {
	return s.sessionId
}

func (s *RtpMediaSession) GetStatus() int {
	return s.status
}

func (s *RtpMediaSession) GetAVPayloadType() uint8 {
	return s.avPayloadNumber
}

func (s *RtpMediaSession) GetAVCodecType() rpc.CodecType {
	return s.avPayloadCodec
}

func (s *RtpMediaSession) GetAVCodecParam() string {
	return s.avCodecParam
}

func (s *RtpMediaSession) GetTelephoneEventPayloadType() uint8 {
	return s.telephoneEventPayloadNumber
}

func (s *RtpMediaSession) GetTelephoneEventCodecType() rpc.CodecType {
	return s.telephoneEventPayloadCodec
}

func (s *RtpMediaSession) GetController() comp.CommandInitiator {
	return s.composer.GetCommandInitiator()
}

func (s *RtpMediaSession) Start() (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.status != sessionStatusCreated {
		err = fmt.Errorf("RtpMediaSession(%v) start at status %v", s.sessionId, s.status)
		logger.Errorf("try to start session(%v) when status is %v", s.sessionId, s.status)
		return
	}

	defer func() {
		// if start failed, stop using this session anymore
		if err != nil {
			logger.Errorf("session(%v) start failed with error(%v), finalize it", s.sessionId, err)
			s.Stop()
		}
	}()

	port := int(s.remotePort)
	if _, err = s.rtpSession.AddRemote(&rtp.Address{
		IPAddr:   s.remoteIp.IP,
		DataPort: port,
		CtrlPort: 1 + port,
		Zone:     "",
	}); err != nil {
		return
	}
	if err = s.rtpSession.StartSession(); err != nil {
		return
	}
	prom.RtpStartedSession.Inc()

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel
	go s.receiveRtcpLoop(ctx)
	go s.receiveRtpLoop(ctx)
	go s.sendRtpLoop(ctx)
	s.status = sessionStatusStarted
	return
}

func (s *RtpMediaSession) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var nbDone int
	if s.status == sessionStatusStopped {
		//logger.Errorf("try to stop already terminated session(%v)", s.sessionId)
		return
	}
	if s.status == sessionStatusCreated {
		// created but not started
		goto cleanup
	}

	if s.cancelFunc != nil {
		s.cancelFunc()
		// for debug purpose, check all loops are finished normally
	done:
		for nbDone < 3 {
			select {
			case <-s.doneC:
				nbDone++
			case <-time.After(5 * time.Second):
				break done
			}
		}
		if nbDone != 3 {
			logger.Errorf("session(%v) loops don't stop normally, finished number:%v", s.sessionId, nbDone)
			prom.RtpAbnormalSession.Inc()
		}
		if s.doneC != nil {
			close(s.doneC)
			s.doneC = nil
		}
	}

cleanup:
	// release all resources this session occupied
	if s.composer != nil {
		s.composer.ExitGraph()
	}
	if s.rtpSession != nil {
		s.rtpSession.CloseSession()
	}
	s.status = sessionStatusStopped
	prom.RtpStartedSession.Dec()
}
