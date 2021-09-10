package server

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
// either generates data by it own or append/change data from previous source
// data from last source in source-list would be used by RTP send loop ultimately
// so be careful to order of sources
// NOTE: the first source would get previousData as nil, previousTs as 0
type Source interface {
	PullData(s *MediaSession, previousData []byte, previousTs uint32) (data []byte, timestampAdvanced uint32)
}

// Sink consumes data from RTP session
// receive loop fetches rtp payload and feeds it to all sink in sink-list
// each sink should return true if the payload can be used by following sinks or
// return false to stop this process (so following sinks can not get the payload)
type Sink interface {
	HandleData(s *MediaSession, previousData []byte) (data []byte, shouldContinue bool)
}

type SourceFactory interface {
	NewSource(s *MediaSession) Source
}
type SinkFactory interface {
	NewSink(s *MediaSession) Sink
}
