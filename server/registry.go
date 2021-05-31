package server

import (
	"fmt"
	"sync"
)

// APIs that allow plugging in method to:
// 1. handle command(take new actions), listen to state change
// 2. handle incoming data(audio,video), which is a sinker
// 3. generate outgoing data, which is a sourcer

func (m *MediaServer) registerCommandExecutor(e CommandExecutor) {
	cmdTrait := e.GetCommandTrait()
	var cm *sync.Map

	for _,trait := range cmdTrait {
		switch trait.CmdTrait {
		case CMD_TRAIT_SIMPLE:
			cm = &m.simpleExecutorMap
		case CMD_TRAIT_STREAM:
			cm = &m.streamExecutorMap
		default:
			fmt.Errorf("register command executor with wrong trait: %v",trait.CmdTrait)
			return
		}
		if cmd,ok := cm.Load(trait.CmdName); ok {
			fmt.Errorf("regsiter executor with command %v already registered",cmd)
			continue
		} else {
			cm.Store(cmd,e)
		}
	}
}

func (m *MediaServer) getExecutorFor(cmd string) (needNotify bool,ce CommandExecutor) {
	if e,ok := m.simpleExecutorMap.Load(cmd); ok {
		needNotify = false
		ce =  e.(CommandExecutor)
	}
	if e,ok := m.streamExecutorMap.Load(cmd); ok {
		needNotify = true
		ce =  e.(CommandExecutor)
	}
	needNotify = false
	ce = nil
	return
}
