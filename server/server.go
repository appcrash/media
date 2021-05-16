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

func (m MediaServer) PrepareMediaStream(ctx context.Context, peer *rpc.Peer) (*rpc.MediaStream, error) {
	fmt.Println("peer is ",peer)
	ms := rpc.MediaStream{}
	ms.StreamId = "some_stream_id"
	ms.PeerIp = peer.GetIp()
	//ms.LocalRtpPort = 2000
	ms.PeerRtpPort = peer.GetPort()
	return &ms,nil
}


func InitServer(port uint16) {
	lis, _ := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	server := MediaServer{listenPort: port}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterMediaApiServer(grpcServer,server)
	grpcServer.Serve(lis)

}





