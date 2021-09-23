package comp

import (
	"errors"
	"fmt"
	"github.com/appcrash/media/server/comp/nmd"
	"github.com/appcrash/media/server/event"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Entry

func InitLogger(gl *logrus.Logger) {
	logger = gl.WithFields(logrus.Fields{"module": "comp"})
}

const (
	TYPE_ENTRY    = "entry"
	TYPE_PUBSUB   = "pubsub"
	TYPE_DISPATCH = "dispatch"
)

type Id struct {
	Name      string
	SessionId string
}

func (id *Id) String() string {
	return id.SessionId + "_" + id.Name
}

func NewId(sessionId, name string) *Id {
	return &Id{SessionId: sessionId, Name: name}
}

// SessionAware enables node to:
// 1. config its static properties before any event flows
// 2. set data output destination
// 3. exit the graph when session ends
type SessionAware interface {
	event.Node

	// ConfigProperties handles props that can not be configured by simple reflection
	ConfigProperties(ci []*nmd.NodeProp)

	// Init do initialization after node is allocated and configured
	Init() error

	// SetPipeOut specify the data endpoint to which this node output
	SetPipeOut(session, name string) error

	SetController(ctrl Controller)

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

// Controller can be used by session commands and nodes in the event graph to invoke actions of other nodes.
// it provides a unified way to send control message to nodes without any knowledge of links (controller would
// establish links to all nodes to communicate beforehand), so simplify the programming pattern
type Controller interface {
	// Call send message and wait for the response (block)
	Call(session, name string, args []string) (resp []string)

	// Cast send message and don't wait (nonblock)
	Cast(session, name string, args []string)
}

// MessageProvider can push data message to event graph
type MessageProvider interface {
	GetName() string
	PushMessage(data DataMessage) error
}

// Registry service for new node type with predefined factories

type SessionNodeFactory func() SessionAware

var sessionNodeRegistry = map[string]SessionNodeFactory{
	TYPE_ENTRY:    newEntryNode,
	TYPE_PUBSUB:   newPubSubNode,
	TYPE_DISPATCH: newDispatch,
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

// DataMessage is used to pass data between nodes in session
type DataMessage []byte
type Cloneable interface {
	Clone() Cloneable
}

func NewDataMessage(d string) DataMessage {
	return DataMessage(d)
}

func (m DataMessage) String() string {
	return string(m)
}
func (m DataMessage) Clone() Cloneable {
	mc := make(DataMessage, len(m))
	copy(mc, m)
	return mc
}

// CtrlMessage is used to invoke or cast function call
type CtrlMessage struct {
	M []string
	C chan []string // used to receive result
}

func NewCallEvent(M []string) *event.Event {
	msg := &CtrlMessage{
		M: M,
		C: make(chan []string, 1),
	}
	return event.NewEvent(CTRL_CALL, msg)
}

func NewCastEvent(M []string) *event.Event {
	msg := &CtrlMessage{
		M: M,
	}
	return event.NewEvent(CTRL_CAST, msg)
}

func NewDataEvent(dm DataMessage) *event.Event {
	return event.NewEvent(DATA_OUTPUT, dm)
}

// public commands
const (
	CTRL_CALL = iota + 10000
	CTRL_CAST
	DATA_OUTPUT
)
