package recordinfo

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

func (info *recordInfo) GetIntValueFrom(fieldName string, record unsafe.Pointer) (value int, isNull bool, err error) {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return 0, false, err
	}

	switch field.Type {
	case Byte:
		raw := *((*[2]byte)(unsafe.Pointer(uintptr(record) + field.location)))
		if raw[1] == 1 {
			return 0, true, nil
		}
		return int(raw[0]), false, nil

	case Int16:
		raw := *((*[3]byte)(unsafe.Pointer(uintptr(record) + field.location)))
		if raw[2] == 1 {
			return 0, true, nil
		}
		return int(binary.LittleEndian.Uint16(raw[:2])), false, nil

	case Int32:
		raw := *((*[5]byte)(unsafe.Pointer(uintptr(record) + field.location)))
		if raw[4] == 1 {
			return 0, true, nil
		}
		return int(binary.LittleEndian.Uint32(raw[:4])), false, nil

	case Int64:
		raw := *((*[9]byte)(unsafe.Pointer(uintptr(record) + field.location)))
		if raw[8] == 1 {
			return 0, true, nil
		}
		return int(binary.LittleEndian.Uint64(raw[:8])), false, nil

	default:
		return 0, false, invalidTypeError(field, `int`)
	}
}

func (info *recordInfo) GetBoolValueFrom(fieldName string, record unsafe.Pointer) (value bool, isNull bool, err error) {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return false, false, err
	}

	if field.Type != Bool {
		return false, false, invalidTypeError(field, `bool`)
	}

	raw := *((*byte)(unsafe.Pointer(uintptr(record) + field.location)))
	if raw == 2 {
		return false, true, nil
	}
	if raw == 0 {
		return false, false, nil
	}
	return true, false, nil
}

func (info *recordInfo) GetFloatValueFrom(fieldName string, record unsafe.Pointer) (value float64, isNull bool, err error) {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return 0, false, err
	}

	switch field.Type {
	case Float:
		raw := *((*[5]byte)(unsafe.Pointer(uintptr(record) + field.location)))
		if raw[4] == 1 {
			return 0, true, nil
		}
		return float64(math.Float32frombits(binary.LittleEndian.Uint32(raw[:4]))), false, nil
	case Double:
		raw := *((*[9]byte)(unsafe.Pointer(uintptr(record) + field.location)))
		if raw[8] == 1 {
			return 0, true, nil
		}
		return math.Float64frombits(binary.LittleEndian.Uint64(raw[:8])), false, nil
	case FixedDecimal:
		raw := info.getRawBytes(field, record)
		if raw[field.fixedLen] == 1 {
			return 0, true, nil
		}
		value, err := strconv.ParseFloat(string(truncateAtNullByte(raw)), 64)
		return value, false, err
	default:
		return 0, false, invalidTypeError(field, `float`)
	}
}

func (info *recordInfo) GetStringValueFrom(fieldName string, record unsafe.Pointer) (value string, isNull bool, err error) {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return ``, false, err
	}

	switch field.Type {
	case String:
		raw := info.getRawBytes(field, record)
		if raw[field.fixedLen] == 1 {
			return ``, true, nil
		}
		return string(truncateAtNullByte(raw)), false, nil

	case V_String:
		raw := info.getRawBytes(field, record)
		if raw == nil {
			return ``, true, nil
		}
		if len(raw) == 0 {
			return ``, false, nil
		}
		return string(raw), false, nil

	case WString:
		raw := info.getRawBytes(field, record)
		if raw[field.fixedLen] == 1 {
			return ``, true, nil
		}
		charLen := len(raw) / 2
		chars := make([]uint16, charLen)
		for charIndex := 0; charIndex < charLen; charIndex++ {
			chars[charIndex] = binary.LittleEndian.Uint16(raw[charIndex*2 : charIndex*2+2])
		}
		return syscall.UTF16ToString(chars), false, nil

	case V_WString:
		raw := info.getRawBytes(field, record)
		if raw == nil {
			return ``, true, nil
		}
		if len(raw) == 0 {
			return ``, false, nil
		}

		charLen := len(raw) / 2
		chars := make([]uint16, charLen)
		for charIndex := 0; charIndex < charLen; charIndex++ {
			chars[charIndex] = binary.LittleEndian.Uint16(raw[charIndex*2 : charIndex*2+2])
		}
		return syscall.UTF16ToString(chars), false, nil

	default:
		return ``, false, invalidTypeError(field, `string`)
	}
}

func (info *recordInfo) GetDateValueFrom(fieldName string, record unsafe.Pointer) (value time.Time, isNull bool, err error) {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return zeroDate, false, err
	}

	raw := info.getRawBytes(field, record)
	if raw[field.fixedLen] == 1 {
		return zeroDate, true, nil
	}

	switch field.Type {
	case Date:
		value, err := time.Parse(dateFormat, string(raw[:field.fixedLen]))
		return value, false, err

	case DateTime:
		value, err := time.Parse(dateTimeFormat, string(raw[:field.fixedLen]))
		return value, false, err

	default:
		return zeroDate, false, invalidTypeError(field, `date`)
	}
}

func (info *recordInfo) GetRawBytesFrom(fieldName string, record unsafe.Pointer) (value []byte, err error) {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return nil, fmt.Errorf(`error getting raw bytes: %v`, err.Error())
	}
	return info.getRawBytes(field, record), nil
}

func (info *recordInfo) getRawBytes(field *fieldInfoEditor, record unsafe.Pointer) []byte {
	switch field.Type {
	case V_String, V_WString:
		return getVarBytes(field, record)
	default:
		return getFixedBytes(field, record)
	}
}

func getFixedBytes(field *fieldInfoEditor, record unsafe.Pointer) []byte {
	totalLen := int(field.fixedLen + field.nullByteLen)
	var raw []byte
	rawHeader := (*reflect.SliceHeader)(unsafe.Pointer(&raw))
	rawHeader.Data = uintptr(record) + field.location
	rawHeader.Len = totalLen
	rawHeader.Cap = totalLen
	return raw
}

func getVarBytes(field *fieldInfoEditor, record unsafe.Pointer) []byte {
	varStart := *((*uint32)(unsafe.Pointer(uintptr(record) + field.location)))
	if varStart == 0 {
		return []byte{}
	}
	if varStart == 1 {
		return nil
	}

	var varLen uint32
	offset := field.location

	// small string optimization, check if high bit is not set and third bit is set
	// small string optimization len is in the 29th and 30th bits
	if (varStart&0x80000000) == 0 && (varStart&0x30000000) != 0 {
		varLen = varStart >> 28
	} else {
		// strip away high bit
		// high bit is set to signal fields larger than int32 bytes to differentiate from small string optimization
		// at this point we have determined there is no small string optimization and so we can strip it away
		varStart &= 0x7fffffff

		varLenFirstByte := *((*byte)(unsafe.Pointer(uintptr(record) + field.location + uintptr(varStart))))
		offset += uintptr(varStart)
		if varLenFirstByte&byte(1) == 1 {
			varLen = uint32(varLenFirstByte >> 1)
			offset += 1
		} else {
			varLen = *((*uint32)(unsafe.Pointer(uintptr(record) + field.location + uintptr(varStart)))) / 2
			offset += 4
		}
	}

	var raw []byte
	rawHeader := (*reflect.SliceHeader)(unsafe.Pointer(&raw))
	rawHeader.Data = uintptr(record) + offset
	rawHeader.Len = int(varLen)
	rawHeader.Cap = int(varLen)
	return raw
}

func truncateAtNullByte(raw []byte) []byte {
	var dataLen int
	for dataLen = 0; dataLen < len(raw); dataLen++ {
		if raw[dataLen] == 0 {
			break
		}
	}
	return raw[:dataLen]
}

func invalidTypeError(field *fieldInfoEditor, expectedType string) error {
	return fmt.Errorf(`[%v]'s type of '%v' is not a valid %v type`, field.Name, fieldTypeMap[field.Type], expectedType)
}
