package comp

const entryNodeDefaultPayloadType = 255

// EntryNode is a basic message provider that simply forward data message to event graph
type EntryNode struct {
	SessionNode

	payloadType uint8
}

//---------------------------------- api & implementation -------------------------------------------

func newEntryNode() SessionAware {
	node := &EntryNode{}
	node.Name = TYPE_ENTRY
	node.payloadType = entryNodeDefaultPayloadType
	return node
}

func (e *EntryNode) PushMessage(data DataMessage) error {
	if data != nil {
		_ = e.SendData(data)
	}
	return nil
}

func (e *EntryNode) Priority() uint32 {
	return uint32(e.payloadType)
}

func (e *EntryNode) GetName() string {
	return e.Name
}

func (e *EntryNode) CanHandlePayloadType(pt uint8) bool {
	if e.payloadType == pt {
		return true
	}
	if e.payloadType == entryNodeDefaultPayloadType {
		return true
	}
	return false
}
