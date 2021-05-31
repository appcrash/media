package server

type TraitEnum uint8
type ExecutorCtrlChan chan string

const (
	CMD_TRAIT_SIMPLE = iota
	CMD_TRAIT_STREAM
)

type CommandTrait struct {
	CmdName string
	CmdTrait TraitEnum
}

type CommandExecutor interface {
	Execute(s *MediaSession,cmd string,args string)
	ExecuteWithNotify(s *MediaSession,cmd string,args string,ctrlIn ExecutorCtrlChan,ctrlOut ExecutorCtrlChan)
	GetCommandTrait() []CommandTrait
}

type Sourcer interface {
	// TODO: add methods
}

type Sinker interface {
	Init(s *MediaSession)
	HandleData(s *MediaSession,data []byte) (shouldContinue bool)
}




