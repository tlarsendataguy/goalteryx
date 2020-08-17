package recordinfo

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"syscall"
	"time"
)

// SetIntField sets an integer field with a value.  The bytes representing the integer are copied into the
// field's buffer for later use by GenerateRecord.  The bytes are stored in little endian order.  Trying to set
// an integer into a non-integer field returns an error.
func (info *recordInfo) SetIntField(fieldName string, value int) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	clearNullFlagFromFixedField(field)
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

// SetBoolField sets a bool field with true/false.  The bytes representing the bool are copied into the
// field's buffer for later use by GenerateRecord.  Trying to set an integer into a non-integer field returns
// an error.
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

// SetFloatField sets a decimal field with a number.  The bytes representing the number are copied into the
// field's buffer for later use by GenerateRecord.  The bytes are stored in little endian order.  Trying to set
// a decimal into a non-decimal field returns an error.
func (info *recordInfo) SetFloatField(fieldName string, value float64) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	clearNullFlagFromFixedField(field)

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

// SetStringField sets an text field with a value.  The bytes representing the text are copied into the
// field's buffer for later use by GenerateRecord.  Trying to set a string into a non-string field returns an error.
//
// WString and V_WString fields are first converted to UTF16 encoding before their bytes are generated and saved.
func (info *recordInfo) SetStringField(fieldName string, value string) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	switch field.Type {
	case String:
		clearNullFlagFromFixedField(field)
		valueBytes := []byte(value)
		size := len(valueBytes)
		copy(field.value[:field.Size], valueBytes)
		if size < int(field.fixedLen) {
			field.value[size] = 0
		}

	case V_String:
		field.isNull = false
		valueBytes := []byte(value)
		actualLen := len(valueBytes)
		maxLen := field.Size
		if actualLen > maxLen {
			actualLen = maxLen
		}

		field.varLen = actualLen
		if actualLen > len(field.value) {
			field.value = make([]byte, actualLen)
		}
		copy(field.value, valueBytes[:actualLen])

	case WString:
		clearNullFlagFromFixedField(field)
		chars, err := syscall.UTF16FromString(value)
		if err != nil {
			return err
		}
		for index, char := range chars {
			if index >= field.Size {
				break
			}
			binary.LittleEndian.PutUint16(field.value[index*2:index*2+2], char)
		}

	case V_WString:
		field.isNull = false
		chars, err := syscall.UTF16FromString(value)
		if err != nil {
			return err
		}

		charsLen := len(chars)
		if charsLen > field.Size {
			chars = chars[:field.Size]
			charsLen = field.Size
		}

		requiredLen := charsLen * 2
		fieldValueLen := len(field.value)

		if requiredLen > fieldValueLen {
			field.value = make([]byte, requiredLen)
		}
		field.varLen = requiredLen
		for index := 0; index < charsLen; index++ {
			binary.LittleEndian.PutUint16(field.value[index*2:index*2+2], chars[index])
		}

	default:
		return invalidTypeError(field, `string`)
	}

	return nil
}

// SetDateField sets a date/datetime field with a value.  The bytes representing the value are copied into the
// field's buffer for later use by GenerateRecord.  Trying to set a date into a non-date field returns an error.
func (info *recordInfo) SetDateField(fieldName string, value time.Time) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	clearNullFlagFromFixedField(field)
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

// SetFieldNull sets the null flag on any field type.
func (info *recordInfo) SetFieldNull(fieldName string) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return nil
	}
	field.varLen = 0
	field.isNull = true
	return nil
}

// SetFromRawBytes sets a field value to the specified bytes.  No validation is performed on the bytes; it is up
// to the caller to provide a valid byte slice for the field.  This is a fast way to set field values, but also
// more dangerous.
func (info *recordInfo) SetFromRawBytes(fieldName string, value []byte) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}
	return info.setFieldFromRawBytes(field, value)
}

// SetIndexFromRawBytes sets a field value to the specified bytes.  No validation is performed on the bytes; it is up
// to the caller to provide a valid byte slice for the field.  This is the fastest way to set field values, but also
// more dangerous.
func (info *recordInfo) SetIndexFromRawBytes(index int, value []byte) error {
	if index < 0 || index > info.numFields {
		return fmt.Errorf(`error setting raw bytes: index was not between 0 and %v`, info.numFields)
	}
	field := info.fields[index]
	return info.setFieldFromRawBytes(field, value)
}

// setFieldFromRawBytes performs the actual work of copying the provided byte slice into the field's buffer.
// TODO: Enforce field length limitations on the incoming byte slice.  We can either throw an error if the length is wrong or truncate the slice
func (info *recordInfo) setFieldFromRawBytes(field *fieldInfoEditor, value []byte) error {
	switch field.Type {
	case V_String, V_WString, Blob, Spatial:
		if value == nil {
			field.isNull = true
			return nil
		}
		valueLen := len(value)
		if valueLen > len(field.value) {
			field.value = make([]byte, valueLen)
		}
		field.varLen = valueLen
		copy(field.value, value)
		field.isNull = false

	case Bool:
		copy(field.value, value)
		field.isNull = value[0] == 2

	default:
		copy(field.value, value)
		field.isNull = value[field.fixedLen] == 1
	}
	return nil
}

// clearNullFlagFromFixedField clears the null flag from any fixed-length, non-bool field type.
func clearNullFlagFromFixedField(field *fieldInfoEditor) {
	field.value[field.fixedLen] = 0
	field.isNull = false
}
