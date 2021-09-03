package event_test

import (
	"github.com/appcrash/media/server/event"
	"time"
)

type onEnterExitFuncType func(t *testNode)
type onEventFuncType func(t *testNode, evt *event.Event)
type onLinkFuncType func(t *testNode, linkId int, scope string, nodeName string)

const (
	cmd_print_self = iota
	cmd_nothing
	cmd_explode
)

type testNode struct {
	scope      string
	name       string
	delegate   *event.NodeDelegate
	onEvent    onEventFuncType
	onEnter    onEnterExitFuncType
	onExit     onEnterExitFuncType
	onLinkUp   onLinkFuncType
	onLinkDown onLinkFuncType

	// optional attributes
	maxLink         int
	dataChannelSize int
	deliveryTimeout time.Duration
}

func (t *testNode) GetNodeName() string {
	return t.name
}

func (t *testNode) GetNodeScope() string {
	return t.scope
}

func (t *testNode) OnEvent(evt *event.Event) {
	if t.onEvent != nil {
		t.onEvent(t, evt)
	}
}

func (t *testNode) OnLinkUp(linkId int, scope string, nodeName string) {
	if t.onLinkUp != nil {
		t.onLinkUp(t, linkId, scope, nodeName)
	}
}

func (t *testNode) OnLinkDown(linkId int, scope string, nodeName string) {
	if t.onLinkDown != nil {
		t.onLinkDown(t, linkId, scope, nodeName)
	}
}

func (t *testNode) OnEnter(delegate *event.NodeDelegate) {
	t.delegate = delegate
	if t.onEnter != nil {
		t.onEnter(t)
	}
}

func (t *testNode) OnExit() {
	if t.onExit != nil {
		t.onExit(t)
	}
}
