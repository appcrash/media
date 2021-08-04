package event

// event network structures

type Event struct {
	cmd int
	obj interface{}
}

func (e *Event) GetCmd() int {
	return e.cmd
}

func (e *Event) GetObj() interface{} {
	return e.obj
}

type Node interface {
	GetNodeName() string
	GetNodeScope() string

	// following methods (On***) are never invoke concurrently
	// all of them are called in multiple separate gorouting sequentially

	// normal event handling
	OnEvent(evt *Event)

	// dlink status change
	OnLinkUp(linkId int, scope string, nodeName string)
	OnLinkDown(linkId int, scope string, nodeName string)

	// after sucessfully added to graph
	OnEnter(delegate *NodeDelegate)

	// the finalizing method after node exits graph
	OnExit()

	// optional attributes if following fields defined in the node
	// -----------------------------------------------------------
	// maxLink int:
	//   override default max output link number
	// dataChannelSize int:
	//   override default buffered event channel size
	// deliveryTimeout time.Duration:
	//   override default event delivery timeout
}
