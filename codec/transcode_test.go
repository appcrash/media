package codec

import (
	"encoding/binary"
	"math"
	"testing"
	"fmt"
)

const (
	encoder = 1
	decoder = 2
)

type codecSpec struct {
	name         string
	sampleRate   []int
	channelCount int
	bitrate      []int
	capability   int
}

type codecConfig struct {
	name string
	sampleRate int
	channelCount int
	bitrate int
}

type codecTask struct {
	from codecConfig
	to codecConfig
}

var codecDb = []codecSpec{
	{"pcm_f64le", []int{8000}, 1, nil, encoder | decoder},
	{"pcm_s16le", []int{8000,16000}, 1, nil, encoder | decoder},
	{"pcm_alaw", []int{8000}, 1, nil, encoder | decoder},
	{"amrnb", []int{8000}, 1, nil, decoder},
	{"amrwb", []int{16000}, 1, nil, decoder},
	{"libopencore_amrnb", []int{8000}, 1, []int{4750, 5150, 5900, 6700, 7400, 7950, 10200, 12200}, encoder | decoder},
	{"libopencore_amrwb", []int{16000}, 1, nil, decoder},
	{"libvo_amrwbenc", []int{16000}, 1, []int{6600, 8850, 12650, 14250, 15850, 18250, 19850, 23050, 23850}, encoder},
}


func (c codecConfig) String() string {
	str := fmt.Sprintf("[%v]:ar(%v):ac(%v)",c.name,c.sampleRate,c.channelCount)
	if c.bitrate != 0 {
		str = str + fmt.Sprintf(":b(%v)",c.bitrate)
	}
	return str
}

func specToConfig(s *codecSpec) (c []codecConfig) {
	for _,rate := range s.sampleRate {
		config := &codecConfig {
			name : s.name,
			sampleRate : rate,
			channelCount : s.channelCount,
		}
		if s.bitrate != nil {
			for _,b := range s.bitrate {
				cb := *config
				cb.bitrate = b
				c = append(c,cb)
			}
		} else {
			c = append(c,*config)
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
	samples := generateSample(1200.0, 10000, 8000)
	var decoders []codecConfig
	var encoders []codecConfig
	for _,spec := range codecDb {
		config := specToConfig(&spec)
		if spec.capability & encoder != 0 {
			encoders = append(encoders,config...)
		}
		if spec.capability & decoder != 0 {
			decoders = append(decoders,config...)
		}
	}
	// convert f64le samples to all formats that encoders support
	for _,config := range encoders {
		param := NewTranscodeParam().
		Decoder("pcm_f64le").SampleRate(8000).ChannelCount(1).
		Encoder(config.name).SampleRate(config.sampleRate).ChannelCount(config.channelCount)
		if config.bitrate != 0 {
			param = param.BitRate(config.bitrate)
		}
		param = param.NewFilter("aresample")
		
		ctx := NewTranscodeContext(param)
		if ctx == nil {
			t.Fatal("create transcode context failed")
		}
		_,ok := ctx.Iterate(samples)
		if ok != 0 {
			t.Fatal("convert failed")
		}
		fmt.Printf("f16le => %v\n",config)
		ctx.Free()
	}
	
}
