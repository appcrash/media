package codec_test

import (
	"fmt"
	"github.com/appcrash/media/codec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestRecord(t *testing.T) {
	_, srcFileName, _, _ := runtime.Caller(0)
	sampleFile := filepath.Join(filepath.Dir(srcFileName), "../assets/sample.wav")
	payload := codec.GetPayloadFromFile(sampleFile)
	if payload == nil {
		t.Fatal("cannot get payload of test file")
	}
	fileName := "/tmp/recording.wav"
	params := fmt.Sprintf("channels=1,sample_rate=8000,codec_id=%v", codec.AV_CODEC_ID_PCM_ALAW)
	ctx := codec.NewRecordContext(fileName, params)
	if ctx == nil {
		t.Fatal("failed to create record context")
	}
	frames := codec.PcmaSplitToFrames(payload, 20)
	ctx.Iterate(frames)
	ctx.Free()
}
