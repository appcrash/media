package comp_test

import (
	"bytes"
	"fmt"
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/event"
	"testing"
	"time"
)

const mtCustom = comp.MtUserMessageBegin

func (m *customMessage) Type() comp.MessageType {
	return mtCustom
}

func (m *customMessage) AsEvent() *event.Event {
	return event.NewEvent(mtCustom, m)
}

func (m *customMessage) Clone() comp.Cloneable {
	return &customMessage{Value: m.Value}
}

type customMessageConvertable interface {
	AscustomMessage() *customMessage
}

type customMessage struct {
	comp.MessageBase
	Value string
}

func (m *customMessage) AsRawByteMessage() *comp.RawByteMessage {
	return &comp.RawByteMessage{
		Data: []byte(m.Value),
	}
}

type fireNode struct {
	comp.SessionNode
	comp.InitiatorNode
}

func (n *fireNode) Offer() []comp.MessageType {
	return []comp.MessageType{mtCustom}
}

func (n *fireNode) OnCall(fromNode string, args []string) (resp []string) {
	msg := &customMessage{
		Value: args[0],
	}
	fmt.Printf("fire call from %v\n", fromNode)
	n.GetLinkPoint(0).SendMessage(msg)
	return comp.WithOk()
}

func (n *fireNode) OnCast(fromNode string, args []string) {
	msg := &customMessage{
		Value: args[0],
	}
	fmt.Printf("fire cast from %v\n", fromNode)
	n.GetLinkPoint(0).SendMessage(msg)
}

func newFireNode() comp.SessionAware {
	n := &fireNode{}
	n.Trait, _ = comp.NodeTraitOfType("fire")
	return n
}

type printNode struct {
	comp.SessionNode
}

func (n *printNode) Accept() []comp.MessageType {
	return []comp.MessageType{
		comp.MtRawByte,
	}
}

func (n *printNode) handleRawByteEvent(evt *event.Event) {
	if msg, ok := comp.EventToMessage[*comp.RawByteMessage](evt); ok {
		fmt.Printf("%v print %v\n", n.GetNodeName(), string(msg.Data))
	}
}

func newPrintNode() comp.SessionAware {
	p := &printNode{}
	p.Trait, _ = comp.NodeTraitOfType("print")
	p.SetMessageHandler(comp.MtRawByte, comp.ChainSetHandler(p.handleRawByteEvent))
	return p
}

func initComposer() {
	comp.AddMessageTrait(comp.MT[customMessage](comp.MetaType[customMessageConvertable]()))
	comp.SetMessageConvertable(mtCustom, comp.MtRawByte)
	comp.RegisterNodeTrait(comp.NT[printNode]("print", newPrintNode))
	comp.RegisterNodeTrait(comp.NT[fireNode]("fire", newFireNode))
}

func composeIt(session, gd string) (*comp.Composer, error) {
	c := comp.NewSessionComposer(session)
	if err := c.ParseGraphDescription(gd); err != nil {
		return nil, fmt.Errorf("parse graph failed: %v", gd)
	}
	graph := event.NewEventGraph()
	if err := c.ComposeNodes(graph); err != nil {
		return nil, err
	}
	return c, nil
}

func TestComposerBasic(t *testing.T) {
	gd := `[input:chan_src] -> [pubsub] -> {[output1:chan_sink],[output2:chan_sink]};`
	c, err := composeIt("test_session", gd)
	if err != nil {
		t.Fatal(err)
	}
	inputC := make(chan []byte)
	outputC1, outputC2 := make(chan []byte, 2), make(chan []byte, 2)

	c.GetNode("input").(*comp.ChanSrc).LinkMe(inputC)
	c.GetNode("output1").(*comp.ChanSink).LinkMe(outputC1)
	c.GetNode("output2").(*comp.ChanSink).LinkMe(outputC2)
	testBytes := []byte("test chan_src,pubsub,chan_sink")
	inputC <- testBytes
	if bytes.Compare(<-outputC1, testBytes) != 0 {
		t.Fatal("send/recv message not equal for output1")
	}
	d := <-outputC2
	if bytes.Compare(d, testBytes) != 0 {
		t.Fatal("send/recv message not equal for output2")
	}
}

func TestComposerWrongNodeType(t *testing.T) {
	gd := "[aaa] -> [bbb]"
	if _, err := composeIt("test_session", gd); err == nil {
		t.Fatal("should fail to prepare wrong type nodes")
	}
}

func ExampleComposerPubSub() {
	gd := `[input:chan_src] -> [pubsub] -> {[p1:print], [ps:pubsub]};
          [ps] -> {[p2:print], [p3:print],[output:chan_sink]}`
	c, err := composeIt("test_session", gd)
	if err != nil {
		panic(err)
	}

	inputC, outputC := make(chan []byte), make(chan []byte, 1)
	c.GetNode("input").(*comp.ChanSrc).LinkMe(inputC)
	c.GetNode("output").(*comp.ChanSink).LinkMe(outputC)
	inputC <- []byte("foobar")
	msg := <-outputC
	fmt.Printf("channel got %v\n", string(msg))
	time.Sleep(50 * time.Millisecond)

	// Unordered OUTPUT:
	// p1 print foobar
	// p2 print foobar
	// p3 print foobar
	// channel got foobar
}

func ExampleComposerMessageConvert() {
	gd := `[fire] -> [pubsub] -> [p1:print];[fire1:fire]`
	c, err := composeIt("test_session", gd)
	if err != nil {
		panic(err)
	}

	fire := c.GetNode("fire1").(*fireNode)
	fire.Call("fire", comp.With("call_action"))
	fire.Cast("fire", comp.With("cast_action"))

	time.Sleep(50 * time.Millisecond)

	// Unordered OUTPUT:
	// fire call from fire1
	// fire cast from fire1
	// p1 print call_action
	// p1 print cast_action
}

//
//func ExampleComposerController() {
//	gd := `[entry payload_type='1'] -> [p1:print] -> [p2:print]`
//	comp.RegisterNodeFactory("print", newPrintNode)
//	c := comp.NewSessionComposer("test_session")
//	graph := event.NewEventGraph()
//	if err := c.ParseGraphDescription(gd); err != nil {
//		fmt.Println("parse graph failed, ", err)
//		return
//	}
//	if err := c.ComposeNodes(graph); err != nil {
//		fmt.Println("prepare node failed, ", err)
//		return
//	}
//	ctrl := c.GetCommandInitiator()
//	ctrl.Cast("", "p1", []string{"foobar"})
//	reply := ctrl.Call("", "p1", []string{})
//	fmt.Printf("%v: %v\n", reply[0], reply[1])
//	time.Sleep(50 * time.Millisecond)
//
//	// Unordered OUTPUT:
//	// p2 print foobar
//	// ok: p1
//}
//
//func ExampleComposerInterSession() {
//	gd := `[entry payload_type='1'] -> [ps:pubsub] -> {[p1:print],[p2:print]}`
//	comp.RegisterNodeFactory("print", newPrintNode)
//	graph := event.NewEventGraph()
//	c1 := comp.NewSessionComposer("test1")
//	c2 := comp.NewSessionComposer("test2")
//	if c1.ParseGraphDescription(gd) != nil || c2.ParseGraphDescription(gd) != nil {
//		fmt.Println("parse graph failed")
//		return
//	}
//	if c1.ComposeNodes(graph) != nil || c2.ComposeNodes(graph) != nil {
//		fmt.Println("prepare node failed")
//		return
//	}
//	mp1, mp2 := c1.GetMessageProvider("entry"), c2.GetMessageProvider("entry")
//	ctrl1 := c1.GetCommandInitiator()
//	mp1.PushMessage(comp.NewRawByteMessage("from_session_1"))
//	mp2.PushMessage(comp.NewRawByteMessage("from_session_2"))
//	connCmd := comp.WithConnect("test1", "p2")
//	ctrl1.Call("test2", "ps", connCmd) // ask "test2:ps" to connect to "test1:p2"
//	mp2.PushMessage(comp.NewRawByteMessage("from_session_2_again"))
//
//	time.Sleep(50 * time.Millisecond)
//	// Unordered OUTPUT:
//	// p1 print from_session_1
//	// p2 print from_session_1
//	// p1 print from_session_2
//	// p2 print from_session_2
//	// p1 print from_session_2_again
//	// p2 print from_session_2_again
//	// p2 print from_session_2_again
//
//}
//
//func ExamplePubSubEnableDisable() {
//	gd := `[e1:entry payload_type='1'] -> [pubsub] -> {[p1:print],[p2:print]}`
//	comp.RegisterNodeFactory("print", newPrintNode)
//	c := comp.NewSessionComposer("test_session")
//	graph := event.NewEventGraph()
//	if err := c.ParseGraphDescription(gd); err != nil {
//		fmt.Println("parse graph failed")
//		return
//	}
//	if err := c.ComposeNodes(graph); err != nil {
//		fmt.Println("prepare node failed")
//		return
//	}
//	mp := c.GetMessageProvider("e1")
//	ctrl := c.GetCommandInitiator()
//	mp.PushMessage(comp.NewRawByteMessage("foo"))
//	ctrl.Call("", "pubsub", comp.With("disable", "node", "test_session", "p1"))
//	mp.PushMessage(comp.NewRawByteMessage("bar"))
//	ctrl.Call("", "pubsub", comp.With("enable", "node", "test_session", "p1"))
//	mp.PushMessage(comp.NewRawByteMessage("foobar"))
//	time.Sleep(50 * time.Millisecond)
//
//	// Unordered OUTPUT:
//	// p1 print foo
//	// p2 print foo
//	// p2 print bar
//	// p1 print foobar
//	// p2 print foobar
//}
//
//func ExampleComposerPushDataMessage() {
//	gd := `[p1:print]; [p2:print];`
//
//	comp.RegisterNodeFactory("print", newPrintNode)
//	c := comp.NewSessionComposer("test_session")
//	graph := event.NewEventGraph()
//	if err := c.ParseGraphDescription(gd); err != nil {
//		fmt.Println("parse graph failed")
//		return
//	}
//	if err := c.ComposeNodes(graph); err != nil {
//		fmt.Println("prepare node failed")
//		return
//	}
//	ctrl := c.GetCommandInitiator()
//	ctrl.PushData("p1", "", []byte("abc"))
//	ctrl.PushData("p2", "", []byte("cba"))
//	time.Sleep(50 * time.Millisecond)
//
//	// Unordered OUTPUT:
//	// p1 print abc
//	// p2 print cba
//}
