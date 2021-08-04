package event_test

import (
	"fmt"
	"github.com/appcrash/media/server/event"
	"math/rand"
	"sync"
	"time"
	"unsafe"
)

type crashNode struct {
	test testNode
	i, n int
}

func ExampleExceptionNode() {
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
