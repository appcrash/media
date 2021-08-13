package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/appcrash/media/server/rpc"
	"google.golang.org/grpc"
	"net"
	"runtime/debug"
	"sync"
)

type MediaServer struct {
	rpc.UnimplementedMediaApiServer
	listenPort        uint16
	simpleExecutorMap sync.Map
	streamExecutorMap sync.Map

	sourceF []SourceFactory
	sinkF []SinkFactory
}

const MYIP = "0.0.0.0"

var local, _ = net.ResolveIPAddr("ip", MYIP)


func (msrv *MediaServer) PrepareMediaStream(ctx context.Context, param *rpc.MediaParam) (*rpc.MediaStream, error) {
	port := getNextPort()
	session := createSession(int(port), param)

	// initialize source/sink list for each session
	// the factory's order is important
	for _,factory := range msrv.sourceF {
		src := factory.NewSource(session)
		session.source = append(session.source,src)
	}
	for _,factory := range msrv.sinkF {
		sink := factory.NewSink(session)
		session.sink = append(session.sink,sink)
	}

	session.AddRemote(param.GetPeerIp(),int(param.GetPeerPort()))

	ms := rpc.MediaStream{}
	ms.StreamId = session.sessionId
	ms.PeerIp = param.GetPeerIp()
	ms.LocalRtpPort = uint32(port)
	ms.LocalIp = MYIP
	ms.PeerRtpPort = param.GetPeerPort()

	return &ms,nil
}

func (msrv *MediaServer) StartSession(ctx context.Context,param *rpc.SessionParam) (*rpc.SessionStatus, error) {
	status := rpc.SessionStatus{Status : "not exist"}
	sessionId := param.GetSessionId()
	if obj,exist := sessionMap.Get(sessionId); exist {
		if session, ok := obj.(*MediaSession); ok {
			session.Start()
			status.Status = "ok"
		} else {
			status.Status = "not a session object"
		}
	}

	return &status,nil
}

func (msrv *MediaServer) StopSession(ctx context.Context,param *rpc.SessionParam) (*rpc.SessionStatus, error) {
	status := rpc.SessionStatus{Status : "not exist"}
	sessionId := param.GetSessionId()
	if obj,exist := sessionMap.Get(sessionId); exist {
		if session, ok := obj.(*MediaSession); ok {
			session.Stop()
			status.Status = "ok"
		} else {
			status.Status = "not a session object"
		}
	}

	return &status,nil
}

func (msrv *MediaServer) ExecuteAction(ctx context.Context,action *rpc.MediaAction) (*rpc.MediaActionResult, error) {
	sessionId := action.StreamId
	s,ok := sessionMap.Get(sessionId)
	result := rpc.MediaActionResult{
		StreamId: sessionId,
	}

	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			fmt.Errorf("ExecuteAction panic(recovered)")
		}
	}()

	if ok {
		session := s.(*MediaSession)
		cmd := action.GetCmd()
		arg := action.GetCmdArg()
		if e,ok1 := msrv.simpleExecutorMap.Load(cmd); ok1 {
			exec := e.(CommandExecute)
			exec.Execute(session,cmd,arg)
			result.State = "ok"
			return &result,nil
		}
	}
	result.State = "cmd execute not exist"
	return &result,nil
}


func (msrv *MediaServer) ExecuteActionWithNotify(action *rpc.MediaAction,stream rpc.MediaApi_ExecuteActionWithNotifyServer) error {
	sessionId := action.StreamId
	s,ok := sessionMap.Get(sessionId)
	eventTemplate := rpc.MediaActionEvent{
		StreamId: sessionId,
	}

	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			fmt.Errorf("ExecuteActionWithNotify panic(recovered)")
		}
	}()

	if ok {
		session := s.(*MediaSession)
		cmd := action.GetCmd()
		arg := action.GetCmdArg()
		if e,ok1 := msrv.streamExecutorMap.Load(cmd); ok1 {
			exec := e.(CommandExecute)
			ctrlIn := make(ExecuteCtrlChan)
			ctrlOut := make(ExecuteCtrlChan)
			go exec.ExecuteWithNotify(session,cmd,arg,ctrlIn,ctrlOut)

			shouldExit := false
outLoop:
			for {
				select {
				case msg,more := <- ctrlOut:
					event := eventTemplate
					event.Event = msg
					if !more {
						shouldExit = true
					}
					if err := stream.Send(&event); err != nil {
						fmt.Errorf("send action event of stream(%v) with event %v error",session,event)
						shouldExit = true
					}
					if shouldExit {
						// either executor runs out of message or send stream with error, exit the loop
						// notify executor by closing the ctrl channel
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

func InitServer(port uint16,executorList []CommandExecute,sourceF []SourceFactory,sinkF []SinkFactory){
	lis, _ := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	server := MediaServer{
		listenPort: port,
		sourceF : sourceF,
		sinkF : sinkF,
	}
	for _,e := range executorList {
		server.registerCommandExecutor(e)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterMediaApiServer(grpcServer,&server)
	grpcServer.Serve(lis)
}









