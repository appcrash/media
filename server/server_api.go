package server

import (
	"fmt"
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/event"
	"github.com/appcrash/media/server/rpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
)

// Initialize all packages logger
func init() {
	gl := logrus.New()
	InitServerLogger(gl)
}

type MediaServer struct {
	rpc.UnimplementedMediaApiServer
	rtpServerIpString string
	rtpServerIpAddr   *net.IPAddr
	portPool          *portPool

	simpleExecutorMap sync.Map
	streamExecutorMap sync.Map

	sourceF []SourceFactory
	sinkF   []SinkFactory

	graph *event.EventGraph
}

type RegisterMore func(s grpc.ServiceRegistrar)

func StartServer(grpcIp string, grpcPort uint16,
	rtpIp string, rtpPortStart uint16, rtpPortEnd uint16,
	regMore RegisterMore,
	executorList []CommandExecute, sourceF []SourceFactory, sinkF []SinkFactory) (err error) {
	lis, _ := net.Listen("tcp", fmt.Sprintf("%s:%d", grpcIp, grpcPort))
	server := MediaServer{
		rtpServerIpString: rtpIp,
		portPool:          new(portPool),
		sourceF:           sourceF,
		sinkF:             sinkF,
	}
	for _, e := range executorList {
		server.registerCommandExecutor(e)
	}
	if server.rtpServerIpAddr, err = net.ResolveIPAddr("ip", rtpIp); err != nil {
		return
	}
	server.portPool.init(rtpPortStart, rtpPortEnd)
	server.graph = event.NewEventGraph()

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterMediaApiServer(grpcServer, &server)
	if regMore != nil {
		regMore(grpcServer)
	}
	grpcServer.Serve(lis)
	return nil
}

// InitServerLogger can be called multiple times before server starts to override default logger
func InitServerLogger(gl *logrus.Logger) {
	comp.InitLogger(gl)
}
