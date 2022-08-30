package comp_test

import (
	"github.com/appcrash/media/server/comp"
	"testing"
)

type testNonCloneableMessage struct {
	comp.BaseMessage
}

func (t *testNonCloneableMessage) Type() comp.MessageType {
	return 100
}

func TestMessageTrait(t *testing.T) {
	comp.RegisterMessageTrait(
		comp.MT[comp.RawByteMessage](),
		comp.MT[testNonCloneableMessage]())
	trait, ok := comp.MessageTraitOfObject(&comp.RawByteMessage{})
	if !ok || !trait.IsCloneable() {
		t.Fatal("raw byte message must be cloneable")
	}
	trait, ok = comp.MessageTraitOfObject(&testNonCloneableMessage{})
	if !ok || trait.IsCloneable() {
		t.Fatal("should not be cloneable")
	}

}

func TestRegisterNodeTrait(t *testing.T) {
	err := comp.RegisterNodeTrait(comp.NT[comp.ChanSink]())
	if err != nil {
		t.Errorf("register failed %v", err)
	}
}
