package server

import (
	"fmt"
	"github.com/appcrash/media/server/channel"
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/event"
	"github.com/appcrash/media/server/rpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"sync"
	"time"
)

var logger *logrus.Entry

func init() {
	gl := logrus.New()
	// Initialize all packages logger
	InitServerLogger(gl)

	// Randomize session id count
	rand.Seed(time.Now().UnixNano())
	sessionIdCounter = rand.Uint32()
}

type MediaServer struct {
	rpc.UnimplementedMediaApiServer
	rtpServerIpString string
	rtpServerIpAddr   *net.IPAddr
	portPool          *portPool
	sessionListener   []SessionListener

	graph *event.Graph

	sessionMutex sync.Mutex
	sessionMap   map[SessionIdType]*MediaSession

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
type StartServerFunc func()
type StopServerFunc func()

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

func NewServer(c *Config) (start StartServerFunc, stop StopServerFunc, err error) {
	var lis net.Listener
	var ip *net.IPAddr

	if lis, err = net.Listen("tcp", fmt.Sprintf("%s:%d", c.GrpcIp, c.GrpcPort)); err != nil {
		logger.Errorf("failed to listen to port(%v) for grpc", c.GrpcPort)
		return
	}
	rtpIp, rtpStartPort, rtpEndPort := c.RtpIp, c.StartPort, c.EndPort
	server := MediaServer{
		rtpServerIpString: rtpIp,
		portPool:          newPortPool(),
		sessionListener:   c.SessionListenerList,
		sessionMap:        make(map[SessionIdType]*MediaSession),

		// read-only maps once executors registered
		simpleExecutorMap: make(map[string]CommandExecute),
		streamExecutorMap: make(map[string]CommandExecute),

		graph: event.NewEventGraph(),
	}
	if ip, err = net.ResolveIPAddr("ip", rtpIp); err != nil {
		return
	}
	server.init(ip, rtpStartPort, rtpEndPort)
	server.registerCommandExecutor(&BuiltinCommandHandler{}) // built-in script executor
	for _, e := range c.ExecutorList {
		server.registerCommandExecutor(e)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterMediaApiServer(grpcServer, &server)
	if c.GrpcRegisterMore != nil {
		c.GrpcRegisterMore(grpcServer)
	}

	start = func() {
		logger.Infof("starting media server")
		grpcServer.Serve(lis)
	}
	stop = func() {
		logger.Infof("try to gracefully stop media server")
		grpcServer.GracefulStop()
		logger.Infof("media server has stopped")
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
