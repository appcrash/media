package event

import (
	"github.com/sirupsen/logrus"
	"time"
)

var logger *logrus.Entry

func InitLogger(gl *logrus.Logger) {
	logger = gl.WithFields(logrus.Fields{"module": "event"})
}

// event network structures

type Callback func()

type Event struct {
	cmd int
	obj interface{}
	cb  Callback
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

	// normal event handling
	OnEvent(evt *Event)

	// methods below (On***) are never invoke concurrently
	// all of them are called in multiple separate goroutine sequentially

	// dlink status change
	OnLinkDown(linkId int, scope string, nodeName string)

	// after sucessfully added to graph
	OnEnter(delegate *NodeDelegate)

	// the finalizing method after node exits graph
	OnExit()
}

// NodeProperty embed it if node needs to be customized
// set properties before adding node to graph
// -----------------------------------------------------------
// maxLink int:
//   override default max output link number
// dataChannelSize int:
//   override default buffered event channel size
// deliveryTimeout time.Duration:
//   override default event delivery timeout
type NodeProperty struct {
	maxLink         int
	dataChannelSize int
	deliveryTimeout time.Duration
}

func (np *NodeProperty) SetMaxLink(m int) {
	np.maxLink = m
}

func (np *NodeProperty) GetMaxLink() int {
	return np.maxLink
}

func (np *NodeProperty) SetDataChannelSize(size int) {
	np.dataChannelSize = size
}

func (np *NodeProperty) GetDataChannelSize() int {
	return np.dataChannelSize
}

func (np *NodeProperty) SetDeliveryTimeout(d time.Duration) {
	np.deliveryTimeout = d
}

func (np *NodeProperty) GetDeliveryTimeout() time.Duration {
	return np.deliveryTimeout
}
