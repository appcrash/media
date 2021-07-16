package codec

import (
	"encoding/binary"
	"math"
	"testing"
)

const (
	encoder = 1
	decoder = 2
)

type codecSpec struct {
	name         string
	sampleRate   int
	channelCount int
	bitrate      []int
	capability   int
}

var codecDb = []codecSpec{
	{"pcm_f64le", 8000, 1, nil, encoder | decoder},
	{"pcm_s16le", 8000, 1, nil, encoder | decoder},
	{"pcm_alaw", 8000, 1, nil, encoder | decoder},
	{"amrnb", 8000, 1, nil, decoder},
	{"amrwb", 16000, 1, nil, decoder},
	{"libopencore_amrnb", 8000, 1, []int{4750, 5150, 5900, 6700, 7400, 7950, 10200, 12200}, encoder | decoder},
	{"libopencore_amrwb", 16000, 1, nil, decoder},
	{"libvo_amrwbenc", 16000, 1, []int{6600, 8850, 12650, 14250, 15850, 18250, 19850, 23050, 23850}, encoder},
}

func generateSample(hz float64, sampleNum int, sampleRate int) (s []byte) {
	s = make([]byte, 8*sampleNum)
	w := 2 * math.Pi * hz
	var step float64 = 1.0 / float64(sampleRate)
	for i := 0; i < sampleNum; i++ {
		v := math.Sin(w * step * float64(i))
		binary.LittleEndian.PutUint64(s[i*8:(i+1)*8], math.Float64bits(v))
	}
	return
}


//func TestTranscode(t *testing.T) {
//	samples := generateSample(1200.0, 10000, 8000)
//	param := NewTranscodeParam().
//		Decoder("pcm_f64le").SampleRate(8000).ChannelCount(1).
//		Encoder("pcm_s16le").SampleRate(8000).ChannelCount(1).
//		NewFilter("aresample")
//	ctx := NewTranscodeContext(param)
//	if ctx == nil {
//		t.Fatal("create transcode context failed")
//	}
//
//	ctx.Free()
//}
