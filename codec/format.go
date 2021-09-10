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

type (
	RecordContext = C.struct_RecordContext
)

func GetPayloadFromFile(fp string) []byte {
	cfp := C.CString(fp)
	payload := (*C.struct_Payload)(C.read_media_file(cfp))
	if payload == nil {
		return nil
	}
	C.free(unsafe.Pointer(cfp))

	logger.Debugf("data is size:%v,bitrate:%v,frame_size:%v\n", payload.size, payload.bitrate,
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

func NewRecordContext(fileName, params string) *RecordContext {
	if params == "" || fileName == "" {
		return nil
	}
	p := C.CString(params)
	fp := C.CString(fileName)
	defer C.free(unsafe.Pointer(p))
	defer C.free(unsafe.Pointer(fp))
	ctx := C.record_init_context(fp, p)
	if ctx != nil {
		return (*RecordContext)(ctx)
	} else {
		return nil
	}
}

// Iterate batch write frames
func (ctx *RecordContext) Iterate(frames [][]byte) {
	var frameDelimits []int32
	var data []byte
	var i int
	for _, frame := range frames {
		l := len(frame)
		if l == 0 {
			continue
		}
		frameDelimits = append(frameDelimits, int32(i+l))
		data = append(data, frame...)
		i += l
	}
	if len(frameDelimits) == 0 {
		return
	}
	C.record_iterate(ctx, (*C.char)(unsafe.Pointer(&data[0])),
		(*C.int)(unsafe.Pointer(&frameDelimits[0])), C.int(len(frameDelimits)))
}

func (ctx *RecordContext) Free() {
	C.record_free(ctx)
}
