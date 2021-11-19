package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/appcrash/media/server/channel"
	"github.com/appcrash/media/server/prom"
	"github.com/appcrash/media/server/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"runtime/debug"
	"sync"
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
	err := srv.startSession(param)
	return &rpc.Status{Status: "ok"}, err
}

func (srv *MediaServer) StopSession(_ context.Context, param *rpc.StopParam) (*rpc.Status, error) {
	if err := srv.stopSession(param); err != nil {
		return nil, err
	} else {
		return &rpc.Status{Status: "ok"}, nil
	}
}

func (srv *MediaServer) ExecuteAction(_ context.Context, action *rpc.Action) (*rpc.ActionResult, error) {
	sessionId := action.SessionId
	result := rpc.ActionResult{
		SessionId: sessionId,
	}
	srv.sessionMutex.Lock()
	session, ok := srv.sessionMap[sessionId]
	srv.sessionMutex.Unlock()

	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			logger.Errorln("ExecuteAction panic(recovered)")
		}
	}()

	if ok {
		cmd := action.GetCmd()
		arg := action.GetCmdArg()
		srv.executorMutex.Lock()
		exec, exist := srv.simpleExecutorMap[cmd]
		srv.executorMutex.Unlock()
		if exist {
			exec.Execute(session, cmd, arg)
			prom.SessionAction.With(prometheus.Labels{"cmd": cmd, "type": "simple"}).Inc()
			result.State = "ok"
			return &result, nil
		}
	}
	result.State = "cmd execute not exist"
	return &result, nil
}

func (srv *MediaServer) ExecuteActionWithNotify(action *rpc.Action, stream rpc.MediaApi_ExecuteActionWithNotifyServer) error {
	sessionId := action.SessionId
	srv.sessionMutex.Lock()
	session, ok := srv.sessionMap[sessionId]
	srv.sessionMutex.Unlock()
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			logger.Errorln("ExecuteActionWithNotify panic(recovered)")
		}
	}()

	if ok {
		cmd := action.GetCmd()
		arg := action.GetCmdArg()
		srv.executorMutex.Lock()
		exec, exist := srv.streamExecutorMap[cmd]
		srv.executorMutex.Unlock()
		if exist {
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
			prom.SessionAction.With(prometheus.Labels{"cmd": cmd, "type": "stream"}).Inc()
		}

		return nil
	}

	return errors.New("cmd not exist")
}

func (srv *MediaServer) SystemChannel(stream rpc.MediaApi_SystemChannelServer) error {
	sendNotifyC := make(chan string, 2)
	recvNotifyC := make(chan string, 2)
	wg := &sync.WaitGroup{}
	var instanceId string
	var errorNotified bool
	var fromC, toC chan *rpc.SystemEvent

	logger.Infof("server's system channel is connected")
	// only work after REGISTER is seen
	for {
		in, err := stream.Recv()
		if err != nil {
			logger.Errorf("system channel rpc error %v", err)
			return err
		}
		if in.Cmd == rpc.SystemCommand_REGISTER {
			instanceId = in.InstanceId
			if instanceId == "" {
				err = fmt.Errorf("system channel got null instance id when registering")
				logger.Error(err)
				return err
			}
			break
		} else {
			if !errorNotified {
				errorNotified = true
				logger.Errorf("system channel got event before client registers itself")
			}
		}
	}

	logger.Infof("instance (%v) enters system channel rpc", instanceId)
	// the client has registered itself
	sc := channel.GetSystemChannel()
	if is, err := sc.RegisterInstance(instanceId); err != nil {
		return err
	} else {
		fromC, toC = is.FromInstanceC, is.ToInstanceC
	}
	logger.Infof("instance (%v) has registered system channel", instanceId)

	wg.Add(2)
	go func() {
		defer func() {
			wg.Done()
			logger.Infof("system channel rpc, exit recv loop")
		}()

		for {
			in, err := stream.Recv()
			if err != nil {
				sendNotifyC <- "recv_error"
				break
			}
			fromC <- in

			// check exit ...
			select {
			case <-recvNotifyC:
				return
			default:
			}
		}
	}()

	go func() {
		defer func() {
			wg.Done()
			logger.Infof("system channel rpc, exit send loop")
		}()

		for {
			select {
			case msg, more := <-toC:
				if !more {
					recvNotifyC <- "send_error"
					return
				}
				if err := stream.Send(msg); err != nil {
					logger.Errorf("system cahnnel rpc, send message error: %v", err)
					recvNotifyC <- "send_error"
					return
				}
			case <-sendNotifyC:
				return
			}
		}
	}()

	wg.Wait()
	logger.Infof("instance (%v) has exited system channel rpc", instanceId)
	return nil
}
