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
	"fmt"
	"unsafe"
)

type (
	DataBuffer       = C.struct_DataBuffer
	TranscodeContext = C.struct_TranscodeContext
)

func GetPayloadFromFile(fp string) []byte {
	cfp := C.CString(fp)
	payload := (*C.struct_Payload)(C.read_media_file(cfp))
	C.free(unsafe.Pointer(cfp))

	fmt.Printf("data is size:%v,bitrate:%v,frame_size:%v\n", payload.size, payload.bitrate,
		BitrateToFrameSize(float64(payload.bitrate), 20))
	if payload != nil {
		data := C.GoBytes(unsafe.Pointer(payload.data), payload.size)
		C.free(unsafe.Pointer(payload.data))
		C.free(unsafe.Pointer(payload))
		return data
	}
	return nil
}

func WritePayloadToFile(payload []byte, fileName string, codecId int, duration int) (ret int) {
	cfileName := C.CString(fileName)
	v := C.write_media_file((*C.char)(unsafe.Pointer(&payload[0])), C.int(len(payload)), (*C.char)(cfileName), C.int(codecId), C.int(duration))
	ret = int(v)
	C.free(unsafe.Pointer(cfileName))
	return
}

func NewTranscodeContext(param *TranscodeParam) *TranscodeContext {
	desc := param.GetDescription()
	if desc == nil {
		return nil
	}
	p := C.CString(*desc)
	defer C.free(unsafe.Pointer(p))
	return (*TranscodeContext)(C.transcode_init_context(p,C.int(len(*desc))))
}

// @param data
// the audio data of source codec, set to nil to get the remaining transcoded data
// it should be of one frame length (normally 20ms)
//
// this method can be called multiple times as long as more data needs transcoding
// you can pass n_frame length data at a time for better performance, but some decoders
// that does not support AV_CODEC_CAP_SUBFRAMES would complain with warning:
// "Multiple frames in a packet"
//
// the data length should always be aligned with deocder's frame size for best compatibility
// and in reasonable size:
// ususally the duration of data should not exceeds 1s beyond which evident lags would occur
// @return transcodedData  the encoded data of destination codec
func (context *TranscodeContext) Iterate(data []byte) (transcodedData []byte, reason int) {
	var dataLen int
	var pdata *byte
	if data != nil {
		pdata = &data[0]
		dataLen = len(data)
	}
	C.transcode_iterate(context, (*C.char)(unsafe.Pointer(pdata)),
		C.int(dataLen), (*C.int)(unsafe.Pointer(&reason)))
	buffer := context.out_buffer
	if buffer.size > 0 {
		transcodedData = C.GoBytes(unsafe.Pointer(buffer.data), buffer.size)
	}
	return
}

func (context *TranscodeContext) Free() {
	C.transcode_free(context)
}
