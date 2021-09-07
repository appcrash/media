package comp_test

import (
	"fmt"
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/event"
	"sync"
	"time"
)

type psTestNode struct {
	delegate  *event.NodeDelegate
	c         chan int
	readyFunc event.Callback
}

type cloneableMsg string

func (m cloneableMsg) Clone() comp.Cloneable {
	return m
}

func (p *psTestNode) GetNodeName() string {
	return "psTest"
}

func (p *psTestNode) GetNodeScope() string {
	return "psTest"
}

func (p *psTestNode) OnEvent(evt *event.Event) {
	switch evt.GetCmd() {
	case comp.DATA_OUTPUT:
		s := evt.GetObj().(cloneableMsg)
		fmt.Printf("node: %v\n", s)
		p.c <- 0
	}
}

func (p *psTestNode) OnLinkUp(linkId int, scope string, nodeName string) {

}

func (p *psTestNode) OnLinkDown(linkId int, scope string, nodeName string) {

}

func (p *psTestNode) OnEnter(delegate *event.NodeDelegate) {
	p.delegate = delegate
	if p.readyFunc != nil {
		p.readyFunc()
	}
}

func (p *psTestNode) OnExit() {

}

func ExampleTwoSubscriber() {
	c1 := make(chan *event.Event, 2)
	c2 := make(chan int)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	tn := &psTestNode{c: c2,
		readyFunc: func() {
			wg.Done()
		},
	}
	ci := comp.ConfigItems{
		"readyFunc": func() {
			wg.Done()
		},
	}
	ps := comp.MakeSessionNode("pubsub", "any_session", ci).(*comp.PubSubNode)

	graph := event.NewEventGraph()
	graph.AddNode(ps)
	graph.AddNode(tn)

	wg.Wait()
	ps.SubscribeNode("psTest", "psTest")
	ps.SubscribeChannel("my_channel", c1)
	<-time.After(100 * time.Millisecond)
	go func() {
		var m1 cloneableMsg = "hello"
		var m2 cloneableMsg = "world"
		ps.Publish(m1)
		ps.Publish(m2)
	}()

	for i := 0; i < 4; i++ {
		select {
		case evt := <-c1:
			if v, ok := evt.GetObj().(cloneableMsg); ok {
				fmt.Printf("channel: %v\n", v)
			} else {
				fmt.Printf("error when converting\n")
			}
		case <-c2:
		}
	}

	ps.UnsubscribeChannel("my_channel")
	<-time.After(100 * time.Millisecond)
	go func() {
		var m cloneableMsg = "after unsubscribe"
		ps.Publish(m)
	}()
	<-c2

	ps.SubscribeChannel("my_channel", c1)
	ps.UnsubscribeNode("psTest", "psTest")
	<-time.After(100 * time.Millisecond)
	go func() {
		var m cloneableMsg = "come back"
		ps.Publish(m)
	}()
	evt := <-c1
	fmt.Printf("channel: %v\n", evt.GetObj().(cloneableMsg))

	// Unordered OUTPUT:
	// node: hello
	// node: world
	// channel: hello
	// channel: world
	// node: after unsubscribe
	// channel: come back
}
