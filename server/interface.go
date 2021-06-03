package server

type TraitEnum uint8
type ExecuteCtrlChan chan string

const (
	CMD_TRAIT_SIMPLE = iota
	CMD_TRAIT_STREAM
)

type CommandTrait struct {
	CmdName string
	CmdTrait TraitEnum
}

type CommandExecute interface {
	Execute(s *MediaSession,cmd string,args string)
	ExecuteWithNotify(s *MediaSession,cmd string,args string,ctrlIn ExecuteCtrlChan,ctrlOut ExecuteCtrlChan)
	GetCommandTrait() []CommandTrait
}

type Source interface {
	PullData(s *MediaSession) (data []byte,timestampAdvanced uint32)
}

type Sink interface {
	HandleData(s *MediaSession,data []byte) (shouldContinue bool)
}

type SourceFactory interface {
	NewSource(s *MediaSession) Source
}
type SinkFactory interface {
	NewSink(s *MediaSession) Sink
}



