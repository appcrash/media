package comp_test

import (
	"github.com/appcrash/media/server/comp"
	"github.com/sirupsen/logrus"
	"testing"
)

func init() {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	//logger.SetLevel(logrus.DebugLevel)
	comp.InitLogger(logger)
}

type testNonCloneableMessage struct {
	comp.BaseMessage
}

func (t *testNonCloneableMessage) Type() comp.MessageType {
	return 100
}

func TestMessageTrait(t *testing.T) {
	comp.RegisterMessageTrait(&comp.RawByteMessage{}, &testNonCloneableMessage{})
	trait, ok := comp.MessageTraitOfObject(&comp.RawByteMessage{})
	if !ok || !trait.IsCloneable() {
		t.Fatal("raw byte message must be cloneable")
	}
	trait, ok = comp.MessageTraitOfObject(&testNonCloneableMessage{})
	if !ok || trait.IsCloneable() {
		t.Fatal("should not be cloneable")
	}

}
