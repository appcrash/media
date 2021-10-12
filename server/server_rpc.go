package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/appcrash/media/server/rpc"
	"runtime/debug"
)

func (srv *MediaServer) PrepareSession(_ context.Context, param *rpc.CreateParam) (*rpc.Session, error) {
	var session *MediaSession
	var err error
	if session, err = srv.createSession(param); err != nil {
		logger.Errorf("fail to prepare session with error:%v", err)
		return nil, err
	}

	if err = session.AddRemote(param.GetPeerIp(), int(param.GetPeerPort())); err != nil {
		session.finalize()
		return nil, err
	}

	logger.Infof("rpc: prepared session %v", session.sessionId)
	rpcSession := rpc.Session{}
	rpcSession.SessionId = session.sessionId
	rpcSession.PeerIp = param.GetPeerIp()
	rpcSession.PeerRtpPort = param.GetPeerPort()
	rpcSession.LocalRtpPort = uint32(session.rtpPort)
	rpcSession.LocalIp = session.localIp.String()

	return &rpcSession, nil
}

func (srv *MediaServer) StartSession(_ context.Context, param *rpc.StartParam) (*rpc.Status, error) {
	sessionId := param.GetSessionId()
	logger.Infof("rpc: start session %v", sessionId)
	if obj, exist := sessionMap.Load(sessionId); exist {
		if session, ok := obj.(*MediaSession); ok {
			if err := session.Start(); err != nil {
				return nil, fmt.Errorf("start session failed: %v", err)
			}
		} else {
			return nil, errors.New("not a session object")
		}
	} else {
		return nil, errors.New("session not exist")
	}

	return &rpc.Status{Status: "ok"}, nil
}

func (srv *MediaServer) StopSession(_ context.Context, param *rpc.StopParam) (*rpc.Status, error) {
	sessionId := param.GetSessionId()
	logger.Infof("rpc: stop session %v", sessionId)
	if obj, exist := sessionMap.Load(sessionId); exist {
		if session, ok := obj.(*MediaSession); ok {
			session.Stop()
			return &rpc.Status{Status: "ok"}, nil
		} else {
			return nil, errors.New("not a session object")
		}
	} else {
		return nil, errors.New("session not exist")
	}
}

func (srv *MediaServer) ExecuteAction(_ context.Context, action *rpc.Action) (*rpc.ActionResult, error) {
	sessionId := action.SessionId
	s, ok := sessionMap.Load(sessionId)
	result := rpc.ActionResult{
		SessionId: sessionId,
	}

	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			logger.Errorln("ExecuteAction panic(recovered)")
		}
	}()

	if ok {
		session := s.(*MediaSession)
		cmd := action.GetCmd()
		arg := action.GetCmdArg()
		if e, ok1 := srv.simpleExecutorMap.Load(cmd); ok1 {
			exec := e.(CommandExecute)
			exec.Execute(session, cmd, arg)
			result.State = "ok"
			return &result, nil
		}
	}
	result.State = "cmd execute not exist"
	return &result, nil
}

func (srv *MediaServer) ExecuteActionWithNotify(action *rpc.Action, stream rpc.MediaApi_ExecuteActionWithNotifyServer) error {
	sessionId := action.SessionId
	s, ok := sessionMap.Load(sessionId)
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			logger.Errorln("ExecuteActionWithNotify panic(recovered)")
		}
	}()

	if ok {
		session := s.(*MediaSession)
		cmd := action.GetCmd()
		arg := action.GetCmdArg()
		if e, ok1 := srv.streamExecutorMap.Load(cmd); ok1 {
			exec := e.(CommandExecute)
			ctrlIn := make(ExecuteCtrlChan, 10)
			ctrlOut := make(ExecuteCtrlChan, 10)
			go exec.ExecuteWithNotify(session, cmd, arg, ctrlIn, ctrlOut)

		outLoop:
			for {
				select {
				case msg, more := <-ctrlOut:
					event := rpc.ActionEvent{
						SessionId: sessionId,
						Event:     msg,
					}
					if !more {
						// executor loop already exits by itself, break without notification
						break outLoop
					}
					if err := stream.Send(&event); err != nil {
						logger.Errorf("send action event of stream(%v) with event %v error", session, event.String())
						// notify executor loop to exit
						close(ctrlIn)
						break outLoop
					}
				}
			}
		}

		return nil
	}

	return errors.New("cmd not exist")
}
