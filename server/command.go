package server

import (
	"context"
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/comp/nmd"
)

// BuiltinCommandHandler provides built-in command to interact with graph by executing nmd script
type BuiltinCommandHandler struct{}

func (sc *BuiltinCommandHandler) Execute(s *MediaSession, _ string, args string) (result []string) {
	if args == "" {
		return
	}
	sessionId := s.GetSessionId().String()
	gt := nmd.NewGraphTopology()
	if err := gt.ParseGraph(sessionId, args); err != nil {
		return
	}
	ctrl := s.GetController()

	// parse rpc command and call/cast to corresponding node, fromNode arg is set to empty string
	for _, call := range gt.GetCallActions() {
		cmd, node := call.Cmd, call.Node
		cmds, err := comp.WithString(cmd)
		if err != nil {
			logger.Errorf("session:%v execute call: node(%v) with %v has error %v", sessionId, node, cmds, err)
			continue
		}
		logger.Infof("session:%v execute call: node %v with %v", sessionId, node, cmds)
		re := ctrl.Call("", node.Name, cmds)
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
		ctrl.Cast("", node.Name, cmds)
	}
	return
}

func (sc *BuiltinCommandHandler) ExecuteWithNotify(s *MediaSession, args string, ctx context.Context, ctrlOut ExecuteCtrlChan) {
	defer func() { close(ctrlOut) }()
	if args == "" {
		return
	}
	sessionId := s.GetSessionId().String()
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
	nodeName := action.NodeName
	outC := make(chan []byte, 32)

	// search the given name node and ensure it is a ChanSink node
	channelSinkNode := s.composer.GetNode(nodeName)
	if channelSinkNode == nil {
		logger.Errorf("execute with notify: link node with name %v failed, no such node", nodeName)
		return
	}
	if node, ok := channelSinkNode.(*comp.ChanSink); ok {
		if err := node.LinkMe(outC); err != nil {
			return
		}
	} else {
		logger.Errorf("node %v is not a channel sink node, can not link", node)
		return
	}

	doneC := ctx.Done()
	for {
		select {
		case rawByte, more := <-outC:
			if !more {
				return
			}
			if rawByte != nil {
				select {
				case ctrlOut <- string(rawByte):
				default:
				}
			}
		case <-doneC:
			return
		}
	}
}

func (sc *BuiltinCommandHandler) ExecuteWithPush(s *MediaSession, dataIn ExecuteDataChan) {
	firstPacket := <-dataIn

	// search the given name node and ensure it is a ChanSrc node
	channelSrcNode := s.composer.GetNode(firstPacket.NodeName)
	if channelSrcNode == nil {
		logger.Errorf("execute with push: link node with name %v failed, no such node", firstPacket.NodeName)
		return
	}
	var inC chan []byte
	defer func() {
		if inC != nil {
			close(inC)
		}
	}()

	if node, ok := channelSrcNode.(*comp.ChanSrc); ok {
		inC = make(chan []byte, 16)
		if err := node.LinkMe(inC); err != nil {
			return
		}
	} else {
		logger.Errorf("node %v is not a channel sink node, can not link", node)
		return
	}

	for {
		select {
		case pushData, more := <-dataIn:
			if !more {
				return
			}
			// convert rpc push data to raw byte and forward to graph
			data := pushData.GetData()
			if len(data) == 0 {
				continue
			}
			select {
			case inC <- data:
			default:
			}
		}
	}
}

func (sc *BuiltinCommandHandler) GetCommandTrait() []CommandTrait {
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
