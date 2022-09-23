package comp

import (
	"github.com/appcrash/media/server/event"
)

// this node accepts input data (from pub, one or many), make multiple copy of it then send them to
// all subscriber (to sub). subscriber can be added or removed dynamically by input commands or api.
// the node will actively remove a subscriber once event is not successfully delivered to it.
//
// the input message type is used as output type when node negotiation, if multiple inputs available, they MUST
// offer the same message type or all of them except the first one would be rejected
//

var cloneableMetaType = MetaType[Cloneable]()

type Pubsub struct {
	SessionNode
	event.NodeProperty

	messageTrait *MessageTrait // all of input&output use this message trait
}

// override default negotiation handler, use the first successfully connecting node's trait as pubub's offer
// succession of connecting nodes must use the same trait as the first one or would be rejected
// NOTE: no message conversion service is provided by pubsub
func (n *Pubsub) handleLinkPoint(msg *LinkPointRequestMessage) {
	agreedTrait := n.messageTrait
	defer func() {
		msg.C <- agreedTrait
	}()
	if len(msg.OfferedTrait) == 0 {
		logger.Errorf("pubsub(%v) get empty offer", n)
		return
	}
	if n.messageTrait != nil {
		// not the first visitor
		for _, trait := range msg.OfferedTrait {
			if n.messageTrait.Match(trait) {
				logger.Infof("pubsub(%v) accept more then one nodes of the same message trait", n)
				return
			}
		}
		agreedTrait = nil
		logger.Errorf("pubsub(%v) reject new comer as it already has incompatible input trait", n)
		return
	}

	// find the first eligible trait
	for _, trait := range msg.OfferedTrait {
		if trait.PtrType.Implements(cloneableMetaType) {
			n.messageTrait = trait.Clone()
			break
		}
	}
	if n.messageTrait == nil {
		logger.Errorf("pubsub(%v) reject the offer as none of them(%v) is cloneable", n, msg.OfferedTrait)
		return
	}
	agreedTrait = n.messageTrait
	n.SetMessageHandler(n.messageTrait.TypeId, ChainSetHandler(n.handleInputStream))
	logger.Debugf("pubsub(%v) accept message type: %v", n, n.messageTrait)
	return
}

func (n *Pubsub) handleInputStream(evt *event.Event) {
	msg, ok := EventToMessage[Message](evt)
	if !ok {
		return
	}
	// no cast check, as in negotiation phase we have inspected message trait, if sender doesn't obey the rule,
	// panic is waiting for you ...
	cloneableMessage := MessageTo[Cloneable](msg)
	for _, lp := range n.linkPoint {
		cloned := cloneableMessage.Clone()
		lp.SendMessage(cloned.(Message))
	}
}

func (n *Pubsub) Offer() []MessageType {
	if n.messageTrait != nil {
		return []MessageType{n.messageTrait.TypeId}
	} else {
		return nil
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
