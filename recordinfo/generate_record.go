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
	if size := info.getRecordSize(); size > len(info.blob) {
		info.blob = make([]byte, size)
	}
	start := 0
	var err error
	for _, field := range info.fields {
		getBytes := field.generator
		if getBytes == nil {
			return nil, fmt.Errorf(`field '%v' does not have a byte generator`, field.Name)
		}
		start, err = getBytes(field, info.blob, start)
		if err != nil {
			return nil, err
		}
	}
	return unsafe.Pointer(&info.blob[0]), nil
}

func generateByte(field *fieldInfoEditor, blob []byte, startAt int) (int, error) {
	dataSize := 2
	putFunc := func(dataSize int, blobSlice []byte) error {
		blobSlice[startAt] = field.value.(byte)
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, startAt, dataSize, putFunc)
}

func generateBool(field *fieldInfoEditor, blob []byte, startAt int) (int, error) {
	if field.value == nil {
		blob[startAt] = 2
	} else if field.value.(bool) {
		blob[startAt] = 1
	} else {
		blob[startAt] = 0
	}
	return startAt + 1, nil
}

func generateInt16(field *fieldInfoEditor, blob []byte, startAt int) (int, error) {
	dataSize := 3
	putFunc := func(dataSize int, blobSlice []byte) error {
		value := field.value.(int16)
		binary.LittleEndian.PutUint16(blobSlice, uint16(value))
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, startAt, dataSize, putFunc)
}

func generateInt32(field *fieldInfoEditor, blob []byte, startAt int) (int, error) {
	dataSize := 5
	putFunc := func(dataSize int, blobSlice []byte) error {
		value := field.value.(int32)
		binary.LittleEndian.PutUint32(blobSlice, uint32(value))
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, startAt, dataSize, putFunc)
}

func generateInt64(field *fieldInfoEditor, blob []byte, startAt int) (int, error) {
	dataSize := 9
	putFunc := func(dataSize int, blobSlice []byte) error {
		value := field.value.(int64)
		binary.LittleEndian.PutUint64(blobSlice, uint64(value))
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, startAt, dataSize, putFunc)
}

func generateFixedDecimal(field *fieldInfoEditor, blob []byte, startAt int) (int, error) {
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
	return putFixedBytesWithNullByte(field, blob, startAt, dataSize, putFunc)
}

func generateFloat32(field *fieldInfoEditor, blob []byte, startAt int) (int, error) {
	dataSize := 5
	putFunc := func(dataSize int, blobSlice []byte) error {
		value := field.value.(float32)
		binary.LittleEndian.PutUint32(blobSlice, math.Float32bits(value))
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, startAt, dataSize, putFunc)
}

func generateFloat64(field *fieldInfoEditor, blob []byte, startAt int) (int, error) {
	dataSize := 9
	putFunc := func(dataSize int, blobSlice []byte) error {
		value := field.value.(float64)
		binary.LittleEndian.PutUint64(blobSlice, math.Float64bits(value))
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, startAt, dataSize, putFunc)
}

func generateString(field *fieldInfoEditor, blob []byte, startAt int) (int, error) {
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
	return putFixedBytesWithNullByte(field, blob, startAt, dataSize, putFunc)
}

func generateWString(field *fieldInfoEditor, blob []byte, startAt int) (int, error) {
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
	return putFixedBytesWithNullByte(field, blob, startAt, dataSize, putFunc)
}

func generateDate(field *fieldInfoEditor, blob []byte, startAt int) (int, error) {
	dataSize := int(field.fixedLen) + 1
	putFunc := func(dataSize int, blobSlice []byte) error {
		value := field.value.(time.Time)
		valueStr := value.Format(dateFormat)
		for index := range valueStr {
			blobSlice[index] = valueStr[index]
		}
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, startAt, dataSize, putFunc)
}

func generateDateTime(field *fieldInfoEditor, blob []byte, startAt int) (int, error) {
	dataSize := int(field.fixedLen) + 1
	putFunc := func(dataSize int, blobSlice []byte) error {
		value := field.value.(time.Time)
		valueStr := value.Format(dateTimeFormat)
		for index := range valueStr {
			blobSlice[index] = valueStr[index]
		}
		return nil
	}
	return putFixedBytesWithNullByte(field, blob, startAt, dataSize, putFunc)
}

func putFixedBytesWithNullByte(field *fieldInfoEditor, blob []byte, startAt int, dataSize int, putData func(dataSize int, blobSlice []byte) error) (int, error) {
	blobSlice := blob[startAt : startAt+dataSize]

	if field.value == nil {
		for index := 0; index < dataSize-1; index++ {
			blobSlice[index] = 0
		}
		blobSlice[dataSize-1] = 1
	} else {
		err := putData(dataSize, blobSlice)
		if err != nil {
			return startAt, err
		}
		blobSlice[dataSize-1] = 0
	}
	return startAt + dataSize, nil
}
