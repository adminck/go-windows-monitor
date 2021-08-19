package windows

// #define WIN32_LEAN_AND_MEAN
// #include <windows.h>
import "C"
import (
	"errors"
	"unsafe"
)

func GetPrivateProfileString(gpath, gApplicationName, gkey string) string {
	path, ApplicationName, key, lpDefault := C.CString(gpath), C.CString(gApplicationName), C.CString(gkey), C.CString("")
	buf := C.CString("")
	C.GetPrivateProfileStringA(ApplicationName, key, lpDefault, buf, 100, path)
	str := C.GoString(buf)
	C.free(unsafe.Pointer(path))
	C.free(unsafe.Pointer(ApplicationName))
	C.free(unsafe.Pointer(key))
	C.free(unsafe.Pointer(lpDefault))
	C.free(unsafe.Pointer(buf))
	return str
}

func WritePrivateProfileString(gpath, gApplicationName, gkey, gval string) error {
	val := C.CString(gval)
	path, ApplicationName, key := C.CString(gpath), C.CString(gApplicationName), C.CString(gkey)
	ret := C.WritePrivateProfileStringA(ApplicationName, key, val, path)
	C.free(unsafe.Pointer(path))
	C.free(unsafe.Pointer(ApplicationName))
	C.free(unsafe.Pointer(key))
	C.free(unsafe.Pointer(val))
	if ret == 0 {
		err := errors.New("write config failed")
		return err
	}
	return nil
}