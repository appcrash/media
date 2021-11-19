package server

import (
	"github.com/appcrash/media/server/utils"
)

type TraitEnum uint8
type ExecuteCtrlChan chan string

const (
	CMD_TRAIT_SIMPLE = iota
	CMD_TRAIT_STREAM
)

type CommandTrait struct {
	CmdName  string
	CmdTrait TraitEnum
}

type CommandExecute interface {
	Execute(s *MediaSession, cmd string, args string)
	ExecuteWithNotify(s *MediaSession, cmd string, args string, ctrlIn ExecuteCtrlChan, ctrlOut ExecuteCtrlChan)
	GetCommandTrait() []CommandTrait
}

// Source provides data for RTP session
// either generates data by it own or append/change data from previous source by
// modifying PacketList object, once it is passed through all sources, RTP send loop
// ultimately create new packet from PacketList then send it
// so be careful the order of sources
type Source interface {
	PullData(s *MediaSession, si **utils.PacketList)
}

// Sink consumes data from RTP session
// receive loop fetches rtp data packet and feeds it to all sinks in sink-list
type Sink interface {
	HandleData(s *MediaSession, si *utils.PacketList)
}

type SourceFactory interface {
	NewSource(s *MediaSession) Source
}
type SinkFactory interface {
	NewSink(s *MediaSession) Sink
}
