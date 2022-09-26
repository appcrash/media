package comp

import (
	"errors"
	"fmt"
	"github.com/appcrash/media/server/comp/nmd"
	"github.com/appcrash/media/server/event"
	"github.com/appcrash/media/server/utils"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

type Id struct {
	Name      string
	SessionId string
}

func NewId(sessionId, name string) *Id {
	return &Id{SessionId: sessionId, Name: name}
}

type MessageHandler func(evt *event.Event)
type MessageHandlerChainer func(previous MessageHandler) MessageHandler // AOP ...

func ChainSetHandler(m MessageHandler) MessageHandlerChainer {
	return func(_ MessageHandler) MessageHandler {
		// replace previous handler any way
		return m
	}
}

func ChainDefaultHandler(m MessageHandler) MessageHandlerChainer {
	return func(previous MessageHandler) MessageHandler {
		if previous != nil {
			// already set a handler, don't override it
			return previous
		} else {
			return m
		}
	}
}

// SessionNode is the base class of all nodes that provide capability in an RTP session
type SessionNode struct {
	Id
	delegate *event.NodeDelegate
	mutex    sync.Mutex

	messageTypeMatch []MessageType
	messageHandler   []MessageHandler
	linkPoint        []LinkPoint // grow only array

	Trait     *NodeTrait       // initialized by gentrait
	Initiator CommandInitiator // initialized by composer
}

//------------------- Base Node Implementation -------------------------

func (s *SessionNode) String() string {
	return fmt.Sprintf("[%v{%v}@%v]", s.GetNodeName(), s.GetNodeTypeName(), s.GetNodeScope())
}

func (s *SessionNode) GetNodeName() string {
	return s.Name
}

func (s *SessionNode) GetNodeScope() string {
	return s.SessionId
}

func (s *SessionNode) GetNodeTypeName() string {
	return s.Trait.NodeType
}

func (s *SessionNode) OnEnter(delegate *event.NodeDelegate) {
	logger.Debugf("node %v enters graph", s.String())
	s.delegate = delegate

	// provide default negotiation behaviour handlers
	s.SetMessageHandler(MtLinkPointRequest, ChainDefaultHandler(s._handleLinkPointRequest))
}

func (s *SessionNode) OnExit() {
	logger.Debugf("node %v exits graph", s)
}

func (s *SessionNode) OnEvent(evt *event.Event) {
	// this function is in the hotspot execution path
	// as each kind of node can accept limited type of message, linear search doesn't add much overhead
	// if this is not true, use fixed-size of handler array for constant time
	msgType := MessageType(evt.GetCmd())
	for i, typ := range s.messageTypeMatch {
		if msgType == typ {
			s.messageHandler[i](evt)
			return
		}
	}
}

func (s *SessionNode) OnLinkDown(linkId int, scope string, nodeName string) {
	logger.Debugf("node got link down (%v:%v) => (%v:%v) ", s.GetNodeScope(), s.GetNodeName(), scope, nodeName)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i, l := range s.linkPoint {
		if l.LinkId() == linkId {
			logger.Debugf("node %v delete link id %v", s, linkId)
			if newLp, err := utils.RemoveElementFromArray(s.linkPoint, i); err != nil {
				logger.Errorf("node %v remove link point with index %v failed", s, linkId)
			} else {
				s.linkPoint = newLp
			}
			return
		}
	}
}

//--------------------------- Base SessionAware Implementation --------------------------------

// this is where negotiation happens, for each offered traits:
// 1. check it can match any accepted trait of this node,
// 2. if none of them matched, test if message can be converted from offered to accepted type
// 3. if no conversion is possible, go to first step with next candidate offer type
func (s *SessionNode) _handleLinkPointRequest(evt *event.Event) {
	linkPointMessage, ok := EventToMessage[*LinkPointRequestMessage](evt)
	if !ok {
		return
	}
	for _, offer := range linkPointMessage.OfferedTrait {
		// try direct match
		for _, answer := range s.Trait.Accept {
			if offer.Match(answer) {
				linkPointMessage.C <- answer
				return
			}
		}

		// not match any one, see if conversion is possible
		for _, answer := range s.Trait.Accept {
			if CanConvertMessage(offer.TypeId, answer.TypeId) {
				// found a conversion path, setup handler for this message type
				// sanity check: ensure a handler for answered message type already exist
				if handler := s.GetMessageHandler(answer.TypeId); handler == nil {
					logger.Errorf("%v has a nil handler for message type %v, conversion is impossible from %v",
						s, answer, offer)
					continue
				}
				// really setup handler for offered message type
				logger.Debugf("node %v create message conversion function from %v -> %v", s, offer.Name(), answer.Name())
				answer = answer.Clone()
				s.SetMessageHandler(offer.TypeId, func(_ MessageHandler) MessageHandler {
					// create message conversion function
					return func(evt *event.Event) {
						msg, ok := EventToMessage[Message](evt)
						if !ok {
							return
						}
						convertedMsg, err := answer.ConvertFrom(msg)
						if err != nil {
							return
						}
						// always retrieve the latest handler, as previous-checked handler may differ at runtime
						handler := s.GetMessageHandler(answer.TypeId)
						handler(convertedMsg.AsEvent())
					}
				})
				linkPointMessage.C <- offer
				return
			}
		}
	}
	linkPointMessage.C <- nil
	return
}

func (s *SessionNode) ExitGraph() {
	if s.delegate != nil {
		_ = s.delegate.RequestNodeExit()
	}
}

func (s *SessionNode) Init() error {
	return nil
}

func (s *SessionNode) OnCall(fromNode string, args []string) (resp []string) {
	return WithOk()
}

func (s *SessionNode) OnCast(fromNode string, args []string) {

}

func (s *SessionNode) Accept() []MessageType {
	return nil
}

func (s *SessionNode) Offer() []MessageType {
	return nil
}

func (s *SessionNode) addLinkPoint(lp LinkPoint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.linkPoint = append(s.linkPoint, lp)
}

func (s *SessionNode) GetLinkPoint(index int) (lp LinkPoint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if index > len(s.linkPoint)-1 || index < 0 {
		return
	}
	return s.linkPoint[index]
}

func (s *SessionNode) GetLinkPointOfType(messageType MessageType) (lp LinkPoint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, l := range s.linkPoint {
		if l.MessageTrait().TypeId == messageType {
			return l
		}
	}
	return
}

// StreamTo sync connect to the other node if negotiation is successful, as it will change the node's state, only call
// it at:
// 1. node initialized but not in work state(i.e. no stream is flowing), such as in composer phase
// 2. outside the node's event goroutine after node has started working
func (s *SessionNode) StreamTo(session, name string, preferredOffer []MessageType) (lp LinkPoint, err error) {
	var linkId int
	var offeredTraits []*MessageTrait

	defer func() {
		if err != nil && linkId >= 0 {
			s.delegate.RequestLinkDown(linkId)
		}
	}()

	for _, o := range preferredOffer {
		if mt, ok := MessageTraitOfType(o); !ok {
			err = fmt.Errorf("message type %v not exist", o)
			return
		} else {
			offeredTraits = append(offeredTraits, mt)
		}
	}
	if len(offeredTraits) == 0 {
		err = fmt.Errorf("(%v:%v) can not set stream target to (%v:%v) due to sender has no offer",
			s.SessionId, s.Name, session, name)
		return
	}

	if linkId = s.delegate.RequestLinkUp(session, name); linkId < 0 {
		err = fmt.Errorf("(%v:%v) can not set stream target to (%v:%v) due to request link up failed",
			s.SessionId, s.Name, session, name)
		return
	}
	sendFunc := func(msg Message) error {
		if !s.delegate.Deliver(linkId, msg.AsEvent()) {
			return errors.New("failed to deliver event")
		}
		return nil
	}
	linkIdentity := MakeLinkIdentity(session, name, linkId)
	newLinkCmd := &LinkPointRequestMessage{
		OfferedTrait: offeredTraits,
		LinkIdentity: linkIdentity,
	}
	newLinkCmd.C = make(chan *MessageTrait, 1)
	evt := newLinkCmd.AsEvent()
	if ok := s.delegate.Deliver(linkId, evt); !ok {
		err = fmt.Errorf("(%v:%v) can not set stream target to (%v:%v) due to deliver link point command failed",
			s.SessionId, s.Name, session, name)
		return
	}
	select {
	case agreedTrait := <-newLinkCmd.C:
		if agreedTrait == nil {
			err = fmt.Errorf("(%v:%v) can not set stream target to (%v:%v) due to offered traits not agreed",
				s.SessionId, s.Name, session, name)
			return
		}
		lp = NewLinkPad(s, linkId, linkIdentity, agreedTrait, sendFunc)
		s.addLinkPoint(lp)
		logger.Infof("new stream connection %v{link:%x} --->[%v]---> (%v@%v)",
			s, linkIdentity, agreedTrait.Name(), name, session)
	case <-time.After(2 * time.Second):
		err = fmt.Errorf("(%v:%v) can not set stream target to (%v:%v) due to link point not retrieved",
			s.SessionId, s.Name, session, name)
	}
	return
}

//--------------------------- Facility methods --------------------------------

// SetMessageHandler calls chainer to get the new handler and replace the previous one if exists
func (s *SessionNode) SetMessageHandler(msgType MessageType, chain MessageHandlerChainer) {
	var i int
	var previousHandler MessageHandler
	for i = 0; i < len(s.messageTypeMatch); i++ {
		if s.messageTypeMatch[i] == msgType {
			previousHandler = s.messageHandler[i]
			break
		}
	}
	newHandler := chain(previousHandler)
	if newHandler == nil {
		panic("session node chain handler error")
	}
	if previousHandler == nil {
		// not found, prepend to head of array
		s.messageTypeMatch = append([]MessageType{msgType}, s.messageTypeMatch...)
		s.messageHandler = append([]MessageHandler{newHandler}, s.messageHandler...)
	} else {
		s.messageHandler[i] = newHandler
	}
}

func (s *SessionNode) GetMessageHandler(msgType MessageType) MessageHandler {
	for i := 0; i < len(s.messageTypeMatch); i++ {
		if s.messageTypeMatch[i] == msgType {
			return s.messageHandler[i]
		}
	}
	return nil
}

// DeliverToStream put message to stream, don't call it in stream event goroutine or deadlock!
func (s *SessionNode) DeliverToStream(msg Message) {
	s.delegate.DeliverSelf(msg.AsEvent())
}

// MakeSessionNode factory method of all session aware nodes
func MakeSessionNode(nodeType string, sessionId string, props []*nmd.NodeProp) SessionAware {
	if nodeType == "" || sessionId == "" {
		return nil
	}
	props = append(props,
		&nmd.NodeProp{
			Key:   "SessionId",
			Type:  "str",
			Value: sessionId,
		})
	if trait, ok := NodeTraitOfType(nodeType); !ok {
		return nil
	} else {
		node := trait.FactoryFunc()
		if node != nil {
			for _, p := range setNodeProperties(node, props) {
				logger.Warnf("node %v doesn't handle property: %v", node, p.Key)
			}
		}
		return node
	}
}

// setNodeProperties use reflection to set fields by NodeType, it is cornerstone of config by scripting
func setNodeProperties(node event.Node, props []*nmd.NodeProp) (newProps []*nmd.NodeProp) {
	ns := reflect.ValueOf(node).Elem()
	for _, p := range props {
		var rvCopy reflect.Value
		k, value := p.Key, p.Value
		field := ns.FieldByName(k)
		rv := reflect.ValueOf(value)
		rvt := rv.Type()
		if !field.IsValid() {
			goto notHandled
		}

		if rvt.AssignableTo(field.Type()) {
			rvCopy = rv
		} else if rvt.ConvertibleTo(field.Type()) {
			rvCopy = rv.Convert(field.Type())
		} else {
			goto notHandled
		}

		if field.CanSet() {
			field.Set(rvCopy)
		} else {
			// forcefully setting unexported variable
			nf := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
			nf.Set(rvCopy)
		}
		continue
	notHandled:
		newProps = append(newProps, p)
	}
	return
}
