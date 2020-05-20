package recordinfo

import (
	"encoding/binary"
	"fmt"
	"goalteryx/convert_strings"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

func (info *recordInfo) GetByteValueFrom(fieldName string, record unsafe.Pointer) (value byte, isNull bool, err error) {
	returnEarly, isNull, err, field := info.shouldReturnEarlyWith(fieldName, record)
	if returnEarly {
		return byte(0), isNull, err
	}
	return *((*byte)(unsafe.Pointer(uintptr(record) + field.location))), false, nil
}

func (info *recordInfo) GetBoolValueFrom(fieldName string, record unsafe.Pointer) (value bool, isNull bool, err error) {
	returnEarly, isNull, err, field := info.shouldReturnEarlyWith(fieldName, record)
	if returnEarly {
		return false, isNull, err
	}
	return *((*bool)(unsafe.Pointer(uintptr(record) + field.location))), false, nil
}

func (info *recordInfo) GetInt16ValueFrom(fieldName string, record unsafe.Pointer) (value int16, isNull bool, err error) {
	returnEarly, isNull, err, field := info.shouldReturnEarlyWith(fieldName, record)
	if returnEarly {
		return 0, isNull, err
	}
	return *((*int16)(unsafe.Pointer(uintptr(record) + field.location))), false, nil
}

func (info *recordInfo) GetInt32ValueFrom(fieldName string, record unsafe.Pointer) (value int32, isNull bool, err error) {
	returnEarly, isNull, err, field := info.shouldReturnEarlyWith(fieldName, record)
	if returnEarly {
		return 0, isNull, err
	}
	return *((*int32)(unsafe.Pointer(uintptr(record) + field.location))), false, nil
}

func (info *recordInfo) GetInt64ValueFrom(fieldName string, record unsafe.Pointer) (value int64, isNull bool, err error) {
	returnEarly, isNull, err, field := info.shouldReturnEarlyWith(fieldName, record)
	if returnEarly {
		return 0, isNull, err
	}
	return *((*int64)(unsafe.Pointer(uintptr(record) + field.location))), false, nil
}

func (info *recordInfo) GetFixedDecimalValueFrom(fieldName string, record unsafe.Pointer) (value float64, isNull bool, err error) {
	returnEarly, isNull, err, field := info.shouldReturnEarlyWith(fieldName, record)
	if returnEarly {
		return 0, isNull, err
	}
	valueStr := convert_strings.CToString(unsafe.Pointer(uintptr(record) + field.location))
	value, err = strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, false, fmt.Errorf(`error converting '%v' to double in '%v' field`, value, fieldName)
	}
	return value, false, nil
}

func (info *recordInfo) GetFloatValueFrom(fieldName string, record unsafe.Pointer) (value float32, isNull bool, err error) {
	returnEarly, isNull, err, field := info.shouldReturnEarlyWith(fieldName, record)
	if returnEarly {
		return 0, isNull, err
	}
	return *((*float32)(unsafe.Pointer(uintptr(record) + field.location))), false, nil
}

func (info *recordInfo) GetDoubleValueFrom(fieldName string, record unsafe.Pointer) (value float64, isNull bool, err error) {
	returnEarly, isNull, err, field := info.shouldReturnEarlyWith(fieldName, record)
	if returnEarly {
		return 0, isNull, err
	}
	return *((*float64)(unsafe.Pointer(uintptr(record) + field.location))), false, nil
}

func (info *recordInfo) GetStringValueFrom(fieldName string, record unsafe.Pointer) (value string, isNull bool, err error) {
	returnEarly, isNull, err, field := info.shouldReturnEarlyWith(fieldName, record)
	if returnEarly {
		return ``, isNull, err
	}
	return convert_strings.CToString(unsafe.Pointer(uintptr(record) + field.location)), false, nil
}

func (info *recordInfo) GetWStringValueFrom(fieldName string, record unsafe.Pointer) (value string, isNull bool, err error) {
	returnEarly, isNull, err, field := info.shouldReturnEarlyWith(fieldName, record)
	if returnEarly {
		return ``, isNull, err
	}
	return convert_strings.WideCToString(unsafe.Pointer(uintptr(record) + field.location)), false, nil
}

func (info *recordInfo) GetV_StringValueFrom(fieldName string, record unsafe.Pointer) (value string, isNull bool, err error) {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return ``, false, err
	}

	if isVarFieldNull(field, record) {
		return ``, true, nil
	}

	varBytes := getVarBytes(field, record)
	return string(varBytes), false, nil
}

func isVarFieldNull(field *fieldInfoEditor, record unsafe.Pointer) bool {
	if *((*int32)(unsafe.Pointer(uintptr(record) + field.location))) == 1 {
		return true
	}
	return false
}

func getVarBytes(field *fieldInfoEditor, record unsafe.Pointer) []byte {
	varStart := *((*uint32)(unsafe.Pointer(uintptr(record) + field.location)))
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

	varBytes := make([]byte, varLen)
	var index uint32
	for index = 0; index < varLen; index++ {
		varBytes[index] = *((*byte)(unsafe.Pointer(uintptr(record) + offset + uintptr(index))))
	}
	return varBytes
}

func (info *recordInfo) GetV_WStringValueFrom(fieldName string, record unsafe.Pointer) (value string, isNull bool, err error) {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return ``, false, err
	}

	if isVarFieldNull(field, record) {
		return ``, true, nil
	}

	varBytes := getVarBytes(field, record)
	varUint16 := make([]uint16, len(varBytes)/2)
	varUint16Index := 0
	for index := 0; index < len(varBytes); index += 2 {
		varUint16[varUint16Index] = binary.LittleEndian.Uint16(varBytes[index : index+2])
		varUint16Index++
	}
	return syscall.UTF16ToString(varUint16), false, nil
}

func (info *recordInfo) GetDateValueFrom(fieldName string, record unsafe.Pointer) (value time.Time, isNull bool, err error) {
	returnEarly, isNull, err, field := info.shouldReturnEarlyWith(fieldName, record)
	if returnEarly {
		return zeroDate, isNull, err
	}
	dateStr := convert_strings.CToString(unsafe.Pointer(uintptr(record) + field.location))
	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return zeroDate, false, fmt.Errorf(`error converting date '%v' in GetDateValueFrom for field [%v], use format yyyy-MM-dd`, dateStr, fieldName)
	}
	return date, false, nil
}

func (info *recordInfo) GetDateTimeValueFrom(fieldName string, record unsafe.Pointer) (value time.Time, isNull bool, err error) {
	returnEarly, isNull, err, field := info.shouldReturnEarlyWith(fieldName, record)
	if returnEarly {
		return zeroDate, isNull, err
	}
	dateStr := convert_strings.CToString(unsafe.Pointer(uintptr(record) + field.location))
	date, err := time.Parse(dateTimeFormat, dateStr)
	if err != nil {
		return zeroDate, false, fmt.Errorf(`error converting datetime '%v' in GetDateValueFrom for field [%v], use format yyyy-MM-dd hh:mm:ss`, dateStr, fieldName)
	}
	return date, false, nil
}

func (info *recordInfo) GetInterfaceValueFrom(fieldName string, record unsafe.Pointer) (value interface{}, isNull bool, err error) {
	returnEarly, isNull, err, field := info.shouldReturnEarlyWith(fieldName, record)
	if returnEarly {
		return nil, isNull, err
	}
	switch field.Type {
	case ByteType:
		return info.GetByteValueFrom(fieldName, record)
	case BoolType:
		return info.GetBoolValueFrom(fieldName, record)
	case Int16Type:
		return info.GetInt16ValueFrom(fieldName, record)
	case Int32Type:
		return info.GetInt32ValueFrom(fieldName, record)
	case Int64Type:
		return info.GetInt64ValueFrom(fieldName, record)
	case FixedDecimalType:
		return info.GetFixedDecimalValueFrom(fieldName, record)
	case FloatType:
		return info.GetFloatValueFrom(fieldName, record)
	case DoubleType:
		return info.GetDoubleValueFrom(fieldName, record)
	case StringType:
		return info.GetStringValueFrom(fieldName, record)
	case WStringType:
		return info.GetWStringValueFrom(fieldName, record)
	case DateType:
		return info.GetDateValueFrom(fieldName, record)
	case DateTimeType:
		return info.GetDateTimeValueFrom(fieldName, record)
	default:
		return nil, false, fmt.Errorf(`field [%v] has invalid type '%v'`, field.Name, field.Type)
	}
}

func (info *recordInfo) shouldReturnEarlyWith(fieldName string, record unsafe.Pointer) (returnEarly bool, isNull bool, err error, field *fieldInfoEditor) {
	field, err = info.getFieldInfo(fieldName)
	if err != nil {
		return true, false, err, nil
	}
	if isValueNull(field, record) {
		return true, true, nil, field
	}
	return false, false, nil, field
}

var nullByteTypes = []string{
	ByteType,
	Int16Type,
	Int32Type,
	Int64Type,
	FixedDecimalType,
	FloatType,
	DoubleType,
	StringType,
	WStringType,
	DateType,
	DateTimeType,
}

func isValueNull(field *fieldInfoEditor, record unsafe.Pointer) bool {
	for _, nullByteType := range nullByteTypes {
		if nullByteType == field.Type {
			nullByte := *((*byte)(unsafe.Pointer(uintptr(record) + field.location + field.fixedLen)))
			return nullByte == byte(1)
		}
	}
	if field.Type == BoolType {
		nullByte := *((*byte)(unsafe.Pointer(uintptr(record) + field.location)))
		return nullByte == byte(2)
	}
	return false
}
