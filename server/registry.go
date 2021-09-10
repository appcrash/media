package server

import (
	"fmt"
	"sync"
)

// APIs that allow plugging in method to:
// 1. handle command(take new actions), listen to state change
// 2. handle incoming data(audio,video), which is a sink
// 3. generate outgoing data, which is a source

func (srv *MediaServer) registerCommandExecutor(e CommandExecute) {
	cmdTrait := e.GetCommandTrait()
	var cm *sync.Map

	for _, trait := range cmdTrait {
		switch trait.CmdTrait {
		case CMD_TRAIT_SIMPLE:
			cm = &srv.simpleExecutorMap
		case CMD_TRAIT_STREAM:
			cm = &srv.streamExecutorMap
		default:
			fmt.Errorf("register command executor with wrong trait: %v", trait.CmdTrait)
			return
		}
		if cmd, ok := cm.Load(trait.CmdName); ok {
			fmt.Errorf("regsiter execute with command %v already registered", cmd)
			continue
		} else {
			fmt.Printf("register execute with command %v\n", trait.CmdName)
			cm.Store(trait.CmdName, e)
		}
	}
}

func (srv *MediaServer) getExecutorFor(cmd string) (needNotify bool, ce CommandExecute) {
	if e, ok := srv.simpleExecutorMap.Load(cmd); ok {
		needNotify = false
		ce = e.(CommandExecute)
	}
	if e, ok := srv.streamExecutorMap.Load(cmd); ok {
		needNotify = true
		ce = e.(CommandExecute)
	}
	needNotify = false
	ce = nil
	return
}
