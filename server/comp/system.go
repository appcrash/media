package comp

import (
	"github.com/appcrash/media/server/event"
)

const SYS_NODE_SCOPE = "sys"

// SystemNode is the base class of all system-wide nodes that keep running throughout graph's lifetime
type SystemNode struct {
	Id
	delegate *event.NodeDelegate
}

func (s *SystemNode) GetNodeName() string {
	return s.Name
}

func (s *SystemNode) GetNodeScope() string {
	return SYS_NODE_SCOPE
}

func (s *SystemNode) OnEvent(evt *event.Event) {
	logger.Debugf("system node(%v) got event with cmd:%v", s.GetNodeName(), evt.GetCmd())
}

func (s *SystemNode) OnLinkDown(linkId int, scope string, nodeName string) {
	logger.Debugf("system node got link down (%v:%v) => (%v:%v) ", s.GetNodeScope(), s.GetNodeName(), scope, nodeName)
}

func (s *SystemNode) OnEnter(delegate *event.NodeDelegate) {
	logger.Debugf("system node(%v) enters graph", s.GetNodeName())
	s.delegate = delegate
}

func (s *SystemNode) OnExit() {
	// system node should never exit
	logger.Errorf("system node(%v) exits graph", s.GetNodeName())
}
