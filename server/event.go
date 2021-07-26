package server

// event network structures

type Event struct {
	cmd string
	obj interface{}
}

type EventLink interface {
	Notify(evt *Event) bool
	TearDown()
}

type EventNode interface {
	GetName() string
	GetFromNode() []string
	GetToNode() []string

	Receive(evt *Event) bool
}

type EventDag interface {
	FindNode(scope string, name string) *EventNode
	Link(from *EventNode, to *EventNode) *EventLink
}
