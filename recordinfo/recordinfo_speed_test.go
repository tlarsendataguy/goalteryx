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
	_ = sourceInfo.SetIntField(`ByteField`, 1)
	_ = sourceInfo.SetBoolField(`BoolField`, true)
	_ = sourceInfo.SetIntField(`Int16Field`, 1)
	_ = sourceInfo.SetIntField(`Int32Field`, 32)
	_ = sourceInfo.SetIntField(`Int64Field`, 64)
	_ = sourceInfo.SetFloatField(`FixedDecimalField`, 123.45)
	_ = sourceInfo.SetFloatField(`FloatField`, 678.9)
	_ = sourceInfo.SetFloatField(`DoubleField`, 0.12345)
	_ = sourceInfo.SetStringField(`StringField`, `A`)
	_ = sourceInfo.SetStringField(`WStringField`, `AB`)
	_ = sourceInfo.SetStringField(`V_StringShortField`, `ABC`)
	_ = sourceInfo.SetStringField(`V_StringLongField`, strings.Repeat(`B`, 500))
	_ = sourceInfo.SetStringField(`V_WStringShortField`, `XYZ`)
	_ = sourceInfo.SetStringField(`V_WStringLongField`, strings.Repeat(`W`, 500))
	_ = sourceInfo.SetDateField(`DateField`, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
	_ = sourceInfo.SetDateField(`DateTimeField`, time.Date(2020, 2, 3, 4, 5, 6, 0, time.UTC))

	record, _ := sourceInfo.GenerateRecord()

	info, _ := recordinfo.FromXml(recordInfoXml)

	fieldCount := info.NumFields()
	results := make(map[string]int64, fieldCount)
	for index := 0; index < fieldCount; index++ {
		field, _ := info.GetFieldByIndex(index)
		start := time.Now()
		for i := 0; i < 100000; i++ {
			value, _ := sourceInfo.GetRawBytesFrom(field.Name, record)
			_ = info.SetFromRawBytes(field.Name, value)
		}
		//_, _ = info.GenerateRecord()
		end := time.Now()
		duration := end.Sub(start)
		results[field.Name] = duration.Milliseconds()
	}
	resultJson, _ := json.MarshalIndent(results, "", "  ")
	t.Logf(string(resultJson))
}
