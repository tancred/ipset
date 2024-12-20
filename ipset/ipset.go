package ipset

/*
#cgo CFLAGS: -W
#cgo LDFLAGS: -lipset
#include "cipset.h"
#include <stdlib.h>
#include <libipset/ipset.h>
*/
import "C"

import (
	"fmt"
	"net"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
)

type IPSet struct {
	ptr *C.struct_ipset
	selfptr unsafe.Pointer
}

func NewIPSet() *IPSet {
	csetptr := C.ipset_init()
	set := &IPSet {
		ptr: csetptr,
		selfptr: nil,
	}
	set.selfptr = gopointer.Save(set)

	C.goips_custom_printf(set.ptr, set.selfptr)

	return set
}

func (set IPSet) Close() {
		fmt.Println("closing ipset")

		r := C.ipset_fini(set.ptr)
		fmt.Println("r=", r)

		fmt.Println("closing ipset: selfptr")
		gopointer.Unref(set.selfptr)
}

func (set IPSet) Save(name string) {
	saveCmd := C.CString(fmt.Sprintf("save %s", name))
	defer C.free(unsafe.Pointer(saveCmd))

	r := C.ipset_parse_line(set.ptr, saveCmd)

	if r == 0 {
		fmt.Println("save OK")
	} else {
		fmt.Println("save NAY")
	}
}

func (set IPSet) Test(name string, addr net.IP) {
	restoreCmd := C.CString(fmt.Sprintf("test %s %s", name, addr.String()))
	defer C.free(unsafe.Pointer(restoreCmd))

	r := C.ipset_parse_line(set.ptr, restoreCmd)

	if r == 0 {
		fmt.Println("test OK")
	} else {
		fmt.Println("test NAY")
	}
}

func (set IPSet) Fail() {
	cmd := C.CString("no command at ALL")
	defer C.free(unsafe.Pointer(cmd))

	r := C.ipset_parse_line(set.ptr, cmd)

	if r == 0 {
		fmt.Println("cmd OK")
	} else {
		fmt.Println("cmd NAY")
	}
}

func init() {
	fmt.Println("ipset: initializing")
	C.ipset_load_types()
}
