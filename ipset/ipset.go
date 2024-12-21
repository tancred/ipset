package ipset
// gen many entries:
//   $ declare i j int
//   $ i=1
//   $ ((j=i+200))
//   $ while (( i < j )); do ipset add bl 1.2.3.$i; ((i++)) ; done

/*
#cgo CFLAGS: -W
#cgo LDFLAGS: -lipset
#include <stdlib.h>
#include <libipset/ipset.h>

int goips_custom_printf(struct ipset *ipset, void *p);
*/
import "C"

import (
	"fmt"
	"net"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
)

type Family int

const (
	FamilyINET = iota
	FamilyINET6
)

type IPSet struct {
	ptr *C.struct_ipset
	selfptr unsafe.Pointer
}

func init() {
	fmt.Println("ipset: initializing")
	C.ipset_load_types()
}

func New() *IPSet {
	csetptr := C.ipset_init()
	set := &IPSet {
		ptr: csetptr,
		selfptr: nil,
	}
	set.selfptr = gopointer.Save(set)

	C.goips_custom_printf(set.ptr, set.selfptr)

	return set
}

func (set *IPSet) Command(command string) int {
	fmt.Println("will run:", command)
	ccmd := C.CString(command)
	defer C.free(unsafe.Pointer(ccmd))

	if set.ptr != nil {
		r := C.ipset_fini(set.ptr)
		fmt.Println("  r =", r)
		set.ptr = C.ipset_init()
		C.goips_custom_printf(set.ptr, set.selfptr)
	}

	return int(C.ipset_parse_line(set.ptr, ccmd))
}

func (set *IPSet) Close() {
		fmt.Println("closing ipset")

		r := C.ipset_fini(set.ptr)
		fmt.Println("r=", r)

		fmt.Println("closing ipset: selfptr")
		gopointer.Unref(set.selfptr)
}

func (set *IPSet) Save(name string) {
	r := set.Command(fmt.Sprintf("save %s", name))
	if r == 0 {
		fmt.Println("save OK")
	} else {
		fmt.Println("save NAY")
	}
}

func (set *IPSet) Test(name string, addr net.IP) {
	r := set.Command(fmt.Sprintf("test %s %s", name, addr.String()))
	if r == 0 {
		fmt.Println("test OK")
	} else {
		fmt.Println("test NAY")
	}
}

func (set *IPSet) Fail() {
	r := set.Command("no command at ALL")
	if r == 0 {
		fmt.Println("cmd OK")
	} else {
		fmt.Println("cmd NAY")
	}
}
