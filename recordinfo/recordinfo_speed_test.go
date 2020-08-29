package recordinfo_test

import (
	"encoding/json"
	"github.com/tlarsen7572/goalteryx/recordinfo"
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
	reader, _ := recordinfo.RecordBlobReaderFromXml(recordInfoXml)

	fieldCount := info.NumFields()
	totalResults := make(map[string]int64, fieldCount)
	doSet := true
	doGenerateRecord := true
	for index := 0; index < fieldCount; index++ {
		field, _ := info.GetFieldByIndex(index)
		start := time.Now()
		for i := 0; i < 100000; i++ {
			value, err := reader.GetRawBytesFrom(field.Name, record)
			if err != nil {
				t.Fatalf(`expected no error getting raw bytes from field %v, but got: %v`, field.Name, err.Error())
			}
			if doSet {
				err = info.SetFromRawBytes(field.Name, value)
				if err != nil {
					t.Fatalf(`expected no error setting raw bytes to field %v, but got: %v`, field.Name, err.Error())
				}
			}
			if doGenerateRecord {
				_, err = info.GenerateRecord()
				if err != nil {
					t.Fatalf(`expected no error setting raw bytes to field %v, but got: %v`, field.Name, err.Error())
				}
			}
		}
		end := time.Now()
		duration := end.Sub(start)
		totalResults[field.Name] = duration.Milliseconds()
	}
	resultJson, _ := json.MarshalIndent(totalResults, "", "  ")
	t.Logf(string(resultJson))
}
