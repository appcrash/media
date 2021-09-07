package comp

import (
	"errors"
	"fmt"
	"github.com/appcrash/media/server/event"
	"github.com/appcrash/media/server/utils"
	"sync"
	"time"
)

// Dispatch is the bridge between every session's nodes and commands
// business commands normally change node's behaviour by send control messages to them with Call() or Cast(),
// nodes react to these messages which may cause more messages(change other node's behaviour,chain reaction) to send.
// all of these messages are sent from Dispatch.
// it is a SessionNode as well as implementing Controller interfaces, so provide api to callers and send messages to
// any other nodes in the event graph.
type Dispatch struct {
	SessionNode
	event.NodeProperty // necessary when the node number exceeds default maxLink of event graph

	mutex     sync.Mutex
	linkMap   map[string]int
	linkSyncC chan bool // sync all link up event
}

func (d *Dispatch) OnLinkUp(linkId int, scope string, nodeName string) {
	if linkId >= 0 {
		d.mutex.Lock()
		id := &Id{SessionId: scope, Name: nodeName}
		d.linkMap[id.String()] = linkId
		d.mutex.Unlock()
		d.linkSyncC <- true
	}
}

//----------------------------------------- api & implementation ---------------------------------------------

// sync connect to other nodes
// TODO: can be called multiple times, skipping existing links
func (d *Dispatch) ConnectTo(nl []*Id) (err error) {
	if d.delegate == nil {
		err = errors.New("send node not ready")
		return
	}
	n := len(nl)
	if n > d.GetMaxLink() {
		err = errors.New("node list exceeds maxLink")
		return
	}
	d.linkSyncC = make(chan bool, n)
	// connect to others and wait until done
	for _, node := range nl {
		d.delegate.RequestLinkUp(node.SessionId, node.Name)
	}
	err = utils.WaitChannelWithTimeout(d.linkSyncC, n, 2*time.Second)
	// not used anymore
	d.linkSyncC = nil
	return
}

func (d *Dispatch) SendTo(sessionId, name string, evt *event.Event) (err error) {
	id := &Id{SessionId: sessionId, Name: name}
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if linkId, ok := d.linkMap[id.String()]; ok {
		d.delegate.Delivery(linkId, evt)
	} else {
		errStr := fmt.Sprintf("no such node:%v", id.String())
		err = errors.New(errStr)
	}
	return
}

func (d *Dispatch) Call(session, name, args string) (resp string) {
	return "ok"
}

func (d *Dispatch) Cast(session, name, args string) {
}

func newDispatch() *Dispatch {
	n := &Dispatch{
		linkMap: make(map[string]int),
	}
	return n
}
