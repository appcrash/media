package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/appcrash/media/server/rpc"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type MediaServer struct {
	rpc.UnimplementedMediaApiServer
	listenPort        uint16
	simpleExecutorMap sync.Map
	streamExecutorMap sync.Map

	sinkerList []Sinker
}

var local, _ = net.ResolveIPAddr("ip", "127.0.0.1")



func (m *MediaServer) PrepareMediaStream(ctx context.Context, peer *rpc.Peer) (*rpc.MediaStream, error) {
	fmt.Println("peer is ",peer)
	port := getNextPort()

	session := createSession(int(port))
	session.sinkerList = m.sinkerList
	for _,s := range session.sinkerList {
		s.Init(session)
	}
	session.AddRemote(peer.GetIp(),int(peer.GetPort()))
	session.StartSession()

	ms := rpc.MediaStream{}
	ms.StreamId = session.sessionId
	ms.PeerIp = peer.GetIp()
	ms.LocalRtpPort = uint32(port)
	ms.PeerRtpPort = peer.GetPort()

	return &ms,nil
}

func (m *MediaServer) ExecuteAction(ctx context.Context,action *rpc.MediaAction) (*rpc.MediaActionResult, error) {
	sessionId := action.StreamId
	s,ok := sessionMap.Get(sessionId)
	result := rpc.MediaActionResult{
		StreamId: sessionId,
	}

	if ok {
		session := s.(*MediaSession)
		cmd := action.GetCmd()
		arg := action.GetCmdArg()
		if e,ok1 := m.simpleExecutorMap.Load(cmd); ok1 {
			exec := e.(CommandExecutor)
			exec.Execute(session,cmd,arg)
			result.State = "ok"
			return &result,nil
		}
	}
	result.State = "cmd executor not exist"
	return &result,nil
}


func (m *MediaServer) ExecuteActionWithNotify(action *rpc.MediaAction,stream rpc.MediaApi_ExecuteActionWithNotifyServer) error {
	sessionId := action.StreamId
	s,ok := sessionMap.Get(sessionId)
	eventTemplate := rpc.MediaActionEvent{
		StreamId: sessionId,
	}

	if ok {
		session := s.(*MediaSession)
		cmd := action.GetCmd()
		arg := action.GetCmdArg()
		if e,ok1 := m.streamExecutorMap.Load(cmd); ok1 {
			exec := e.(CommandExecutor)
			ctrlIn := make(ExecutorCtrlChan)
			ctrlOut := make(ExecutorCtrlChan)
			go exec.ExecuteWithNotify(session,cmd,arg,ctrlIn,ctrlOut)

			shouldExit := false
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
						break
					}
				}
			}
		}

		return nil
	}

	return errors.New("cmd not exist")
}


func InitServer(port uint16,executorList []CommandExecutor,sinkerList []Sinker){
	lis, _ := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	server := MediaServer{
		listenPort: port,
	}
	if sinkerList != nil {
		server.sinkerList = sinkerList
	}
	for _,e := range executorList {
		server.registerCommandExecutor(e)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterMediaApiServer(grpcServer,&server)
	grpcServer.Serve(lis)
}









