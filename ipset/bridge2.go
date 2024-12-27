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
	return status
}

//export goipsStandardErrorFn
func goipsStandardErrorFn(cset *C.struct_ipset, p unsafe.Pointer, errType C.int, msg *C.char) {
	set := gopointer.Restore(p).(*IPSet)
	gomsg := C.GoString(msg)
	set.stdError(cset, int(errType), strings.TrimSpace(gomsg))
}

//export goipsPrintOutFn
func goipsPrintOutFn(p unsafe.Pointer, msg *C.char) {
	set := gopointer.Restore(p).(*IPSet)
	gomsg := C.GoString(msg)
	set.printOut(strings.TrimSpace(gomsg))
}

func (set *IPSet) customError(cset *C.struct_ipset, status int, msg string) {
	fmt.Printf("  customError: %d (%s) msg: `%s'\n", status, status2String(status), msg)
}

func (set *IPSet) stdError(cset *C.struct_ipset, errType int, msg string) {
	set.recentError = &cmdError{
		Level:   errType2Level(errType),
		Message: msg,
	}
}

func (set *IPSet) printOut(msg string) {
	set.recentMessage = set.recentMessage + msg
}

func errType2Level(errType int) errorLevel {
	switch errType {
	case C.IPSET_NO_ERROR:
		return errorLevelNoError
	case C.IPSET_NOTICE:
		return errorLevelNotice
	case C.IPSET_WARNING:
		return errorLevelWarning
	case C.IPSET_ERROR:
		return errorLevelError
	}
	return errorLevelUnknown
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
