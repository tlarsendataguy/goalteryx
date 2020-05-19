package recordinfo

import (
	"time"
)

var ByteType = `Byte`
var BoolType = `Bool`
var Int16Type = `Int16`
var Int32Type = `Int32`
var Int64Type = `Int64`
var FixedDecimalType = `FixedDecimal`
var FloatType = `Float`
var DoubleType = `Double`
var StringType = `String`
var WStringType = `WString`
var V_StringType = `V_String`
var V_WStringType = `V_WString`
var DateType = `Date`
var DateTimeType = `DateTime`

var zeroDate = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

type generateBytes func(field *fieldInfoEditor, blob []byte, startAt int) (int, error)

var dateFormat = `2006-01-02`
var dateTimeFormat = `2006-01-02 15:04:05`
