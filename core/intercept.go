package core

/*
#include <stddef.h>
#include <stdlib.h>
#include "cgo_helpers.h"
*/
import "C"

import (
	"unsafe"
)

//export diveCallback
func diveCallback(cdata *C.uchar, csize C.uint, cfingerprint *C.uchar, cfsize C.uint, cuserdata unsafe.Pointer) C.int {
	if diveCallback1A75734AFunc != nil {
		data1a75734a := string(C.GoBytes(unsafe.Pointer(cdata), C.int(csize)))
		size1a75734a := (uint32)(csize)
		fingerprint1a75734a := string(C.GoBytes(unsafe.Pointer(cfingerprint), C.int(cfsize)))
		fsize1a75734a := (uint32)(cfsize)
		userdata1a75734a := (unsafe.Pointer)(unsafe.Pointer(cuserdata))
		ret1a75734a := diveCallback1A75734AFunc(data1a75734a, size1a75734a, fingerprint1a75734a, fsize1a75734a, userdata1a75734a)
		ret, _ := (C.int)(ret1a75734a), cgoAllocsUnknown
		return ret
	}
	panic("callback func has not been set (race?)")
}
