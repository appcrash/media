package comp_test

import (
	"github.com/appcrash/media/server/comp"
	"testing"
)

func TestRegisterNodeTrait(t *testing.T) {
	err := comp.RegisterNodeTrait(comp.NT[comp.ChanSink]())
	if err != nil {
		t.Errorf("register node failed %v", err)
	}
}
