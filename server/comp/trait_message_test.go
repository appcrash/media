package comp_test

import (
	"github.com/appcrash/media/server/comp"
	"testing"
)

type testNonCloneableMessage struct {
	comp.MessageBase
}

func (t *testNonCloneableMessage) Type() comp.MessageType {
	return 100
}

func TestMessageTrait(t *testing.T) {
	comp.InitBuiltinMessage()
	comp.RegisterMessageTrait(comp.MT[testNonCloneableMessage]())
	trait, ok := comp.MessageTraitOfObject(&comp.RawByteMessage{})
	if !ok || !trait.IsCloneable() {
		t.Fatal("raw byte message must be cloneable")
	}
	trait, ok = comp.MessageTraitOfObject(&testNonCloneableMessage{})
	if !ok || trait.IsCloneable() {
		t.Fatal("should not be cloneable")
	}

}
