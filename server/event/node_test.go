package event_test

import (
	"fmt"
	"github.com/appcrash/media/server/event"
	"math/rand"
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

func ExampleExceptionNodeOnEvent() {
	rand.Seed(time.Now().UTC().UnixNano())
	wg := &sync.WaitGroup{}
	wg.Add(6)
	initDone := make(chan int)
	excpNode1 := testNode{scope: "exception", name: "node1",
		onEnter: func(t *testNode) {
			initDone <- 0
		},
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
	cn1.n = 10000 + rand.Intn(20000)
	cn2.n = 10000 + rand.Intn(20000)

	normalNode := testNode{scope: "test", name: "normal",
		onEnter: func(t *testNode) {
			t.delegate.RequestLinkUp("exception", "node1")
			t.delegate.RequestLinkUp("exception", "node2")
		},
		onLinkUp: func(t *testNode, linkId int, scope string, nodeName string) {
			go func(id int, name string) {
				for {
					evt := event.NewEvent(cmd_nothing, 0)
					if ok := t.delegate.Delivery(id, evt); !ok {
						fmt.Printf("delivery stopped: %v\n", name)
						wg.Done()
						return
					}
				}
			}(linkId, nodeName)
		},
		onLinkDown: func(t *testNode, linkId int, scope string, nodeName string) {
			fmt.Printf("link down with %v\n", nodeName)
			wg.Done()
		},
	}

	graph := event.NewEventGraph()
	graph.AddNode(&cn1.test)
	graph.AddNode(&cn2.test)
	_, _ = <-initDone, <-initDone
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

func ExampleExceptionNodeOnControl() {
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
	connPanicNode := &testNode{scope:"panic",name:"conn",
		onEnter: func(t *testNode) {
			t.delegate.RequestLinkUp("normal","normal")
		},
		onLinkUp: func(t *testNode, linkId int, scope string, nodeName string) {
			panic("I am also nervous")
		},
		onExit: func(t *testNode) {
			fmt.Println("conn panic exit")
			wg.Done()
		},
	}
	normalNode := &testNode{scope: "normal",name:"normal"}
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
	wg := sync.WaitGroup{}
	wg.Add(3)
	tn := &testNode{
		onEnter: func(tn *testNode) {
			tn.delegate.RequestLinkUp("test", "node1")
			tn.delegate.RequestLinkUp("test", "node2")
			tn.delegate.RequestLinkUp("test", "node3")
		},
		onLinkUp: func(tn *testNode, linkId int, scope string, nodeName string) {
			if linkId >= 0 {
				atomic.AddInt32(&count, 1)
			}
			wg.Done()
		},
		maxLink: 2,
	}
	graph := event.NewEventGraph()
	for i := 1; i <= 3; i++ {
		node := &testNode{scope: "test", name: "node" + strconv.Itoa(i)}
		graph.AddNode(node)
	}
	graph.AddNode(tn)
	wg.Wait()
	if count != 2 {
		t.Fatal("link more than maxLink")
	}
}

func TestDataChannelSizeAndDeliveryTimeout(t *testing.T) {
	var timeout = time.Second * 1
	blockC := make(chan int)
	initDone := make(chan int)
	done := make(chan int)
	blockNode := &testNode{scope: "block", name: "block", dataChannelSize: 1,
		onEnter: func(t *testNode) {
			initDone <- 0
		},
		onEvent: func(t *testNode, evt *event.Event) {
			// we are here when already consumed one event
			<-blockC
		},
	}
	sendNode := &testNode{scope: "send", name: "send", deliveryTimeout: timeout,
		onEnter: func(t *testNode) {
			t.delegate.RequestLinkUp("block", "block")
		},
		onLinkUp: func(tn *testNode, linkId int, scope string, nodeName string) {
			<-initDone
			// delivered and consumed, make the OnEvent stuck
			err1 := tn.delegate.Delivery(linkId, event.NewEvent(cmd_nothing, 1))

			// delivered to queue, make the queue full, not consumed
			err2 := tn.delegate.Delivery(linkId, event.NewEvent(cmd_nothing, 2))

			// can not be delivered as data channel is full, test the timeout
			now := time.Now()
			err3 := tn.delegate.Delivery(linkId, event.NewEvent(cmd_nothing, 3))
			duration := time.Since(now)
			if err1 != true || err2 != true || err3 != false {
				t.Errorf("should block the third event")
			}
			if duration < timeout {
				t.Errorf("delivery returns too early")
			}
			done <- 0
		},
	}
	graph := event.NewEventGraph()
	graph.AddNode(blockNode)
	graph.AddNode(sendNode)
	<-done
	blockC <- 0
	blockC <- 0
}
