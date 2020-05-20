package recordinfo

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

func (info *recordInfo) GenerateRecord() (unsafe.Pointer, error) {
	fixed, variable := info.getRecordSizes()
	totalSize := fixed + 4 + variable
	if totalSize > len(info.blob) {
		info.blob = make([]byte, totalSize)
	}
	binary.LittleEndian.PutUint32(info.blob[fixed:fixed+4], uint32(variable))

	fixedStart := 0
	varStart := fixed + 4
	var err error
	for _, field := range info.fields {
		getBytes := field.generator
		if getBytes == nil {
			return nil, fmt.Errorf(`field '%v' does not have a byte generator`, field.Name)
		}
		fixedStart, varStart, err = getBytes(field, info.blob, fixedStart, varStart)
		if err != nil {
			return nil, err
		}
	}
	return unsafe.Pointer(&info.blob[0]), nil
}

func generateByte(field *fieldInfoEditor, blob []byte, fixedStartAt int, varStartAt int) (int, int, error) {
	dataSize := 2
	putFunc := func(dataSize int, blobSlice []byte) error {
		blobSlice[fixedStartAt] = field.value.(byte)
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, fixedStartAt, dataSize, putFunc, varStartAt)
}

func generateBool(field *fieldInfoEditor, blob []byte, startAt int, varStartAt int) (int, int, error) {
	if field.value == nil {
		blob[startAt] = 2
	} else if field.value.(bool) {
		blob[startAt] = 1
	} else {
		blob[startAt] = 0
	}
	return startAt + 1, varStartAt, nil
}

func generateInt16(field *fieldInfoEditor, blob []byte, fixedStartAt int, varStartAt int) (int, int, error) {
	dataSize := 3
	putFunc := func(dataSize int, blobSlice []byte) error {
		value := field.value.(int16)
		binary.LittleEndian.PutUint16(blobSlice, uint16(value))
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, fixedStartAt, dataSize, putFunc, varStartAt)
}

func generateInt32(field *fieldInfoEditor, blob []byte, fixedStartAt int, varStartAt int) (int, int, error) {
	dataSize := 5
	putFunc := func(dataSize int, blobSlice []byte) error {
		value := field.value.(int32)
		binary.LittleEndian.PutUint32(blobSlice, uint32(value))
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, fixedStartAt, dataSize, putFunc, varStartAt)
}

func generateInt64(field *fieldInfoEditor, blob []byte, fixedStartAt int, varStartAt int) (int, int, error) {
	dataSize := 9
	putFunc := func(dataSize int, blobSlice []byte) error {
		value := field.value.(int64)
		binary.LittleEndian.PutUint64(blobSlice, uint64(value))
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, fixedStartAt, dataSize, putFunc, varStartAt)
}

func generateFixedDecimal(field *fieldInfoEditor, blob []byte, fixedStartAt int, varStartAt int) (int, int, error) {
	dataSize := int(field.fixedLen) + 1
	putFunc := func(dataSize int, blobSlice []byte) error {
		value := field.value.(float64)
		format := `%` + fmt.Sprintf(`%v.%vf`, field.Size, field.Precision)
		valueStr := strings.TrimSpace(fmt.Sprintf(format, value))
		for index := 0; index < dataSize-1; index++ {
			if index >= len(valueStr) {
				blobSlice[index] = 0
			} else {
				blobSlice[index] = valueStr[index]
			}
		}
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, fixedStartAt, dataSize, putFunc, varStartAt)
}

func generateFloat32(field *fieldInfoEditor, blob []byte, fixedStartAt int, varStartAt int) (int, int, error) {
	dataSize := 5
	putFunc := func(dataSize int, blobSlice []byte) error {
		value := field.value.(float32)
		binary.LittleEndian.PutUint32(blobSlice, math.Float32bits(value))
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, fixedStartAt, dataSize, putFunc, varStartAt)
}

func generateFloat64(field *fieldInfoEditor, blob []byte, fixedStartAt int, varStartAt int) (int, int, error) {
	dataSize := 9
	putFunc := func(dataSize int, blobSlice []byte) error {
		value := field.value.(float64)
		binary.LittleEndian.PutUint64(blobSlice, math.Float64bits(value))
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, fixedStartAt, dataSize, putFunc, varStartAt)
}

func generateString(field *fieldInfoEditor, blob []byte, fixedStartAt int, varStartAt int) (int, int, error) {
	dataSize := int(field.fixedLen) + 1
	putFunc := func(dataSize int, blobSlice []byte) error {
		value := field.value.(string)
		for index := 0; index < dataSize-1; index++ {
			if index >= len(value) {
				blobSlice[index] = 0
			} else {
				blobSlice[index] = value[index]
			}
		}
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, fixedStartAt, dataSize, putFunc, varStartAt)
}

func generateWString(field *fieldInfoEditor, blob []byte, fixedStartAt int, varStartAt int) (int, int, error) {
	dataSize := int(field.fixedLen) + 1
	putFunc := func(dataSize int, blobSlice []byte) error {
		value := field.value.(string)
		valueChars, err := syscall.UTF16FromString(value)
		if err != nil {
			return err
		}
		for index := 0; index < int(field.fixedLen)/2; index++ {
			byteIndex := index * 2
			if index > len(valueChars) {
				blobSlice[byteIndex] = 0
				blobSlice[byteIndex+1] = 0
				continue
			}
			char := valueChars[index]
			if char == 0 {
				break
			}
			binary.LittleEndian.PutUint16(blobSlice[byteIndex:byteIndex+2], char)
		}
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, fixedStartAt, dataSize, putFunc, varStartAt)
}

func generateDate(field *fieldInfoEditor, blob []byte, fixedStartAt int, varStartAt int) (int, int, error) {
	dataSize := int(field.fixedLen) + 1
	putFunc := func(dataSize int, blobSlice []byte) error {
		value := field.value.(time.Time)
		valueStr := value.Format(dateFormat)
		for index := range valueStr {
			blobSlice[index] = valueStr[index]
		}
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, fixedStartAt, dataSize, putFunc, varStartAt)
}

func generateDateTime(field *fieldInfoEditor, blob []byte, fixedStartAt int, varStartAt int) (int, int, error) {
	dataSize := int(field.fixedLen) + 1
	putFunc := func(dataSize int, blobSlice []byte) error {
		value := field.value.(time.Time)
		valueStr := value.Format(dateTimeFormat)
		for index := range valueStr {
			blobSlice[index] = valueStr[index]
		}
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, fixedStartAt, dataSize, putFunc, varStartAt)
}

func generateV_String(field *fieldInfoEditor, blob []byte, fixedStartAt int, varStartAt int) (int, int, error) {
	value := field.value.(string)
	valueBytes := []byte(value)
	return putVarData(field, blob, valueBytes, fixedStartAt, varStartAt)
}

func generateV_WString(field *fieldInfoEditor, blob []byte, fixedStartAt int, varStartAt int) (int, int, error) {
	value := field.value.(string)
	valueUtf16, err := syscall.UTF16FromString(value)
	if err != nil {
		return fixedStartAt, varStartAt, err
	}

	// remove the null terminator from the UTF16 string
	valueUtf16 = valueUtf16[0 : len(valueUtf16)-1]

	valueBytes := make([]byte, len(valueUtf16)*2)
	for index, valueChar := range valueUtf16 {
		binary.LittleEndian.PutUint16(valueBytes[index*2:(index*2)+2], valueChar)
	}
	return putVarData(field, blob, valueBytes, fixedStartAt, varStartAt)
}

func putFixedBytesWithNullByte(field *fieldInfoEditor, blob []byte, fixedStartAt int, dataSize int, putData func(dataSize int, blobSlice []byte) error, varStartAt int) (int, int, error) {
	blobSlice := blob[fixedStartAt : fixedStartAt+dataSize]

	if field.value == nil {
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
