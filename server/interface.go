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
	Execute(s *RtpMediaSession, cmd string, args string) (result []string, err error)
	ExecuteWithNotify(s *RtpMediaSession, args string, ctx context.Context, ctrlOut ExecuteCtrlChan)
	ExecuteWithPush(s *RtpMediaSession, dataIn ExecuteDataChan)
	GetCommandTrait() []CommandTrait
}

// RtpPacketProvider provides data for RTP session
// RTP send loop creates new packets from RtpPacketList then send them.
type RtpPacketProvider interface {
	comp.NodeTraitTag
	PullPacketChannel() <-chan *utils.RtpPacketList
}

// RtpPacketConsumer consumes data from RTP session
// receive loop fetches rtp data packet and feeds it to consumer
type RtpPacketConsumer interface {
	comp.NodeTraitTag
	HandlePacketChannel() chan<- *utils.RtpPacketList
}

// RtpPacketInterceptor can intercept packets bidirectional, that is on the way of graph -> socket or socket -> graph
type RtpPacketInterceptor interface {
	InterceptRtpPacket(pl *utils.RtpPacketList)
}
