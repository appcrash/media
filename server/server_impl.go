package server

import (
	"errors"
	"fmt"
	"github.com/appcrash/media/server/channel"
	"github.com/appcrash/media/server/prom"
	"github.com/appcrash/media/server/rpc"
	"net"
)

func (srv *GrpcServer) init(ip *net.IPAddr, portStart, portEnd uint16) {
	srv.rtpServerIpAddr = ip
	srv.portPool.Init(portStart, portEnd)
	channel.GetSystemChannel().AddListener(srv)
}

func (srv *GrpcServer) addToSessionMap(session *RtpMediaSession) {
	srv.sessionMutex.Lock()
	defer srv.sessionMutex.Unlock()
	srv.sessionMap[session.sessionId] = session
	prom.RtpCreatedSession.Inc()
}

func (srv *GrpcServer) removeFromSessionMap(session *RtpMediaSession) {
	srv.sessionMutex.Lock()
	defer srv.sessionMutex.Unlock()
	delete(srv.sessionMap, session.sessionId)
	prom.RtpCreatedSession.Dec()
}

func (srv *GrpcServer) getNextAvailableRtpPort() uint16 {
	return srv.portPool.Get()
}

func (srv *GrpcServer) reclaimRtpPort(port uint16) {
	if port != 0 {
		srv.portPool.Put(port)
	}
}

func (srv *GrpcServer) invokeSessionListener(session *RtpMediaSession, status int) {
	switch status {
	case sessionStatusCreated:
		for _, listener := range srv.sessionListener {
			listener.OnSessionCreated(session)
		}
	case sessionStatusStarted:
		for _, listener := range srv.sessionListener {
			listener.OnSessionStarted(session)
		}
	case sessionStatusStopped:
		for _, listener := range srv.sessionListener {
			listener.OnSessionStopped(session)
		}
	case sessionStatusUpdated:
		for _, listener := range srv.sessionListener {
			listener.OnSessionUpdated(session)
		}
	}
}

func (srv *GrpcServer) createSession(param *rpc.CreateParam) (session *RtpMediaSession, err error) {
	var localIp, remoteIp *net.IPAddr
	var localPort, remotePort uint16
	defer func() {
		if err != nil {
			if localPort > 0 {
				// if create session failed, avoid port leaking
				srv.reclaimRtpPort(localPort)
			}
			if session != nil {
				session.Stop()
			}
		}
	}()

	logger.Infof("create rtp session param: %v", param)

	if localPort = srv.getNextAvailableRtpPort(); localPort == 0 {
		err = errors.New("grpc server runs out of port resource")
		return
	}
	if param.GetPeerPort()&0xffff0000 != 0 {
		// not an uint16 port number
		err = fmt.Errorf("invalid peer port: %v", param.GetPeerPort())
		return
	}
	remotePort = uint16(param.GetPeerPort())
	if remoteIp, err = net.ResolveIPAddr("ip", param.GetPeerIp()); err != nil {
		return
	}
	codecInfos := param.GetCodecs()
	if len(codecInfos) == 0 {
		err = errors.New("create session without any codec info")
		return
	}
	localIp = srv.rtpServerIpAddr
	gd := param.GetGraphDesc()

	if session, err = NewRtpMediaSession(localIp, remoteIp, localPort, remotePort, codecInfos, gd, srv.graph); err != nil {
		return
	}

	// connect source/sink into event graph of this session
	// then listen on udp messages
	if err = session.activate(); err != nil {
		return
	}
	prom.RtpAllSession.Inc()
	srv.addToSessionMap(session)
	srv.invokeSessionListener(session, sessionStatusCreated)
	return session, nil
}

func (srv *GrpcServer) updateSession(param *rpc.UpdateParam) (err error) {
	var sessionId SessionIdType
	if sessionId, err = SessionIdFromString(param.GetSessionId()); err != nil {
		err = errors.New("invalid session id")
		return
	}
	var remoteIp *net.IPAddr
	var err1 error
	if remoteIp, err1 = net.ResolveIPAddr("ip", param.GetPeerIp()); err1 != nil {
		return fmt.Errorf("update with invalid peer ip address: %v", param.GetPeerIp())
	}
	if param.GetPeerPort()&0xffff0000 != 0 {
		// not a uint16 port number
		return fmt.Errorf("invalid peer port: %v", param.GetPeerPort())
	}

	srv.sessionMutex.Lock()
	session, exist := srv.sessionMap[sessionId]
	srv.sessionMutex.Unlock()
	if exist {
		if session.status != sessionStatusCreated {
			// if session already started or stopped, no updating is done
			err = fmt.Errorf("try to update already started/stopped session(%v)", sessionId)
			return
		}
		logger.Infof("update session(%v) with param:%v", sessionId, param)
		session.remoteIp = remoteIp
		session.remotePort = uint16(param.GetPeerPort())

		//update rtp params when necessary
		pt := param.GetPayloadNumber()
		if pt > 0 { // ignore pt==0(PCMU) static payload type/default value
			_pt := uint8(pt)
			if _pt != session.avPayloadNumber { //whether update
				logger.Infof("update payload number from previous=%v,to current=%v", session.avPayloadNumber, pt)
				session.avPayloadNumber = _pt
				//session.UpdateRtpParams()
			}

		}

		srv.invokeSessionListener(session, sessionStatusUpdated)
	} else {
		err = errors.New("session not exist")
	}
	return
}

func (srv *GrpcServer) startSession(param *rpc.StartParam) (err error) {
	var sessionId SessionIdType
	if sessionId, err = SessionIdFromString(param.GetSessionId()); err != nil {
		err = errors.New("invalid session id")
		return
	}
	logger.Infof("rpc: start session %v", sessionId)
	srv.sessionMutex.Lock()
	session, exist := srv.sessionMap[sessionId]
	srv.sessionMutex.Unlock()
	if exist {
		err = session.Start()
	} else {
		err = fmt.Errorf("session:%v not exist", sessionId)
	}
	if err == nil {
		srv.invokeSessionListener(session, sessionStatusStarted)
	} else {
		if exist {
			// session exists but start failed, reclaim resource
			srv.sessionMutex.Lock()
			delete(srv.sessionMap, sessionId)
			srv.sessionMutex.Unlock()
			srv.reclaimRtpPort(session.localPort)
			srv.invokeSessionListener(session, sessionStatusStopped)
		}
	}
	return
}

func (srv *GrpcServer) stopSession(param *rpc.StopParam) (err error) {
	var sessionId SessionIdType
	if sessionId, err = SessionIdFromString(param.GetSessionId()); err != nil {
		err = errors.New("invalid session id")
		return
	}
	logger.Infof("rpc: stop session %v", sessionId)
	srv.sessionMutex.Lock()
	session, exist := srv.sessionMap[sessionId]
	srv.sessionMutex.Unlock()
	if exist {
		session.Stop()
		srv.sessionMutex.Lock()
		delete(srv.sessionMap, sessionId)
		srv.sessionMutex.Unlock()
		srv.invokeSessionListener(session, sessionStatusStopped)
	} else {
		err = errors.New("session not exist")
	}
	return
}

// APIs that allow plugging in method to:
// 1. handle command(take new actions), listen to state change
// 2. pull data from media server
// 3. push data to media server
func (srv *GrpcServer) registerCommandExecutor(e CommandExecute) (err error) {
	cmdTrait := e.GetCommandTrait()
	var cm map[string]CommandExecute

	for _, trait := range cmdTrait {
		switch trait.CmdTrait {
		case CmdTraitSimple:
			cm = srv.simpleExecutorMap
		case CmdTraitPullStream, CmdTraitPushStream:
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

func (srv *GrpcServer) getExecutorFor(cmd string) (needNotify bool, ce CommandExecute) {
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

// OnChannelEvent just forwards system event to session
func (srv *GrpcServer) OnChannelEvent(e *rpc.SystemEvent) {
	var exist bool
	var session *RtpMediaSession
	var sessionId SessionIdType
	var err error
	if sessionId, err = SessionIdFromString(e.SessionId); err != nil {
		logger.Errorf("OnChannelEvent got invalid session id")
		return
	}
	srv.sessionMutex.Lock()
	session, exist = srv.sessionMap[sessionId]
	srv.sessionMutex.Unlock()
	if !exist {
		logger.Debugf("OnChannelEvent got an non-existent session report: %v", sessionId)
		return
	}
	session.onSystemEvent(e)
}
