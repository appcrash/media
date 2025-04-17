package codec_test

import (
	"github.com/appcrash/media/codec"
	"math/rand/v2"
	"testing"
)

func TestMix(t *testing.T) {
	testWithSampleRate(t, 8000, 8000, 8000)
	testWithSampleRate(t, 8000, 8000, 16000)
	testWithSampleRate(t, 16000, 16000, 8000)

	// TODO: support different sample rates
	//testWithSampleRate(t,8000,16000,16000)
	//testWithSampleRate(t,16000,16000,8000)

}

func testWithSampleRate(t *testing.T, rate1 int, rate2 int, outRate int) {
	hz1 := rand.IntN(18_000) + 440
	hz2 := rand.IntN(18_000) + 440
	sampleNum := rand.IntN(1000) + 10000
	samples1 := generateSample(float64(hz1), sampleNum, rate1)
	samples2 := generateSample(float64(hz2), sampleNum, rate2)

	mixParam := codec.NewMixParam().
		Input1().SampleRate(rate1).SampleFormat(codec.AV_SAMPLE_FMT_DBL).ChannelLayout(codec.AV_CH_LAYOUT_MONO).
		Input2().SampleRate(rate2).SampleFormat(codec.AV_SAMPLE_FMT_DBL).ChannelLayout(codec.AV_CH_LAYOUT_MONO).
		Output().SampleRate(outRate).SampleFormat(codec.AV_SAMPLE_FMT_S16).ChannelLayout(codec.AV_CH_LAYOUT_MONO).
		GetDescription()

	// sample1 and sample2 passed to iterate must have equal audio time length
	sampleLen1, sampleLen2 := sampleNum, sampleNum
	//if rate1 > rate2 {
	//	sampleLen1 = int(math.Floor(float64(sampleNum) / float64(rate1) * float64(rate2)))
	//} else if rate2 > rate1 {
	//	sampleLen2 = int(math.Floor(float64(sampleNum) / float64(rate2) * float64(rate1)))
	//	sampleLen2 -= 2000
	//}
	mixctx := codec.NewMixContext(*mixParam)
	mixed, ok := mixctx.Iterate(samples1[:sampleLen1*8], samples2[:sampleLen2*8], sampleLen1, sampleLen2)
	if ok != 0 || len(mixed) == 0 {
		t.Fatalf("mix failed from rate: (%v,%v) to rate: %v, with samples: %v", rate1, rate2, outRate, sampleNum)
	}
	mixctx.Free()
}
