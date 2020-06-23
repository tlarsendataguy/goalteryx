package recordinfo

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/binary"
	"github.com/tlarsen7572/goalteryx/recordblob"
	"reflect"
	"unsafe"
)

// GenerateRecord creates a record blob from the current values saved to RecordInfo.  The recordInfo struct contains
// a buffer of bytes where the record is stored.  The return value is an unsafe pointer pointing to the first
// element of this blob.  For performance reasons, we try to minimize memory allocation and so we reuse the
// buffer each time GenerateRecord is called.  Allocations only occur when the record being generated would exceed
// the size of the buffer.  In these cases, a new buffer is allocated with a bit of extra padding.  Buffers will
// only grow; we do not shrink them.
func (info *recordInfo) GenerateRecord() (recordblob.RecordBlob, error) {
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
			if field.isNull {
				if field.Type == Bool {
					info.blob[fixedWriteIndex] = 2
				} else {
					info.blob[fixedWriteIndex+fixedLen-1] = 1
				}
			} else {
				copy(info.blob[fixedWriteIndex:fixedWriteIndex+fixedLen], field.value)
			}
			fixedWriteIndex += fixedLen
		}
	}
	return recordblob.NewRecordBlob(info.blobHandle), nil
}

// allocateNewBuffer performs the actual allocation, as well as frees the old buffer.  These buffers are
// malloc'd from C to ensure we don't copy data when passing from Go to C.  The byte slice that sits on top
// of the malloc'd memory simply allows us to interact with the memory from Go.
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

// putVarData saves variable-length field data into the blob.  There are a few optimizations that are made.  If
// the variable-length data is smaller than 4 bytes, it gets put into the fixed-length portion of the record blob
// and no data is stored in the variable-length portion.  If the data is between 4 and 127 bytes, its length can
// be stored into a single byte rather than the normal 4-byte integer.  Data with lengths exceeding 127 bytes are
// stored with the full 4-byte integer containing the length of the data, followed by the data itself.
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
