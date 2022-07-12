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

	sourceF         []SourceFactory
	sinkF           []SinkFactory
	sessionListener []SessionListener

	graph *event.Graph

	sessionMutex sync.Mutex
	sessionMap   map[string]*MediaSession

	simpleExecutorMap map[string]CommandExecute
	streamExecutorMap map[string]CommandExecute
}

type Config struct {
	RtpIp               string
	StartPort, EndPort  uint16
	ExecutorList        []CommandExecute
	SourceFactoryList   []SourceFactory
	SinkFactoryList     []SinkFactory
	SessionListenerList []SessionListener

	GrpcIp           string
	GrpcPort         uint16
	GrpcRegisterMore RegisterMore
}

type RegisterMore func(s grpc.ServiceRegistrar)

// SessionListener methods are called concurrently, so it must be goroutine safe
type SessionListener interface {
	OnSessionCreated(s *MediaSession)
	OnSessionUpdated(s *MediaSession)
	OnSessionStarted(s *MediaSession)
	OnSessionStopped(s *MediaSession)
}

// BaseSessionListener facilitates building listener if not interested in all events
type BaseSessionListener struct{}

func (b *BaseSessionListener) OnSessionCreated(s *MediaSession) {}
func (b *BaseSessionListener) OnSessionUpdated(s *MediaSession) {}
func (b *BaseSessionListener) OnSessionStarted(s *MediaSession) {}
func (b *BaseSessionListener) OnSessionStopped(s *MediaSession) {}

func StartServer(c *Config) (err error) {
	var lis net.Listener
	var ip *net.IPAddr
	if lis, err = net.Listen("tcp", fmt.Sprintf("%s:%d", c.GrpcIp, c.GrpcPort)); err != nil {
		logger.Errorf("failed to listen to port(%v) for grpc", c.GrpcPort)
		return
	}
	rtpIp, rtpStartPort, rtpEndPort := c.RtpIp, c.StartPort, c.EndPort
	server := MediaServer{
		rtpServerIpString: rtpIp,
		portPool:          new(portPool),
		sourceF:           c.SourceFactoryList,
		sinkF:             c.SinkFactoryList,
		sessionListener:   c.SessionListenerList,
		sessionMap:        make(map[string]*MediaSession),

		// read-only maps once executors registered
		simpleExecutorMap: make(map[string]CommandExecute),
		streamExecutorMap: make(map[string]CommandExecute),

		graph: event.NewEventGraph(),
	}
	if ip, err = net.ResolveIPAddr("ip", rtpIp); err != nil {
		return
	}
	server.init(ip, rtpStartPort, rtpEndPort)
	server.registerCommandExecutor(&ScriptCommandHandler{}) // built-in script executor
	for _, e := range c.ExecutorList {
		server.registerCommandExecutor(e)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterMediaApiServer(grpcServer, &server)
	if c.GrpcRegisterMore != nil {
		c.GrpcRegisterMore(grpcServer)
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
