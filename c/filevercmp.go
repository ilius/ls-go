package main

/*

#cgo LDFLAGS: -lconfig

#include "filevercmp.h"
#include <stdlib.h>
#include <stdio.h>

int filevercmp (const char *s1, const char *s2);

*/
import "C"

import (
	"unsafe"
)

func filevercmp(s1 string, s2 string) int {
	cs1 := C.CString(s1)
	cs2 := C.CString(s2)
	res := C.filevercmp(cs1, cs2)
	C.free(unsafe.Pointer(cs1))
	C.free(unsafe.Pointer(cs2))
	return int(res)
}
