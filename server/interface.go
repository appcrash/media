package server

import (
	"context"
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/rpc"
	"github.com/appcrash/media/server/utils"
)

type TraitEnum uint8
type ExecuteCtrlChan chan string
type ExecuteDataChan chan *rpc.PushData

const (
	CmdTraitSimple = iota
	CmdTraitPullStream
	CmdTraitPushStream
)

type CommandTrait struct {
	CmdName  string
	CmdTrait TraitEnum
}

type CommandExecute interface {
	Execute(s *MediaSession, cmd string, args string) (result []string)
	ExecuteWithNotify(s *MediaSession, args string, ctx context.Context, ctrlOut ExecuteCtrlChan)
	ExecuteWithPush(s *MediaSession, dataIn ExecuteDataChan)
	GetCommandTrait() []CommandTrait
}

// RtpPacketProvider provides data for RTP session
// either generates data by it own or append/change data from previous provider by
// modifying RtpPacketList object, once it is passed through all sources, RTP send loop
// ultimately create new packet from RtpPacketList then send it.
// so be careful the order of providers
type RtpPacketProvider interface {
	comp.NodeTrait
	PullPacketChannel() <-chan *utils.RtpPacketList
}

// RtpPacketConsumer consumes data from RTP session
// receive loop fetches rtp data packet and feeds it to all consumers in consumer-list
type RtpPacketConsumer interface {
	comp.NodeTrait
	HandlePacketChannel() chan<- *utils.RtpPacketList
}

// RtpPacketInterceptor can intercept packets bidirectional, that is on the way of graph -> socket or socket -> graph
type RtpPacketInterceptor interface {
	InterceptRtpPacket(pl *utils.RtpPacketList)
}

type SourceFactory interface {
	NewSource(s *MediaSession) RtpPacketProvider
}
type SinkFactory interface {
	NewSink(s *MediaSession) RtpPacketConsumer
}
