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
	TypeENTRY    = "entry"
	TypePUBSUB   = "pubsub"
	TypeDISPATCH = "dispatch"
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

// RawByteMessage is used to pass byte streaming data between nodes intra/inter session, for efficiency
type RawByteMessage []byte

// CtrlMessage is used to invoke or cast function call
type CtrlMessage struct {
	M []string
	C chan []string // used to receive result if not nil
}

// GenericMessage is used to encapsulate custom object in message with tagged type name
// nodes can only communicate with each other who can build/handle GenericMessage of the same tagged type
type GenericMessage struct {
	Subtype string
	Obj     interface{}
}

// Message is the base interface of all kinds of message
type Message interface {
	AsEvent() *event.Event
}

type Cloneable interface {
	Clone() Cloneable
}

// CloneableMessage is a must if message need to pass through pubsub node
type CloneableMessage interface {
	Message
	Clone() CloneableMessage
}

func NewRawByteMessage(d string) RawByteMessage {
	return RawByteMessage(d)
}

func (m RawByteMessage) String() string {
	return string(m)
}
func (m RawByteMessage) Clone() CloneableMessage {
	mc := make(RawByteMessage, len(m))
	copy(mc, m)
	return mc
}

func (m RawByteMessage) AsEvent() *event.Event {
	return event.NewEvent(RawByte, m)
}

func (cm *CtrlMessage) AsEvent() *event.Event {
	if cm.C != nil {
		return event.NewEvent(CtrlCall, cm)
	} else {
		return event.NewEvent(CtrlCast, cm)
	}
}

func (gm *GenericMessage) AsEvent() *event.Event {
	return event.NewEvent(Generic, gm)
}

func (gm *GenericMessage) String() string {
	return fmt.Sprintf("GenericMessage type:%v value:%v", gm.Subtype, gm.Obj)
}

// Clone returns non-nil object only if internal object is also cloneable
func (gm *GenericMessage) Clone() (obj CloneableMessage) {
	if gm.Obj == nil || gm.Obj == gm {
		// prevent recursive clone
		return
	}
	var cloned interface{}
	switch gm.Obj.(type) {
	case Cloneable:
		cloned = gm.Obj.(Cloneable).Clone()
	case CloneableMessage:
		cloned = gm.Obj.(CloneableMessage).Clone()
	default:
		return
	}
	obj = &GenericMessage{
		Subtype: gm.Subtype,
		Obj:     cloned,
	}
	return
}

// public commands
const (
	CtrlCall = iota + 10000
	CtrlCast
	RawByte
	Generic
)
