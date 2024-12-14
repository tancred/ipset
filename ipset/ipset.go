package ipset

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

type IPSet struct {
	ptr *C.struct_ipset
}

func NewIPSet() IPSet {
	return IPSet {
		ptr: C.ipset_init(),
	}
}

func (i IPSet) Close() {
		fmt.Println("closing ipset")
		r := C.ipset_fini(i.ptr)
		fmt.Println("r=", r)
}

func init() {
	fmt.Println("ipset: initializing")
	C.ipset_load_types()
}
