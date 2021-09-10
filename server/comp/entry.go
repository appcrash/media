package comp

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
	if data != nil {
		e.SendData(data)
	}
	return nil
}

func (e *EntryNode) GetName() string {
	return e.Name
}
