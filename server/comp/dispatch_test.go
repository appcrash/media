package comp_test

import (
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/event"
	"sync"
	"testing"
)

func TestSendCommand(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	ci := make(comp.ConfigItems)
	ci.Set("readyFunc", func() { wg.Done() })
	ci["maxLink"] = 1
	sn := comp.MakeSessionNode("send", "some_session", ci).(*comp.Dispatch)
	ps := comp.MakeSessionNode("pubsub", "some_session", ci).(*comp.PubSubNode)
	graph := event.NewEventGraph()

	graph.AddNode(ps)
	graph.AddNode(sn)

	wg.Wait()
	var msg cloneableMsg = "amsg"
	evt := event.NewEvent(comp.DATA_OUTPUT, msg)
	c := make(chan *event.Event, 1)
	nl := []*comp.Id{comp.NewId("some_session", comp.TYPE_PUBSUB)}
	ps.SubscribeChannel("testC", c)
	sn.ConnectTo(nl)
	sn.SendTo("some_session", comp.TYPE_PUBSUB, evt)
	revt := <-c
	rmsg := revt.GetObj().(cloneableMsg)
	if msg != rmsg {
		t.Fatal("received msg not equal to sent msg")
	}
}
