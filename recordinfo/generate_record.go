package recordinfo

import (
	"encoding/binary"
	"unsafe"
)

func (info *recordInfo) GenerateRecord() (unsafe.Pointer, error) {
	fixed, variable := info.getRecordSizes()
	totalSize := fixed + 4 + variable
	if totalSize > len(info.blob) {
		info.blob = make([]byte, totalSize)
	}
	binary.LittleEndian.PutUint32(info.blob[fixed:fixed+4], uint32(variable))
	return unsafe.Pointer(&info.blob[0]), nil
}

func putFixedBytesWithNullByte(field *fieldInfoEditor, blob []byte, fixedStartAt int, dataSize int, putData func(dataSize int, blobSlice []byte) error, varStartAt int) (int, int, error) {
	blobSlice := blob[fixedStartAt : fixedStartAt+dataSize]

	if field.fixedValue == nil {
		for index := 0; index < dataSize-1; index++ {
			blobSlice[index] = 0
		}
		blobSlice[dataSize-1] = 1
	} else {
		err := putData(dataSize, blobSlice)
		if err != nil {
			return fixedStartAt, varStartAt, err
		}
		blobSlice[dataSize-1] = 0
	}
	return fixedStartAt + dataSize, varStartAt, nil
}

func putVarNull(blob []byte, fixedStartAt int, varStartAt int) (int, int, error) {
	binary.LittleEndian.PutUint32(blob[fixedStartAt:fixedStartAt+4], 1)
	return fixedStartAt + 4, varStartAt, nil
}

func putVarData(field *fieldInfoEditor, blob []byte, data []byte, fixedStartAt int, varStartAt int) (int, int, error) {
	varDataLen := uint32(len(data))

	// Small string optimization
	if varDataLen < 4 {
		varDataLen <<= 28
		fixedBytes := make([]byte, 4)
		for index, fixedByte := range data {
			fixedBytes[index] = fixedByte
		}
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

	var index uint32
	for index = 0; index < varDataLen; index++ {
		blob[varStartAt] = data[index]
		varStartAt++
	}
	return fixedStartAt, varStartAt, nil
}
