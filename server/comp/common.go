package comp

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Entry

func InitLogger(gl *logrus.Logger) {
	logger = gl.WithFields(logrus.Fields{"module": "comp"})
}

const (
	TYPE_ENTRY  = "entry"
	TYPE_PUBSUB = "pubsub"
	TYPE_SEND   = "send"
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

// Controller can be used by session commands and nodes in the event graph to invoke actions of other nodes.
// it provides a unified way to send control message to nodes without any knowledge of links (controller would
// establish links to all nodes to communicate beforehand), so simplify the programming pattern
type Controller interface {
	// Call send message and wait for the response (block)
	Call(session, name, args string) (resp string)

	// Cast send message and don't wait (nonblock)
	Cast(session, name, args string)
}

type Cloneable interface {
	Clone() Cloneable
}

type DataMessage []byte

func (m DataMessage) Clone() Cloneable {
	mc := make(DataMessage, len(m))
	copy(mc, m)
	return mc
}

func getNodeByName(typeName string) SessionAware {
	switch typeName {
	case TYPE_ENTRY:
		return newEntryNode()
	case TYPE_PUBSUB:
		return newPubSubNode()
	case TYPE_SEND:
		return newDispatch()
	default:
		logger.Errorf("unknown node type:%v\n", typeName)
		return nil
	}
}

type GenericRouteCommand struct {
	Id
}

// public commands & responses
const (
	CMD_GENERIC_SET_ROUTE = iota // set the destination to which a node output
	DATA_OUTPUT

	CMD_PUBSUB_ADD_NODE_SUBSCRIBER
	CMD_PUBSUB_ADD_CHANNEL_SUBSCRIBER
	CMD_PUBSUB_REMOVE_NODE_SUBSCRIBER
	CMD_PUBSUB_REMOVE_CHANNEL_SUBSCRIBER
)
