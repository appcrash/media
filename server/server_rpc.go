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
	"time"
)

func (srv *GrpcServer) GetVersion(_ context.Context, _ *rpc.Empty) (*rpc.VersionNumber, error) {
	return &rpc.VersionNumber{Ver: rpc.Version_DEFAULT}, nil
}

func (srv *GrpcServer) PrepareSession(_ context.Context, param *rpc.CreateParam) (*rpc.Session, error) {
	var session *RtpMediaSession
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

func (srv *GrpcServer) UpdateSession(_ context.Context, param *rpc.UpdateParam) (*rpc.Status, error) {
	// only remote (ip, port) can be updated
	err := srv.updateSession(param)
	return &rpc.Status{Status: "ok"}, err
}

func (srv *GrpcServer) StartSession(_ context.Context, param *rpc.StartParam) (*rpc.Status, error) {
	err := srv.startSession(param)
	return &rpc.Status{Status: "ok"}, err
}

func (srv *GrpcServer) StopSession(_ context.Context, param *rpc.StopParam) (*rpc.Status, error) {
	if err := srv.stopSession(param); err != nil {
		return nil, err
	} else {
		return &rpc.Status{Status: "ok"}, nil
	}
}

func (srv *GrpcServer) ExecuteAction(_ context.Context, action *rpc.Action) (*rpc.ActionResult, error) {
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
			re, err1 := exec.Execute(session, cmd, arg)
			prom.GrpcSessionAction.With(prometheus.Labels{"cmd": cmd, "type": "simple"}).Inc()
			result.State = strings.Join(re, " ")
			return &result, err1
		}
	}
	result.State = "error execute not exist"
	return &result, nil
}

func (srv *GrpcServer) ExecuteActionWithNotify(action *rpc.Action, stream rpc.MediaApi_ExecuteActionWithNotifyServer) error {
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
				prom.GrpcSessionAction.With(prometheus.Labels{"cmd": cmd, "type": "pull_stream"}).Inc()
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

func (srv *GrpcServer) ExecuteActionWithPush(stream rpc.MediaApi_ExecuteActionWithPushServer) error {
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
			return errors.New("push action with invalid session id")
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
		prom.GrpcSessionAction.With(prometheus.Labels{"type": "push_stream"}).Inc()
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

// SystemChannel is long-keepalive connection to ease bidirectional system-level message exchange
func (srv *GrpcServer) SystemChannel(stream rpc.MediaApi_SystemChannelServer) error {
	wg := &sync.WaitGroup{}
	var instanceId string
	var errorLogged bool
	var fromC, toC chan *rpc.SystemEvent

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
			if !errorLogged {
				errorLogged = true
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
	logger.Infof("instance:%v has registered system channel", instanceId)
	wg.Add(2)

	// the receive&send goroutines are independent, one exit wouldn't end the other one
	// only network exception or instance state close would end them
	now := time.Now().Format("2006-01-02 15:04:05")
	go func() {
		defer func() {
			wg.Done()
			logger.Infof("instance:%v at %v system channel rpc, exit recv loop", instanceId, now)
		}()

		for {
			in, err := stream.Recv()
			if err != nil {
				break
			}
			_, ok := <-fromC
			if !ok {
				// Channel is closed
				break
			}
			select {
			case fromC <- in:
			default:
			}

		}
	}()

	go func() {
		defer func() {
			wg.Done()
			logger.Infof("instance:%v at %v system channel rpc, exit send loop", instanceId, now)
		}()

		for {
			select {
			case msg, more := <-toC:
				if !more {
					return
				}
				if err := stream.Send(msg); err != nil {
					logger.Errorf("instance:%v system channel rpc, send message error: %v", instanceId, err)
					return
				}
			}
		}
	}()

	wg.Wait()
	logger.Infof("instance:%v at %v has exited system channel rpc normally", instanceId, now)
	return nil
}
