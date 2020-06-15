package recordinfo

import (
	"encoding/binary"
	"fmt"
	"github.com/tlarsen7572/goalteryx/recordblob"
	"math"
	"reflect"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

// GetIntValueFrom retrieves integers from integer fields.  Each type of integer field uses a different fixed
// length, and so we must treat each separately.  The storage size for each integer field is the number of bytes
// needed to store each integer, plus 1.  The last byte is used as a null flag: 0 means the field has a value and 1
// means the field is null.
func (info *recordInfo) GetIntValueFrom(fieldName string, record *recordblob.RecordBlob) (value int, isNull bool, err error) {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return 0, false, err
	}

	switch field.Type {
	case Byte:
		raw := *((*[2]byte)(unsafe.Pointer(uintptr(record.Blob) + field.location)))
		if raw[1] == 1 {
			return 0, true, nil
		}
		return int(raw[0]), false, nil

	case Int16:
		raw := *((*[3]byte)(unsafe.Pointer(uintptr(record.Blob) + field.location)))
		if raw[2] == 1 {
			return 0, true, nil
		}
		return int(binary.LittleEndian.Uint16(raw[:2])), false, nil

	case Int32:
		raw := *((*[5]byte)(unsafe.Pointer(uintptr(record.Blob) + field.location)))
		if raw[4] == 1 {
			return 0, true, nil
		}
		return int(binary.LittleEndian.Uint32(raw[:4])), false, nil

	case Int64:
		raw := *((*[9]byte)(unsafe.Pointer(uintptr(record.Blob) + field.location)))
		if raw[8] == 1 {
			return 0, true, nil
		}
		return int(binary.LittleEndian.Uint64(raw[:8])), false, nil

	default:
		return 0, false, invalidTypeError(field, `int`)
	}
}

// GetBoolValueFrom extracts a boolean value from a boolean field.  Bool fields are the only fields without
// a byte for the null flag.  Bool fields can either be 0 (false), 1 (true), or 2 (null).
func (info *recordInfo) GetBoolValueFrom(fieldName string, record *recordblob.RecordBlob) (value bool, isNull bool, err error) {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return false, false, err
	}

	if field.Type != Bool {
		return false, false, invalidTypeError(field, `bool`)
	}

	raw := *((*byte)(unsafe.Pointer(uintptr(record.Blob) + field.location)))
	if raw == 2 {
		return false, true, nil
	}
	if raw == 0 {
		return false, false, nil
	}
	return true, false, nil
}

// GetFloatValueFrom retrieves float values from decimal fields.  The Float and Double fields are fixed-size fields
// (4 bytes and 8 bytes, respectively).  The size of FixedDecimal fields is specified in the field definition.  All
// decimal fields have an additional byte at the end for a null flag: 0 means the field has a value and 1
// means the field is null.
func (info *recordInfo) GetFloatValueFrom(fieldName string, record *recordblob.RecordBlob) (value float64, isNull bool, err error) {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return 0, false, err
	}

	switch field.Type {
	case Float:
		raw := *((*[5]byte)(unsafe.Pointer(uintptr(record.Blob) + field.location)))
		if raw[4] == 1 {
			return 0, true, nil
		}
		return float64(math.Float32frombits(binary.LittleEndian.Uint32(raw[:4]))), false, nil
	case Double:
		raw := *((*[9]byte)(unsafe.Pointer(uintptr(record.Blob) + field.location)))
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

// GetStringValueFrom extracts text values from string fields.  String and WString fields are fixed-length and
// contain an extra byte for a null flag: 0 means the field has a value and 1 means the field is null.
//
// V_String and V_WString fields are variable-length fields and are more complicated to retrieve.  4 bytes are
// stored in the fixed portion of the record blob.  If the 4 bytes are an integer equal to 0, the value is a
// zero-length string.  If the 4 bytes are an integer equal to 1, the field is null.  Otherwise, the 4 bytes
// are an integer telling how many bytes you must skip until you reach the actual variable-length data.  The
// variable-length data itself contains an integer (1 or 4 bytes long, depending on the size of the text) that
// describes how long the text is.  Following this second integer are the actual bytes that make up the text value.
//
// String and V_String are narrow strings (such as ASCII) whereas WString and V_WString are wide strings.  Wide
// strings are encoded in little-endian UTF16.
func (info *recordInfo) GetStringValueFrom(fieldName string, record *recordblob.RecordBlob) (value string, isNull bool, err error) {
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

// GetDateValueFrom extracts date/datetime values from date/datetime fields.  Date fields are 10 bytes long,
// representing a date string formatted as yyyy-MM-dd.  DateTime fields are 19 bytes long, representing a datetime
// string formatted as yyyy-MM-dd hh:mm:ss.  There is an extra byte at the end of both types of fields for a
// null flag: 0 means the field has a value and 1 means the field is null.
func (info *recordInfo) GetDateValueFrom(fieldName string, record *recordblob.RecordBlob) (value time.Time, isNull bool, err error) {
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

// GetRawBytesFrom extracts the raw bytes from any field.  It is a universal value getter as it can be
// used on any field.  For fixed-length fields the return value contains the bytes from the fixed portion
// of the record blob, including the trailing byte used for the null flag.  For variable-length fields, the
// return value is either nil (if the field is null), a zero-length byte array (if the field is empty), or
// the bytes containing the actual stored data.
func (info *recordInfo) GetRawBytesFrom(fieldName string, record *recordblob.RecordBlob) (value []byte, err error) {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return nil, fmt.Errorf(`error getting raw bytes: %v`, err.Error())
	}
	return info.getRawBytes(field, record), nil
}

// GetRawBytesFrom extracts the raw bytes from any field position.  It is a universal value getter as it can be
// used on any field.  For fixed-length fields the return value contains the bytes from the fixed portion
// of the record blob, including the trailing byte used for the null flag.  For variable-length fields, the
// return value is either nil (if the field is null), a zero-length byte array (if the field is empty), or
// the bytes containing the actual stored data.
func (info *recordInfo) GetRawBytesFromIndex(index int, record *recordblob.RecordBlob) (value []byte, err error) {
	if index < 0 || index > info.numFields {
		return nil, fmt.Errorf(`error getting raw bytes: index was not between 0 and %v`, info.numFields)
	}
	return info.getRawBytes(info.fields[index], record), nil
}

// getRawBytes returns the raw bytes of a field.  For fixed-length fields the return value contains the bytes from the fixed portion
// of the record blob, including the trailing byte used for the null flag.  For variable-length fields, the
// return value is either nil (if the field is null), a zero-length byte array (if the field is empty), or
// the bytes containing the actual stored data.
func (info *recordInfo) getRawBytes(field *fieldInfoEditor, record *recordblob.RecordBlob) []byte {
	switch field.Type {
	case V_String, V_WString, Blob, Spatial:
		return getVarBytes(field, record)
	default:
		return getFixedBytes(field, record)
	}
}

// getFixedBytes gets the bytes for a field from the fixed portion of the record blob.  This includes any trailing
// null flag byte.
func getFixedBytes(field *fieldInfoEditor, record *recordblob.RecordBlob) []byte {
	totalLen := int(field.fixedLen + field.nullByteLen)
	var raw []byte
	rawHeader := (*reflect.SliceHeader)(unsafe.Pointer(&raw))
	rawHeader.Data = uintptr(record.Blob) + field.location
	rawHeader.Len = totalLen
	rawHeader.Cap = totalLen
	return raw
}

// getVarBytes gets the data bytes of variable-length fields.  The return value is either nil (if the field is null),
// a zero-length byte array (if the field is empty), or the bytes containing the actual stored data.
func getVarBytes(field *fieldInfoEditor, record *recordblob.RecordBlob) []byte {
	varStart := *((*uint32)(unsafe.Pointer(uintptr(record.Blob) + field.location)))
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

		varLenFirstByte := *((*byte)(unsafe.Pointer(uintptr(record.Blob) + field.location + uintptr(varStart))))
		offset += uintptr(varStart)
		if varLenFirstByte&byte(1) == 1 {
			varLen = uint32(varLenFirstByte >> 1)
			offset += 1
		} else {
			varLen = *((*uint32)(unsafe.Pointer(uintptr(record.Blob) + field.location + uintptr(varStart)))) / 2
			offset += 4
		}
	}

	var raw []byte
	rawHeader := (*reflect.SliceHeader)(unsafe.Pointer(&raw))
	rawHeader.Data = uintptr(record.Blob) + offset
	rawHeader.Len = int(varLen)
	rawHeader.Cap = int(varLen)
	return raw
}

// truncateAtNullByte is used to truncate bytes after the end of a narrow string.  For example, if a String field has
// a length of 10 but the value stored is 'ABC', there will be 6 garbage bytes at the end of the string.  Strings
// that are shorter than their field length will end with a null byte, so when we extract the raw bytes of certain
// fields, we stop when we reach a null byte so the return value does not get garbled.
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
