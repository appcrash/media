package server

import (
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/comp/nmd"
	"github.com/appcrash/media/server/event"
)

// ScriptCommandHandler provides built-in command to execute nmd script
type ScriptCommandHandler struct{}

func (sc *ScriptCommandHandler) Execute(s *MediaSession, _ string, args string) (result []string) {
	if args == "" {
		return
	}
	sessionId := s.GetSessionId()
	gt := nmd.NewGraphTopology()
	if err := gt.ParseGraph(sessionId, args); err != nil {
		return
	}
	ctrl := s.GetController()
	for _, call := range gt.GetCallActions() {
		cmd, node := call.Cmd, call.Node
		cmds, err := comp.WithString(cmd)
		if err != nil {
			logger.Errorf("session:%v execute call: node(%v) with %v has error %v", sessionId, node, cmds, err)
			continue
		}
		logger.Infof("session:%v execute call: node(%v) with %v", sessionId, node, cmds)
		re := ctrl.Call(node.Scope, node.Name, cmds)
		if len(result) != 0 {
			result = append(result, "\n")
		}
		result = append(result, re...)
	}
	for _, cast := range gt.GetCastActions() {
		cmd, node := cast.Cmd, cast.Node
		cmds, err := comp.WithString(cmd)
		if err != nil {
			logger.Errorf("session:%v execute cast: node(%v) with %v has error %v", sessionId, node, cmds, err)
			continue
		}
		logger.Infof("session:%v execute cast: node(%v) with %v", sessionId, node, cmds)
		ctrl.Cast(node.Scope, node.Name, cmds)
	}
	return
}

func (sc *ScriptCommandHandler) ExecuteWithNotify(s *MediaSession, cmd string, args string, ctrlIn ExecuteCtrlChan, ctrlOut ExecuteCtrlChan) {
	defer func() { close(ctrlOut) }()
	if args == "" {
		return
	}
	sessionId := s.GetSessionId()
	gt := nmd.NewGraphTopology()
	if err := gt.ParseGraph(sessionId, args); err != nil {
		return
	}
	sinkAction := gt.GetSinkActions()
	if sinkAction == nil || len(sinkAction) != 1 {
		// only support one sink action
		return
	}
	action := sinkAction[0]
	chName := action.ChannelName
	ch := make(chan *event.Event, 10)
	if err := s.composer.LinkChannel(chName, ch); err != nil {
		logger.Errorf("execute with notify: link channel error: %v", err)
		return
	}
	defer func() {
		// tell graph to stop sending event to this channel
		if err := s.composer.UnlinkChannel(chName); err != nil {
			logger.Errorln(err)
		}
	}()

	for {
		select {
		case evt, more := <-ch:
			if !more {
				return
			}
			if evt == nil || evt.GetObj() == nil {
				continue
			}
			select {
			case ctrlOut <- string(evt.GetObj().(comp.RawByteMessage)):
			default:
			}
		case _, more := <-ctrlIn:
			if !more {
				return
			}
		}
	}
}

func (sc *ScriptCommandHandler) ExecuteWithPush(s *MediaSession, dataIn ExecuteDataChan) {
	controller := s.GetController()
	for {
		select {
		case pushData, more := <-dataIn:
			if !more {
				return
			}
			// sanity check before pushing
			nodeName := pushData.GetNodeName()
			data := pushData.GetData()
			if len(nodeName) == 0 || len(data) == 0 {
				continue
			}
			controller.PushData(nodeName, pushData.GetMsgType(), data)
		}
	}
}

func (sc *ScriptCommandHandler) GetCommandTrait() []CommandTrait {
	return []CommandTrait{
		{
			"exec",
			CmdTraitSimple,
		},
		{
			"pull_stream",
			CmdTraitPullStream,
		},
		{
			"push_stream",
			CmdTraitPushStream,
		},
	}
}
