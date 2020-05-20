package recordinfo_test

import (
	"goalteryx/recordinfo"
	"strings"
	"testing"
	"time"
)

func TestSetValuesAndGenerateRecord(t *testing.T) {
	recordInfo := generateTestRecordInfo()
	setRecordInfoTestData(recordInfo)

	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	byteVal, isNull, err := recordInfo.GetByteValueFrom(`ByteField`, record)
	checkExpectedGetValueFrom(t, byteVal, byte(1), isNull, false, err, nil)

	boolVal, isNull, err := recordInfo.GetBoolValueFrom(`BoolField`, record)
	checkExpectedGetValueFrom(t, boolVal, true, isNull, false, err, nil)

	int16Val, isNull, err := recordInfo.GetInt16ValueFrom(`Int16Field`, record)
	checkExpectedGetValueFrom(t, int16Val, int16(2), isNull, false, err, nil)

	int32Val, isNull, err := recordInfo.GetInt32ValueFrom(`Int32Field`, record)
	checkExpectedGetValueFrom(t, int32Val, int32(3), isNull, false, err, nil)

	int64Val, isNull, err := recordInfo.GetInt64ValueFrom(`Int64Field`, record)
	checkExpectedGetValueFrom(t, int64Val, int64(4), isNull, false, err, nil)

	fixedDecimalVal, isNull, err := recordInfo.GetFixedDecimalValueFrom(`FixedDecimalField`, record)
	checkExpectedGetValueFrom(t, fixedDecimalVal, 123.45, isNull, false, err, nil)

	floatVal, isNull, err := recordInfo.GetFloatValueFrom(`FloatField`, record)
	checkExpectedGetValueFrom(t, floatVal, float32(654.321), isNull, false, err, nil)

	doubleVal, isNull, err := recordInfo.GetDoubleValueFrom(`DoubleField`, record)
	checkExpectedGetValueFrom(t, doubleVal, 909.33, isNull, false, err, nil)

	stringVal, isNull, err := recordInfo.GetStringValueFrom(`StringField`, record)
	checkExpectedGetValueFrom(t, stringVal, `ABCDEFG`, isNull, false, err, nil)

	wstringVal, isNull, err := recordInfo.GetWStringValueFrom(`WStringField`, record)
	checkExpectedGetValueFrom(t, wstringVal, `CXVY`, isNull, false, err, nil)

	dateVal, isNull, err := recordInfo.GetDateValueFrom(`DateField`, record)
	expectedDate := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	checkExpectedGetValueFrom(t, dateVal, expectedDate, isNull, false, err, nil)

	dateTimeVal, isNull, err := recordInfo.GetDateTimeValueFrom(`DateTimeField`, record)
	expectedDate = time.Date(2021, 3, 4, 5, 6, 7, 0, time.UTC)
	checkExpectedGetValueFrom(t, dateTimeVal, expectedDate, isNull, false, err, nil)
}

func TestCachedRecords(t *testing.T) {
	recordInfo := generateTestRecordInfo()
	setRecordInfoTestData(recordInfo)

	record1, _ := recordInfo.GenerateRecord()
	record2, _ := recordInfo.GenerateRecord()
	if record1 != record2 {
		t.Fatalf(`record1 and record2 are 2 different pointers`)
	}
}

func TestSetLongVarDataFieldsAndGenerateRecord(t *testing.T) {
	recordInfo := recordinfo.New()
	recordInfo.AddByteField(`ByteField`, ``)
	recordInfo.AddV_WStringField(`V_WStringField`, ``, 250)
	recordInfo.AddV_StringField(`V_StringField`, ``, 250)

	_ = recordInfo.SetByteField(`ByteField`, 1)
	_ = recordInfo.SetV_StringField(`V_StringField`, strings.Repeat(`B`, 200))
	_ = recordInfo.SetV_WStringField(`V_WStringField`, strings.Repeat(`A`, 100))

	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull, err := recordInfo.GetV_StringValueFrom(`V_StringField`, record)
	checkExpectedGetValueFrom(t, value, strings.Repeat(`B`, 200), isNull, false, err, nil)

	value, isNull, err = recordInfo.GetV_WStringValueFrom(`V_WStringField`, record)
	checkExpectedGetValueFrom(t, value, strings.Repeat(`A`, 100), isNull, false, err, nil)
}

func TestSetShortVarDataFieldsAndGenerateRecord(t *testing.T) {
	recordInfo := recordinfo.New()
	recordInfo.AddByteField(`ByteField`, ``)
	recordInfo.AddV_WStringField(`V_WStringField`, ``, 250)
	recordInfo.AddV_StringField(`V_StringField`, ``, 250)

	_ = recordInfo.SetByteField(`ByteField`, 1)
	_ = recordInfo.SetV_StringField(`V_StringField`, strings.Repeat(`B`, 100))
	_ = recordInfo.SetV_WStringField(`V_WStringField`, strings.Repeat(`A`, 50))

	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull, err := recordInfo.GetV_StringValueFrom(`V_StringField`, record)
	checkExpectedGetValueFrom(t, value, strings.Repeat(`B`, 100), isNull, false, err, nil)

	value, isNull, err = recordInfo.GetV_WStringValueFrom(`V_WStringField`, record)
	checkExpectedGetValueFrom(t, value, strings.Repeat(`A`, 50), isNull, false, err, nil)
}

func TestSetTinyVarDataFieldsAndGenerateRecord(t *testing.T) {
	recordInfo := recordinfo.New()
	recordInfo.AddByteField(`ByteField`, ``)
	recordInfo.AddV_WStringField(`V_WStringField`, ``, 250)
	recordInfo.AddV_StringField(`V_StringField`, ``, 250)

	_ = recordInfo.SetByteField(`ByteField`, 1)
	_ = recordInfo.SetV_StringField(`V_StringField`, `B`)
	_ = recordInfo.SetV_WStringField(`V_WStringField`, `A`)

	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	value, isNull, err := recordInfo.GetV_StringValueFrom(`V_StringField`, record)
	checkExpectedGetValueFrom(t, value, `B`, isNull, false, err, nil)

	value, isNull, err = recordInfo.GetV_WStringValueFrom(`V_WStringField`, record)
	checkExpectedGetValueFrom(t, value, `A`, isNull, false, err, nil)
}

func generateTestRecordInfo() recordinfo.RecordInfo {
	recordInfo := recordinfo.New()
	recordInfo.AddByteField(`ByteField`, ``)
	recordInfo.AddBoolField(`BoolField`, ``)
	recordInfo.AddInt16Field(`Int16Field`, ``)
	recordInfo.AddInt32Field(`Int32Field`, ``)
	recordInfo.AddInt64Field(`Int64Field`, ``)
	recordInfo.AddFixedDecimalField(`FixedDecimalField`, ``, 19, 2)
	recordInfo.AddFloatField(`FloatField`, ``)
	recordInfo.AddDoubleField(`DoubleField`, ``)
	recordInfo.AddStringField(`StringField`, ``, 10)
	recordInfo.AddWStringField(`WStringField`, ``, 5)
	recordInfo.AddDateField(`DateField`, ``)
	recordInfo.AddDateTimeField(`DateTimeField`, ``)
	return recordInfo
}

func setRecordInfoTestData(recordInfo recordinfo.RecordInfo) {
	_ = recordInfo.SetByteField(`ByteField`, byte(1))
	_ = recordInfo.SetBoolField(`BoolField`, true)
	_ = recordInfo.SetInt16Field(`Int16Field`, 2)
	_ = recordInfo.SetInt32Field(`Int32Field`, 3)
	_ = recordInfo.SetInt64Field(`Int64Field`, 4)
	_ = recordInfo.SetFixedDecimalField(`FixedDecimalField`, 123.45)
	_ = recordInfo.SetFloatField(`FloatField`, float32(654.321))
	_ = recordInfo.SetDoubleField(`DoubleField`, 909.33)
	_ = recordInfo.SetStringField(`StringField`, `ABCDEFG`)
	_ = recordInfo.SetWStringField(`WStringField`, `CXVY`)
	_ = recordInfo.SetDateField(`DateField`, time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC))
	_ = recordInfo.SetDateTimeField(`DateTimeField`, time.Date(2021, 3, 4, 5, 6, 7, 0, time.UTC))
}
