package recordinfo

import (
	"time"
)

type FieldType byte

const (
	Invalid      FieldType = 0
	Bool         FieldType = 1
	Byte         FieldType = 2
	Int16        FieldType = 3
	Int32        FieldType = 4
	Int64        FieldType = 5
	FixedDecimal FieldType = 6
	Float        FieldType = 7
	Double       FieldType = 8
	String       FieldType = 9
	WString      FieldType = 10
	V_String     FieldType = 11
	V_WString    FieldType = 12
	Date         FieldType = 13
	DateTime     FieldType = 14

	byteType         = `Byte`
	boolType         = `Bool`
	int16Type        = `Int16`
	int32Type        = `Int32`
	int64Type        = `Int64`
	fixedDecimalType = `FixedDecimal`
	floatType        = `Float`
	doubleType       = `Double`
	stringType       = `String`
	wStringType      = `WString`
	v_StringType     = `V_String`
	v_WStringType    = `V_WString`
	dateType         = `Date`
	dateTimeType     = `DateTime`
)

var fieldTypeMap = []string{
	Invalid:      `invalid`,
	Bool:         boolType,
	Byte:         byteType,
	Int16:        int16Type,
	Int32:        int32Type,
	Int64:        int64Type,
	FixedDecimal: fixedDecimalType,
	Float:        floatType,
	Double:       doubleType,
	String:       stringType,
	WString:      wStringType,
	V_String:     v_StringType,
	V_WString:    v_WStringType,
	Date:         dateType,
	DateTime:     dateTimeType,
}

var zeroDate = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

const dateFormat = `2006-01-02`
const dateTimeFormat = `2006-01-02 15:04:05`
