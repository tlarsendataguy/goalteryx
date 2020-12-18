package api_new

import (
	"reflect"
	"unsafe"
)

func generateGetFixedBytes(startAt int, length int) func(Record) []byte {
	startAtUint := uintptr(startAt)
	return func(data Record) []byte {
		var raw []byte
		rawHeader := (*reflect.SliceHeader)(unsafe.Pointer(&raw))
		rawHeader.Data = uintptr(data) + startAtUint
		rawHeader.Len = length
		rawHeader.Cap = length
		return raw
	}
}

func generateGetVarBytes(startAt int) func(Record) []byte {
	startAtUint := uintptr(startAt)
	return func(data Record) []byte {
		varStart := *((*uint32)(unsafe.Pointer(uintptr(data) + startAtUint)))
		if varStart == 0 {
			return []byte{}
		}
		if varStart == 1 {
			return nil
		}

		var varLen uint32
		offset := startAtUint

		// small string optimization, check if high bit is not set and third bit is set
		// small string optimization len is in the 29th and 30th bits
		if (varStart&0x80000000) == 0 && (varStart&0x30000000) != 0 {
			varLen = varStart >> 28
		} else {
			// strip away high bit
			// high bit is set to signal fields larger than int32 bytes to differentiate from small string optimization
			// at this point we have determined there is no small string optimization and so we can strip it away
			varStart &= 0x7fffffff
			varStartUint := uintptr(varStart)

			varLenFirstByte := *((*byte)(unsafe.Pointer(uintptr(data) + startAtUint + varStartUint)))
			offset += varStartUint
			if varLenFirstByte&byte(1) == 1 {
				varLen = uint32(varLenFirstByte >> 1)
				offset += 1
			} else {
				varLen = *((*uint32)(unsafe.Pointer(uintptr(data) + startAtUint + varStartUint))) / 2
				offset += 4
			}
		}

		var raw []byte
		rawHeader := (*reflect.SliceHeader)(unsafe.Pointer(&raw))
		rawHeader.Data = uintptr(data) + offset
		rawHeader.Len = int(varLen)
		rawHeader.Cap = int(varLen)
		return raw
	}
}
