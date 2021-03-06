package server

import (
	"context"
	"github.com/appcrash/GoRTP/rtp"
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/event"
	"github.com/appcrash/media/server/prom"
	"github.com/appcrash/media/server/rpc"
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

type MediaSession struct {
	server                *MediaServer
	sessionId             string
	localIp, remoteIp     *net.IPAddr
	localPort, remotePort uint16
	rtpSession            *rtp.Session
	instanceId            string // which instance created this session

	avPayloadNumber uint8
	avPayloadCodec  rpc.CodecType
	avCodecParam    string

	telephoneEventPayloadNumber uint8
	telephoneEventPayloadCodec  rpc.CodecType
	telephoneEventCodecParam    string

	mutex                sync.Mutex
	createTimestamp      time.Time
	activeCheckTimestamp time.Time // last time we send session info state to instance, updated by server
	activeEchoTimestamp  time.Time // last time we recv session info state from instance, updated by server
	status               int
	cancelFunc           context.CancelFunc
	doneC                chan string // notify this channel when loop is done

	source   []Source
	sink     []Sink
	composer *comp.Composer
}

func (s *MediaSession) GetSessionId() string {
	return s.sessionId
}

func (s *MediaSession) GetStatus() int {
	return s.status
}

func (s *MediaSession) GetAVPayloadType() uint8 {
	return s.avPayloadNumber
}

func (s *MediaSession) GetAVCodecType() rpc.CodecType {
	return s.avPayloadCodec
}

func (s *MediaSession) GetAVCodecParam() string {
	return s.avCodecParam
}

func (s *MediaSession) GetTelephoneEventPayloadType() uint8 {
	return s.telephoneEventPayloadNumber
}

func (s *MediaSession) GetTelephoneEventCodecType() rpc.CodecType {
	return s.telephoneEventPayloadCodec
}

func (s *MediaSession) GetEventGraph() *event.Graph {
	return s.server.graph
}

func (s *MediaSession) GetController() comp.Controller {
	return s.composer.GetController()
}

func (s *MediaSession) Start() (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.status != sessionStatusCreated {
		logger.Errorf("try to start session(%v) when status is %v", s.sessionId, s.status)
		return
	}

	defer func() {
		// if start failed, stop using this session anymore
		if err != nil {
			logger.Errorf("session(%v) start failed with error(%v), finalize it", s.sessionId, err)
			s.finalize()
			s.status = sessionStatusStopped
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
	prom.StartedSession.Inc()

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel
	go s.receiveCtrlLoop(ctx)
	go s.receivePacketLoop(ctx)
	go s.sendPacketLoop(ctx)
	s.status = sessionStatusStarted
	return
}

func (s *MediaSession) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var nbDone int
	if s.status == sessionStatusStopped {
		logger.Errorf("try to stop already terminated session(%v)", s.sessionId)
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
			case <-time.After(10 * time.Second):
				break done
			}
		}
		if nbDone != 3 {
			logger.Errorf("session(%v) loops don't stop normally, finished number:%v", s.sessionId, nbDone)
		}
	}

cleanup:
	s.finalize()
}

// AddNode add an event node to server-wide event graph
func (s *MediaSession) AddNode(node event.Node) {
	s.server.graph.AddNode(node)
}

// LinkChannel links *ch* to the channel registered as *name* in graph description
func (s *MediaSession) LinkChannel(name string, ch chan<- *event.Event) {
	_ = s.composer.LinkChannel(name, ch)
}
