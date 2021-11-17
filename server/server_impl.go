package server

import (
	"fmt"
	"github.com/appcrash/media/server/channel"
	"github.com/appcrash/media/server/prom"
	"github.com/appcrash/media/server/rpc"
	"net"
	"time"
)

func (srv *MediaServer) init(ip *net.IPAddr, portStart, portEnd uint16) {
	srv.rtpServerIpAddr = ip
	srv.portPool.init(portStart, portEnd)
	//go srv.healthCheck()
	//channel.GetSystemChannel().AddListener(srv)
}

func (srv *MediaServer) addToSessionMap(session *MediaSession) {
	srv.sessionMutex.Lock()
	defer srv.sessionMutex.Unlock()
	srv.sessionMap[session.sessionId] = session
	prom.CreatedSession.Inc()
}

func (srv *MediaServer) removeFromSessionMap(session *MediaSession) {
	srv.sessionMutex.Lock()
	defer srv.sessionMutex.Unlock()
	delete(srv.sessionMap, session.sessionId)
	prom.CreatedSession.Dec()
}

func (srv *MediaServer) getNextAvailableRtpPort() uint16 {
	return srv.portPool.get()
}

func (srv *MediaServer) reclaimRtpPort(port uint16) {
	if port != 0 {
		srv.portPool.put(port)
	}
}

func (srv *MediaServer) createSession(param *rpc.CreateParam) (session *MediaSession, err error) {
	defer func() {
		if err != nil && session != nil {
			session.finalize()
		}
	}()

	if session, err = newSession(srv, param); err != nil {
		return
	}
	// initialize source/sink list for each session
	// the factory's order is important
	for _, factory := range srv.sourceF {
		src := factory.NewSource(session)
		session.source = append(session.source, src)
	}
	for _, factory := range srv.sinkF {
		sink := factory.NewSink(session)
		session.sink = append(session.sink, sink)
	}

	// connect source/sink into event graph of this session
	// then listen on udp messages
	if err = session.activate(); err != nil {
		return
	}
	prom.AllSession.Inc()
	srv.addToSessionMap(session)
	return session, nil
}

// APIs that allow plugging in method to:
// 1. handle command(take new actions), listen to state change
// 2. handle incoming data(audio,video), which is a sink
// 3. generate outgoing data, which is a source
func (srv *MediaServer) registerCommandExecutor(e CommandExecute) (err error) {
	cmdTrait := e.GetCommandTrait()
	var cm map[string]CommandExecute

	srv.executorMutex.Lock()
	defer srv.executorMutex.Unlock()
	for _, trait := range cmdTrait {
		switch trait.CmdTrait {
		case CMD_TRAIT_SIMPLE:
			cm = srv.simpleExecutorMap
		case CMD_TRAIT_STREAM:
			cm = srv.streamExecutorMap
		default:
			err = fmt.Errorf("register command executor with wrong trait: %v", trait.CmdTrait)
			return
		}
		if cmd, ok := cm[trait.CmdName]; ok {
			err = fmt.Errorf("regsiter execute with command %v already registered", cmd)
			return
		} else {
			logger.Infof("register execute with command %v\n", trait.CmdName)
			cm[trait.CmdName] = e
		}
	}
	return
}

func (srv *MediaServer) getExecutorFor(cmd string) (needNotify bool, ce CommandExecute) {
	srv.executorMutex.Lock()
	defer srv.executorMutex.Unlock()
	if e, ok := srv.simpleExecutorMap[cmd]; ok {
		needNotify = false
		ce = e.(CommandExecute)
	}
	if e, ok := srv.streamExecutorMap[cmd]; ok {
		needNotify = true
		ce = e.(CommandExecute)
	}
	needNotify = false
	ce = nil
	return
}

const (
	healthCheckPeriod    = 10 * time.Second
	sessionCheckPeriod   = 30 * time.Second
	sessionTimeoutPeriod = 2 * time.Minute
)

// healthCheck periodically check sessions' state
func (srv *MediaServer) healthCheck() {
	ticker := time.NewTicker(healthCheckPeriod)
	sysChannel := channel.GetSystemChannel()
	for {
		select {
		case <-ticker.C:
			var ms []*MediaSession
			srv.sessionMutex.Lock()
			for _, s := range srv.sessionMap {
				ms = append(ms, s)
			}
			srv.sessionMutex.Unlock()
			for _, session := range ms {
				session.mutex.Lock()
				checkTs := session.activeCheckTimestamp
				echoTs := session.activeEchoTimestamp
				status := session.status
				session.mutex.Unlock()
				sessionId := session.sessionId
				instanceId := session.instanceId
				if time.Since(echoTs) > sessionTimeoutPeriod {
					logger.Infof("session(%v) is stopped by health check due to timeout", sessionId)
					session.Stop()
				}
				if time.Since(checkTs) > sessionCheckPeriod {
					var evt string
					switch status {
					case sessionStatusCreated:
						evt = "create"
					case sessionStatusStarted:
						evt = "start"
					case sessionStatusStopped:
						evt = "stop"
					default:
						logger.Errorf("session(%v) has unknown state(%v)", sessionId, status)
					}
					logger.Debugf("session(%v) state is queried by health check", sessionId)
					msg := rpc.SystemEvent{
						Cmd:        rpc.SystemCommand_SESSION_INFO,
						InstanceId: instanceId,
						SessionId:  sessionId,
						Event:      evt,
					}
					if instanceId == "" {
						goto broadcast
					}
					if err := sysChannel.NotifyInstance(&msg); err != nil {
						logger.Error(err)
						goto broadcast
					}
					continue
				broadcast:
					// instance of session is unknown or disconnected from server, so broadcast it
					// in hope that somebody echos this session's state
					logger.Infof("session(%v)'s instance(%v) is not available, broadcast its info",
						sessionId, instanceId)
					sysChannel.BroadcastInstance(&msg)
				}
			}
		}
	}
}

func (srv *MediaServer) OnChannelEvent(e *rpc.SystemEvent) {
	var exist bool
	var session *MediaSession
	sessionId := e.SessionId
	event := e.Event
	if sessionId == "" {
		logger.Errorf("OnChannelEvent got empty session id")
		return
	}
	switch e.Cmd {
	case rpc.SystemCommand_SESSION_INFO:
		// report from signalling server about the session state
		srv.sessionMutex.Lock()
		session, exist = srv.sessionMap[sessionId]
		srv.sessionMutex.Unlock()
		if !exist {
			logger.Debugf("OnChannelEvent got an non-existent session report: %v", sessionId)
			return
		}
		switch event {
		case "ok":
			// update session timestamp
			logger.Debugf("update session(%v) active timestamp", sessionId)
			session.mutex.Lock()
			session.activeEchoTimestamp = time.Now()
			session.mutex.Unlock()
		case "stop":
			// the session has been terminated in signalling server, so must be done in media server too
			logger.Infof("session(%v) is stopped by channel event", sessionId)
			session.Stop()
		default:
			logger.Errorf("unknown event(%v) for SESSION_INFO", event)
		}

	}
}
