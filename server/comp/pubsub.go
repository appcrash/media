package comp

import (
	"fmt"
	"github.com/appcrash/media/server/event"
	"sync"
	"time"
)

// this node accepts input data (from pub, one or many), make multiple copy of it then send them to
// all subscriber (to sub). subscriber can be added or removed dynamically by input commands or api.
// the node will actively remove a subscriber once event is not successfully delivered to it.

// for input, one of :
// 1. other node send event to pubsub node (inter-node communication)
// 2. call this node's Publish method (feed event to event graph)
//
// for output(subscriber), one of:
// 1. other node that receives event from pubsub node (inter-node communication)
// 2. provider a channel of type (chan<- *event.Event) to which pubsub deliveries (consume event from event graph),
// the channel must be buffered channel, i.e. cap(c) != 0
//
// PubSub can have only one input and one output, then it becomes like 'tee' in linux command line
// which read from stdin and write to stdout. so PubSub can be a bridge between event graph and outside world
//
// |outside| ------> pubsub -----> event graph
//            feed
//
// event graph -----> pubsub -----> |outside|
//                          consume
//
// command details:
// #add subscriber:
//   for node:   receiving_session receiving_name
//   for channel:   Name  channel_reference (of type chan<- *event.Event)
// #remove subscriber:
//   for node:   receiving_scope receiving_name
//   for channel:   Name
// #publish data
//   <any object>
//
// response details:
// #subscribed_data
//   <any object from publish data command>
//

const PUBSUB_DEFAULT_DELIVERY_TIMEOUT = 20 * time.Millisecond

const (
	psSubscribeTypeNode = iota
	psSubscribeTypeChannel
)

type cmdPsAddNode struct {
	session, name string
}

type cmdPsAddChannel struct {
	name string
	c    chan<- *event.Event
}

type cmdPsRemoveNode struct {
	session, name string
}

type cmdPsRemoveChannel struct {
	name string
}

type psSubscriberInfo struct {
	subType int
	linkId  int                 // if subscriber is a node
	channel chan<- *event.Event // if subscriber is a chan
	name    string
}

type PubSubNode struct {
	SessionNode

	deliveryTimeout time.Duration
	mutex           sync.Mutex
	subscribers     []*psSubscriberInfo
}

func (p *PubSubNode) OnEvent(evt *event.Event) {
	obj := evt.GetObj()
	if obj == nil {
		return
	}
	switch evt.GetCmd() {
	case DATA_OUTPUT:
		if c, ok := obj.(Cloneable); ok {
			p.Publish(c)
		}
	case CMD_GENERIC_SET_ROUTE:
		if c, ok := obj.(*GenericRouteCommand); ok {
			p.doSubscribeNode(c.SessionId, c.Name)
		}
	case CMD_PUBSUB_ADD_NODE_SUBSCRIBER:
		if i, ok := obj.(*cmdPsAddNode); ok {
			p.doSubscribeNode(i.session, i.name)
		}
	case CMD_PUBSUB_ADD_CHANNEL_SUBSCRIBER:
		if i, ok := obj.(*cmdPsAddChannel); ok {
			p.doSubscribeChannel(i.name, i.c)
		}
	case CMD_PUBSUB_REMOVE_NODE_SUBSCRIBER:
		if i, ok := obj.(*cmdPsRemoveNode); ok {
			p.doUnsubscribeNode(i.session, i.name)
		}
	case CMD_PUBSUB_REMOVE_CHANNEL_SUBSCRIBER:
		if i, ok := obj.(*cmdPsRemoveChannel); ok {
			p.doUnsubscribeChannel(i.name)
		}
	}
}

func (p *PubSubNode) OnLinkUp(linkId int, scope string, nodeName string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if index, si := p.findNodeSubscriber(scope, nodeName); si != nil {
		// if link up request failed, remove the subscriber, or set the link id
		if linkId < 0 {
			p.deleteSubscriber(index)
		} else {
			si.linkId = linkId
		}
	}
}

func (p *PubSubNode) OnLinkDown(_linkId int, scope string, nodeName string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if index, si := p.findNodeSubscriber(scope, nodeName); si != nil {
		// a node subscriber is down, just remove it from subscribers
		p.deleteSubscriber(index)
	}
}

//------------------------------- api & implementation --------------------------------------

func (p *PubSubNode) Publish(obj Cloneable) {
	if obj == nil {
		return
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// publish message to all subscribers
	// for node: delivery timeout by field 'deliveryTimeout'
	// for channel: nonblock sending without timeout
	for _, s := range p.subscribers {
		switch s.subType {
		case psSubscribeTypeNode:
			if s.linkId < 0 {
				continue
			}
			evt := event.NewEvent(DATA_OUTPUT, obj.Clone())
			p.delegate.Delivery(s.linkId, evt)
		case psSubscribeTypeChannel:
			if s.channel == nil {
				continue
			}
			evt := event.NewEvent(DATA_OUTPUT, obj.Clone())
			select {
			case s.channel <- evt:
			default:
			}
		}
	}
}

// SubscribeNode add a node as a subscriber of this pubsub node
func (p *PubSubNode) SubscribeNode(session, name string) bool {
	evt := event.NewEvent(CMD_PUBSUB_ADD_NODE_SUBSCRIBER, &cmdPsAddNode{session, name})
	return p.delegate.DeliverSelf(evt)
}

// UnsubscribeNode remove a node subscriber of this pubsub node
func (p *PubSubNode) UnsubscribeNode(session, name string) bool {
	evt := event.NewEvent(CMD_PUBSUB_REMOVE_NODE_SUBSCRIBER, &cmdPsRemoveNode{session, name})
	return p.delegate.DeliverSelf(evt)
}

// SubscribeChannel add a channel as a subscriber of this pubsub node
func (p *PubSubNode) SubscribeChannel(name string, c chan<- *event.Event) bool {
	evt := event.NewEvent(CMD_PUBSUB_ADD_CHANNEL_SUBSCRIBER, &cmdPsAddChannel{name, c})
	return p.delegate.DeliverSelf(evt)
}

func (p *PubSubNode) UnsubscribeChannel(name string) bool {
	evt := event.NewEvent(CMD_PUBSUB_REMOVE_CHANNEL_SUBSCRIBER, &cmdPsRemoveChannel{name})
	return p.delegate.DeliverSelf(evt)
}

func newPubSubNode() *PubSubNode {
	node := new(PubSubNode)
	node.Name = TYPE_PUBSUB
	node.deliveryTimeout = PUBSUB_DEFAULT_DELIVERY_TIMEOUT
	return node
}

func psNewNodeSubscriber(scope, nodeName string) *psSubscriberInfo {
	name := psMakeNodeName(scope, nodeName)
	si := new(psSubscriberInfo)
	si.subType = psSubscribeTypeNode
	si.name = name
	si.linkId = -1 // set when link up
	return si
}

func psNewChannelSubscriber(chName string, c chan<- *event.Event) *psSubscriberInfo {
	name := psMakeChannelName(chName)
	si := new(psSubscriberInfo)
	si.subType = psSubscribeTypeChannel
	si.name = name
	si.channel = c
	si.linkId = -1
	return si
}

func psMakeNodeName(scope, name string) string {
	return fmt.Sprintf("node_%v_%v", scope, name)
}

func psMakeChannelName(chName string) string {
	return "chan_" + chName
}

func (p *PubSubNode) findSubInfo(name string) (index int, si *psSubscriberInfo) {
	if name == "" {
		return -1, nil
	}
	for i, n := range p.subscribers {
		if n.name == name {
			return i, n
		}
	}
	return -1, nil
}

func (p *PubSubNode) findNodeSubscriber(scope, name string) (index int, si *psSubscriberInfo) {
	nodeName := psMakeNodeName(scope, name)
	return p.findSubInfo(nodeName)
}

func (p *PubSubNode) findChannelSubscriber(chName string) (index int, si *psSubscriberInfo) {
	name := psMakeChannelName(chName)
	return p.findSubInfo(name)
}

func (p *PubSubNode) doSubscribeNode(scope, name string) {
	if _, s := p.findNodeSubscriber(scope, name); s != nil {
		return
	}

	si := psNewNodeSubscriber(scope, name)
	p.mutex.Lock()
	p.subscribers = append(p.subscribers, si)
	p.mutex.Unlock()
	p.delegate.RequestLinkUp(scope, name)
}

func (p *PubSubNode) doUnsubscribeNode(scope, name string) {
	if index, si := p.findNodeSubscriber(scope, name); si == nil {
		return
	} else {
		if si.linkId >= 0 {
			// delete the subscriber now
			p.deleteSubscriber(index)
			p.delegate.RequestLinkDown(si.linkId)
		}
	}
}

func (p *PubSubNode) doSubscribeChannel(chName string, c chan<- *event.Event) {
	if _, s := p.findChannelSubscriber(chName); s != nil {
		return
	}
	if cap(c) == 0 {
		// must be buffered channel
		return
	}
	si := psNewChannelSubscriber(chName, c)
	p.mutex.Lock()
	p.subscribers = append(p.subscribers, si)
	p.mutex.Unlock()
}

func (p *PubSubNode) doUnsubscribeChannel(chName string) {
	if index, si := p.findChannelSubscriber(chName); si == nil {
		return
	} else {
		p.deleteSubscriber(index)
	}
}

func (p *PubSubNode) deleteSubscriber(index int) {
	p.mutex.Lock()
	siLen := len(p.subscribers)
	p.subscribers[index] = p.subscribers[siLen-1]
	p.subscribers = p.subscribers[:siLen-1]
	p.mutex.Unlock()
}
