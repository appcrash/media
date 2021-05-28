package server

import (
	"context"
	"fmt"
	"github.com/appcrash/media/server/rpc"
	"google.golang.org/grpc"
	"net"
)

type MediaServer struct {
	rpc.UnimplementedMediaApiServer
	listenPort uint16
}

var local, _ = net.ResolveIPAddr("ip", "127.0.0.1")



func (m *MediaServer) PrepareMediaStream(ctx context.Context, peer *rpc.Peer) (*rpc.MediaStream, error) {
	fmt.Println("peer is ",peer)
	port := getNextPort()

	session := createSession(int(port))
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
		Result:   "fail",
	}
	if ok {
		fileName := action.GetAction()
		session := s.(*MediaSession)
		session.sndCtrlC <- fileName
		result.Result = "ok"
	}
	return &result,nil
}

func InitServer(port uint16) {
	lis, _ := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	server := MediaServer{listenPort: port}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterMediaApiServer(grpcServer,&server)
	grpcServer.Serve(lis)
}









