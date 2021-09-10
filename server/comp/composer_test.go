package comp_test

import (
	"fmt"
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/event"
	"testing"
	"time"
)

type printNode struct {
	comp.SessionNode
}

func newPrintNode() comp.SessionAware {
	return &printNode{}
}

func (p *printNode) OnEvent(e *event.Event) {
	switch e.GetCmd() {
	case comp.DATA_OUTPUT:
		msg := e.GetObj().(comp.DataMessage)
		fmt.Printf("%v print %v\n", p.Name, msg)
	case comp.CTRL_CALL:
		msg := e.GetObj().(*comp.CtrlMessage)
		reply := comp.WithOk(p.Name)
		msg.C <- reply
	case comp.CTRL_CAST:
		msg := e.GetObj().(*comp.CtrlMessage)
		data := msg.M[0]
		p.SendData(comp.NewDataMessage(data))
	}
}

func TestComposerBasic(t *testing.T) {
	gd := "[pubsub]: channel=src1,src2 \n" +
		"[entry] -> [pubsub]"
	c := comp.NewSessionComposer("test_session")
	graph := event.NewEventGraph()
	if err := c.ParseGraphDescription(gd); err != nil {
		t.Fatal("parse graph failed")
	}
	ch1 := make(chan *event.Event, 2)
	ch2 := make(chan *event.Event, 2)
	c.RegisterChannel("src1", ch1) // statically register
	if err := c.PrepareNodes(graph); err != nil {
		t.Fatal("prepare node failed", err)
	}

	mp := c.GetMessageProvider(comp.TYPE_ENTRY)
	mp.PushMessage(comp.NewDataMessage("hello"))
	evt := <-ch1
	if evt.GetObj().(comp.DataMessage).String() != "hello" {
		t.Fatal("send/recv message not equal for src1")
	}
	c.RegisterChannel("src2", ch2)
	mp.PushMessage(comp.NewDataMessage("hello again"))
	evt = <-ch2
	if evt.GetObj().(comp.DataMessage).String() != "hello again" {
		t.Fatal("send/recv message not equal for src2")
	}
	evt = <-ch1
	if evt.GetObj().(comp.DataMessage).String() != "hello again" {
		t.Fatal("send/recv message not equal for src1 (again)")
	}
}

func TestComposerWrongNodeType(t *testing.T) {
	gd := "[aaa] -> [bbb]"
	c := comp.NewSessionComposer("test_session")
	if err := c.ParseGraphDescription(gd); err != nil {
		t.Fatal("parse graph failed")
	}
	if err := c.PrepareNodes(event.NewEventGraph()); err == nil {
		t.Fatal("should fail to prepare wrong type nodes")
	}
}

func ExampleComposerPubSub() {
	gd := "[pubsub]: channel=src \n" +
		"[entry] -> [pubsub] \n" +
		"[pubsub] -> [p1:print] \n" +
		"[pubsub] -> [ps:pubsub] \n" +
		"[ps:pubsub] -> [p2:print] \n" +
		"[ps:pubsub] -> [p3:print]"
	comp.RegisterNodeFactory("print", newPrintNode)
	c := comp.NewSessionComposer("test_session")
	graph := event.NewEventGraph()
	if err := c.ParseGraphDescription(gd); err != nil {
		fmt.Println("parse graph failed")
		return
	}
	ch := make(chan *event.Event, 2)
	c.RegisterChannel("src", ch)
	if err := c.PrepareNodes(graph); err != nil {
		fmt.Println("prepare node failed")
		return
	}
	mp := c.GetMessageProvider(comp.TYPE_ENTRY)
	mp.PushMessage(comp.NewDataMessage("foobar"))
	evt := <-ch
	msg := evt.GetObj().(comp.DataMessage).String()
	fmt.Printf("channel got %v\n", msg)
	time.Sleep(50 * time.Millisecond)

	// Unordered OUTPUT:
	// p1 print foobar
	// p2 print foobar
	// p3 print foobar
	// channel got foobar
}

func ExampleComposerMultipleEntry() {
	gd := "[e1:entry] -> [pubsub] \n" +
		"[pubsub] -> [p1:print] \n" +
		"[pubsub] -> [p2:print] \n" +
		"[e2:entry] -> [ps:pubsub] \n" +
		"[ps:pubsub] -> [p2:print] \n" +
		"[ps:pubsub] -> [p3:print]"
	comp.RegisterNodeFactory("print", newPrintNode)
	c := comp.NewSessionComposer("test_session")
	graph := event.NewEventGraph()
	if err := c.ParseGraphDescription(gd); err != nil {
		fmt.Println("parse graph failed")
		return
	}
	if err := c.PrepareNodes(graph); err != nil {
		fmt.Println("prepare node failed")
		return
	}
	mp1 := c.GetMessageProvider("e1")
	mp2 := c.GetMessageProvider("e2")
	mp1.PushMessage(comp.NewDataMessage("foo"))
	mp2.PushMessage(comp.NewDataMessage("bar"))
	time.Sleep(50 * time.Millisecond)

	// Unordered OUTPUT:
	// p1 print foo
	// p2 print foo
	// p2 print bar
	// p3 print bar
}

func ExampleComposerController() {
	gd := "[entry] -> [p1:print] \n" +
		"[p1] -> [p2:print]"
	comp.RegisterNodeFactory("print", newPrintNode)
	c := comp.NewSessionComposer("test_session")
	graph := event.NewEventGraph()
	if err := c.ParseGraphDescription(gd); err != nil {
		fmt.Println("parse graph failed, ", err)
		return
	}
	if err := c.PrepareNodes(graph); err != nil {
		fmt.Println("prepare node failed, ", err)
		return
	}
	ctrl := c.GetController()
	ctrl.Cast("", "p1", []string{"foobar"})
	reply := ctrl.Call("", "p1", []string{})
	fmt.Printf("%v: %v\n", reply[0], reply[1])
	time.Sleep(50 * time.Millisecond)

	// Unordered OUTPUT:
	// p2 print foobar
	// ok: p1
}

func ExampleComposerInterSession() {
	gd := "[entry] -> [ps:pubsub] \n" +
		"[ps] -> [p1:print] \n" +
		"[ps] -> [p2:print]"
	comp.RegisterNodeFactory("print", newPrintNode)
	graph := event.NewEventGraph()
	c1 := comp.NewSessionComposer("test1")
	c2 := comp.NewSessionComposer("test2")
	if c1.ParseGraphDescription(gd) != nil || c2.ParseGraphDescription(gd) != nil {
		fmt.Println("parse graph failed")
		return
	}
	if c1.PrepareNodes(graph) != nil || c2.PrepareNodes(graph) != nil {
		fmt.Println("prepare node failed")
		return
	}
	mp1, mp2 := c1.GetMessageProvider("entry"), c2.GetMessageProvider("entry")
	ctrl1 := c1.GetController()
	mp1.PushMessage(comp.NewDataMessage("from_session_1"))
	mp2.PushMessage(comp.NewDataMessage("from_session_2"))
	connCmd := comp.WithConnect("test1", "p2")
	ctrl1.Call("test2", "ps", connCmd) // ask "test2:ps" to connect to "test1:p2"
	mp2.PushMessage(comp.NewDataMessage("from_session_2_again"))

	time.Sleep(50 * time.Millisecond)
	// Unordered OUTPUT:
	// p1 print from_session_1
	// p2 print from_session_1
	// p1 print from_session_2
	// p2 print from_session_2
	// p1 print from_session_2_again
	// p2 print from_session_2_again
	// p2 print from_session_2_again

}
