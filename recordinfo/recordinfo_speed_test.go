package recordinfo_test

import (
	"encoding/json"
	"goalteryx/recordinfo"
	"strings"
	"testing"
	"time"
)

func TestSpeed(t *testing.T) {
	sourceInfo, _ := recordinfo.FromXml(recordInfoXml)
	_ = sourceInfo.SetByteField(`ByteField`, byte(1))
	_ = sourceInfo.SetBoolField(`BoolField`, true)
	_ = sourceInfo.SetInt16Field(`Int16Field`, int16(16))
	_ = sourceInfo.SetInt32Field(`Int32Field`, int32(32))
	_ = sourceInfo.SetInt64Field(`Int64Field`, int64(64))
	_ = sourceInfo.SetFixedDecimalField(`FixedDecimalField`, 123.45)
	_ = sourceInfo.SetFloatField(`FloatField`, float32(678.9))
	_ = sourceInfo.SetDoubleField(`DoubleField`, 0.12345)
	_ = sourceInfo.SetStringField(`StringField`, `A`)
	_ = sourceInfo.SetWStringField(`WStringField`, `AB`)
	_ = sourceInfo.SetV_StringField(`V_StringShortField`, `ABC`)
	_ = sourceInfo.SetV_StringField(`V_StringLongField`, strings.Repeat(`B`, 500))
	_ = sourceInfo.SetV_WStringField(`V_WStringShortField`, `XYZ`)
	_ = sourceInfo.SetV_WStringField(`V_WStringLongField`, strings.Repeat(`W`, 500))
	_ = sourceInfo.SetDateField(`DateField`, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
	_ = sourceInfo.SetDateTimeField(`DateTimeField`, time.Date(2020, 2, 3, 4, 5, 6, 0, time.UTC))

	record, _ := sourceInfo.GenerateRecord()

	info, _ := recordinfo.FromXml(recordInfoXml)

	fieldCount := info.NumFields()
	results := make(map[string]int64, fieldCount)
	for index := 0; index < fieldCount; index++ {
		field, _ := info.GetFieldByIndex(index)
		start := time.Now()
		for i := 0; i < 100000; i++ {
			value, isNull, _ := sourceInfo.GetInterfaceValueFrom(field.Name, record)
			if isNull {
				_ = info.SetFieldNull(field.Name)
			} else {
				_ = info.SetFromInterface(field.Name, value)
			}
		}
		//_, _ = info.GenerateRecord()
		end := time.Now()
		duration := end.Sub(start)
		results[field.Name] = duration.Milliseconds()
	}
	resultJson, _ := json.MarshalIndent(results, "", "  ")
	t.Logf(string(resultJson))
}
