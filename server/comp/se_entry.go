package comp

import (
	"fmt"
	"strconv"
	"strings"
)

// EntryNode is a basic message provider that simply forward data message to event graph
type EntryNode struct {
	SessionNode

	priority       uint32
	payloadType    string  // if multiple payload types be set, separated them by "," , e.g "payload_type=96,97"
	allowedPayload []uint8 // don't set it directly in nmd
}

//---------------------------------- api & implementation -------------------------------------------

func newEntryNode() SessionAware {
	node := &EntryNode{}
	node.Name = TypeENTRY

	return node
}

func (e *EntryNode) Init() error {
	if len(e.payloadType) == 0 {
		return fmt.Errorf("entry node without payload type")
	}

	payloads := strings.Split(e.payloadType, ",")
	for _, p := range payloads {
		if pt, err := strconv.Atoi(p); err != nil {
			return fmt.Errorf("entry node with invalid payload type:%v", e.payloadType)
		} else {
			if pt < 0 || pt > 255 {
				return fmt.Errorf("entry node: payload type must in range (0,255) but got %v", pt)
			}
			e.allowedPayload = append(e.allowedPayload, uint8(pt))
		}
	}
	logger.Infof("entry node(%v:%v) allow payload types: %v", e.GetNodeScope(), e.GetNodeName(), e.allowedPayload)
	return nil
}

func (e *EntryNode) PushMessage(msg Message) error {
	if msg != nil {
		return e.SendMessage(msg)
	}
	return fmt.Errorf("push nil message")
}

func (e *EntryNode) Priority() uint32 {
	return e.priority
}

func (e *EntryNode) GetName() string {
	return e.Name
}

func (e *EntryNode) CanHandlePayloadType(pt uint8) bool {
	for _, p := range e.allowedPayload {
		if p == pt {
			return true
		}
	}
	return false
}
