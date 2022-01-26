package comp_test

import (
	"github.com/appcrash/media/server/comp"
	"testing"
)

func TestGraphSorting(t *testing.T) {
	desc := `[sink_entry] -> [transcode] -> [asr] -> [pubsub];
		[transcode from_name=amrwb from_samplerate='16k' to_name=pcm_s16le to_samplerate='8k']; 
		[pubsub channel=source_exit]`

	sc := comp.NewSessionComposer("test")
	err := sc.ParseGraphDescription(desc)
	if err != nil {
		t.Fatal("should not happen")
	}
	expected := []string{"pubsub", "asr", "transcode", "sink_entry"}
	for i, n := range sc.GetSortedNodes() {
		if expected[i] != n.Name {
			t.Fatal("wrong order of sorted node")
		}
	}
}

func TestGraphLoop(t *testing.T) {
	desc := `[a] -> [b] -> [c] -> [e] -> [a];
             [d] -> [c]`
	sc := comp.NewSessionComposer("test")
	err := sc.ParseGraphDescription(desc)
	//t.Logf("correctly detect: %v", err)
	if err == nil {
		t.Fatal("loop not detected")
	}
}
