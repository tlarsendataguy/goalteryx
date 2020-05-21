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
		field.fixedValue[0] = byte(value)
	case Int16Type:
		binary.LittleEndian.PutUint16(field.fixedValue[0:2], uint16(value))
	case Int32Type:
		binary.LittleEndian.PutUint32(field.fixedValue[0:4], uint32(value))
	case Int64Type:
		binary.LittleEndian.PutUint64(field.fixedValue[0:8], uint64(value))
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
		field.fixedValue[0] = 1
	} else {
		field.fixedValue[0] = 0
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
				field.fixedValue[index] = 0
			} else {
				field.fixedValue[index] = valueStr[index]
			}
		}

	case FloatType:
		data := math.Float32bits(float32(value))
		binary.LittleEndian.PutUint32(field.fixedValue[0:field.fixedLen], data)

	case DoubleType:
		data := math.Float64bits(value)
		binary.LittleEndian.PutUint64(field.fixedValue[0:field.fixedLen], data)

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
		valueStr := []byte(value)
		size := int(field.fixedLen)
		for index := 0; index < size; index++ {
			if index >= len(valueStr) {
				field.fixedValue[index] = 0
			} else {
				field.fixedValue[index] = valueStr[index]
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
			binary.LittleEndian.PutUint16(field.fixedValue[index*2:index*2+2], char)
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
			field.fixedValue[index] = 0
		} else {
			field.fixedValue[index] = valueStr[index]
		}
	}

	return nil
}

func (info *recordInfo) SetFieldNull(fieldName string) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return nil
	}
	field.fixedValue = nil
	return nil
}

func (info *recordInfo) SetFromRawBytes(fieldName string, value []byte) error {
	_, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	return nil
}

func clearNullFlag(field *fieldInfoEditor) {
	nullByteLocation := len(field.fixedValue) - 1
	field.fixedValue[nullByteLocation] = 0
}
