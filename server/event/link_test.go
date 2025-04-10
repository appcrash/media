package event_test

import (
	"fmt"
	"github.com/appcrash/media/server/event"
	"math/rand"
	"strconv"
	"sync"
)

func Example_sendEvent() {
	done := make(chan int)
	node1 := &testNode{scope: "scope1", name: "node1",
		onEvent: func(t *testNode, evt *event.Event) {
			if evt.GetCmd() == cmd_print_self {
				go func() {
					fmt.Printf("scope:%v,name:%v\n", t.scope, t.name)
					done <- 0
				}()
			}
		}}
	node2 := &testNode{scope: "scope2", name: "node2",
		onEnter: func(tn *testNode) {
			go func() {
				linkId := tn.delegate.RequestLinkUp("scope1", "node1")
				if linkId >= 0 {
					fmt.Println("got link scope1:node1")
				}
				evt := event.NewEvent(cmd_print_self, nil)
				tn.delegate.Deliver(linkId, evt)
				linkId = tn.delegate.RequestLinkUp("scope1", "node1")
				if linkId < 0 {
					fmt.Println("duplicated link scope1:node1")
					done <- 0
				}
			}()
		},
	}
	graph := event.NewEventGraph()
	graph.AddNode(node1)
	graph.AddNode(node2)
	_, _ = <-done, <-done

	// Unordered OUTPUT:
	// scope:scope1,name:node1
	// got link scope1:node1
	// duplicated link scope1:node1
}

func Example_linkDownNonExistent() {
	const loopNum = 5
	wg := &sync.WaitGroup{}
	wg.Add(loopNum)
	downnode := &testNode{scope: "downscope", name: "downnode"}
	baseNode := &testNode{scope: "ask_linkdown",
		onEnter: func(tn *testNode) {
			linkId := tn.delegate.RequestLinkUp("downscope", "downnode")
			tn.delegate.RequestLinkDown(linkId)
		},
		onLinkDown: func(tn *testNode, linkId int, _scope string, _nodeName string) {
			if linkId >= 0 {
				fmt.Printf("linkdown: %v\n", tn.name)
				if err := tn.delegate.RequestLinkDown(linkId + 1); err != nil {
					fmt.Printf("linkdown: non-existent link\n")
				} // request wrong link id
				wg.Done()
			}
		},
	}

	graph := event.NewEventGraph()
	graph.AddNode(downnode)
	for i := 0; i < loopNum; i++ {
		node := *baseNode
		node.name = "node" + strconv.Itoa(i)
		graph.AddNode(&node)
	}
	wg.Wait()

	// Unordered output:
	// linkdown: node0
	// linkdown: node1
	// linkdown: node2
	// linkdown: node3
	// linkdown: node4
	// linkdown: non-existent link
	// linkdown: non-existent link
	// linkdown: non-existent link
	// linkdown: non-existent link
	// linkdown: non-existent link
}

func Example_linkDuplicated() {
	done := make(chan int)
	downnode := &testNode{scope: "downscope", name: "downnode"}
	baseNode := &testNode{scope: "ask_linkdown", name: "down_same_link",
		onEnter: func(tn *testNode) {
			linkId := tn.delegate.RequestLinkUp("downscope", "downnode")
			dupId := tn.delegate.RequestLinkUp("downscope", "downnode")
			if dupId == -1 {
				fmt.Printf("linkup: request failed\n")
			}

			tn.delegate.RequestLinkDown(linkId)
		},
		onLinkDown: func(tn *testNode, linkId int, _scope string, _nodeName string) {
			if linkId >= 0 {
				fmt.Printf("linkdown: %v\n", tn.name)
				if err := tn.delegate.RequestLinkDown(linkId); err != nil {
					fmt.Printf("linkdown: link already down\n")
					done <- 0
				}
			}
		},
	}

	graph := event.NewEventGraph()
	graph.AddNode(downnode)
	graph.AddNode(baseNode)
	_ = <-done

	// Unordered output:
	// linkup: request failed
	// linkdown: down_same_link
	// linkdown: link already down
}

func Example_moreLink() {
	done := make(chan int)
	n1 := &testNode{scope: "target", name: "node1"}
	n2 := &testNode{scope: "target", name: "node2"}
	shootNode := &testNode{scope: "shoot", name: "shooter",
		onEnter: func(tn *testNode) {
			id1 := tn.delegate.RequestLinkUp("target", "node1")
			id2 := tn.delegate.RequestLinkUp("target", "node2")

			f := func(nd *event.NodeDelegate, id int) {
				for {
					// fire!
					bullets := 10000 + rand.Intn(50000)
					for bullets > 0 {
						bullets--
						evt := event.NewEvent(cmd_nothing, nil)
						nd.Deliver(id, evt)
					}
					nd.RequestLinkDown(id)
				}
			}
			fmt.Println("aiming target:node1")
			go f(tn.delegate, id1)
			fmt.Println("aiming target:node2")
			go f(tn.delegate, id2)

		},
		onLinkDown: func(tn *testNode, linkId int, scope string, nodeName string) {
			fmt.Printf("shot down %v:%v\n", scope, nodeName)
			done <- 0
		},
	}

	graph := event.NewEventGraph()
	graph.AddNode(n1)
	graph.AddNode(n2)
	graph.AddNode(shootNode)
	_, _ = <-done, <-done

	// Unordered output:
	// aiming target:node1
	// aiming target:node2
	// shot down target:node1
	// shot down target:node2
}

// test when receiving node exit, all senders would fail to deliver events
// and get link-down notification, the receiving node also gets OnExit callback
func Example_receiverExit() {
	wg := &sync.WaitGroup{}
	wg.Add(5) // 5 = aim x 2 + shotdown x 2 + bombExit x 1
	s1 := testNode{scope: "shoot", name: "shooter1",
		onEnter: func(tn *testNode) {
			linkId := tn.delegate.RequestLinkUp("target", "bomb")
			fmt.Printf("%v aiming target:bomb\n", tn.name)
			go func(nd *event.NodeDelegate) {
				for {
					// fire!
					var evt *event.Event
					criticalHit := 10000 + rand.Intn(50000)
					i := 0
					for {
						if i == criticalHit {
							evt = event.NewEvent(cmd_explode, nil)
						} else {
							evt = event.NewEvent(cmd_nothing, nil)
						}
						if ok := nd.Deliver(linkId, evt); !ok {
							fmt.Printf("%v target missing\n", tn.name)
							wg.Done()
							return
						}
						i++
					}
				}
			}(tn.delegate)
		},

		onLinkDown: func(tn *testNode, linkId int, scope string, nodeName string) {
			fmt.Printf("%v shot down %v:%v\n", tn.name, scope, nodeName)
			wg.Done()
		},
	}
	s2 := s1
	s2.name = "shooter2"

	initDone := make(chan int)
	earlyExitNode := testNode{scope: "target", name: "bomb",
		onEnter: func(t *testNode) {
			initDone <- 0
		},
		onEvent: func(t *testNode, evt *event.Event) {
			if evt.GetCmd() == cmd_explode {
				t.delegate.RequestNodeExit()
			}
		},
		onExit: func(t *testNode) {
			fmt.Println("bomb: exited")
			wg.Done()
		},
	}

	graph := event.NewEventGraph()
	graph.AddNode(&earlyExitNode)
	<-initDone
	graph.AddNode(&s1)
	graph.AddNode(&s2)
	wg.Wait()

	// Unordered output:
	// shooter1 aiming target:bomb
	// shooter2 aiming target:bomb
	// bomb: exited
	// shooter1 target missing
	// shooter2 target missing
	// shooter1 shot down target:bomb
	// shooter2 shot down target:bomb
}
