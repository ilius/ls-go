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
	"fmt"
	"os"
	"unsafe"
)

func main() {
	s1 := os.Args[1]
	s2 := os.Args[2]
	cs1 := C.CString(s1)
	cs2 := C.CString(s2)
	res := C.filevercmp(cs1, cs2)
	C.free(unsafe.Pointer(cs1))
	C.free(unsafe.Pointer(cs2))
	fmt.Printf("res = %v\n", res)
	switch {
	case res == 0:
		fmt.Printf("%#v = %#v\n", s1, s2)
	case res < 0:
		fmt.Printf("%#v < %#v\n", s1, s2)
	case res > 0:
		fmt.Printf("%#v > %#v\n", s1, s2)
	}
}
