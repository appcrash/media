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

	inLinkPoint, outLinkPoint []LinkPoint

	eventTypeMatch []int
	eventHandler   []EventHandler
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
	s.AddHandler(NewLinkPoint, s.handleLinkPoint)
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

}

func (s *SessionNode) ExitGraph() {
	if s.delegate != nil {
		_ = s.delegate.RequestNodeExit()
	}
}

func (s *SessionNode) Init() error {
	return nil
}

func (s *SessionNode) OnCall(fromSession, fromNode string, args []string) (resp []string) {
	return WithOk()
}

func (s *SessionNode) OnCast(fromSession, fromNode string, args []string) {

}

func (s *SessionNode) StreamTo(session, name string, offer []*MessageTrait) (err error, lp LinkPoint) {
	var linkId int
	if linkId = s.delegate.RequestLinkUp(session, name); linkId < 0 {
		err = fmt.Errorf("(%v:%v) can not set stream target to %v:%v due to request link up failed",
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
		OfferedTrait: offer,
		LinkIdentity: linkIdentity,
		C:            make(chan *MessageTrait, 1),
	}
	evt := event.NewEvent(NewLinkPoint, newLinkCmd)
	if ok := s.delegate.Deliver(linkId, evt); !ok {
		err = fmt.Errorf("(%v:%v) can not set stream target to %v:%v due to deliver event failed",
			s.SessionId, s.Name, session, name)
		return
	}
	select {
	case agreedTrait := <-newLinkCmd.C:
		if agreedTrait == nil {
			err = fmt.Errorf("(%v:%v) can not set stream target to %v:%v due to offered traits not agreed",
				s.SessionId, s.Name, session, name)
		}
		lp = &LinkPad{
			owner:        s,
			linkId:       linkId,
			identity:     linkIdentity,
			messageTrait: agreedTrait,
			sendFunc:     sendFunc,
		}
		s.outLinkPoint = append(s.outLinkPoint, lp)
		logger.Infof("new stream connection (%v|%v) {%x} --->[%v]---> (%v|%v)",
			s.GetNodeScope(), s.GetNodeName(), linkIdentity, agreedTrait.Name, session, name)
	case <-time.After(10 * time.Second):
		err = fmt.Errorf("(%v:%v) can not set stream target to %v:%v due to link point not retrieved",
			s.SessionId, s.Name, session, name)
	}
	return
}

func (s *SessionNode) StreamBy(offer []*MessageTrait, linkIdentity uint64) *MessageTrait {

}

func (s *SessionNode) SetController(ctrl CommandInitiator) {
	s.ctrl = ctrl
}

//--------------------------- Facility methods --------------------------------

func (s *SessionNode) AddHandler(msgType int, handler EventHandler) {
	s.eventTypeMatch = append([]int{msgType}, s.eventTypeMatch...)
	s.eventHandler = append([]EventHandler{handler}, s.eventHandler...)
}

// Call forward to controller
func (s *SessionNode) Call(session, name string, args []string) (resp []string) {
	return s.ctrl.Call(session, name, args)
}

// Cast forward to controller
func (s *SessionNode) Cast(session, name string, args []string) {
	s.ctrl.Cast(session, name, args)
}

// ToMessage convert event object back to concrete message
func ToMessage[K Message](evt *event.Event) K {
	obj := evt.GetObj()
	if obj == nil {
		return nil
	}
	if msg, ok := obj.(K); ok {
		return msg
	}
	return nil
}

// MakeSessionNode factory method of all session aware nodes
func MakeSessionNode(nodeType string, sessionId string, props []*nmd.NodeProp) SessionAware {
	if nodeType == "" || sessionId == "" {
		logger.Errorln("make session node failed")
		return nil
	}
	props = append(props,
		&nmd.NodeProp{
			Key:   "SessionId",
			Type:  "str",
			Value: sessionId,
		})
	node := getNodeByName(nodeType)
	if node != nil {
		props = setNodeProperties(node, props)
		//node.ConfigProperties(props)
	}
	return node
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
