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

type cloneableObj struct {
	data string
}

func (c *cloneableObj) Clone() comp.Cloneable {
	return &cloneableObj{c.data}
}

func (c *cloneableObj) String() string {
	return c.data
}

func newPrintNode() comp.SessionAware {
	return &printNode{}
}

func (p *printNode) OnEvent(e *event.Event) {
	switch e.GetCmd() {
	case comp.RawByte, comp.Generic:
		msg := e.GetObj()
		fmt.Printf("%v print %v\n", p.Name, msg)
	case comp.CtrlCall:
		msg := e.GetObj().(*comp.CtrlMessage)
		reply := comp.WithOk(p.Name)
		msg.C <- reply
	case comp.CtrlCast:
		msg := e.GetObj().(*comp.CtrlMessage)
		data := msg.M[0]
		p.SendMessage(comp.NewRawByteMessage(data))
	}
}

func TestComposerBasic(t *testing.T) {
	gd := `[entry payload_type='1'] -> [pubsub channel='src1,src2'];`
	c := comp.NewSessionComposer("test_session")
	graph := event.NewEventGraph()
	if err := c.ParseGraphDescription(gd); err != nil {
		t.Fatal("parse graph failed")
	}
	ch1 := make(chan *event.Event, 2)
	ch2 := make(chan *event.Event, 2)
	if err := c.ComposeNodes(graph); err != nil {
		t.Fatal("prepare node failed", err)
	}
	c.LinkChannel("src1", ch1)
	mp := c.GetMessageProvider(comp.TypeENTRY)
	mp.PushMessage(comp.NewRawByteMessage("hello"))
	evt := <-ch1
	if evt.GetObj().(comp.RawByteMessage).String() != "hello" {
		t.Fatal("send/recv message not equal for src1")
	}
	c.LinkChannel("src2", ch2)
	mp.PushMessage(comp.NewRawByteMessage("hello again"))
	evt = <-ch2
	if evt.GetObj().(comp.RawByteMessage).String() != "hello again" {
		t.Fatal("send/recv message not equal for src2")
	}
	evt = <-ch1
	if evt.GetObj().(comp.RawByteMessage).String() != "hello again" {
		t.Fatal("send/recv message not equal for src1 (again)")
	}
}

func TestComposerWrongNodeType(t *testing.T) {
	gd := "[aaa] -> [bbb]"
	c := comp.NewSessionComposer("test_session")
	if err := c.ParseGraphDescription(gd); err != nil {
		t.Fatal("parse graph failed")
	}
	if err := c.ComposeNodes(event.NewEventGraph()); err == nil {
		t.Fatal("should fail to prepare wrong type nodes")
	}
}

func ExampleComposerPubSub() {
	gd := `[entry payload_type='1'] -> [pubsub channel='src'] -> {[p1:print], [ps:pubsub]};
           [ps] -> {[p2:print], [p3:print]}`
	comp.RegisterNodeFactory("print", newPrintNode)
	c := comp.NewSessionComposer("test_session")
	graph := event.NewEventGraph()
	if err := c.ParseGraphDescription(gd); err != nil {
		fmt.Println("parse graph failed")
		return
	}
	ch := make(chan *event.Event, 2)
	if err := c.ComposeNodes(graph); err != nil {
		fmt.Println("prepare node failed")
		return
	}
	c.LinkChannel("src", ch)
	mp := c.GetMessageProvider(comp.TypeENTRY)
	mp.PushMessage(comp.NewRawByteMessage("foobar"))
	evt := <-ch
	msg := evt.GetObj().(comp.RawByteMessage).String()
	fmt.Printf("channel got %v\n", msg)
	time.Sleep(50 * time.Millisecond)

	// Unordered OUTPUT:
	// p1 print foobar
	// p2 print foobar
	// p3 print foobar
	// channel got foobar
}

func ExampleComposerGenericMessage() {
	gd := `[entry payload_type='1'] -> [pubsub channel='out'] -> [p1:print];`

	comp.RegisterNodeFactory("print", newPrintNode)
	c := comp.NewSessionComposer("test_session")
	graph := event.NewEventGraph()
	if err := c.ParseGraphDescription(gd); err != nil {
		fmt.Println("parse graph failed")
		return
	}
	ch := make(chan *event.Event, 2)
	if err := c.ComposeNodes(graph); err != nil {
		fmt.Println("prepare node failed")
		return
	}
	c.LinkChannel("out", ch)
	msg := &comp.GenericMessage{
		Subtype: "cloneable",
		Obj:     &cloneableObj{data: "cloneMe"},
	}
	mp := c.GetMessageProvider(comp.TypeENTRY)
	mp.PushMessage(msg)
	evt := <-ch
	obj := evt.GetObj()
	fmt.Printf("channel got %v\n", obj)
	time.Sleep(50 * time.Millisecond)

	// Unordered OUTPUT:
	// p1 print GenericMessage type:cloneable value:cloneMe
	// channel got GenericMessage type:cloneable value:cloneMe
}

func ExampleComposerMultipleEntry() {
	gd := `[e1:entry payload_type='1'] -> [pubsub] -> {[p1:print],[p2:print]};
           [e2:entry payload_type='1'] -> [ps:pubsub] -> {[p2:print],[p3:print]}`
	comp.RegisterNodeFactory("print", newPrintNode)
	c := comp.NewSessionComposer("test_session")
	graph := event.NewEventGraph()
	if err := c.ParseGraphDescription(gd); err != nil {
		fmt.Println("parse graph failed")
		return
	}
	if err := c.ComposeNodes(graph); err != nil {
		fmt.Println("prepare node failed")
		return
	}
	mp1 := c.GetMessageProvider("e1")
	mp2 := c.GetMessageProvider("e2")
	mp1.PushMessage(comp.NewRawByteMessage("foo"))
	mp2.PushMessage(comp.NewRawByteMessage("bar"))
	time.Sleep(50 * time.Millisecond)

	// Unordered OUTPUT:
	// p1 print foo
	// p2 print foo
	// p2 print bar
	// p3 print bar
}

func ExampleComposerController() {
	gd := `[entry payload_type='1'] -> [p1:print] -> [p2:print]`
	comp.RegisterNodeFactory("print", newPrintNode)
	c := comp.NewSessionComposer("test_session")
	graph := event.NewEventGraph()
	if err := c.ParseGraphDescription(gd); err != nil {
		fmt.Println("parse graph failed, ", err)
		return
	}
	if err := c.ComposeNodes(graph); err != nil {
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
	gd := `[entry payload_type='1'] -> [ps:pubsub] -> {[p1:print],[p2:print]}`
	comp.RegisterNodeFactory("print", newPrintNode)
	graph := event.NewEventGraph()
	c1 := comp.NewSessionComposer("test1")
	c2 := comp.NewSessionComposer("test2")
	if c1.ParseGraphDescription(gd) != nil || c2.ParseGraphDescription(gd) != nil {
		fmt.Println("parse graph failed")
		return
	}
	if c1.ComposeNodes(graph) != nil || c2.ComposeNodes(graph) != nil {
		fmt.Println("prepare node failed")
		return
	}
	mp1, mp2 := c1.GetMessageProvider("entry"), c2.GetMessageProvider("entry")
	ctrl1 := c1.GetController()
	mp1.PushMessage(comp.NewRawByteMessage("from_session_1"))
	mp2.PushMessage(comp.NewRawByteMessage("from_session_2"))
	connCmd := comp.WithConnect("test1", "p2")
	ctrl1.Call("test2", "ps", connCmd) // ask "test2:ps" to connect to "test1:p2"
	mp2.PushMessage(comp.NewRawByteMessage("from_session_2_again"))

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

func ExamplePubSubEnableDisable() {
	gd := `[e1:entry payload_type='1'] -> [pubsub] -> {[p1:print],[p2:print]}`
	comp.RegisterNodeFactory("print", newPrintNode)
	c := comp.NewSessionComposer("test_session")
	graph := event.NewEventGraph()
	if err := c.ParseGraphDescription(gd); err != nil {
		fmt.Println("parse graph failed")
		return
	}
	if err := c.ComposeNodes(graph); err != nil {
		fmt.Println("prepare node failed")
		return
	}
	mp := c.GetMessageProvider("e1")
	ctrl := c.GetController()
	mp.PushMessage(comp.NewRawByteMessage("foo"))
	ctrl.Call("", "pubsub", comp.With("disable", "node", "test_session", "p1"))
	mp.PushMessage(comp.NewRawByteMessage("bar"))
	ctrl.Call("", "pubsub", comp.With("enable", "node", "test_session", "p1"))
	mp.PushMessage(comp.NewRawByteMessage("foobar"))
	time.Sleep(50 * time.Millisecond)

	// Unordered OUTPUT:
	// p1 print foo
	// p2 print foo
	// p2 print bar
	// p1 print foobar
	// p2 print foobar
}

func ExampleComposerPushDataMessage() {
	gd := `[p1:print]; [p2:print];`

	comp.RegisterNodeFactory("print", newPrintNode)
	c := comp.NewSessionComposer("test_session")
	graph := event.NewEventGraph()
	if err := c.ParseGraphDescription(gd); err != nil {
		fmt.Println("parse graph failed")
		return
	}
	if err := c.ComposeNodes(graph); err != nil {
		fmt.Println("prepare node failed")
		return
	}
	ctrl := c.GetController()
	ctrl.PushData("p1", "", []byte("abc"))
	ctrl.PushData("p2", "", []byte("cba"))
	time.Sleep(50 * time.Millisecond)

	// Unordered OUTPUT:
	// p1 print abc
	// p2 print cba
}
