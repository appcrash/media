package comp

import (
	"errors"
	"github.com/appcrash/media/server/event"
)

// EntryNode is a basic message provider that simply forward data message to event graph
type EntryNode struct {
	SessionNode
}

//---------------------------------- api & implementation -------------------------------------------

func newEntryNode() SessionAware {
	node := &EntryNode{}
	node.Name = TYPE_ENTRY
	return node
}

func (e *EntryNode) PushMessage(data DataMessage) error {
	if e.dataLinkId >= 0 && data != nil {
		evt := event.NewEvent(DATA_OUTPUT, data)
		if !e.delegate.Deliver(e.dataLinkId, evt) {
			return errors.New("failed to deliver message")
		}
	}
	return nil
}

func (e *EntryNode) GetName() string {
	return e.Name
}
