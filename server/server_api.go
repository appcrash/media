package server

import (
	"fmt"
	"github.com/appcrash/media/server/channel"
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/event"
	"github.com/appcrash/media/server/rpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
)

var logger *logrus.Entry

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

	sourceF []SourceFactory
	sinkF   []SinkFactory

	graph *event.Graph

	sessionMutex      sync.Mutex
	sessionMap        map[string]*MediaSession
	executorMutex     sync.Mutex
	simpleExecutorMap map[string]CommandExecute
	streamExecutorMap map[string]CommandExecute
}

type RegisterMore func(s grpc.ServiceRegistrar)

func StartServer(grpcIp string, grpcPort uint16,
	rtpIp string, rtpPortStart uint16, rtpPortEnd uint16,
	regMore RegisterMore,
	executorList []CommandExecute, sourceF []SourceFactory, sinkF []SinkFactory) (err error) {
	var lis net.Listener
	var ip *net.IPAddr
	if lis, err = net.Listen("tcp", fmt.Sprintf("%s:%d", grpcIp, grpcPort)); err != nil {
		logger.Errorf("failed to listen to port(%v) for grpc", grpcPort)
		return
	}
	server := MediaServer{
		rtpServerIpString: rtpIp,
		portPool:          new(portPool),
		sourceF:           sourceF,
		sinkF:             sinkF,
		sessionMap:        make(map[string]*MediaSession),
		simpleExecutorMap: make(map[string]CommandExecute),
		streamExecutorMap: make(map[string]CommandExecute),
		graph:             event.NewEventGraph(),
	}
	if ip, err = net.ResolveIPAddr("ip", rtpIp); err != nil {
		return
	}
	server.init(ip, rtpPortStart, rtpPortEnd)
	server.registerCommandExecutor(&ScriptCommandHandler{}) // built-in script executor
	for _, e := range executorList {
		server.registerCommandExecutor(e)
	}

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
	logger = gl.WithFields(logrus.Fields{"module": "server"})
	event.InitLogger(gl)
	comp.InitLogger(gl)
	channel.InitLogger(gl)
}
