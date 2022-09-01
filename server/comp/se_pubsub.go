package comp

import (
	"github.com/appcrash/media/server/event"
	"time"
)

// this node accepts input data (from pub, one or many), make multiple copy of it then send them to
// all subscriber (to sub). subscriber can be added or removed dynamically by input commands or api.
// the node will actively remove a subscriber once event is not successfully delivered to it.
//
// the input message type is used as output type when node negotiation, if multiple inputs available, they MUST
// offer the same message type or all of them except the first one would be rejected
//
// the pubsub node is the only one that explicitly implement 'conn' command provided by framework, which means
// if cross scope interaction is demanded, a pubsub node must kick in

const PubsubDefaultDeliveryTimeout = 20 * time.Millisecond

type Pubsub struct {
	SessionNode
	event.NodeProperty

	messageTrait *MessageTrait // all of input&output use this message trait
}

func (p *Pubsub) OnEnter(delegate *event.NodeDelegate) {
	p.SessionNode.OnEnter(delegate)
	p.SetHandler(MtNewLinkPoint, p.handleLinkPoint)
}

func (p *Pubsub) handleLinkPoint(evt *event.Event) {
	msg, ok := ToMessage[*LinkPointMessage](evt)
	if !ok {
		return
	}
	defer func() {
		msg.C <- p.messageTrait
	}()
	if len(msg.OfferedTrait) == 0 {
		logger.Errorf("pubsub(%v) get empty offer", p)
		return
	}

	// find the first eligible trait
	for _, trait := range msg.OfferedTrait {
		if trait.IsCloneable() {
			p.messageTrait = trait.Clone()
		}
	}
	if p.messageTrait == nil {
		logger.Errorf("pubsub(%v) reject the offer as none of them(%v) is cloneable", p, msg.OfferedTrait)
		return
	}
	p.SetHandler(p.messageTrait.TypeId, p.handleInputStream)
	logger.Debugf("pubsub(%v) accept message type: %v", p, p.messageTrait)
	return
}

func (p *Pubsub) handleInputStream(evt *event.Event) {
	msg, ok := ToMessage[Message](evt)
	if !ok {
		return
	}
	// no cast check, as in negotiation phase we have inspected message trait, if sender doesn't obey the rule,
	// panic is waiting for you ...
	cloneableMessage := msg.(Cloneable)
	for _, lp := range p.linkPoint {
		cloned := cloneableMessage.Clone()
		lp.SendMessage(cloned.(Message))
	}
}

//------------------------------- api & implementation --------------------------------------

//func (p *Pubsub) handleCall(msg *CtrlMessage) {
//	m := msg.M
//	if ml := len(m); ml > 0 {
//		action := m[0]
//		switch action {
//		case "conn":
//			// conn {sessionName} {nodeName}
//			if ml == 3 {
//				toSession, toName := m[1], m[2]
//				if err := p.SetPipeOut(toSession, toName); err == nil {
//					msg.C <- WithOk()
//					return
//				}
//			}
//			goto error
//		case "enable":
//			// enable node {sessionName} {nodeName}
//			// enable channel {channelName}
//			fallthrough
//		case "disable":
//			// disable node {sessionName} {nodeName}
//			// disable channel {channelName}
//			var si *psSubscriberInfo
//			switch ml {
//			case 3:
//				if m[1] != "channel" {
//					goto error
//				}
//				channelName := m[2]
//				_, si = p.findChannelSubscriber(channelName)
//				goto error
//			case 4:
//				if m[1] != "node" {
//					goto error
//				}
//				sessionName, nodeName := m[2], m[3]
//				_, si = p.findNodeSubscriber(sessionName, nodeName)
//			}
//			if si == nil {
//				goto error
//			}
//			p.mutex.Lock()
//			if action == "enable" {
//				si.enabled = true
//			} else {
//				si.enabled = false
//			}
//			p.mutex.Unlock()
//			msg.C <- WithOk()
//			return
//		}
//	}
//
//error:
//	msg.C <- WithError()
//}
