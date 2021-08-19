package windows

import (
	"unsafe"
)
import (
	"syscall"
)

var (
	modkernel32     = syscall.NewLazyDLL("kernel32.dll")
	procCreateEvent = modkernel32.NewProc("CreateEventW")
	procSetEvent    = modkernel32.NewProc("SetEvent")
)

func CreateEvent(sa *syscall.SecurityAttributes, bManualReset uint32, bInitialState uint32, name string) (handle syscall.Handle, err error) {
	n, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}

	r0, _, e1 := procCreateEvent.Call(uintptr(unsafe.Pointer(sa)), uintptr(bManualReset), uintptr(bInitialState), uintptr(unsafe.Pointer(n)))
	handle = syscall.Handle(r0)
	if handle == syscall.InvalidHandle {
		if e1 != nil {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}

	return
}

func SetEvent(handle syscall.Handle) bool {
	r1, _, _ := procSetEvent.Call(uintptr(handle))
	if 1 == r1 {
		return true
	}

	return false
}
