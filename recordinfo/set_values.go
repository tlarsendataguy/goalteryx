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
	case ByteType:
		field.value[0] = byte(value)
	case Int16Type:
		binary.LittleEndian.PutUint16(field.value[0:2], uint16(value))
	case Int32Type:
		binary.LittleEndian.PutUint32(field.value[0:4], uint32(value))
	case Int64Type:
		binary.LittleEndian.PutUint64(field.value[0:8], uint64(value))
	default:
		return fmt.Errorf(`[%v]'s type of '%v' is not a valid int type`, field.Name, field.Type)
	}
	return nil
}

func (info *recordInfo) SetBoolField(fieldName string, value bool) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	if field.Type != BoolType {
		return fmt.Errorf(`[%v]'s type of '%v' is not a valid bool type`, field.Name, field.Type)
	}

	if value {
		field.value[0] = 1
	} else {
		field.value[0] = 0
	}
	return nil
}

func (info *recordInfo) SetFloatField(fieldName string, value float64) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	clearNullFlag(field)

	switch field.Type {
	case FixedDecimalType:
		format := `%` + fmt.Sprintf(`%v.%vf`, field.Size, field.Precision)
		valueStr := []byte(strings.TrimSpace(fmt.Sprintf(format, value)))
		size := int(field.fixedLen)
		for index := 0; index < size; index++ {
			if index >= len(valueStr) {
				field.value[index] = 0
			} else {
				field.value[index] = valueStr[index]
			}
		}

	case FloatType:
		data := math.Float32bits(float32(value))
		binary.LittleEndian.PutUint32(field.value[0:field.fixedLen], data)

	case DoubleType:
		data := math.Float64bits(value)
		binary.LittleEndian.PutUint64(field.value[0:field.fixedLen], data)

	default:
		return fmt.Errorf(`[%v]'s type of '%v' is not a valid float type`, field.Name, field.Type)
	}

	return nil
}

func (info *recordInfo) SetStringField(fieldName string, value string) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	switch field.Type {
	case StringType:
		valueBytes := []byte(value)
		size := int(field.fixedLen)
		for index := 0; index < size; index++ {
			if index >= len(valueBytes) {
				field.value[index] = 0
			} else {
				field.value[index] = valueBytes[index]
			}
		}

	case V_StringType:
		valueBytes := []byte(value)
		if len(valueBytes) >= len(field.value) {
			field.value = make([]byte, len(valueBytes)+20) // arbitrary padding to try and minimize memory allocation
		}
		varLenIsUnset := true
		size := len(field.value)
		for index := 0; index < size; index++ {
			if index >= len(valueBytes) {
				if varLenIsUnset {
					field.varLen = index
					varLenIsUnset = false
				}
				field.value[index] = 0
			} else {
				field.value[index] = valueBytes[index]
			}
		}

	case WStringType:
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

	case V_WStringType:
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
		return fmt.Errorf(`[%v]'s type of '%v' is not a valid string type`, field.Name, field.Type)
	}

	return nil
}

func (info *recordInfo) SetDateField(fieldName string, value time.Time) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	var valueStr string
	switch field.Type {
	case DateType:
		valueStr = value.Format(dateFormat)
	case DateTimeType:
		valueStr = value.Format(dateTimeFormat)
	default:
		return fmt.Errorf(`[%v]'s type of '%v' is not a valid date type`, field.Name, field.Type)
	}

	size := int(field.fixedLen)
	for index := 0; index < size; index++ {
		if index >= len(valueStr) {
			field.value[index] = 0
		} else {
			field.value[index] = valueStr[index]
		}
	}

	return nil
}

func (info *recordInfo) SetFieldNull(fieldName string) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return nil
	}
	field.value = nil
	return nil
}

func (info *recordInfo) SetFromRawBytes(fieldName string, value []byte) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	switch field.Type {
	case V_StringType, V_WStringType:
		if len(value) >= len(field.value) {
			field.value = make([]byte, len(value)+20)
		}
		field.varLen = len(value)
		for index := 0; index < len(field.value); index++ {
			if index >= len(value) {
				field.value[index] = 0
			} else {
				field.value[index] = value[index]
			}
		}
	default:
		for index := 0; index < int(field.fixedLen); index++ {
			if index >= len(value) {
				field.value[index] = 0
			} else {
				field.value[index] = value[index]
			}
		}
	}
	return nil
}

func clearNullFlag(field *fieldInfoEditor) {
	nullByteLocation := len(field.value) - 1
	field.value[nullByteLocation] = 0
}
