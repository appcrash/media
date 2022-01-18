package codec

//#cgo pkg-config: libavformat libavcodec libavutil libswresample libavfilter freetype2
//
//#include <stdio.h>
//#include <stdlib.h>
//#include <inttypes.h>
//#include <stdint.h>
//#include <string.h>
//#include <libavformat/avformat.h>
//#include "codec.h"
import "C"
import "unsafe"

func EncodeText() {
	buff := C.video_render()
	payload := C.GoBytes(unsafe.Pointer(buff.data), buff.size)
	analyzeH264(payload)
}
