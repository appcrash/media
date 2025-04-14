package comp

import (
	"github.com/appcrash/media/server/event"
	"github.com/sirupsen/logrus"
	"reflect"
)

// CAVEAT: always generate message trait first, as node analysis depends on message types
//go:generate go run ../../cmd/gentrait -gen-root -t message -o trait_message_generated.go -v
//go:generate go run ../../cmd/gentrait -gen-root -t node -o trait_node_generated.go -v

var logger *logrus.Entry

func InitLogger(gl *logrus.Logger) {
	logger = gl.WithFields(logrus.Fields{"module": "comp"})
}

func InitBuiltIn() {
	InitMessage()
	InitNode()
}

// MetaType return the type of type/interface object
func MetaType[T any]() reflect.Type {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	return typ
}

type MessageType int
type MessagePostProcessor func(msg Message)

// Message is the base interface of all kinds of message
type Message interface {
	AsEvent() *event.Event
	GetHeader(name string) []byte
	SetHeader(name string, data []byte)
	Type() MessageType
}

type LinkIdentityType uint64 // it differs from linkId as it is unique among whole graph instead of node scope

// LinkPoint is an affiliated output gateway of a node. it is used to communicate with other node as long as a link
// keeps persistent, and gets destroyed once link down. A node can have many LinkPoint as needed
type LinkPoint interface {
	LinkId() int
	Identity() LinkIdentityType
	Owner() SessionAware
	SendMessage(msg Message) error
	MessageTrait() *MessageTrait
	SetEnabled(e bool)
}

// Command is used to send instant info to nodes. it differs from message in:
// --------------------------------------------------------------------------------------------
// |   /      |                 Command              |               Message                  |
// | tx-path  | signalling channel(out of band)      | data channel,in-band (Streamable)      |
// |addressing| no link is created, invoke directly  | create link before data is transmitted |
// |direction | uni(Cast) or bi (Call)               | uni-only, from one node to the other   |
// |  peer    | anyone implements CommandInitiator   |only nodes that have been added to graph|
// |  scope   | only send to local scope(in session) | can cross scope, i.e. a-leg to b-leg   |
// --------------------------------------------------------------------------------------------

// CommandInitiator can be used by session commands and nodes in the event graph to invoke actions of other nodes.
// it provides a unified way to send control message to nodes without any knowledge of links
type CommandInitiator interface {
	// Call send message and wait for the response (block)
	Call(fromNode, toNode string, args []string) (resp []string)

	// Cast send message and don't wait (nonblock)
	Cast(fromNode, toNode string, args []string)
}

// CommandReceiver receives commands from the other components in sync or async manners
type CommandReceiver interface {
	OnCall(fromNode string, args []string) (resp []string)
	OnCast(fromNode string, args []string)
}

// Streamable can create/accept link to/from the other Streamable
type Streamable interface {
	Accept() []MessageType
	Offer() []MessageType
	StreamTo(session, name string, preferredOffer []MessageType) (LinkPoint, error)
	GetLinkPoint(index int) (lp LinkPoint)
	GetLinkPointOfType(messageType MessageType) (lp LinkPoint)

	// GetNodeTypeName return the node trait name
	GetNodeTypeName() string
}

// GraphPhaseAop enable aop calls at each graph-related phase, each sub-interface in it is optional to all nodes
type GraphPhaseAop interface {
	InitializingNode
	PreComposer
	PostComposer
	UnInitializingNode
}

// SessionAware enables node to:
// 1. config its static properties before any event starts to flow
// 2. react to commands when event flowing
// 3. stream data to any other nodes after negotiation
// 4. exit the graph when session ends
type SessionAware interface {
	event.Node
	CommandReceiver
	Streamable
}
