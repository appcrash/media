package codec

//#cgo pkg-config: libavformat libavcodec libavutil
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
	DecodedFrame     = C.struct_DecodedFrame
	TranscodeContext = C.struct_TranscodeContext
)

func ConvertFormat(payload []byte) []byte {
	cpayload := C.CBytes(payload)
	defer C.free(unsafe.Pointer(cpayload))
	var df *DecodedFrame = (*DecodedFrame)(C.convert_format((*C.char)(cpayload), C.int(len(payload))))

	var converted []byte = C.GoBytes(unsafe.Pointer(df.data), df.size)
	C.free(unsafe.Pointer(df.data))
	C.free(unsafe.Pointer(df))
	return converted
}

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

func WritePayloadToFile(payload []byte, fileName string, codecId int,duration int) (ret int) {
	cfileName := C.CString(fileName)
	v := C.write_media_file((*C.char)(unsafe.Pointer(&payload[0])), C.int(len(payload)), (*C.char)(cfileName), C.int(codecId),C.int(duration))
	ret = int(v)
	C.free(unsafe.Pointer(cfileName))
	return
}

func (frame *DecodedFrame) ToBytes() []byte {
	return C.GoBytes(unsafe.Pointer(frame.data), frame.size)
}
func (frame *DecodedFrame) Free() {
	C.free(unsafe.Pointer(frame))
}


func TranscodeNew(fromCodecName string, toCodecName string) *TranscodeContext{
	fname := C.CString(fromCodecName)
	tname := C.CString(toCodecName)
	defer C.free(unsafe.Pointer(fname))
	defer C.free(unsafe.Pointer(tname))
	return (*TranscodeContext)(C.transcode_init_context(fname,tname))
}

func (context *TranscodeContext) Iterate(data []byte) (frame *DecodedFrame,reason int) {
	dataLen := len(data)
	frame = C.transcode_iterate(context,(*C.char)(unsafe.Pointer(&data[0])),
		C.int(dataLen),(*C.int)(unsafe.Pointer(&reason)))
	return
}

func (context *TranscodeContext) Free() {
	C.transcode_free(context)
}
