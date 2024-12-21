package ipset

/*
#include <libipset/ipset.h>
*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
)

//export goipsCustomErrorFn
func goipsCustomErrorFn(cset *C.struct_ipset, p unsafe.Pointer, status C.int, msg *C.char) C.int {
	set := gopointer.Restore(p).(*IPSet)
	gomsg := C.GoString(msg)
	set.customError(cset, int(status), strings.TrimSpace(gomsg))
	// set.customError(cset, int(status), gomsg)
	return status
}

//export goipsStandardErrorFn
func goipsStandardErrorFn(cset *C.struct_ipset, p unsafe.Pointer, errType C.int, msg *C.char) {
	set := gopointer.Restore(p).(*IPSet)
	gomsg := C.GoString(msg)
	set.stdError(cset, int(errType), strings.TrimSpace(gomsg))
	// set.stdError(cset, int(errType), gomsg)
}

//export goipsPrintOutFn
func goipsPrintOutFn(p unsafe.Pointer, msg *C.char) {
	set := gopointer.Restore(p).(*IPSet)
	gomsg := C.GoString(msg)
	set.printOut(strings.TrimSpace(gomsg))
	// set.printOut(gomsg)
}

func (set *IPSet) customError(cset *C.struct_ipset, status int, msg string) {
	fmt.Printf("---- customError: %d (%s) msg: `%s'\n", status, status2String(status), msg)

}

func (set *IPSet) stdError(cset *C.struct_ipset, errType int, msg string) {
	fmt.Printf("---- stdError: %d (%s) msg: `%s'\n", errType, errType2String(errType), msg)
}

func (set *IPSet) printOut(msg string) {
	fmt.Printf("---- printOut: msg: `%s'\n", msg)
}

func errType2String(errType int) string {
	switch errType {
	case C.IPSET_NO_ERROR:
		return "no error"
	case C.IPSET_NOTICE:
		return "NOTE"
	case C.IPSET_WARNING:
		return "WARN"
	case C.IPSET_ERROR:
		return "ERRR"
	default:
		return fmt.Sprintf("unknown error %d", errType)
	}
}

func status2String(status int) string {
	switch status {
	case C.IPSET_NO_PROBLEM:
		return "nemas problemas"
	case C.IPSET_OTHER_PROBLEM:
		return "other problem"
	case C.IPSET_PARAMETER_PROBLEM:
		return "parameter problem"
	case C.IPSET_VERSION_PROBLEM:
		return "version problem"
	case C.IPSET_SESSION_PROBLEM:
		return "session problem"
	default:
		return fmt.Sprintf("unknown problem %d", status)
	}
}
