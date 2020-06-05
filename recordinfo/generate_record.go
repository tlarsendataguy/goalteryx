package recordinfo

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/binary"
	"reflect"
	"unsafe"
)

func (info *recordInfo) GenerateRecord() (unsafe.Pointer, error) {
	fixed, variable := info.getRecordSizes()
	totalSize := fixed + 4 + variable
	if totalSize >= info.blobLen {
		info.allocateNewBuffer(totalSize + 20) // arbitrary padding to try and minimize memory allocation
	}
	binary.LittleEndian.PutUint32(info.blob[fixed:fixed+4], uint32(variable))

	fixedWriteIndex := 0
	varWriteIndex := fixed + 4
	for _, field := range info.fields {
		switch field.Type {
		case V_String, V_WString, Blob, Spatial:
			if field.isNull {
				binary.LittleEndian.PutUint32(info.blob[fixedWriteIndex:fixedWriteIndex+4], 1)
				fixedWriteIndex += 4
				continue
			}
			if field.varLen == 0 {
				binary.LittleEndian.PutUint32(info.blob[fixedWriteIndex:fixedWriteIndex+4], 0)
				fixedWriteIndex += 4
				continue
			}
			var err error
			fixedWriteIndex, varWriteIndex, err = putVarData(info.blob, field, fixedWriteIndex, varWriteIndex)
			if err != nil {
				return nil, err
			}
		default:
			fixedLen := len(field.value)
			copy(info.blob[fixedWriteIndex:fixedWriteIndex+fixedLen], field.value)
			fixedWriteIndex += fixedLen
		}
	}
	return info.blobHandle, nil
}

func (info *recordInfo) allocateNewBuffer(newSize int) {
	if info.blobHandle != nil {
		C.free(info.blobHandle)
	}
	info.blobHandle = C.malloc(C.ulonglong(newSize))

	info.blob = []byte{}
	blobSlice := (*reflect.SliceHeader)(unsafe.Pointer(&info.blob))
	blobSlice.Data = uintptr(info.blobHandle)
	blobSlice.Len = newSize
	blobSlice.Cap = newSize

	info.blobLen = newSize
}

func putVarData(blob []byte, field *fieldInfoEditor, fixedStartAt int, varStartAt int) (int, int, error) {
	varDataLen := uint32(field.varLen)

	// Small string optimization
	if varDataLen < 4 {
		varDataLen <<= 28
		fixedBytes := make([]byte, 4)
		copy(fixedBytes, field.value[0:field.varLen])
		varDataUint32 := binary.LittleEndian.Uint32(fixedBytes) | varDataLen
		binary.LittleEndian.PutUint32(blob[fixedStartAt:fixedStartAt+4], varDataUint32)
		return fixedStartAt + 4, varStartAt, nil
	}

	binary.LittleEndian.PutUint32(blob[fixedStartAt:fixedStartAt+4], uint32(varStartAt-fixedStartAt))
	fixedStartAt += 4

	if varDataLen < 128 {
		blob[varStartAt] = byte(varDataLen*2) | 1 // Alteryx seems to multiply all var lens by 2
		varStartAt += 1
	} else {
		binary.LittleEndian.PutUint32(blob[varStartAt:varStartAt+4], varDataLen*2) // Alteryx seems to multiply all var lens by 2
		varStartAt += 4
	}

	copy(blob[varStartAt:varStartAt+field.varLen], field.value)
	return fixedStartAt, varStartAt + field.varLen, nil
}
