package comp

import (
	"fmt"
	"github.com/appcrash/media/server/comp/nmd"
	"github.com/appcrash/media/server/event"
	"reflect"
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

type EventHandler func(evt *event.Event)

func NewId(sessionId, name string) *Id {
	return &Id{SessionId: sessionId, Name: name}
}

// SessionNode is the base class of all nodes that provide capability in an RTP session
type SessionNode struct {
	Id
	delegate *event.NodeDelegate
	ctrl     CommandInitiator

	eventTypeMatch []int
	eventHandler   []EventHandler

	Trait     *NodeTrait
	LinkPoint []LinkPoint
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
	s.SetHandler(NewLinkPoint, s.handleLinkPoint)
}

func (s *SessionNode) OnExit() {
	logger.Debugf("node(%v) exits graph", s.GetNodeName())
}

func (s *SessionNode) OnEvent(evt *event.Event) {
	cmd := evt.GetCmd()
	for i, typ := range s.eventTypeMatch {
		if cmd == typ {
			s.eventHandler[i](evt)
			return
		}
	}
}

func (s *SessionNode) OnLinkDown(linkId int, scope string, nodeName string) {
	logger.Debugf("node got link down (%v:%v) => (%v:%v) ", s.GetNodeScope(), s.GetNodeName(), scope, nodeName)
}

//--------------------------- Base SessionAware Implementation --------------------------------

func (s *SessionNode) handleLinkPoint(evt *event.Event) {
	msg, ok := ToMessage[*LinkPointCommand](evt)
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

func (s *SessionNode) StreamTo(session, name string) (err error, lp LinkPoint) {
	var linkId int
	if linkId = s.delegate.RequestLinkUp(session, name); linkId < 0 {
		err = fmt.Errorf("(%v:%v) can not set stream target to (%v:%v) due to request link up failed",
			s.SessionId, s.Name, session, name)
		return
	}
	sendFunc := func(msg Message) error {
		if !s.delegate.Deliver(linkId, msg.AsEvent()) {
			return fmt.Errorf("failed to deliver event")
		}
		return nil
	}
	linkIdentity := MakeLinkIdentity(session, name, linkId)
	newLinkCmd := &LinkPointCommand{
		//OfferedTrait: s.offer,
		LinkIdentity: linkIdentity,
		C:            make(chan *MessageTrait, 1),
	}
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
		}
		lp = &LinkPad{
			owner:        s,
			linkId:       linkId,
			identity:     linkIdentity,
			messageTrait: agreedTrait,
			sendFunc:     sendFunc,
		}
		s.LinkPoint = append(s.LinkPoint, lp)
		logger.Infof("new stream connection (%v|%v){%x} --->[%v]---> (%v|%v)",
			s.GetNodeScope(), s.GetNodeName(), linkIdentity, agreedTrait.Type.String(), session, name)
	case <-time.After(10 * time.Second):
		err = fmt.Errorf("(%v:%v) can not set stream target to (%v:%v) due to link point not retrieved",
			s.SessionId, s.Name, session, name)
	}
	return
}

//--------------------------- Facility methods --------------------------------

func (s *SessionNode) SetHandler(msgType int, handler EventHandler) {
	for i := 0; i < len(s.eventTypeMatch); i++ {
		if s.eventTypeMatch[i] == msgType {
			s.eventHandler[i] = handler
			return
		}
	}
	// not found, prepend to head of array
	s.eventTypeMatch = append([]int{msgType}, s.eventTypeMatch...)
	s.eventHandler = append([]EventHandler{handler}, s.eventHandler...)
}

// SendMessage use the first link point as in most use case
func (s *SessionNode) SendMessage(msg Message) error {
	if len(s.LinkPoint) == 0 {
		return fmt.Errorf("(%v:%v) has no link point", s.GetNodeScope(), s.GetNodeName())
	}
	return s.LinkPoint[0].SendMessage(msg)
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
