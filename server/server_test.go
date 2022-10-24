package server_test

import (
	"context"
	"fmt"
	"github.com/appcrash/media/server"
	"github.com/appcrash/media/server/channel"
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/rpc"
	"github.com/appcrash/media/server/utils"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	grpcIp   = "127.0.0.1"
	grpcPort = 5678
)

type echo struct {
	comp.SessionNode
	comp.ChannelNode

	channel chan *utils.RtpPacketList
}

func (n *echo) PullPacketChannel() <-chan *utils.RtpPacketList {
	return n.channel
}

func (n *echo) HandlePacketChannel() chan<- *utils.RtpPacketList {
	return n.channel
}

func (n *echo) OnCast(_ string, args []string) {
	n.NotifyInstance(strings.Join(args, "#"))
}

type recvFunc func(event *rpc.SystemEvent)

type client struct {
	instanceId  string
	conn        *grpc.ClientConn
	mediaClient rpc.MediaApiClient
	sysStream   rpc.MediaApi_SystemChannelClient
}

func (c *client) connect(onReceive recvFunc) {
	var opts = []grpc.DialOption{grpc.WithInsecure()}
	var callOpts []grpc.CallOption
	conn, err1 := grpc.Dial(fmt.Sprintf("%v:%v", grpcIp, grpcPort), opts...)
	if err1 != nil {
		panic(err1)
	}
	c.conn = conn
	c.mediaClient = rpc.NewMediaApiClient(conn)

	// register myself first
	stream, err2 := c.mediaClient.SystemChannel(context.Background(), callOpts...)
	if err2 != nil {
		panic(err2)
	}
	stream.Send(&rpc.SystemEvent{
		Cmd:        rpc.SystemCommand_REGISTER,
		InstanceId: c.instanceId,
	})
	c.sysStream = stream

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				return
			}
			if err != nil {
				log.Fatalf("receive error: %v", err)
				return
			}
			onReceive(in)
		}
	}()
}

func (c *client) keepalive(ctx context.Context) {
	ticker := time.NewTicker(channel.KeepAliveCheckDuration)
	for {
		select {
		case <-ticker.C:
			c.sysStream.Send(&rpc.SystemEvent{
				Cmd:        rpc.SystemCommand_KEEPALIVE,
				InstanceId: c.instanceId,
			})
		case <-ctx.Done():
			return
		}
	}
}

func (c *client) close() {
	c.conn.Close()
}

func startServer() {
	config := &server.Config{
		RtpIp:     "127.0.0.1",
		StartPort: 10000,
		EndPort:   20000,
		GrpcIp:    grpcIp,
		GrpcPort:  grpcPort,
	}
	if err := server.StartServer(config); err != nil {
		panic(err)
	}
}

func initComposer() {
	comp.InitBuiltIn()
	comp.RegisterNodeTrait(comp.NT[echo]("echo", func() comp.SessionAware {
		n := &echo{channel: make(chan *utils.RtpPacketList, 2)}
		n.Trait, _ = comp.NodeTraitOfType("echo")
		return n
	}))
}

func TestMain(m *testing.M) {
	initComposer()
	go startServer()
	time.Sleep(1 * time.Second)
	os.Exit(m.Run())
}

func TestInstanceKeepalive(t *testing.T) {
	ch := make(chan struct{})
	instanceId := "test"
	c := &client{instanceId: instanceId}
	c.connect(func(event *rpc.SystemEvent) {
		close(ch)
	})
	c.sysStream.Send(&rpc.SystemEvent{
		Cmd:        rpc.SystemCommand_KEEPALIVE,
		InstanceId: instanceId,
	})
	c.sysStream.CloseSend()
	<-ch
}

func TestInstanceReconnect(t *testing.T) {
	instanceId := "reconnect"
	c := &client{instanceId: instanceId}
	c.connect(func(event *rpc.SystemEvent) {
		t.Logf("received %v", event)
	})
	time.Sleep(channel.KeepAliveTimeout)
	c = &client{instanceId: instanceId}
	c.connect(func(event *rpc.SystemEvent) {
		t.Logf("after reconnect, receive %v", event)
	})
	ctx, cancel := context.WithCancel(context.Background())
	go c.keepalive(ctx)
	time.Sleep(10 * time.Second)
	cancel()
	time.Sleep(channel.KeepAliveTimeout)
}

func TestOperateSession(t *testing.T) {
	instanceId := "new_session"
	c := &client{instanceId: instanceId}
	c.connect(func(event *rpc.SystemEvent) {
		t.Logf("recv %v", event)
	})
	ctx, cancel := context.WithCancel(context.Background())
	go c.keepalive(ctx)
	var opts []grpc.CallOption
	session, err := c.mediaClient.PrepareSession(ctx, &rpc.CreateParam{
		PeerIp:   "127.0.0.1",
		PeerPort: 2000,
		Codecs: []*rpc.CodecInfo{{
			PayloadNumber: 8,
			PayloadType:   rpc.CodecType_PCM_ALAW,
			CodecParam:    "",
		}},
		GraphDesc:  "[echo]",
		InstanceId: instanceId,
	}, opts...)
	if err != nil {
		panic(err)
	}
	if _, err = c.mediaClient.UpdateSession(ctx, &rpc.UpdateParam{SessionId: session.SessionId, PeerPort: 3000}, opts...); err != nil {
		panic(err)
	}
	if _, err = c.mediaClient.StartSession(ctx, &rpc.StartParam{SessionId: session.SessionId}, opts...); err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)
	if _, err = c.mediaClient.ExecuteAction(ctx, &rpc.Action{
		SessionId: session.SessionId,
		Cmd:       "exec",
		CmdArg:    "[echo] <-- 'a b c d'",
	}, opts...); err != nil {
		panic(err)
	}
	time.Sleep(server.SessionAuditPeriod)
	if _, err = c.mediaClient.StopSession(ctx, &rpc.StopParam{SessionId: session.SessionId}, opts...); err != nil {
		panic(err)
	}
	cancel()
	time.Sleep(1 * time.Second)
}
