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
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
)

type Family int

const (
	FamilyINET = iota
	FamilyINET6
)

type errorLevel int

const (
	errorLevelNoError = iota
	errorLevelNotice
	errorLevelWarning
	errorLevelError
	errorLevelUnknown
)

type cmdError struct {
	Level errorLevel
	Message string
}

func (err *cmdError) Error() string {
	var lvl string
	switch err.Level {
	case errorLevelNoError:
		lvl = ""
	case errorLevelNotice:
		lvl = "notice"
	case errorLevelWarning:
		lvl = "warning"
	case errorLevelError:
		lvl = "error"
	}

	return fmt.Sprintf("%s: %s", lvl, err.Message)
}


type IPSet struct {
	ptr     *C.struct_ipset
	selfptr unsafe.Pointer
	recentError *cmdError
	recentMessage string
}

type Info struct {
	Name string
	Type string
	Family string
	Timeout *int
}

func init() {
	C.ipset_load_types()
}

func New() *IPSet {
	csetptr := C.ipset_init()
	set := &IPSet{
		ptr:     csetptr,
		selfptr: nil,
	}
	set.selfptr = gopointer.Save(set)

	C.goips_custom_printf(set.ptr, set.selfptr)

	return set
}

func (set *IPSet) Close() {
	_ = C.ipset_fini(set.ptr)
	gopointer.Unref(set.selfptr)
}

type CreateOption func (i Info) Info

func CreateOptionTimeout(timeout int) CreateOption {
	return func (i Info) Info {
		i.Timeout = &timeout
		return i
	}
}

func CreateOptionFamily(family string) CreateOption {
	return func (i Info) Info {
		i.Family = family
		return i
	}
}

func (set *IPSet) Create(name string, options ...CreateOption) error {
	info := Info{
		Name: name,
		Type: "hash:ip",
		Family: "inet",
		Timeout: nil,
	}

	for _, o := range options {
		info = o(info)
	}

	cmd := fmt.Sprintf("create %s %s", info.Name, info.Type)
	cmd = cmd + fmt.Sprintf(" family %s", info.Family)
	if info.Timeout != nil {
		cmd = cmd + fmt.Sprintf(" timeout %d", *info.Timeout)
	}
	_, _, err := set.Command(cmd)

	if err != nil {
		return err
	}

	return nil
	// var cmderr *cmdError
	// if errors.As(err, &cmderr) {
	// 	if cmderr.Level >= errorLevelError {
	// 		return false, err
	// 	}
	// } else if err != nil {
	// 	return false, err
	// }

	// return r == 0, nil
}

func (set *IPSet) Destroy(name string) error {
	_, _, err := set.Command(fmt.Sprintf("destroy %s", name))

	if err != nil {
		return err
	}

	return nil
	// var cmderr *cmdError
	// if errors.As(err, &cmderr) {
	// 	if cmderr.Level >= errorLevelError {
	// 		return false, err
	// 	}
	// } else if err != nil {
	// 	return false, err
	// }

	// return r == 0, nil
}

func (set *IPSet) Info(name string) (Info, error) {
	_, msg, err := set.Command(fmt.Sprintf("save %s", name))

	var cmderr *cmdError
	if errors.As(err, &cmderr) {
		return Info{}, err
	} else if err != nil {
		return Info{}, err
	}

	// create bl hash:ip family inet hashsize 1024 maxelem 65536 bucketsize 12 initval 0xd263dc02
	// ...
	lines := strings.Split(msg, "\n")
	fields := strings.Fields(lines[0])

	info := Info{}

	info.Name = fields[1]
	info.Type = fields[2]

	for i := 3; i + 1 < len(fields); i++ {
		key := fields[i]
		val := fields[i+1]

		switch key {
		case "family":
			info.Family = val
		case "timeout":
			if n, err := strconv.Atoi(val); err == nil {
				info.Timeout = &n
			}
		}
	}

	return info, nil
}

func (set *IPSet) Add(name string, addr net.IP) (bool, error) {
	_, _, err := set.Command(fmt.Sprintf("add %s %s", name, addr.String()))

	if err != nil {
		return false, err
	}

	return true, nil
}

func (set *IPSet) Add6(name string, addr net.IP) (bool, error) {
	var addrString string

	if addr.To4() != nil {
		addrString = fmt.Sprintf("::ffff:%s", addr.String())
	} else {
		addrString = addr.String()
	}

	_, _, err := set.Command(fmt.Sprintf("add %s %s", name, addrString))

	if err != nil {
		return false, err
	}

	return true, nil
}

func (set *IPSet) Test(name string, addr net.IP) (bool, error) {
	r, _, err := set.Command(fmt.Sprintf("test %s %s", name, addr.String()))

	var cmderr *cmdError
	if errors.As(err, &cmderr) {
		if cmderr.Level >= errorLevelError {
			return false, err
		}
	} else if err != nil {
		return false, err
	}

	return r == 0, nil
}

func (set *IPSet) Test6(name string, addr net.IP) (bool, error) {
	var addrString string

	if addr.To4() != nil {
		addrString = fmt.Sprintf("::ffff:%s", addr.String())
	} else {
		addrString = addr.String()
	}

	r, _, err := set.Command(fmt.Sprintf("test %s %s", name, addrString))

	var cmderr *cmdError
	if errors.As(err, &cmderr) {
		if cmderr.Level >= errorLevelError {
			return false, err
		}
	} else if err != nil {
		return false, err
	}

	return r == 0, nil
}

func (set *IPSet) Command(command string) (int, string, error) {
	ccmd := C.CString(command)
	defer C.free(unsafe.Pointer(ccmd))

	if set.ptr != nil {
		_ = C.ipset_fini(set.ptr)
		set.ptr = C.ipset_init()
		C.goips_custom_printf(set.ptr, set.selfptr)
	}

	set.recentError = nil
	set.recentMessage = ""

	r := int(C.ipset_parse_line(set.ptr, ccmd))

	if set.recentError != nil {
		err := set.recentError
		set.recentError = nil
		return r, "", err
	}

	msg := set.recentMessage
	set.recentMessage = ""

	return r, msg, nil
}

func (set *IPSet) Save(name string) {
	r, _, _ := set.Command(fmt.Sprintf("save %s", name))
	if r == 0 {
		fmt.Fprintln(os.Stderr, "save OK")
	} else {
		fmt.Fprintln(os.Stderr, "save NAY")
	}
}

func (set *IPSet) Fail() {
	r, _, _ := set.Command("no command at ALL")
	if r == 0 {
		fmt.Fprintln(os.Stderr, "cmd OK")
	} else {
		fmt.Fprintln(os.Stderr, "cmd NAY")
	}
}
