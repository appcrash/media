package comp_test

import (
	"github.com/appcrash/media/server/comp"
	"testing"
)

func TestGraphSorting(t *testing.T) {
	desc := "[sink_entry] -> [transcode] \n" +
		"[transcode] -> [lxyasr] \n" +
		"[lxyasr] -> [pubsub] \n" +
		"[transcode]: from_name=amrwb;from_samplerate=16k;to_name=pcm_s16le;to_samplerate=8k \n" +
		"[pubsub]: channel=source_exit"

	sc := comp.NewSessionComposer("test")
	err := sc.ParseGraphDescription(desc)
	if err != nil {
		t.Fatal("should not happen")
	}
	expected := []string{"pubsub", "lxyasr", "transcode", "sink_entry"}
	for i, n := range sc.GetSortedNodes() {
		if expected[i] != n.Name {
			t.Fatal("wrong order of sorted node")
		}
	}
}

func TestGraphLoop(t *testing.T) {
	desc := "[a] -> [b] \n" +
		"[b] -> [c] \n" +
		"[d] -> [c] \n" +
		"[c] -> [e] \n" +
		"[e] -> [a] \n"
	sc := comp.NewSessionComposer("test")
	err := sc.ParseGraphDescription(desc)
	if err == nil {
		t.Fatal("loop not detected")
	}
}
