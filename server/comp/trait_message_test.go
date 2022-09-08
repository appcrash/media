package comp_test

import (
	"github.com/appcrash/media/server/comp"
	"testing"
)

func TestMessageTrait(t *testing.T) {
	cm := &comp.ChannelLinkMessage{}
	trait, ok := comp.MessageTraitOfObject(&comp.RawByteMessage{})
	if !ok {
		t.Fatal("not found")
	}
	converted, err := trait.ConvertFrom(cm)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("converted is %v\n", converted)
}
