package util

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"reflect"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

func Encrypt(value string) (string, error) {
	dll := syscall.NewLazyDLL(`SrcLib.dll`)
	proc := dll.NewProc(`EncryptPassword`)

	dest := C.malloc(C.ulonglong(1024 * 2))
	nMode := uintptr(0)
	size := uintptr(1024)

	sourcePtr, err := syscall.UTF16PtrFromString(value)
	if err != nil {
		return ``, err
	}

	written, _, err := proc.Call(nMode, uintptr(unsafe.Pointer(sourcePtr)), uintptr(dest), size)
	if err.(syscall.Errno) != 0 {
		return ``, err
	}
	written, _, err = proc.Call(nMode, uintptr(dest), uintptr(dest), size)
	if err.(syscall.Errno) != 0 {
		return ``, err
	}

	var destSlice []uint16
	destHeader := (*reflect.SliceHeader)(unsafe.Pointer(&destSlice))
	destHeader.Data = uintptr(dest)
	destHeader.Len = int(written)
	destHeader.Cap = int(written)
	output := string(utf16.Decode(destSlice))

	C.free(dest)
	return output, nil
}
