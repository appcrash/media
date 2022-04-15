package event

import (
	"context"
	"errors"
	"github.com/appcrash/media/server/prom"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// NodeDelegate is delegated to interact with event graph directly
// it is the only way for Node to send events by NodeDelegate
type NodeDelegate struct {
	nodeImpl        Node
	id              string
	ctrlC           chan *Event
	dataC           chan *Event
	userEventDoneC  chan int
	graph           *Graph
	inExit          atomic.Value
	deliveryTimeout time.Duration // in milliseconds

	// how many concurrent delivering on the way, it is important for safe exiting
	deliveryCount int32

	invokeMutex sync.Mutex
	// lock-free dlink array with fixed size(maxLink)
	// even though links field and its pointers(*atomic.Value) is not
	// volatile(that means they can be cpu-local cached), it doesn't matter.
	// in fact, all of them are read-only once created,
	// accessing element by atomic.Load always gets the latest value it points to
	// the limitation of max links are the trade-off between performance and
	// flexibility. if links can dynamically grow and event delivery has to atomic
	// read it every time(instead of cached) because we have to rewrite this variable
	// when enlarging array. in real world, output link number is fixed in most case,
	// so putting this limitation doesn't hurt much
	links []*atomic.Value
	// recycled dlink slots are saved here
	freeLinkSlot []int
}

// cope with atomic(Store/Load) that doesn't permit nil value
var nullLink = &dlink{}

const (
	defaultMaxLink         = 4
	defaultDataChannelSize = 128
	defaultDeliveryTimeout = 100 * time.Millisecond
	defaultExitDelay       = 50 * time.Millisecond
)

func newNodeDelegate(graph *Graph, node Node, maxLink int) *NodeDelegate {
	delegate := &NodeDelegate{
		nodeImpl:       node,
		ctrlC:          make(chan *Event),
		userEventDoneC: make(chan int),
		graph:          graph,
	}
	delegate.id = node.GetNodeScope() + ":" + node.GetNodeName()
	delegate.inExit.Store(false)
	delegate.links = make([]*atomic.Value, maxLink)
	delegate.freeLinkSlot = make([]int, maxLink)
	for i := 0; i < maxLink; i++ {
		val := new(atomic.Value)
		val.Store(nullLink)
		delegate.links[i] = val
		delegate.freeLinkSlot[i] = maxLink - 1 - i
	}

	dataSize := defaultDataChannelSize
	deliveryTimeout := defaultDeliveryTimeout
	ps := reflect.ValueOf(node)
	elem := ps.Elem()
	if elem.Kind() == reflect.Struct {
		field := elem.FieldByName("dataChannelSize")
		if field.IsValid() {
			switch field.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if ds := int(field.Int()); ds > 0 {
					dataSize = ds
				}
			}
		}
		field = elem.FieldByName("deliveryTimeout")
		if field.IsValid() && field.Type().String() == "time.Duration" {
			to := (*time.Duration)(unsafe.Pointer(field.UnsafeAddr()))
			if *to > 0 {
				deliveryTimeout = *to
			}
		}
	}

	// only buffered channel can satisfy nonblock sending in most case
	delegate.dataC = make(chan *Event, dataSize)
	delegate.deliveryTimeout = deliveryTimeout
	return delegate
}

func (nd *NodeDelegate) getNodeName() string {
	return nd.nodeImpl.GetNodeName()
}

func (nd *NodeDelegate) getNodeScope() string {
	return nd.nodeImpl.GetNodeScope()
}

// unique id in event graph
func (nd *NodeDelegate) getId() string {
	return nd.id
}

func (nd *NodeDelegate) isExiting() bool {
	return nd.inExit.Load().(bool)
}

func (nd *NodeDelegate) setExiting() {
	nd.inExit.Store(true)
}

func (nd *NodeDelegate) receiveCtrl(evt *Event) {
	// sync receive, assume event graph is bug free
	nd.ctrlC <- evt
}

func (nd *NodeDelegate) receiveData(evt *Event, timeoutMs time.Duration) (ok bool) {
	if nd.isExiting() {
		// stop receiving any event if we are exiting
		return false
	}
	atomic.AddInt32(&nd.deliveryCount, 1)
	defer func() {
		atomic.AddInt32(&nd.deliveryCount, -1)
	}()

	// always try nonblock delivery first, take chance to avoid creating timer
	select {
	case nd.dataC <- evt:
		ok = true
		return
	default:
	}

	// receiver doesn't catch up with me, wait until timeout
	select {
	case nd.dataC <- evt:
		ok = true
	case <-time.After(timeoutMs):
	}
	return
}

func (nd *NodeDelegate) startEventLoop() {
	ctx, cancel := context.WithCancel(context.Background())
	go func(n *NodeDelegate, cancelUserLoop context.CancelFunc) {
		err := nd.systemEventLoop()
		if err != nil {
			// come to here as graph has a bug
			if logger != nil {
				logger.Errorf("[graph]: fatal error: %v\n", err)
			}
			nd.finalize(err, cancelUserLoop)
			return
		}
		// notify user event loop to drain events then stop
		cancelUserLoop()
	}(nd, cancel)
	go func(n *NodeDelegate, ctx context.Context, cancelUserLoop context.CancelFunc) {
		err := n.userEventLoop(ctx)
		if err != nil {
			// normal execution should not come here, we are out of loop
			// because either normal exit(err == nil) or exception
			// thrown(err != nil), it means this node is in abnormal state,
			// report to graph and finalize me
			n.finalize(err, cancelUserLoop)
		}
	}(nd, ctx, cancel)
}

func (nd *NodeDelegate) systemEventLoop() (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			prom.NodeSystemEventException.Inc()
			if logger != nil {
				logger.Errorf("[graph]: bug %v\n", r)
			}
			if err, ok = r.(error); ok {
				return
			}
			err = errors.New("unknown error type in event loop")
		}
	}()

	for {
		evt := <-nd.ctrlC
		nd.handleSystemEvent(evt)
		if evt.cmd == respNodeExit {
			// node-exit is the last ctrl message in node's lifecycle
			// end this goroutine now
			return
		}
	}
}

func (nd *NodeDelegate) preSystemInvoke() {
	nd.invokeMutex.Lock()
}
func (nd *NodeDelegate) postSystemInvoke() {
	if r := recover(); r != nil {
		nd.invokeMutex.Unlock()
		_ = nd.RequestNodeExit()
		return
	}
	nd.invokeMutex.Unlock()
}

func (nd *NodeDelegate) handleSystemEvent(evt *Event) {
	switch evt.cmd {
	case respNodeAdd:
		if resp, ok := evt.obj.(*nodeAddResponse); !ok {
			panic(errors.New("[graph]:node add response with wrong event object"))
		} else {
			go func(n *NodeDelegate, cb Callback) {
				n.preSystemInvoke()
				defer n.postSystemInvoke()
				n.nodeImpl.OnEnter(n)
				if cb != nil {
					cb()
				}
			}(resp.delegate, resp.cb)
		}
	case respNodeExit:
		if _, ok := evt.obj.(*nodeExitResponse); !ok {
			panic(errors.New("[graph]:node exit response with wrong event object"))
		} else {
			go func(n *NodeDelegate) {
				n.preSystemInvoke()
				defer n.postSystemInvoke()
				// block here until all user events handled then call OnExit
				<-nd.userEventDoneC
				if len(nd.dataC) > 0 && logger != nil {
					logger.Errorf("node(%v) exit with dataC size:%v\n", nd.getNodeName(), len(nd.dataC))
				}
				n.nodeImpl.OnExit()
			}(nd)
		}
	case respLinkUp:
		if resp, ok := evt.obj.(*linkUpResponse); !ok {
			panic(errors.New("[graph]:dlink up response with wrong event object"))
		} else {
			if resp.state != 0 {
				resp.c <- -1
				return
			}
			// if free list is not empty, recycle the slot
			var newLinkId int
			link := resp.link
			freeLen := len(nd.freeLinkSlot)
			if freeLen > 0 {
				var val *atomic.Value
				newLinkId = nd.freeLinkSlot[freeLen-1]
				nd.freeLinkSlot = nd.freeLinkSlot[:freeLen-1]
				val = nd.links[newLinkId]
				if val == nil {
					val = &atomic.Value{}
					nd.links[newLinkId] = val
				}
				// atomic rewrite old dlink info
				val.Store(link)
			}

			resp.c <- newLinkId
		}
	case respLinkDown:
		// we receive link_down either toNode exits the graph or we actively requested previously
		if resp, ok := evt.obj.(*linkDownResponse); !ok {
			panic(errors.New("[graph]:dlink down response with wrong event object"))
		} else {
			link := resp.link
			linkId := -1
			for i, v := range nd.links {
				l := v.Load().(*dlink)
				if l == link {
					linkId = i
					break
				}
			}
			if linkId < 0 {
				// wrong link passed to node
				if logger != nil {
					logger.Errorf("[node]: request link down with wrong linkId:%v\n", linkId)
				}
				return
			}
			// put index to free list, set the slot to null link so Deliver would fail
			nd.freeLinkSlot = append(nd.freeLinkSlot, linkId)
			nd.links[linkId].Store(nullLink)
			scope := link.toNode.getNodeScope()
			nodeName := link.toNode.getNodeName()
			go func(n *NodeDelegate, id int, s string, name string) {
				n.preSystemInvoke()
				defer n.postSystemInvoke()
				n.nodeImpl.OnLinkDown(id, s, name)
			}(nd, linkId, scope, nodeName)
		}
	default:
		panic(errors.New("[graph]: received unknown command"))
	}
}

func (nd *NodeDelegate) userEventLoop(ctx context.Context) (err error) {
	defer func() {
		// sync with node exit, notify that loop is done
		close(nd.userEventDoneC)
		if r := recover(); r != nil {
			var ok bool
			prom.NodeUserEventException.Inc()
			if err, ok = r.(error); ok {
				return
			}
			err = errors.New("unknown error type in event loop")
		}
	}()

	doneC := ctx.Done()
	dataC := nd.dataC
outerLoop:
	for {
		select {
		case evt := <-dataC:
			nd.handleUserEvent(evt)
		case <-doneC:
			// close instead of return immediately to let events buffered in dataC being drained by
			// OnEvent before OnExit invoked
			doneC = nil
			break outerLoop
		}
	}
	// we need to consume all user events before exit
	// simply set nd.dataC to nil is not a good way to stop receiving more events as:
	// 1. the new value of dataC may be not seen by other goroutines due to caches
	// 2. race detector would complain when running test
	// 3. this operation is not atomic and some goroutines may have already passed inExiting() check to
	// start writing to dataC even after dataC set to nil
	// we introduce extra overhead to record delivery count that is still on its way, once it reaches to
	// zero, node can assert no more event would come in as inExiting() is true right now, i.e. delivery
	// count can only be decreased to 0 or not changed(keep 0) at this moment
	ticker := time.NewTicker(defaultExitDelay)
	for {
		select {
		case evt, more := <-dataC:
			if !more {
				return
			}
			nd.handleUserEvent(evt)
		case <-ticker.C:
			if atomic.LoadInt32(&nd.deliveryCount) == 0 {
				// no more delivery can succeed now, it is safe to close dataC now
				ticker.Stop()
				close(nd.dataC)
				nd.dataC = nil
			}
		}
	}
}

func (nd *NodeDelegate) handleUserEvent(evt *Event) {
	nd.nodeImpl.OnEvent(evt)
	// check call back and invoke it if necessary
	if evt.cb != nil {
		evt.cb()
	}
}

func (nd *NodeDelegate) finalize(err error, cancel context.CancelFunc) {
	// exception handling
	// user code error happened and caught by node delegate, we will
	// 1. notify user event loop to exit, all senders' delivery would fail;
	// 2. request node exit to graph, the node would be notified
	// output links are down and moved out of the graph, the node's
	// system event loop would terminate when node-exit response received

	if err != nil && logger != nil {
		logger.Errorf("finalizing node with error:%v", err)
	}
	cancel()
	_ = nd.RequestNodeExit()
}

/*****************************************************
 ****************** APIs for end user ****************
 *****************************************************/

// RequestLinkUp [SYNC] node request dlink to other node of @param scope and name @param nodeName
// the request is passed to graph, then graph would create the dlink and notify the node
// delegate.  if the requested node doesn't exist or any error happened, linkId == -1
func (nd *NodeDelegate) RequestLinkUp(scope string, nodeName string) (linkId int) {
	// graph will give error response if duplication happened
	// i.e. (fromScope,fromName,toScope,toName) quaternion is unique across graph
	if scope == "" || nodeName == "" {
		return -1
	}
	c := make(chan int, 1)
	evt := newLinkUpRequest(nd, scope, nodeName, c)
	nd.graph.deliveryEvent(evt)
	linkId = <-c
	return
}

// RequestLinkDown [ASYNC] node request tearing down an output dlink, and node's
// OnLinkDown would be invoked once successfully tearing down
func (nd *NodeDelegate) RequestLinkDown(linkId int) (err error) {
	// what if linkId >= len(nd.links) when nd.links itself not protected by
	// mutex and not being atomic? nd.links can only increase, so linkId assigned
	// before is always safe to access. however, the value may have been changed
	// due to slot recycling, sanity checking is necessary. a normal-behavior node
	// should never see such case (i.e. pass wrong linkId in)
	if linkId >= len(nd.links) {
		err = errors.New("linkId out of range")
		return
	}
	link := nd.links[linkId].Load().(*dlink)
	if link.fromNode != nd {
		err = errors.New("wrong linkId with different fromNode(not you)")
		return
	}
	evt := newLinkDownRequest(link)
	nd.graph.deliveryEvent(evt)
	return
}

// RequestNodeExit [ASYNC] ask graph to remove this node, and OnExit would be invoked
func (nd *NodeDelegate) RequestNodeExit() (err error) {
	if nd.isExiting() {
		err = errors.New("node is already in exiting state")
		return
	}

	// ensure all input/output links are torn down before the node gets
	// out of the graph, send exit request to graph, and graph will handle
	// all of this
	nd.setExiting()
	nd.graph.deliveryEvent(newNodeExitRequest(nd))
	return
}

// DeliverWithTimeout [SYNC] return true if successfully delivered
func (nd *NodeDelegate) DeliverWithTimeout(linkId int, evt *Event, timeout time.Duration) bool {
	link := nd.links[linkId].Load().(*dlink)
	if link == nullLink {
		return false
	}
	if link.fromNode != nd {
		// if it is an input dlink
		return false
	}
	return link.toNode.receiveData(evt, timeout)
}

func (nd *NodeDelegate) Deliver(linkId int, evt *Event) bool {
	return nd.DeliverWithTimeout(linkId, evt, nd.deliveryTimeout)
}

// DeliverSelf [SYNC] directly puts event to this node's event loop
// it is a convenient way to talk to the node, and node can choose to expose api to let
// caller who has a reference to this node directly sending message to node
// from the node's perspective, it doesn't care about the source of every event
func (nd *NodeDelegate) DeliverSelf(evt *Event) bool {
	return nd.receiveData(evt, nd.deliveryTimeout)
}
