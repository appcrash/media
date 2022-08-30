package comp_test

import (
	"github.com/appcrash/media/server/comp"
	"testing"
)

func TestRegisterNodeTrait(t *testing.T) {
	err := comp.RegisterNodeTrait(comp.NT[comp.ChannelSink]())
	if err != nil {
		t.Errorf("register failed %v", err)
	}
}
