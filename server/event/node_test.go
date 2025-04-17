package event_test

import (
	"fmt"
	"github.com/appcrash/media/server/event"
	"math/rand/v2"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

type crashNode struct {
	test testNode
	i, n int
}

func Example_exceptionNodeOnEvent() {
	wg := &sync.WaitGroup{}
	wg.Add(6)
	excpNode1 := testNode{scope: "exception", name: "node1",
		onEvent: func(t *testNode, evt *event.Event) {
			cn := (*crashNode)(unsafe.Pointer(t))
			cn.i++
			if cn.i == cn.n {
				// access null pointer, panic...
				(*crashNode)(unsafe.Pointer(nil)).i++
			}
		},
		onExit: func(t *testNode) {
			fmt.Printf("%v exited\n", t.name)
			wg.Done()
		},
	}
	excpNode2 := excpNode1
	excpNode2.name = "node2"
	cn1, cn2 := crashNode{excpNode1, 0, 0}, crashNode{excpNode2, 0, 0}
	cn1.n = 10000 + rand.IntN(20000)
	cn2.n = 10000 + rand.IntN(20000)

	normalNode := testNode{scope: "test", name: "normal",
		onEnter: func(t *testNode) {
			id1 := t.delegate.RequestLinkUp("exception", "node1")
			id2 := t.delegate.RequestLinkUp("exception", "node2")
			f := func(id int, name string) {
				for {
					evt := event.NewEvent(cmd_nothing, 0)
					if ok := t.delegate.Deliver(id, evt); !ok {
						fmt.Printf("delivery stopped: %v\n", name)
						wg.Done()
						return
					}
				}
			}
			go f(id1, "node1")
			go f(id2, "node2")
		},
		onLinkDown: func(t *testNode, linkId int, scope string, nodeName string) {
			fmt.Printf("link down with %v\n", nodeName)
			wg.Done()
		},
	}

	graph := event.NewEventGraph()
	graph.AddNode(&cn1.test)
	graph.AddNode(&cn2.test)
	graph.AddNode(&normalNode)
	wg.Wait()

	// Unordered output:
	// delivery stopped: node1
	// delivery stopped: node2
	// link down with node1
	// link down with node2
	// node1 exited
	// node2 exited
}

func Example_exceptionNodeOnControl() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	panicNode := &testNode{scope: "panic", name: "panic",
		onEnter: func(t *testNode) {
			panic("I am nervous!")
		},
		onExit: func(t *testNode) {
			fmt.Println("panic exit")
			wg.Done()
		},
	}
	connPanicNode := &testNode{scope: "panic", name: "conn",
		onEnter: func(t *testNode) {
			t.delegate.RequestLinkUp("normal", "normal")
			panic("I am also nervous")
		},
		onExit: func(t *testNode) {
			fmt.Println("conn panic exit")
			wg.Done()
		},
	}
	normalNode := &testNode{scope: "normal", name: "normal"}
	graph := event.NewEventGraph()
	graph.AddNode(panicNode)
	graph.AddNode(normalNode)
	graph.AddNode(connPanicNode)
	wg.Wait()

	// Unordered output:
	// panic exit
	// conn panic exit
}

func TestMaxLink(t *testing.T) {
	var count int32
	//wg := sync.WaitGroup{}
	//wg.Add(3)
	tn := &testNode{
		onEnter: func(tn *testNode) {
			if tn.delegate.RequestLinkUp("test", "node1") >= 0 {
				count++
			}
			if tn.delegate.RequestLinkUp("test", "node2") >= 0 {
				count++
			}
			if tn.delegate.RequestLinkUp("test", "node3") >= 0 {
				count++
			}
		},
	}
	tn.SetMaxLink(2)
	graph := event.NewEventGraph()
	for i := 1; i <= 3; i++ {
		node := &testNode{scope: "test", name: "node" + strconv.Itoa(i)}
		graph.AddNode(node)
	}
	graph.AddNode(tn)
	//wg.Wait()
	if count != 2 {
		t.Fatal("link more than maxLink")
	}
}

func TestDataChannelSizeAndDeliveryTimeout(t *testing.T) {
	var timeout = time.Second * 1
	blockC := make(chan int)
	blockNode := &testNode{scope: "block", name: "block",
		onEvent: func(t *testNode, evt *event.Event) {
			// we are here when already consumed one event
			<-blockC
		},
	}
	blockNode.SetDataChannelSize(1)
	sendNode := &testNode{scope: "send", name: "send",
		onEnter: func(tn *testNode) {
			linkId := tn.delegate.RequestLinkUp("block", "block")
			// delivered and consumed, make the OnEvent stuck
			err1 := tn.delegate.Deliver(linkId, event.NewEvent(cmd_nothing, 1))

			// delivered to queue, make the queue full, not consumed
			err2 := tn.delegate.Deliver(linkId, event.NewEvent(cmd_nothing, 2))

			// can not be delivered as data channel is full, test the timeout
			now := time.Now()
			err3 := tn.delegate.Deliver(linkId, event.NewEvent(cmd_nothing, 3))
			duration := time.Since(now)
			if err1 != true || err2 != true || err3 != false {
				t.Errorf("should block the third event")
			}
			if duration < timeout {
				t.Errorf("delivery returns too early")
			}
			blockC <- 0
			blockC <- 0
		},
	}
	sendNode.SetDeliveryTimeout(timeout)
	graph := event.NewEventGraph()
	graph.AddNode(blockNode)
	graph.AddNode(sendNode)
}

type countNode testNode

func (c *countNode) Count(i int) {
	evt := event.NewEvent(cmd_nothing, i)
	c.delegate.DeliverSelf(evt)
}

func TestSelfDelivery(t *testing.T) {
	c := make(chan int)
	node := &countNode{scope: "count", name: "count",
		onEvent: func(t *testNode, evt *event.Event) {
			c <- evt.GetObj().(int)
		},
	}

	graph := event.NewEventGraph()
	graph.AddNode((*testNode)(node))

	go func() {
		for i := 1; i <= 500; i++ {
			node.Count(i)
		}
	}()

	for i := 1; i <= 500; i++ {
		v := <-c
		if v != i {
			t.Errorf("delivered data wrong")
		}
	}
}

func TestSyncEvent(t *testing.T) {
	ready := make(chan int)
	c := make(chan int)
	node := &countNode{scope: "count", name: "count",
		onEnter: func(t *testNode) {
			linkId := t.delegate.RequestLinkUp("nothing", "nothing")
			for i := 1; i <= 500; i++ {
				evt := event.NewEventWithCallback(cmd_nothing, nil, func() {
					c <- 0
				})
				t.delegate.Deliver(linkId, evt)
			}
		},
	}
	nothing := &testNode{scope: "nothing", name: "nothing"}

	go func() {
		for i := 1; i <= 500; i++ {
			<-c
		}
		ready <- 0
	}()

	graph := event.NewEventGraph()
	graph.AddNode(nothing)
	graph.AddNode((*testNode)(node))

	select {
	case <-ready:
	case <-time.After(2 * time.Second):
		t.Error("event callback not working")
	}
}

func TestOnExit(t *testing.T) {
	n := 50000 + rand.IntN(50000)
	var count int
	done := make(chan int)
	tn := &testNode{scope: "exit", name: "exit",
		onEvent: func(t *testNode, evt *event.Event) {
			count++
		},
		onExit: func(n *testNode) {
			done <- 0
		},
	}
	tn.SetDataChannelSize(n)
	graph := event.NewEventGraph()
	graph.AddNode(tn)
	go func() {
		for i := 0; i < n; i++ {
			ok := tn.delegate.DeliverSelf(event.NewEvent(cmd_nothing, 0))
			if !ok {
				t.Errorf("delivered failed")
			}
		}
		tn.delegate.RequestNodeExit()
	}()
	<-done
	if n != count {
		t.Fatal("count must equal to n, onExit must be called after all events handled")
	}
}

func TestOnExitUnderConcurrentDeliver(t *testing.T) {
	concurrent := 500
	done := make(chan int)
	wg := &sync.WaitGroup{}
	wg.Add(500)
	var deliveredEvents, handledEvents int32

	graph := event.NewEventGraph()
	tn := &testNode{scope: "exitConcurrent", name: "exitConcurrent",
		onEvent: func(t *testNode, evt *event.Event) {
			atomic.AddInt32(&handledEvents, 1)
		},
		onExit: func(n *testNode) {
			done <- 0
		},
	}
	tn.SetDeliveryTimeout(time.Duration(5) * time.Second)
	graph.AddNode(tn)
	for i := 0; i < concurrent; i++ {
		go func() {
			for {
				ok := tn.delegate.DeliverSelf(event.NewEvent(cmd_nothing, 0))
				if !ok {
					wg.Done()
					return
				} else {
					atomic.AddInt32(&deliveredEvents, 1)
				}
			}
		}()
	}
	go func() {
		delay := 5 + rand.IntN(5)
		<-time.After(time.Duration(delay) * time.Second)
		tn.delegate.RequestNodeExit()
	}()

	<-done
	wg.Wait()
	if handledEvents != deliveredEvents {
		t.Fatal("deliveredEvents must equal to handledEvents, onExit must be called after all events handled")
	}
}
