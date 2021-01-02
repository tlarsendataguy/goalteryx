package api_new

import (
	"reflect"
	"unsafe"
)

func bytesToUtf16(value []byte) []uint16 {
	utf16Len := len(value) / 2
	var utf16Bytes []uint16
	rawHeader := (*reflect.SliceHeader)(unsafe.Pointer(&utf16Bytes))
	rawHeader.Data = uintptr(unsafe.Pointer(&value[0]))
	rawHeader.Len = utf16Len
	rawHeader.Cap = utf16Len
	return utf16Bytes
}

func utf16ToBytes(value []uint16) []byte {
	bytesLen := len(value) * 2
	var bytes []byte
	rawHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	rawHeader.Data = uintptr(unsafe.Pointer(&value[0]))
	rawHeader.Len = bytesLen
	rawHeader.Cap = bytesLen
	return bytes
}

func ptrToBytes(value unsafe.Pointer, start int, length int) []byte {
	var bytes []byte
	rawHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	rawHeader.Data = uintptr(value) + uintptr(start)
	rawHeader.Len = length
	rawHeader.Cap = length
	return bytes
}
