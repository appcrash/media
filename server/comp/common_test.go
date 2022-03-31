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

func TestGenericMessage_Clone(t *testing.T) {
	rbm := comp.RawByteMessage("some_byte")
	gm := &comp.GenericMessage{
		Subtype: "test_generic",
		Obj:     rbm,
	}
	cgm := gm.Clone()
	if cgm == nil {
		t.Fatal("generic message does not clone its internal object")
	}
	gm.Obj = gm
	cgm = gm.Clone()
	if cgm != nil {
		t.Fatal("generic message allow recursive clone")
	}
	gm.Obj = nil
	cgm = gm.Clone()
	if cgm != nil {
		t.Fatal("generic message cloned when object is nil")
	}
}
