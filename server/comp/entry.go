package comp

import (
	"github.com/appcrash/media/server/event"
)

// EntryNode is the place where input data flow into event graph in the session
type EntryNode struct {
	SessionNode

	linkId int
}

func (e *EntryNode) OnEvent(evt *event.Event) {
	obj := evt.GetObj()
	if obj == nil {
		return
	}
	switch evt.GetCmd() {
	case CMD_GENERIC_SET_ROUTE:
		if c, ok := obj.(*GenericRouteCommand); ok {
			e.setOutputNode(c.SessionId, c.Name)
		}

	}
}

func (e *EntryNode) OnLinkUp(linkId int, scope string, nodeName string) {
	if linkId >= 0 {
		e.linkId = linkId
	}
}

//---------------------------------- api & implementation -------------------------------------------

func newEntryNode() *EntryNode {
	node := &EntryNode{
		linkId: -1,
	}
	node.Name = TYPE_ENTRY
	return node
}

func (e *EntryNode) setOutputNode(sessionId, name string) {
	e.delegate.RequestLinkUp(sessionId, name)
}

func (e *EntryNode) Input(data []byte) bool {
	if e.linkId >= 0 && data != nil {
		evt := event.NewEvent(DATA_OUTPUT, data)
		return e.delegate.Delivery(e.linkId, evt)
	}
	return false
}
