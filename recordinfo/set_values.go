package recordinfo

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"syscall"
	"time"
)

func (info *recordInfo) SetIntField(fieldName string, value int) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	clearNullFlag(field)
	switch field.Type {
	case Byte:
		field.value[0] = byte(value)
	case Int16:
		binary.LittleEndian.PutUint16(field.value[0:2], uint16(value))
	case Int32:
		binary.LittleEndian.PutUint32(field.value[0:4], uint32(value))
	case Int64:
		binary.LittleEndian.PutUint64(field.value[0:8], uint64(value))
	default:
		return invalidTypeError(field, `int`)
	}
	return nil
}

func (info *recordInfo) SetBoolField(fieldName string, value bool) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	if field.Type != Bool {
		return invalidTypeError(field, `bool`)
	}

	if value {
		field.value[0] = 1
	} else {
		field.value[0] = 0
	}
	field.isNull = false
	return nil
}

func (info *recordInfo) SetFloatField(fieldName string, value float64) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	clearNullFlag(field)

	switch field.Type {
	case FixedDecimal:
		format := `%` + fmt.Sprintf(`%v.%vf`, field.Size, field.Precision)
		valueStr := []byte(strings.TrimSpace(fmt.Sprintf(format, value)))
		size := int(field.fixedLen)
		copy(field.value, valueStr)
		if size < int(field.fixedLen) {
			field.value[size] = 0
		}

	case Float:
		data := math.Float32bits(float32(value))
		binary.LittleEndian.PutUint32(field.value[0:field.fixedLen], data)

	case Double:
		data := math.Float64bits(value)
		binary.LittleEndian.PutUint64(field.value[0:field.fixedLen], data)

	default:
		return invalidTypeError(field, `float`)
	}

	return nil
}

func (info *recordInfo) SetStringField(fieldName string, value string) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	switch field.Type {
	case String:
		clearNullFlag(field)
		valueBytes := []byte(value)
		size := int(field.fixedLen)
		copy(field.value, valueBytes)
		if size < int(field.fixedLen) {
			field.value[size] = 0
		}

	case V_String:
		valueBytes := []byte(value)
		field.varLen = len(valueBytes)
		if field.varLen >= len(field.value) {
			field.value = make([]byte, len(valueBytes)+20) // arbitrary padding to try and minimize memory allocation
		}
		copy(field.value, valueBytes)

	case WString:
		clearNullFlag(field)
		chars, err := syscall.UTF16FromString(value)
		if err != nil {
			return err
		}
		for index, char := range chars {
			if index*2 > int(field.fixedLen) {
				break
			}
			binary.LittleEndian.PutUint16(field.value[index*2:index*2+2], char)
		}

	case V_WString:
		chars, err := syscall.UTF16FromString(value)
		if err != nil {
			return err
		}
		requiredLen := len(chars) * 2
		if requiredLen >= len(field.value) {
			field.value = make([]byte, requiredLen+20) // arbitrary padding to try and minimize memory allocation
		}
		varLenIsUnset := true
		size := len(field.value) / 2
		for index := 0; index < size; index++ {
			if index >= len(chars) {
				if varLenIsUnset {
					field.varLen = index * 2
					varLenIsUnset = false
				}
				field.value[index*2] = 0
				field.value[index*2+1] = 0
			} else {
				binary.LittleEndian.PutUint16(field.value[index*2:index*2+2], chars[index])
			}
		}

	default:
		return invalidTypeError(field, `string`)
	}

	return nil
}

func (info *recordInfo) SetDateField(fieldName string, value time.Time) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	clearNullFlag(field)
	var valueStr string
	switch field.Type {
	case Date:
		valueStr = value.Format(dateFormat)
	case DateTime:
		valueStr = value.Format(dateTimeFormat)
	default:
		return invalidTypeError(field, `date`)
	}

	copy(field.value, valueStr)
	return nil
}

func (info *recordInfo) SetFieldNull(fieldName string) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return nil
	}
	field.varLen = 0
	field.isNull = true
	return nil
}

func (info *recordInfo) SetFromRawBytes(fieldName string, value []byte) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	switch field.Type {
	case V_String, V_WString:
		if value == nil {
			field.isNull = true
			return nil
		}
		valueLen := len(value)
		if valueLen >= len(field.value) {
			field.value = make([]byte, valueLen+20)
		}
		field.varLen = valueLen
		copy(field.value, value)

	default:
		copy(field.value, value)
	}
	field.isNull = false
	return nil
}

func clearNullFlag(field *fieldInfoEditor) {
	nullByteLocation := field.fixedLen - 1
	field.value[nullByteLocation] = 0
	field.isNull = false
}
