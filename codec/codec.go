package codec
//#cgo pkg-config: libavformat libavcodec libavutil
//#include <stdio.h>
//#include <stdlib.h>
//#include <inttypes.h>
//#include <stdint.h>
//#include <string.h>
//#include "codec.h"
import "C"
import "unsafe"

type (
	DecodedFrame = C.struct_DecodedFrame
)

func ConvertFormat(payload []byte) []byte {
	cpayload := C.CBytes(payload)
	//defer C.free(unsafe.Pointer(cpayload))
	var df *DecodedFrame = (*DecodedFrame)(C.convert_format((*C.char)(cpayload),C.int(len(payload))))

	var converted []byte = C.GoBytes(unsafe.Pointer(df.data),df.size)
	C.free(unsafe.Pointer(df.data))
	C.free(unsafe.Pointer(df))
	return converted
}
