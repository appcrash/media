package codec

//#cgo pkg-config: libavformat libavcodec libavutil libswresample libavfilter
//
//#include <stdio.h>
//#include <stdlib.h>
//#include <inttypes.h>
//#include <stdint.h>
//#include <string.h>
//#include <libavformat/avformat.h>
//#include "codec.h"
import "C"
import (
	"unsafe"
)

type DataBuffer struct {
	cobj *C.struct_DataBuffer
}

type TranscodeContext struct {
	cobj *C.struct_TranscodeContext
}

type MixContext struct {
	cobj *C.struct_MixContext
}

func NewTranscodeContext(param *TranscodeParam) *TranscodeContext {
	desc := param.GetDescription()
	if desc == nil {
		return nil
	}
	p := C.CString(*desc)
	defer C.free(unsafe.Pointer(p))
	cobj := C.transcode_init_context(p, C.int(len(*desc)))
	if cobj == nil {
		return nil
	}
	return &TranscodeContext{cobj}
}

// Iterate
// @param data
// the audio data of source codec, set to nil to get the remaining transcoded data
// it should be of one frame length (normally 20ms)
//
// this method can be called multiple times as long as more data needs transcoding
// you can pass n_frame length data at a time for better performance, but some decoders
// that does not support AV_CODEC_CAP_SUBFRAMES would complain with warning:
// "Multiple frames in a packet"
//
// the data length should always be aligned with decoder's frame size for best compatibility
// and in reasonable size:
// usually the duration of data should not exceed 1s beyond which evident lags would occur
// @return transcodedData  the encoded data of destination codec
func (context *TranscodeContext) Iterate(data []byte) (transcodedData []byte, reason int) {
	var dataLen int
	var pdata *byte
	if data != nil {
		pdata = &data[0]
		dataLen = len(data)
	}
	C.transcode_iterate(context.cobj, (*C.char)(unsafe.Pointer(pdata)),
		C.int(dataLen), (*C.int)(unsafe.Pointer(&reason)))
	buffer := context.cobj.out_buffer
	if buffer.size > 0 {
		transcodedData = C.GoBytes(unsafe.Pointer(buffer.data), buffer.size)
	}
	return
}

func (context *TranscodeContext) Free() {
	C.transcode_free(context.cobj)
}

func NewMixContext(param string) *MixContext {
	p := C.CString(param)
	defer C.free(unsafe.Pointer(p))
	cobj := C.mix_init_context(p, C.int(len(param)))
	if cobj == nil {
		return nil
	}
	return &MixContext{cobj}
}

func (context *MixContext) Iterate(data1 []byte, data2 []byte, samples1 int, samples2 int) (mixed []byte, reason int) {
	C.mix_iterate(context.cobj, (*C.char)(unsafe.Pointer(&data1[0])), C.int(len(data1)),
		(*C.char)(unsafe.Pointer(&data2[0])), C.int(len(data2)),
		C.int(samples1), C.int(samples2), (*C.int)(unsafe.Pointer(&reason)))
	buffer := context.cobj.out_buffer
	if buffer.size > 0 {
		mixed = C.GoBytes(unsafe.Pointer(buffer.data), buffer.size)
	}
	return
}

func (context *MixContext) Free() {
	C.mix_free(context.cobj)
}
