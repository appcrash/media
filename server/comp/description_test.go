package comp_test

import (
	"github.com/appcrash/media/server/comp"
	"testing"
)

func TestGraphSorting(t *testing.T) {
	desc := `[a_gateway] -> [sink_entry] -> [transcode] -> [asr] -> [pubsub];
		[transcode from_name=amrwb from_samplerate='16k' to_name=pcm_s16le to_samplerate='8k'];`

	sc := comp.NewSessionComposer("test", "")
	err := sc.ParseGraphDescription(desc)
	if err != nil {
		t.Fatal("should not happen")
	}
	expected := []string{"pubsub", "asr", "transcode", "sink_entry", "a_gateway"}
	for i, n := range sc.GetSortedNodes() {
		if expected[i] != n.Name {
			t.Fatal("wrong order of sorted node")
		}
	}
}

func TestGraphLoop(t *testing.T) {
	desc := `[a] -> [b] -> [c] -> [e] -> [a];
             [d] -> [c]`
	sc := comp.NewSessionComposer("test", "")
	err := sc.ParseGraphDescription(desc)
	//t.Logf("correctly detect: %v", err)
	if err == nil {
		t.Fatal("loop not detected")
	}
}

func TestAllowGraphLoop(t *testing.T) {
	desc1 := `[a_gateway] -> [b] -> [c] -> [e] -> [a_gateway];
             [d] -> [c]`
	sc := comp.NewSessionComposer("test", "")
	err := sc.ParseGraphDescription(desc1)
	if err != nil {
		t.Fatal("loop with gateway should be valid")
	}
	desc2 := `[a_gateway] -> [b_gateway];[b_gateway] -> [a_gateway]`
	sc = comp.NewSessionComposer("test", "")
	err = sc.ParseGraphDescription(desc2)
	if err != nil {
		t.Fatal("gateway nodes should talk to each other")
	}
}
