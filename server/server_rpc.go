package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/appcrash/media/server/channel"
	"github.com/appcrash/media/server/prom"
	"github.com/appcrash/media/server/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"runtime/debug"
	"strings"
	"sync"
)

func (srv *MediaServer) GetVersion(_ context.Context, _ *rpc.Empty) (*rpc.VersionNumber, error) {
	return &rpc.VersionNumber{Ver: rpc.Version_DEFAULT}, nil
}

func (srv *MediaServer) PrepareSession(_ context.Context, param *rpc.CreateParam) (*rpc.Session, error) {
	var session *MediaSession
	var err error
	if session, err = srv.createSession(param); err != nil {
		logger.Errorf("fail to prepare session with error:%v", err)
		return nil, err
	}

	logger.Infof("rpc: prepared session %v", session.sessionId)
	rpcSession := rpc.Session{}
	rpcSession.SessionId = session.sessionId.String()
	rpcSession.PeerIp = param.GetPeerIp()
	rpcSession.PeerRtpPort = param.GetPeerPort()
	rpcSession.LocalRtpPort = uint32(session.localPort)
	rpcSession.LocalIp = session.localIp.String()

	return &rpcSession, nil
}

func (srv *MediaServer) UpdateSession(_ context.Context, param *rpc.UpdateParam) (*rpc.Status, error) {
	// only remote (ip, port) can be updated
	err := srv.updateSession(param)
	return &rpc.Status{Status: "ok"}, err
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
	var sessionId SessionIdType
	var err error
	if sessionId, err = SessionIdFromString(action.SessionId); err != nil {
		err = errors.New("invalid session id")
		return nil, err
	}
	result := rpc.ActionResult{
		SessionId: action.SessionId,
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
		exec, exist := srv.simpleExecutorMap[cmd]
		if exist {
			re := exec.Execute(session, cmd, arg)
			prom.SessionAction.With(prometheus.Labels{"cmd": cmd, "type": "simple"}).Inc()
			result.State = strings.Join(re, " ")
			return &result, nil
		}
	}
	result.State = "error execute not exist"
	return &result, nil
}

func (srv *MediaServer) ExecuteActionWithNotify(action *rpc.Action, stream rpc.MediaApi_ExecuteActionWithNotifyServer) error {
	var sessionId SessionIdType
	var err error
	if sessionId, err = SessionIdFromString(action.SessionId); err != nil {
		err = errors.New("invalid session id")
		return err
	}
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
		exec, exist := srv.streamExecutorMap[cmd]
		if exist {
			ctx, cancel := context.WithCancel(context.Background())
			ctrlOut := make(ExecuteCtrlChan, 32)
			defer func() {
				prom.SessionAction.With(prometheus.Labels{"cmd": cmd, "type": "pull_stream"}).Inc()
				// notify executor loop to exit
				cancel()
			}()
			go exec.ExecuteWithNotify(session, arg, ctx, ctrlOut)

		outLoop:
			for {
				select {
				case msg, more := <-ctrlOut:
					event := rpc.ActionEvent{
						SessionId: action.SessionId,
						Event:     msg,
					}
					if !more {
						// executor loop already exits by itself, break without notification
						break outLoop
					}
					if err := stream.Send(&event); err != nil {
						logger.Errorf("send action event of stream(%v) with event %v error", session, event.String())
						break outLoop
					}
				}
			}
		}

		return nil
	}

	return errors.New("cmd not exist")
}

func (srv *MediaServer) ExecuteActionWithPush(stream rpc.MediaApi_ExecuteActionWithPushServer) error {
	var sessionId SessionIdType
	var dataIn chan *rpc.PushData

	// receive the first data to retrieve the session id
	if data, err := stream.Recv(); err != nil {
		return err
	} else {
		sidStr := data.GetSessionId()
		if sidStr == "" {
			return errors.New("push action with empty session id")
		}
		if sessionId, err = SessionIdFromString(sidStr); err != nil {
			return errors.New("push action with invalid sessin id")
		}
		srv.sessionMutex.Lock()
		session, ok := srv.sessionMap[sessionId]
		srv.sessionMutex.Unlock()
		if !ok {
			return fmt.Errorf("push action with session(%v) that is not exist", sessionId)
		}
		exec := srv.streamExecutorMap[data.Cmd]
		if exec == nil {
			logger.Errorf("no push executor cmd: %v registered", data.Cmd)
			return fmt.Errorf("push cmd %v not found", data.Cmd)
		}
		dataIn = make(chan *rpc.PushData, 32)
		go exec.ExecuteWithPush(session, dataIn)
		// don't forget to send first packet
		dataIn <- data
	}

	defer func() {
		prom.SessionAction.With(prometheus.Labels{"type": "push_stream"}).Inc()
		// let push loop stop
		if dataIn != nil {
			close(dataIn)
		}
	}()

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&rpc.ActionResult{
				SessionId: sessionId.String(),
				State:     "ok",
			})
		}
		if err != nil {
			return err
		}
		select {
		case dataIn <- data:
		default:
		}
	}
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
