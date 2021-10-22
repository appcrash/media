package comp

import (
	"errors"
	"fmt"
	"github.com/appcrash/media/server/event"
	"sync"
)

// Dispatch is a bridge between every session's nodes and commands
// business commands normally change node's behaviour by send control messages to them with Call() or Cast(),
// nodes react to these messages which may cause more messages(change other node's behaviour,chain reaction) to send.
// all of these messages are sent from Dispatch.
// it is a SessionNode as well as implementing Controller interfaces, so provide api to callers while sending
// messages to any other nodes in the event graph.
type Dispatch struct {
	SessionNode
	event.NodeProperty // necessary when the node number exceeds default maxLink of event graph

	mutex   sync.Mutex
	linkMap map[string]int
}

//----------------------------------------- api & implementation ---------------------------------------------

// sync connect to other nodes, if the provided node list contains already-linked nodes, skip them
func (d *Dispatch) connectTo(nl []*Id) (err error) {
	if d.delegate == nil {
		err = errors.New("send node not ready")
		return
	}
	n := len(nl)
	if n > d.GetMaxLink() {
		err = errors.New("node list exceeds maxLink")
		return
	}

	// skip existing links
	var newNl []*Id
	d.mutex.Lock()
	for _, node := range nl {
		if _, ok := d.linkMap[node.String()]; !ok {
			newNl = append(newNl, node)
		}
	}
	d.mutex.Unlock()

	// connect to others and wait until done
	for _, node := range newNl {
		if linkId := d.delegate.RequestLinkUp(node.SessionId, node.Name); linkId >= 0 {
			d.mutex.Lock()
			d.linkMap[node.String()] = linkId
			d.mutex.Unlock()
		} else {
			err = fmt.Errorf("dispatch connect to %v failed", node.Name)
			return
		}
	}
	return
}

// getLinkId retrieve requested node with given sessionId and name, if the link to it
// has not established yet, connect to that node on the fly, then return created link id
func (d *Dispatch) getLinkId(sessionId, name string) (linkId int, err error) {
	var ok bool
	id := NewId(sessionId, name)
	d.mutex.Lock()
	linkId, ok = d.linkMap[id.String()]
	d.mutex.Unlock()
	if !ok {
		// has no link to the requested node, establish on the fly
		if err = d.connectTo([]*Id{id}); err != nil {
			return
		}
		d.mutex.Lock()
		linkId, _ = d.linkMap[id.String()]
		d.mutex.Unlock()
	}
	return
}

// Call send control message to a node in the graph and wait for its reply
// if session is "", it means sending to local session nodes
func (d *Dispatch) Call(session, name string, args []string) (resp []string) {
	if session == "" {
		session = d.SessionId
	}
	linkId, err := d.getLinkId(session, name)
	if err != nil {
		resp = WithError("can not connect to requested node")
		return
	}
	evt := NewCallEvent(args)
	if !d.delegate.Deliver(linkId, evt) {
		resp = []string{"err"}
		return
	}
	ctrlMsg := evt.GetObj().(*CtrlMessage)
	resp = <-ctrlMsg.C
	return
}

// Cast send control message to a node in the graph
// if session is "", it means send to local session nodes
func (d *Dispatch) Cast(session, name string, args []string) {
	if session == "" {
		session = d.SessionId
	}
	linkId, err := d.getLinkId(session, name)
	if err != nil {
		return
	}
	evt := NewCastEvent(args)
	if !d.delegate.Deliver(linkId, evt) {
		return
	}
	return
}

func newDispatch() SessionAware {
	n := &Dispatch{
		linkMap: make(map[string]int),
	}
	n.Name = TYPE_DISPATCH
	return n
}
