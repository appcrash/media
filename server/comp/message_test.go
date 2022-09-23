package comp_test

import (
	"github.com/appcrash/media/server/comp"
	"github.com/appcrash/media/server/utils"
	"testing"
)

//func TestMessageTrait(t *testing.T) {
//	cm := &comp.ChannelLinkRequestMessage{}
//	trait, ok := comp.MessageTraitOfObject(&comp.RawByteMessage{})
//	if !ok {
//		t.Fatal("not found")
//	}
//	converted, err := trait.ConvertFrom(cm)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Logf("converted is %v\n", converted)
//}

func TestNodeTrait(t *testing.T) {
	node := comp.MakeSessionNode("rtp_src", "abc", nil)
	composer := comp.NewSessionComposer("abc")
	metaType := comp.MetaType[comp.PreComposer]()

	utils.AopCall(node, nil, comp.MetaType[comp.PreInitializer](), "PreInit")
	utils.AopCall(node, []interface{}{composer, node}, metaType, "BeforeCompose")
}
