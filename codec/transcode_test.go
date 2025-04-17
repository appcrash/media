package codec_test

import (
	"encoding/binary"
	"fmt"
	"github.com/appcrash/media/codec"
	"math"
	"math/rand/v2"
	"testing"
)

const (
	encoder = 1
	decoder = 2
)

type codecSpec struct {
	name         string
	codecId      int
	sampleRate   []int
	channelCount int
	bitrate      []int
	capability   int
}

type codecConfig struct {
	name         string
	codecId      int
	sampleRate   int
	channelCount int
	bitrate      int
}

var codecDb = []codecSpec{
	{"pcm_f64le", codec.AV_CODEC_ID_PCM_F64LE,
		[]int{8000}, 1, nil, encoder | decoder},
	{"pcm_s16le", codec.AV_CODEC_ID_PCM_S16LE,
		[]int{8000, 16000}, 1, nil, encoder | decoder},
	{"pcm_alaw", codec.AV_CODEC_ID_PCM_ALAW,
		[]int{8000}, 1, nil, encoder | decoder},
	{"amrnb", codec.AV_CODEC_ID_AMR_NB,
		[]int{8000}, 1, nil, decoder},
	{"amrwb", codec.AV_CODEC_ID_AMR_WB,
		[]int{16000}, 1, nil, decoder},
	{"libopencore_amrnb", codec.AV_CODEC_ID_AMR_NB,
		[]int{8000}, 1, []int{4750, 5150, 5900, 6700, 7400, 7950, 10200, 12200}, encoder | decoder},
	{"libopencore_amrwb", codec.AV_CODEC_ID_AMR_WB,
		[]int{16000}, 1, nil, decoder},
	{"libvo_amrwbenc", codec.AV_CODEC_ID_AMR_WB,
		[]int{16000}, 1, []int{6600, 8850, 12650, 14250, 15850, 18250, 19850, 23050, 23850}, encoder},
}

func (c codecConfig) String() string {
	str := fmt.Sprintf("[%v]:ar(%v):ac(%v)", c.name, c.sampleRate, c.channelCount)
	if c.bitrate != 0 {
		str = str + fmt.Sprintf(":b(%v)", c.bitrate)
	}
	return str
}

func specToConfig(s *codecSpec) (c []codecConfig) {
	for _, rate := range s.sampleRate {
		config := &codecConfig{
			name:         s.name,
			codecId:      s.codecId,
			sampleRate:   rate,
			channelCount: s.channelCount,
		}
		if s.bitrate != nil {
			for _, b := range s.bitrate {
				cb := *config
				cb.bitrate = b
				c = append(c, cb)
			}
		} else {
			c = append(c, *config)
		}
	}
	return
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

func TestTranscode(t *testing.T) {
	hz := rand.IntN(18_000) + 440
	sampleNum := rand.IntN(10_000) + 10_000
	samples := generateSample(float64(hz), sampleNum, 8000)
	t.Logf("sample Hz: %v,  number is %v", hz, sampleNum)
	var decoders []codecConfig
	var encoders []codecConfig
	for _, spec := range codecDb {
		config := specToConfig(&spec)
		if spec.capability&encoder != 0 {
			encoders = append(encoders, config...)
		}
		if spec.capability&decoder != 0 {
			decoders = append(decoders, config...)
		}
	}
	// convert f64le samples to all formats that encoders support
	payloadMap := make(map[int][]byte)
	for _, config := range encoders {
		param := codec.NewTranscodeParam().
			Decoder("pcm_f64le").SampleRate(8000).ChannelCount(1).
			Encoder(config.name).SampleRate(config.sampleRate).ChannelCount(config.channelCount)
		if config.bitrate != 0 {
			param = param.BitRate(config.bitrate)
		}
		param = param.NewFilter("aresample")

		ctx := codec.NewTranscodeContext(param)
		if ctx == nil {
			t.Fatal("create transcode context failed")
		}
		payload, ok := ctx.Iterate(samples)
		if ok != 0 {
			t.Fatal("convert failed")
		}
		payloadMap[config.codecId] = payload
		t.Logf("samples:f16le => %v\n", config)
		ctx.Free()
	}

	// encoded kinds of payloads now, then transcode them to all other codecs
	for codecId, payload := range payloadMap {
		// search all decoders that can decode this payload
		for _, c := range codecDb {
			if c.codecId == codecId && (c.capability&decoder != 0) {
				for _, decodeConfig := range specToConfig(&c) {
					// 1:many transcode
					for _, encodec := range codecDb {
						if encodec.capability&encoder == 0 {
							continue
						}
						for _, encodeConfig := range specToConfig(&encodec) {
							param := codec.NewTranscodeParam().
								Decoder(decodeConfig.name).SampleRate(decodeConfig.sampleRate).ChannelCount(decodeConfig.channelCount).
								Encoder(encodeConfig.name).SampleRate(encodeConfig.sampleRate).ChannelCount(encodeConfig.channelCount)
							if encodeConfig.bitrate != 0 {
								param = param.BitRate(encodeConfig.bitrate)
							}
							param = param.NewFilter("aresample")

							t.Logf("transcode:%v => %v", decodeConfig, encodeConfig)
							ctx := codec.NewTranscodeContext(param)
							if ctx == nil {
								t.Fatal("create transcode context failed")
							}
							_, ok := ctx.Iterate(payload)
							if ok != 0 {
								t.Fatal("convert failed")
							}
							// drain remaining data
							_, ok = ctx.Iterate(nil)
							if ok != 0 {
								t.Fatal("convert with last draining failed")
							}
							ctx.Free()
						}
					}
				}
			}
		}
	}

}
