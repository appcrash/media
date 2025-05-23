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
	lp comp.LinkPoint `comp:"type=custom,nullable"`
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
	n.Self = n
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

type printHeaderNode struct {
	printNode
}

func (n *printHeaderNode) handleRawByteEvent(evt *event.Event) {
	msg, _ := comp.EventToMessage[comp.Message](evt)
	if v := msg.GetHeader(comp.Origin); v != nil {
		fmt.Printf("%v print origin %v", n.GetNodeName(), string(v))
	}
}

func newPrintNode() comp.SessionAware {
	p := &printNode{}
	p.Self = p
	p.Trait, _ = comp.NodeTraitOfType("print")
	p.SetMessageHandler(comp.MtRawByte, comp.ChainSetHandler(p.handleRawByteEvent))
	return p
}

func newPrintHeaderNode() comp.SessionAware {
	p := &printHeaderNode{}
	p.Self = p
	p.Trait, _ = comp.NodeTraitOfType("print_header")
	p.SetMessageHandler(comp.MtRawByte, comp.ChainSetHandler(p.handleRawByteEvent))
	return p
}

type fakeGateway struct {
	comp.SessionNode
}

func (n *fakeGateway) Accept() []comp.MessageType {
	return []comp.MessageType{
		comp.MtRawByte,
	}
}

func (n *fakeGateway) Offer() []comp.MessageType {
	return []comp.MessageType{comp.MtRawByte}
}

func (n *fakeGateway) handleRawByteEvent(evt *event.Event) {
}

func newFakeGatewayNode() comp.SessionAware {
	n := &fakeGateway{}
	n.Self = n
	n.Trait, _ = comp.NodeTraitOfType("fake_gateway")
	n.SetMessageHandler(comp.MtRawByte, comp.ChainSetHandler(n.handleRawByteEvent))
	return n
}

func initComposer() {
	comp.AddMessageTrait(comp.MT[customMessage](comp.MetaType[customMessageConvertable]()))
	comp.SetMessageConvertable(mtCustom, comp.MtRawByte)
	comp.RegisterNodeTrait(comp.NT[printNode]("print", newPrintNode))
	comp.RegisterNodeTrait(comp.NT[printHeaderNode]("print_header", newPrintHeaderNode))
	comp.RegisterNodeTrait(comp.NT[fireNode]("fire", newFireNode))
	comp.RegisterNodeTrait(comp.NT[fakeGateway]("fake_gateway", newFakeGatewayNode))
}

func composeIt(session, gd string) (*comp.Composer, error) {
	c := comp.NewSessionComposer(session, "")
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
	gd := `[input:chan_src] -> [pubsub] <raw_byte> {[output1:chan_sink],[output2:chan_sink]};`
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

func TestLoop(t *testing.T) {
	gd := `[input:chan_src] -> [abc:fake_gateway] ->[cba:fake_gateway] -> {[output1:chan_sink],[output2:chan_sink]};`
	_, err := composeIt("test_session", gd)
	if err != nil {
		t.Fatal(err)
	}
}

func TestComposerWrongNodeType(t *testing.T) {
	gd := "[aaa] -> [bbb]"
	if _, err := composeIt("test_session", gd); err == nil {
		t.Fatal("should fail to prepare wrong type nodes")
	}
}

func TestChanSrcSink(t *testing.T) {
	gd := "[src:chan_src] -> [sink:chan_sink]"
	if c, e := composeIt("test_session", gd); e != nil {
		t.Fatal(e)
	} else {
		inputC, outputC := make(chan []byte), make(chan []byte)
		i := c.GetNode("src").(*comp.ChanSrc)
		o := c.GetNode("sink").(*comp.ChanSink)
		if err := i.LinkMe(inputC); err != nil {
			t.Fatal(err)
		}
		if err := i.LinkMe(inputC); err == nil {
			t.Fatal("should not link again")
		}
		if err := o.LinkMe(outputC); err != nil {
			t.Fatal(err)
		}
		if err := o.LinkMe(outputC); err == nil {
			t.Fatal("should not link again")
		}
	}
}

func Example_messagePostProcessor() {
	gd := `[input:chan_src trackable=true] -> [pubsub] -> [p1:print_header];`
	c, err := composeIt("test_session", gd)
	if err != nil {
		panic(err)
	}
	inputC := make(chan []byte)
	c.GetNode("input").(*comp.ChanSrc).LinkMe(inputC)
	inputC <- []byte("abcdefg")
	time.Sleep(50 * time.Millisecond)

	// Unordered OUTPUT:
	// p1 print origin input
}

func Example_composerPubSub() {
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

func Example_composerMessageConvert() {
	gd := `[fire] -> [pubsub] -> [p1:print];[fire1:fire]`
	c, err := composeIt("test_session", gd)
	if err != nil {
		panic(err)
	}

	if c.GetNode("fire").(*fireNode).lp == nil {
		panic("fire node link point not injected")
	}
	fire1 := c.GetNode("fire1").(*fireNode)
	fire1.Call("fire", comp.With("call_action"))
	fire1.Cast("fire", comp.With("cast_action"))

	time.Sleep(50 * time.Millisecond)

	// Unordered OUTPUT:
	// fire call from fire1
	// fire cast from fire1
	// p1 print call_action
	// p1 print cast_action
}

func Example_composerBuiltInCommand() {
	gd1 := `[fire] -> [pubsub] -> [p1:print];`
	gd2 := `[src:chan_src] -> [pubsub] -> {[p2:print],[output:chan_sink]};`
	c1 := comp.NewSessionComposer("session1", "")
	c2 := comp.NewSessionComposer("session2", "")
	graph := event.NewEventGraph()
	if err := c1.ParseGraphDescription(gd1); err != nil {
		panic(err)
	}
	if err := c2.ParseGraphDescription(gd2); err != nil {
		panic(err)
	}
	if err := c1.ComposeNodes(graph); err != nil {
		panic(err)
	}
	if err := c2.ComposeNodes(graph); err != nil {
		panic(err)
	}
	inputC := make(chan []byte)
	outputC := make(chan []byte, 2)
	c2.GetNode("src").(*comp.ChanSrc).LinkMe(inputC)
	c2.GetNode("output").(*comp.ChanSink).LinkMe(outputC)
	initiator1 := c1.GetCommandInitiator()
	initiator2 := c2.GetCommandInitiator()

	inputC <- []byte("src before conn")
	<-outputC // wait message passwd through pubsub
	args, _ := comp.WithString("conn session1 p1")
	resp := initiator2.Call("", "pubsub", args)
	inputC <- []byte("src after conn")
	<-outputC

	args, _ = comp.WithString("disable_link " + resp[1])
	initiator2.Call("", "pubsub", args)
	inputC <- []byte("src after disable_link")
	<-outputC

	args, _ = comp.WithString("enable_link " + resp[1])
	initiator2.Call("", "pubsub", args)
	inputC <- []byte("src after enable_link")
	<-outputC

	args, _ = comp.WithString("conn session2 p2")
	initiator1.Call("", "pubsub", args)
	initiator1.Cast("rpc", "fire", []string{"fire in the hole"})

	time.Sleep(500 * time.Millisecond)
	c1.ExitGraph()
	c2.ExitGraph()

	// Unordered OUTPUT:
	// p2 print src before conn
	// p1 print src after conn
	// p2 print src after conn
	// p2 print src after disable_link
	// p1 print src after enable_link
	// p2 print src after enable_link
	// fire cast from rpc
	// p1 print fire in the hole
	// p2 print fire in the hole
}
