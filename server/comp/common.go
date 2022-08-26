package comp

import (
	"errors"
	"fmt"
	"github.com/appcrash/media/server/event"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Entry

func InitLogger(gl *logrus.Logger) {
	logger = gl.WithFields(logrus.Fields{"module": "comp"})
}

const (
	TypeENTRY    = "entry"
	TypePUBSUB   = "pubsub"
	TypeDISPATCH = "dispatch"
)

// LinkPoint is an affiliated in/out gateway of a node. it is used to communicate with other node as long as a link
// keeps persistent, and gets destroyed once link down. A node can have many LinkPoint as needed
type LinkPoint interface {
	LinkId() int
	Owner() SessionAware
	SetPeer(point LinkPoint)
	Peer() LinkPoint
	SendMessage(msg Message) error
	MessageTrait() MessageTrait
}

// Command is used to send instant info  to nodes. it differs from message in:
// --------------------------------------------------------------------------------------------
// |   /      |                 Command              |               Message                  |
// | tx-path  | signalling channel                   | data channel (implements Streamable)   |
// |addressing| no link is created, invoke directly  | create link before data is transmitted |
// |direction | uni(Cast) or bi (Call)               | uni-only, from one node to the other   |
// |  peer    | anyone implements CommandInitiator   | only nodes that has been added to graph|
// |  scope   | only send to local scope(in session) | can cross scope, i.e. a-leg to b-leg   |

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
	OnCall(fromSession, fromNode string, args []string) (resp []string)
	OnCast(fromSession, fromNode string, args []string)
}

// Streamable can negotiate message traits before link is established and transfer data after that
type Streamable interface {
	ProvideOffer() (mt []*MessageTrait)
	AnswerOffer(mt []*MessageTrait) (filteredMt []*MessageTrait)
	SetStreamTarget(session, name string, mt *MessageTrait) (error, LinkPoint)
	OnSetStream(session, name string, mt *MessageTrait) (error, LinkPoint)
}

// SessionAware enables node to:
// 1. config its static properties before any event flows
// 2. set data output destination
// 3. exit the graph when session ends
type SessionAware interface {
	event.Node

	CommandReceiver
	Streamable

	// Init do initialization after node is allocated and configured
	Init() error

	SetController(ctrl CommandInitiator)

	// ExitGraph is used when initialization failed or session terminated
	ExitGraph()
}

// ComposerAware enable class to interact with composer at each phase
type ComposerAware interface {
	// PreSetup called after composer parsed graph description, before creating node instances
	PreSetup(c *Composer) error

	// PostSetup called after composer created and setuped nodes
	PostSetup(c *Composer) error
}

// MessageProvider can push data message to event-graph
type MessageProvider interface {
	GetName() string
	PushMessage(data Message) error
	CanHandlePayloadType(pt uint8) bool
	Priority() uint32 // multiple message providers can be ordered by priority
}

// Registry service for new node type with predefined factories

type SessionNodeFactory func() SessionAware

var sessionNodeRegistry = map[string]SessionNodeFactory{
	TypeENTRY:    newEntryNode,
	TypePUBSUB:   newPubSubNode,
	TypeDISPATCH: newDispatch,
}

func getNodeByName(typeName string) SessionAware {
	if f, ok := sessionNodeRegistry[typeName]; ok {
		return f()
	}
	logger.Errorf("unknown node type:%v\n", typeName)
	return nil
}

func RegisterNodeFactory(typeName string, f SessionNodeFactory) error {
	if typeName == "" || f == nil {
		return errors.New("wrong typename or nil factory")
	}
	if _, ok := sessionNodeRegistry[typeName]; ok {
		return errors.New(fmt.Sprintf("node of type:%v already registered", typeName))
	}
	sessionNodeRegistry[typeName] = f
	return nil
}

// public commands
const (
	CtrlCall = iota + 10000
	CtrlCast
	RawByte
	Generic
)
