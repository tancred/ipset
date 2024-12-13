package main

/*
#cgo CFLAGS: -W
#cgo LDFLAGS: -lipset
#include "cipset.h"
#include <libipset/ipset.h>
*/
import "C"

import (
	"fmt"
)

func init() {
	fmt.Println("init ipset")
	C.goips_init()

	var ptr *C.struct_ipset
	ptr = C.ipset_init()

	if ptr == nil {
		fmt.Println("unable to get ipset if")
	}

	fmt.Printf("if %v\n", ptr)
	fmt.Println("finishing")
	r := C.ipset_fini(ptr)
	ptr = nil
	fmt.Println("r=", r)
}
