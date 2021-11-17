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
	sessionStatusStarted
	sessionStatusStopped
)

type MediaSession struct {
	server     *MediaServer
	sessionId  string
	localIp    *net.IPAddr
	localPort  int
	rtpPort    uint16
	rtpSession *rtp.Session
	instanceId string // which instance created this session

	audioPayloadNumber uint8
	audioPayloadCodec  rpc.CodecType
	audioCodecParam    string

	telephoneEventPayloadNumber uint8

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

func (s *MediaSession) GetAudioPayloadType() uint8 {
	return s.audioPayloadNumber
}

func (s *MediaSession) GetAudioCodecType() rpc.CodecType {
	return s.audioPayloadCodec
}

func (s *MediaSession) GetAudioCodecParam() string {
	return s.audioCodecParam
}

func (s *MediaSession) GetTelephoneEventPayloadType() uint8 {
	return s.telephoneEventPayloadNumber
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
	nbDone := 0
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
			logger.Errorf("s(%v) loops don't stop normally, finished number:%v", s.sessionId, nbDone)
		}
	}

cleanup:
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

// LinkChannel links *ch* to the channel registered as *name* in graph description
func (s *MediaSession) LinkChannel(name string, ch chan<- *event.Event) {
	_ = s.composer.LinkChannel(name, ch)
}
