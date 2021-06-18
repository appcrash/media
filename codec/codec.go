package codec
//#cgo pkg-config: libavformat libavcodec libavutil
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
	DecodedFrame = C.struct_DecodedFrame
)

func ConvertFormat(payload []byte) []byte {
	cpayload := C.CBytes(payload)
	defer C.free(unsafe.Pointer(cpayload))
	var df *DecodedFrame = (*DecodedFrame)(C.convert_format((*C.char)(cpayload),C.int(len(payload))))

	var converted []byte = C.GoBytes(unsafe.Pointer(df.data),df.size)
	C.free(unsafe.Pointer(df.data))
	C.free(unsafe.Pointer(df))
	return converted
}

func GetPayloadFromFile(fp string) []byte {
	cfp := C.CString(fp)
	payload := (*C.struct_Payload)(C.read_media_file(cfp))
	C.free(unsafe.Pointer(cfp))

	fmt.Printf("data is size:%v,bitrate:%v,frame_size:%v\n",payload.size,payload.bitrate,
		BitrateToFrameSize(float64(payload.bitrate),20))
	if payload != nil {
		data := C.GoBytes(unsafe.Pointer(payload.data),payload.size)
		C.free(unsafe.Pointer(payload.data))
		C.free(unsafe.Pointer(payload))
		return data
	}
	return nil
}

func WritePayloadToFile(payload []byte,fileName string,codecId int) (ret int){
	cfileName := C.CString(fileName)
	v := C.write_media_file((*C.char)(unsafe.Pointer(&payload[0])),C.int(len(payload)),(*C.char)(cfileName),C.int(codecId))
	ret = int(v)
	C.free(unsafe.Pointer(cfileName))
	return
}
