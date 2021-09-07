package comp

import (
	"errors"
	"fmt"
	"github.com/appcrash/media/server/event"
	"reflect"
	"regexp"
	"strings"
	"unsafe"
)

type ConfigItems map[string]interface{}

// SessionNode is the base class of all nodes that provide capability in an RTP session
type SessionNode struct {
	Id
	delegate *event.NodeDelegate
	ctrl     Controller

	// record where the pipe output to
	dataLinkId                    int
	dataSessionName, dataNodeName string
}

// SessionAware enables node to:
// 1. config its static properties before any event flows
// 2. set data output destination
// 3. exit the graph when session ends
type SessionAware interface {
	event.Node

	// ConfigProperties handles props that can not be configured by simple reflection
	ConfigProperties(ci ConfigItems)

	// SetPipeOut specify the data endpoint to which this node output
	SetPipeOut(session, name string) error

	SetController(ctrl Controller)

	// ExitGraph is used when initialization failed or session terminated
	ExitGraph()
}

//------------------------- ConfigItems -------------------------

const regCamelCasePattern = `_+[a-z]`

var regCamelCase = regexp.MustCompile(regCamelCasePattern)

// Set converts key in form of foo_bar or foo__bar ... into fooBar if possible
// normal keys with camel case remain intact
func (ci ConfigItems) Set(key string, val interface{}) {
	newKey := regCamelCase.ReplaceAllStringFunc(key, func(match string) string {
		last := string(match[len(match)-1])
		return strings.ToUpper(last)
	})
	ci[newKey] = val
}

//------------------- Base Node Implementation -------------------------

func (s *SessionNode) GetNodeName() string {
	return s.Name
}

func (s *SessionNode) GetNodeScope() string {
	return s.SessionId
}

func (s *SessionNode) OnEnter(delegate *event.NodeDelegate) {
	logger.Debugf("node(%v) enters graph\n", s.GetNodeName())
	s.delegate = delegate
}

func (s *SessionNode) OnExit() {
	logger.Debugf("node(%v) exits graph", s.GetNodeName())
}

func (s *SessionNode) OnEvent(evt *event.Event) {
	logger.Debugf("node(%v) got event with cmd:%v\n", s.GetNodeName(), evt.GetCmd())
}

func (s *SessionNode) OnLinkDown(linkId int, scope string, nodeName string) {
	logger.Debugf("node(%v) got link to (%v:%v) down\n", s.GetNodeName(), scope, nodeName)
}

//--------------------------- Base SessionAware Implementation --------------------------------

func (s *SessionNode) ExitGraph() {
	if s.delegate != nil {
		s.delegate.RequestNodeExit()
	}
}

func (s *SessionNode) ConfigProperties(ci ConfigItems) {
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

// SendData utility method to put data message to next node
func (s *SessionNode) SendData(msg DataMessage) (err error) {
	if s.dataLinkId >= 0 {
		evt := NewDataEvent(msg)
		s.delegate.Deliver(s.dataLinkId, evt)
	} else {
		err = errors.New("data link is not established")
	}
	return
}

// MakeSessionNode factory method of all session aware nodes
func MakeSessionNode(nodeType string, sessionId string, ci ConfigItems) SessionAware {
	if ci == nil || nodeType == "" || sessionId == "" {
		logger.Errorln("make session node failed")
		return nil
	}
	ci.Set("SessionId", sessionId) // always set it
	node := getNodeByName(nodeType)
	if node != nil {
		ci = setNodeProperties(node, ci)
		node.ConfigProperties(ci)
	}
	return node
}

// setNodeProperties use reflection to set fields by Name, it is cornerstone of config by scripting
func setNodeProperties(node event.Node, ci ConfigItems) (nci ConfigItems) {
	ns := reflect.ValueOf(node).Elem()
	for k, v := range ci {
		if v == nil {
			continue
		}
		field := ns.FieldByName(k)
		rv := reflect.ValueOf(v)
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
		if nci == nil {
			nci = make(ConfigItems)
		}
		nci[k] = v
	}
	return
}
