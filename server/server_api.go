package server

import (
	"fmt"
	"github.com/appcrash/media/server/channel"
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/event"
	"github.com/appcrash/media/server/rpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"math/rand/v2"
	"net"
	"sync"
)

var logger *logrus.Entry

func init() {
	gl := logrus.New()
	// Initialize all packages logger
	InitServerLogger(gl)

	sessionIdCounter = rand.Uint32()
}

type GrpcServer struct {
	rpc.UnimplementedMediaApiServer
	rtpServerIpString string
	rtpServerIpAddr   *net.IPAddr
	portPool          *PortPool
	sessionListener   []SessionListener

	graph *event.Graph

	sessionMutex sync.Mutex
	sessionMap   map[SessionIdType]*RtpMediaSession

	simpleExecutorMap map[string]CommandExecute
	streamExecutorMap map[string]CommandExecute
}

type Config struct {
	RtpIp               string
	StartPort, EndPort  uint16
	ExecutorList        []CommandExecute
	SessionListenerList []SessionListener

	GrpcIp           string
	GrpcPort         uint16
	GrpcRegisterMore RegisterMore
}

type RegisterMore func(s grpc.ServiceRegistrar)
type StartServerFunc func() error
type StopServerFunc func()

// SessionListener methods are called concurrently, so it must be goroutine safe
type SessionListener interface {
	OnSessionCreated(s *RtpMediaSession)
	OnSessionUpdated(s *RtpMediaSession)
	OnSessionStarted(s *RtpMediaSession)
	OnSessionStopped(s *RtpMediaSession)
}

// BaseSessionListener facilitates building listener if not interested in all events
type BaseSessionListener struct{}

func (b *BaseSessionListener) OnSessionCreated(_ *RtpMediaSession) {}
func (b *BaseSessionListener) OnSessionUpdated(_ *RtpMediaSession) {}
func (b *BaseSessionListener) OnSessionStarted(_ *RtpMediaSession) {}
func (b *BaseSessionListener) OnSessionStopped(_ *RtpMediaSession) {}

func NewGrpcServer(c *Config) (start StartServerFunc, stop StopServerFunc, err error) {
	var lis net.Listener
	var ip *net.IPAddr

	if lis, err = net.Listen("tcp", fmt.Sprintf("%s:%d", c.GrpcIp, c.GrpcPort)); err != nil {
		logger.Errorf("failed to listen to port(%v) for grpc", c.GrpcPort)
		return
	}
	rtpIp, rtpStartPort, rtpEndPort := c.RtpIp, c.StartPort, c.EndPort
	server := GrpcServer{
		rtpServerIpString: rtpIp,
		portPool:          NewPortPool(),
		sessionListener:   c.SessionListenerList,
		sessionMap:        make(map[SessionIdType]*RtpMediaSession),

		// read-only maps once executors registered
		simpleExecutorMap: make(map[string]CommandExecute),
		streamExecutorMap: make(map[string]CommandExecute),

		graph: event.NewEventGraph(),
	}
	if ip, err = net.ResolveIPAddr("ip", rtpIp); err != nil {
		return
	}
	server.init(ip, rtpStartPort, rtpEndPort)
	_ = server.registerCommandExecutor(&BuiltinCommandHandler{}) // built-in script executor
	for _, e := range c.ExecutorList {
		if err = server.registerCommandExecutor(e); err != nil {
			logger.Errorf("failed to register command executor: %v", err)
			return
		}
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterMediaApiServer(grpcServer, &server)
	if c.GrpcRegisterMore != nil {
		c.GrpcRegisterMore(grpcServer)
	}

	start = func() error {
		logger.Infof("starting GRPC/RTP endpoint")
		return grpcServer.Serve(lis)
	}
	stop = func() {
		logger.Infof("try to gracefully stop GRPC/RTP endpoint")
		grpcServer.GracefulStop()
		logger.Infof("GRPC/RTP endpoint has stopped")
	}
	return
}

// InitServerLogger can be called multiple times before server starts to override default logger
func InitServerLogger(gl *logrus.Logger) {
	logger = gl.WithFields(logrus.Fields{"module": "server"})
	event.InitLogger(gl)
	comp.InitLogger(gl)
	channel.InitLogger(gl)
}
