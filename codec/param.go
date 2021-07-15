package codec

import (
	"strconv"
	"strings"
)

// utility classes to config transcode parameters
// refer to:  https://ffmpeg.org/ffmpeg-filters.html  for detailed filter description

type configable interface {
	setName(name string)
	getName() string
	setOption(key string, value string)
	getOption() *kvoption
}

type kvoption map[string]string

type codecParam struct {
	codecName string
	option    kvoption
}

type filterParam struct {
	filterName string

	// some filter option comes first without key, like aresample=8000
	// to set immediateOption just set key to "" when configuring option
	immediateOption string
	option          kvoption
}

type TranscodeParam struct {
	encoder       *codecParam
	decoder       *codecParam
	currentConfig configable
	filterGraph   []*filterParam
	hasError      bool
}

func (cp *codecParam) setName(name string) {
	cp.codecName = name
}

func (cp *codecParam) getName() string {
	return cp.codecName
}

func (cp *codecParam) setOption(key string, value string) {
	cp.option[key] = value
}

func (cp *codecParam) getOption() *kvoption{
	return &cp.option
}

func (fp *filterParam) setName(name string) {
	fp.filterName = name
}

func (fp *filterParam) getName() string {
	return fp.filterName
}

func (fp *filterParam) setOption(key string, value string) {
	if key == "" {
		fp.immediateOption = value
	} else {
		fp.option[key] = value
	}
}

func (fp *filterParam) getOption() *kvoption{
	return &fp.option
}

func NewTranscodeParam() *TranscodeParam {
	return &TranscodeParam{
		encoder: &codecParam{option: make(kvoption)},
		decoder: &codecParam{option: make(kvoption)}}
}

func (tp *TranscodeParam) Encoder(name string) *TranscodeParam {
	tp.saveFilter()
	tp.encoder.setName(name)
	tp.currentConfig = tp.encoder
	return tp
}

func (tp *TranscodeParam) Decoder(name string) *TranscodeParam {
	tp.saveFilter()
	tp.decoder.setName(name)
	tp.currentConfig = tp.decoder
	return tp
}

func (tp *TranscodeParam) With(key string, value string) *TranscodeParam {
	tp.currentConfig.setOption(key, value)
	return tp
}

func (tp *TranscodeParam) saveFilter() {
	if tp.currentConfig == nil {
		return
	}
	if oldFilter, ok := tp.currentConfig.(*filterParam); ok {
		// save current filter
		if oldFilter.getName() == "" {
			// name is not set
			tp.hasError = true
		}
		tp.filterGraph = append(tp.filterGraph, oldFilter)
	}
	tp.currentConfig = nil
}

func (tp *TranscodeParam) NewFilter(name string) *TranscodeParam{
	tp.saveFilter()
	tp.currentConfig = &filterParam{
		filterName: name,
		option:     make(kvoption),
	}
	return tp
}

// combine all params together to a string which is used by transcode initilization
// final string will in form of:
// ****************************************************************
// encoder:some_encoder_name param1=value1,param2=value2,... \n
// decoder:some_decoder_name param1=value1,param2=value2,... \n
// filter_graph: filter1=<filter options>,filter2=<filter options>
// ****************************************************************
// NOTE: encoder/decoder lines are mandatory, filter_graph line is optional
// @return nil if the required parameters are not enough or error exist
func (tp *TranscodeParam) GetDescription() *string {
	if tp.hasError {
		return nil
	}
	encoderLen := len(tp.encoder.option)
	decoderLen := len(tp.decoder.option)
	if encoderLen == 0 ||
		decoderLen == 0 ||
		tp.encoder.getName() == "" ||
		tp.decoder.getName() == "" {
		return nil
	}
	tp.saveFilter()

	var sb strings.Builder
	writeCodec := func(typeName string,codec configable,optLen int) {
		sb.WriteString(typeName + ":")
		sb.WriteString(codec.getName())
		sb.WriteString(" ")
		i := 0
		for k,v := range *codec.getOption() {
			sb.WriteString(k)
			sb.WriteString("=")
			sb.WriteString(v)
			if i == optLen - 1 {
				sb.WriteString("\n")
			} else {
				sb.WriteString(",")
			}
			i++
		}
	}
	writeCodec("encoder",tp.encoder,encoderLen)
	writeCodec("decoder",tp.decoder,decoderLen)

	fgLen := len(tp.filterGraph)
	i := 0
	if fgLen > 0 {
		sb.WriteString("filter_graph: ")
		for _,filter := range tp.filterGraph {
			if filter.immediateOption == "" && len(filter.option) == 0 {
				// a filter without any value and option, just append its name
				sb.WriteString(filter.getName())
			} else {
				sb.WriteString(filter.getName() + "=")
				if len(filter.immediateOption) > 0 {
					sb.WriteString(filter.immediateOption)
					if len(filter.option) > 0 {
						sb.WriteString(":")        // filter options separated by colon, parsed by libavfilter
					}
				}
				for k,v := range filter.option {
					sb.WriteString(":")
					sb.WriteString(k)
					sb.WriteString("=")
					sb.WriteString(v)
				}
			}
			if i != fgLen - 1 {
				sb.WriteString(",") // filters are separated by comma, parsed by libavfilter
			}
			i++
		}
	}

	result := sb.String()
	return &result
}

// following methods are for convenience, used for codecParam only
func (tp *TranscodeParam) checkCodec() {
	if _,ok := tp.currentConfig.(*codecParam); !ok {
		tp.hasError = true
	}
}

// option names for AVCodecContext, see "libavcodec/options_table.h" in FFMPEG
func (tp *TranscodeParam) SampleRate(rate int) *TranscodeParam {
	tp.checkCodec()
	tp.currentConfig.setOption("ar",strconv.Itoa(rate))
	return tp
}

func (tp *TranscodeParam) ChannelCount(count int) *TranscodeParam{
	tp.checkCodec()
	tp.currentConfig.setOption("ac",strconv.Itoa(count))
	return tp
}

func (tp *TranscodeParam) BitRate(bitrate int) *TranscodeParam{
	tp.checkCodec()
	tp.currentConfig.setOption("b",strconv.Itoa(bitrate))
	return tp
}
