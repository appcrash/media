package comp

import (
	"errors"
	"fmt"
	"github.com/appcrash/media/server/comp/nmd"
	"github.com/appcrash/media/server/event"
	"reflect"
	"unsafe"
)

// SessionNode is the base class of all nodes that provide capability in an RTP session
type SessionNode struct {
	Id
	delegate *event.NodeDelegate
	ctrl     Controller

	// record where the pipe output to
	dataLinkId                    int
	dataSessionName, dataNodeName string
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
}

func (s *SessionNode) OnExit() {
	logger.Debugf("node(%v) exits graph", s.GetNodeName())
}

func (s *SessionNode) OnEvent(evt *event.Event) {
	logger.Debugf("node(%v) got event with cmd:%v", s.GetNodeName(), evt.GetCmd())
}

func (s *SessionNode) OnLinkDown(linkId int, scope string, nodeName string) {
	logger.Debugf("node got link down (%v:%v) => (%v:%v) ", s.GetNodeScope(), s.GetNodeName(), scope, nodeName)
	if linkId >= 0 && s.dataLinkId == linkId {
		s.dataLinkId = -1
	}
}

//--------------------------- Base SessionAware Implementation --------------------------------

func (s *SessionNode) ExitGraph() {
	if s.delegate != nil {
		_ = s.delegate.RequestNodeExit()
	}
}

func (s *SessionNode) ConfigProperties(_ []*nmd.NodeProp) {
}

func (s *SessionNode) Init() error {
	return nil
}

func (s *SessionNode) SetPipeOut(session, name string) error {
	if s.delegate == nil {
		return errors.New("delegate not ready when set pipe")
	}
	s.dataSessionName, s.dataNodeName = session, name
	if s.dataLinkId = s.delegate.RequestLinkUp(session, name); s.dataLinkId < 0 {
		return errors.New(fmt.Sprintf("can not set pipe to %v:%v", session, name))
	}
	return nil
}

func (s *SessionNode) SetController(ctrl Controller) {
	s.ctrl = ctrl
}

//--------------------------- Facility methods --------------------------------

// DataPipeReady return whether data link is established
func (s *SessionNode) DataPipeReady() bool {
	return s.dataLinkId >= 0
}

// SendMessage utility method to put data message to next node
func (s *SessionNode) SendMessage(msg Message) (err error) {
	if s.DataPipeReady() {
		evt := msg.AsEvent()
		s.delegate.Deliver(s.dataLinkId, evt)
	} else {
		err = errors.New("data link is not established")
	}
	return
}

// Call forward to controller
func (s *SessionNode) Call(session, name string, args []string) (resp []string) {
	return s.ctrl.Call(session, name, args)
}

func (s *SessionNode) CallSys(name string, args []string) (resp []string) {
	return s.Call(SYS_NODE_SCOPE, name, args)
}

// Cast forward to controller
func (s *SessionNode) Cast(session, name string, args []string) {
	s.ctrl.Cast(session, name, args)
}

func (s *SessionNode) CastSys(name string, args []string) {
	s.Cast(SYS_NODE_SCOPE, name, args)
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
		}, &nmd.NodeProp{
			Key:   "dataLinkId",
			Type:  "int",
			Value: -1,
		})
	node := getNodeByName(nodeType)
	if node != nil {
		props = setNodeProperties(node, props)
		node.ConfigProperties(props)
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
