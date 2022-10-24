package server

import (
	"context"
	"github.com/appcrash/media/server/rpc"
	"github.com/appcrash/media/server/utils"
	"sync"
	"time"
)

const (
	SessionAuditPeriod       = 30 * time.Second
	SessionTimeoutPeriod     = 5 * time.Minute
	ReportErrorThreshold     = 10
	ReportInfoPacketInterval = 200
)

const (
	sendLoop = iota
	receiveLoop
	rtcpLoop

	nbLoopReporter
)

// WatchDog is used to detect sessions in abnormal state such as zombie session and end it if necessary
// it detects state by:
// 1. send/recv loops actively report info or error
// 2. periodically check signalling server reported session info(if it is capable)
//
// if any of above reported timestamp timeout, watchdog will end this session
type WatchDog struct {
	session *MediaSession

	mutex                  sync.Mutex
	started                bool
	errorLogged            *utils.Set[int]
	createTimestamp        time.Time
	instanceAliveTimestamp time.Time // last time we recv session info state from instance
	loopAliveTimestamp     [nbLoopReporter]time.Time
	nbError                int32
	cancel                 context.CancelFunc
}

func newWatchDog(s *MediaSession) *WatchDog {
	now := time.Now()
	return &WatchDog{
		session:         s,
		createTimestamp: now,
		errorLogged:     utils.NewSet[int](),
	}
}

func (wd *WatchDog) start() {
	wd.mutex.Lock()
	defer wd.mutex.Unlock()
	if wd.started {
		return
	}
	wd.started = true
	ctx, cancel := context.WithCancel(context.Background())
	wd.cancel = cancel
	go wd.healthCheck(ctx)
}

func (wd *WatchDog) reportLoopInfo(loopId int) {
	wd.mutex.Lock()
	defer wd.mutex.Unlock()
	wd.loopAliveTimestamp[loopId] = time.Now()
}

func (wd *WatchDog) reportLoopError(loopId int, err error) {
	wd.mutex.Lock()
	defer wd.mutex.Unlock()
	if !wd.errorLogged.Contain(loopId) {
		logger.Error(err)
	}

	wd.nbError++
	if wd.nbError > ReportErrorThreshold {
		logger.Errorf("watchdog(%v): stop session due to too many errors", wd.session.GetSessionId())
		wd.stop()
	}
}

func (wd *WatchDog) reportSessionInfo(_ *rpc.SystemEvent) {
	wd.mutex.Lock()
	defer wd.mutex.Unlock()
	wd.instanceAliveTimestamp = time.Now()
}

// stop session as well as watchdog itself
func (wd *WatchDog) stop() {
	wd.session.Stop()
	if wd.cancel != nil {
		wd.cancel()
	}
}

// healthCheck periodically check session's state
func (wd *WatchDog) healthCheck(ctx context.Context) {
	ticker := time.NewTicker(SessionAuditPeriod)
	session := wd.session
	for {
		select {
		case <-ticker.C:
			sessionId := session.sessionId
			// start a new round audit ...
			var stopSession bool
			wd.mutex.Lock()
			switch session.status {
			case sessionStatusStarted:
				// currently only check that is any packet still received
				// open question: how to use send loop's info to determine zombie session
				recvTs := wd.loopAliveTimestamp[receiveLoop]
				if !recvTs.IsZero() && time.Since(recvTs) > SessionTimeoutPeriod {
					logger.Errorf("session(%v) has not received any packet in timeout period, stop it", sessionId)
					stopSession = true
				}
				fallthrough // more checks
			case sessionStatusCreated:
				// created session has no running loops, check instance aliveness and if create timestamp too far away
				if wd.instanceAliveTimestamp.IsZero() {
					// instance has not reported any info yet, so examine session's creation moment
					if session.status == sessionStatusCreated && time.Since(wd.createTimestamp) > SessionTimeoutPeriod {
						logger.Errorf("session(%v) created but not started until timeout, stop it", sessionId)
						stopSession = true
					}
				} else {
					// the instance is able to report its session info, check whether disconnected
					if time.Since(wd.instanceAliveTimestamp) > SessionTimeoutPeriod {
						logger.Errorf("session(%v) has no update from instance since %v, timeout, stop it",
							wd.instanceAliveTimestamp, sessionId)
						stopSession = true
					}
				}
			case sessionStatusStopped:
				stopSession = true
			default:
				logger.Errorf("session(%v) has unknown state(%v)", sessionId, session.status)
			}
			wd.mutex.Unlock()
			if stopSession {
				wd.stop()
			}
		case <-ctx.Done():
			return
		}
	}
}
