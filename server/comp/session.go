package comp

import (
	"errors"
	"fmt"
	"github.com/appcrash/media/server/comp/nmd"
	"github.com/appcrash/media/server/event"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

type Id struct {
	Name      string
	SessionId string
}

func (id *Id) String() string {
	return id.SessionId + "_" + id.Name
}

type MessageHandler func(evt *event.Event)

func NewId(sessionId, name string) *Id {
	return &Id{SessionId: sessionId, Name: name}
}

// SessionNode is the base class of all nodes that provide capability in an RTP session
type SessionNode struct {
	Id
	delegate *event.NodeDelegate
	ctrl     CommandInitiator
	mutex    sync.Mutex

	messageTypeMatch []MessageType
	messageHandler   []MessageHandler

	Trait     *NodeTrait
	linkPoint []LinkPoint // grow only array
}

//------------------- Base Node Implementation -------------------------

func (s *SessionNode) GetNodeName() string {
	return s.Name
}

func (s *SessionNode) GetNodeScope() string {
	return s.SessionId
}

func (s *SessionNode) OnEnter(delegate *event.NodeDelegate) {
	logger.Debugf("node(%v) enters graph", s.GetNodeName())
	s.delegate = delegate
	s.SetHandler(MtNewLinkPoint, s.handleLinkPoint)
	s.SetHandler(MtConnectNode, s.handleConnectNode)
}

func (s *SessionNode) OnExit() {
	logger.Debugf("node(%v) exits graph", s.GetNodeName())
}

func (s *SessionNode) OnEvent(evt *event.Event) {
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
}

//--------------------------- Base SessionAware Implementation --------------------------------

func (s *SessionNode) handleLinkPoint(evt *event.Event) {
	msg, ok := ToMessage[*LinkPointMessage](evt)
	if !ok {
		return
	}
	for _, offer := range msg.OfferedTrait {
		for _, answer := range s.Trait.Accept {
			if offer.Match(answer) {
				msg.C <- answer
				return
			}
		}
	}
	msg.C <- nil
	return
}

// DESIGN DRAWBACK: connect to the other node may cause jitter when stream flow is under heavy load,
// because this is a sync operation
// ALTERNATIVE: put the StreamTo to other goroutine?
func (s *SessionNode) handleConnectNode(evt *event.Event) {
	var connected bool
	msg, ok := ToMessage[*ConnectNodeMessage](evt)
	if !ok {
		return
	}
	defer func() { msg.C <- connected }()
	if len(msg.Session) == 0 || len(msg.NodeName) == 0 {
		logger.Errorf("node(%v) got invalid connect node request %v", s, msg)
		return
	}
	if lp, err := s.StreamTo(msg.Session, msg.NodeName); err != nil {
		return
	} else {
		connected = true
		s.linkPoint = append(s.linkPoint, lp)
	}
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

func (s *SessionNode) AddLinkPoint(lp LinkPoint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.linkPoint = append(s.linkPoint, lp)
}

func (s *SessionNode) GetLinkPoint(index int) (lp LinkPoint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if index > len(s.linkPoint)-1 {
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
// 2. within the node's event goroutine after node has started working
func (s *SessionNode) StreamTo(session, name string) (lp LinkPoint, err error) {
	var linkId int

	defer func() {
		if err != nil && linkId >= 0 {
			s.delegate.RequestLinkDown(linkId)
		}
	}()
	offer := s.Trait.Offer
	if len(offer) == 0 {
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
	newLinkCmd := &LinkPointMessage{
		OfferedTrait: offer,
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
		s.AddLinkPoint(lp)
		logger.Infof("new stream connection (%v|%v){%x} --->[%v]---> (%v|%v)",
			s.GetNodeScope(), s.GetNodeName(), linkIdentity, agreedTrait.Type.String(), session, name)
	case <-time.After(2 * time.Second):
		err = fmt.Errorf("(%v:%v) can not set stream target to (%v:%v) due to link point not retrieved",
			s.SessionId, s.Name, session, name)
	}
	return
}

//--------------------------- Facility methods --------------------------------

func (s *SessionNode) SetHandler(msgType MessageType, handler MessageHandler) {
	for i := 0; i < len(s.messageTypeMatch); i++ {
		if s.messageTypeMatch[i] == msgType {
			s.messageHandler[i] = handler
			return
		}
	}
	// not found, prepend to head of array
	s.messageTypeMatch = append([]MessageType{msgType}, s.messageTypeMatch...)
	s.messageHandler = append([]MessageHandler{handler}, s.messageHandler...)
}

// DeliverToStream put message to stream, don't call it in stream event goroutine or deadlock!
func (s *SessionNode) DeliverToStream(msg Message) {
	s.delegate.DeliverSelf(msg.AsEvent())
}

// ToMessage convert event object back to concrete message
func ToMessage[M Message](evt *event.Event) (msg M, ok bool) {
	obj := evt.GetObj()
	if obj == nil {
		return
	}
	msg, ok = obj.(M)
	return
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
	if trait, ok := NodeTraitOfName(nodeType); !ok {
		return nil
	} else {
		node := trait.NewFunc()
		if node != nil {
			props = setNodeProperties(node, props)
		}
		return node
	}
}

// setNodeProperties use reflection to set fields by Name, it is cornerstone of config by scripting
func setNodeProperties(node event.Node, props []*nmd.NodeProp) (newProps []*nmd.NodeProp) {
	ns := reflect.ValueOf(node).Elem()
	for _, p := range props {
		k, value := p.Key, p.Value
		field := ns.FieldByName(k)
		rv := reflect.ValueOf(value)
		if !field.IsValid() {
			continue
		}
		rvt := rv.Type()
		var rvCopy reflect.Value
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
