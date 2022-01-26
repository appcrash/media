package comp

import (
	"errors"
	"fmt"
	"github.com/appcrash/media/server/event"
	"sync"
	"time"
)

// this node accepts input data (from pub, one or many), make multiple copy of it then send them to
// all subscriber (to sub). subscriber can be added or removed dynamically by input commands or api.
// the node will actively remove a subscriber once event is not successfully delivered to it.

// for input, one of :
// 1. other node send data message to pubsub node (inter-node communication)
// 2. call this node's Publish method (feed event to event graph)
//
// for output(subscriber), one of:
// 1. other node that receives event from pubsub node (inter-node communication)
// 2. provider a channel of type (chan<- *event.Event) to which pubsub delivers (consume event from event graph),
// the channel must be buffered channel, i.e. cap(c) != 0, and in the SAME session of this node, so communication
// across sessions must take inter-node measures
//
// PubSub can have only one input and one output, then it becomes like 'tee' in linux command line
// which read from stdin and write to stdout. so PubSub can be a bridge between event graph and outside world
//
// |outside| ------> pubsub -----> event graph
//            feed
//
// event graph -----> pubsub -----> |outside|
//                          consume

const PubsubDefaultDeliveryTimeout = 20 * time.Millisecond

const (
	psSubscribeTypeNode = iota
	psSubscribeTypeChannel
)

type psSubscriberInfo struct {
	subType int
	linkId  int                 // if subscriber is a node
	channel chan<- *event.Event // if subscriber is a chan
	name    string
	enabled bool
}

type PubSubNode struct {
	SessionNode
	event.NodeProperty

	mutex       sync.Mutex
	subscribers []*psSubscriberInfo
}

func (p *PubSubNode) OnEvent(evt *event.Event) {
	obj := evt.GetObj()
	if obj == nil {
		return
	}
	switch evt.GetCmd() {
	case RawByte, Generic:
		if c, ok := obj.(CloneableMessage); ok {
			p.Publish(c)
		}
	case CtrlCall:
		if msg, ok := obj.(*CtrlMessage); ok {
			p.handleCall(msg)
		}
	case CtrlCast:
		// none
	}
}

func (p *PubSubNode) OnLinkDown(_ int, scope string, nodeName string) {
	if index, si := p.findNodeSubscriber(scope, nodeName); si != nil {
		// a node subscriber is down, just remove it from subscribers
		p.deleteSubscriber(index)
	}
}

// OnExit close all channel subscribers
func (p *PubSubNode) OnExit() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, s := range p.subscribers {
		if s.subType == psSubscribeTypeChannel {
			if s.channel != nil {
				close(s.channel)
				s.channel = nil
			}
		}
	}
}

// SetPipeOut overrides default session node's behaviour, it allows multiple pipes simultaneously
// (that's what "pubsub" stands for) instead of only one data output pipe
func (p *PubSubNode) SetPipeOut(session, name string) error {
	if p.delegate == nil {
		return errors.New("delegate not ready when set pipe")
	}
	return p.SubscribeNode(session, name)
}

//------------------------------- api & implementation --------------------------------------

func (p *PubSubNode) Publish(obj CloneableMessage) {
	var subscribers []*psSubscriberInfo
	if obj == nil {
		return
	}

	p.mutex.Lock()
	size := len(p.subscribers)
	if size == 0 {
		p.mutex.Unlock()
		return
	}
	subscribers = make([]*psSubscriberInfo, size)
	// as subscribers may change during event delivering,
	// copy the array, then release lock
	copy(subscribers, p.subscribers)
	p.mutex.Unlock()

	// publish message to all subscribers
	// for node: delivery timeout by field 'deliveryTimeout'
	// for channel: nonblock sending without timeout, so must use buffered channel to avoid losing message
	for _, s := range subscribers {
		if !s.enabled {
			continue
		}
		switch s.subType {
		case psSubscribeTypeNode:
			if s.linkId < 0 {
				continue
			}
			msg := obj.Clone()
			if msg == nil {
				logger.Debugf("pubsub got message that clone to nil")
				continue
			}
			evt := msg.AsEvent()
			p.delegate.Deliver(s.linkId, evt)
		case psSubscribeTypeChannel:
			if s.channel == nil {
				continue
			}
			msg := obj.Clone()
			if msg == nil {
				logger.Debugf("pubsub got message that clone to nil")
				continue
			}
			evt := msg.AsEvent()
			select {
			case s.channel <- evt:
			default:
			}
		}
	}
}

func newPubSubNode() SessionAware {
	node := new(PubSubNode)
	node.Name = TypePUBSUB
	node.SetDeliveryTimeout(PubsubDefaultDeliveryTimeout)
	return node
}

func psNewNodeSubscriber(scope, nodeName string, linkId int) *psSubscriberInfo {
	name := psMakeNodeName(scope, nodeName)
	si := new(psSubscriberInfo)
	si.subType = psSubscribeTypeNode
	si.name = name
	si.linkId = linkId
	si.enabled = true
	return si
}

func psNewChannelSubscriber(chName string, c chan<- *event.Event) *psSubscriberInfo {
	name := psMakeChannelName(chName)
	si := new(psSubscriberInfo)
	si.subType = psSubscribeTypeChannel
	si.name = name
	si.channel = c
	si.linkId = -1
	si.enabled = true
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
	p.mutex.Lock()
	defer p.mutex.Unlock()
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

func (p *PubSubNode) deleteSubscriber(index int) {
	p.mutex.Lock()
	siLen := len(p.subscribers)
	p.subscribers[index] = p.subscribers[siLen-1]
	p.subscribers = p.subscribers[:siLen-1]
	p.mutex.Unlock()
}

func (p *PubSubNode) handleCall(msg *CtrlMessage) {
	m := msg.M
	if ml := len(m); ml > 0 {
		action := m[0]
		switch action {
		case "conn":
			// conn {sessionName} {nodeName}
			if ml == 3 {
				toSession, toName := m[1], m[2]
				if err := p.SetPipeOut(toSession, toName); err == nil {
					msg.C <- WithOk()
					return
				}
			}
			goto error
		case "enable":
			// enable node {sessionName} {nodeName}
			// enable channel {channelName}
			fallthrough
		case "disable":
			// disable node {sessionName} {nodeName}
			// disable channel {channelName}
			var si *psSubscriberInfo
			switch ml {
			case 3:
				if m[1] != "channel" {
					goto error
				}
				channelName := m[2]
				_, si = p.findChannelSubscriber(channelName)
				goto error
			case 4:
				if m[1] != "node" {
					goto error
				}
				sessionName, nodeName := m[2], m[3]
				_, si = p.findNodeSubscriber(sessionName, nodeName)
			}
			if si == nil {
				goto error
			}
			p.mutex.Lock()
			if action == "enable" {
				si.enabled = true
			} else {
				si.enabled = false
			}
			p.mutex.Unlock()
			msg.C <- WithOk()
			return
		}
	}

error:
	msg.C <- WithError()
}

// SubscribeNode add a node as a subscriber of this pubsub node
func (p *PubSubNode) SubscribeNode(scope, name string) error {
	if _, s := p.findNodeSubscriber(scope, name); s != nil {
		return errors.New(fmt.Sprintf("node %v:%v is already a subscriber", scope, name))
	}
	if linkId := p.delegate.RequestLinkUp(scope, name); linkId >= 0 {
		si := psNewNodeSubscriber(scope, name, linkId)
		p.mutex.Lock()
		p.subscribers = append(p.subscribers, si)
		p.mutex.Unlock()
		return nil
	} else {
		return errors.New(fmt.Sprintf("node(%v:%v) subscribes failed", scope, name))
	}
}

// UnsubscribeNode remove a node subscriber from this pubsub node
func (p *PubSubNode) UnsubscribeNode(scope, name string) error {
	if index, si := p.findNodeSubscriber(scope, name); si == nil {
		return nil
	} else {
		// delete the subscriber now
		p.deleteSubscriber(index)
		return p.delegate.RequestLinkDown(si.linkId)
	}
}

// SubscribeChannel add a channel as a subscriber of this pubsub node
func (p *PubSubNode) SubscribeChannel(chName string, c chan<- *event.Event) error {
	if _, s := p.findChannelSubscriber(chName); s != nil {
		return errors.New(fmt.Sprintf("channel %v is already subscribed", chName))
	}
	if cap(c) == 0 {
		// must be buffered channel
		return errors.New("must use buffered channel to subscribe")
	}
	si := psNewChannelSubscriber(chName, c)
	p.mutex.Lock()
	p.subscribers = append(p.subscribers, si)
	p.mutex.Unlock()
	return nil
}

// UnsubscribeChannel remove a channel subscriber with given name from this pubsub node
func (p *PubSubNode) UnsubscribeChannel(chName string) error {
	if index, si := p.findChannelSubscriber(chName); si == nil {
		return errors.New(fmt.Sprintf("channel %v is not subscribed, so unsubscribe fails", chName))
	} else {
		p.deleteSubscriber(index)
		return nil
	}
}
