package recordinfo_test

import (
	"github.com/tlarsen7572/goalteryx/recordcopier"
	"github.com/tlarsen7572/goalteryx/recordinfo"
	"strings"
	"testing"
	"time"
)

func TestSetValuesAndGenerateRecord(t *testing.T) {
	recordInfo, reader := generateTestRecordInfo()
	setRecordInfoTestData(recordInfo)

	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	byteVal, isNull, err := reader.GetIntValueFrom(`ByteField`, record)
	checkExpectedGetValueFrom(t, byteVal, 1, isNull, false, err, nil, `error setting byte:`)

	boolVal, isNull, err := reader.GetBoolValueFrom(`BoolField`, record)
	checkExpectedGetValueFrom(t, boolVal, true, isNull, false, err, nil, `error setting bool:`)

	int16Val, isNull, err := reader.GetIntValueFrom(`Int16Field`, record)
	checkExpectedGetValueFrom(t, int16Val, 2, isNull, false, err, nil, `error setting int16:`)

	int32Val, isNull, err := reader.GetIntValueFrom(`Int32Field`, record)
	checkExpectedGetValueFrom(t, int32Val, 3, isNull, false, err, nil, `error setting int32:`)

	int64Val, isNull, err := reader.GetIntValueFrom(`Int64Field`, record)
	checkExpectedGetValueFrom(t, int64Val, 4, isNull, false, err, nil, `error setting int64:`)

	fixedDecimalVal, isNull, err := reader.GetFloatValueFrom(`FixedDecimalField`, record)
	checkExpectedGetValueFrom(t, fixedDecimalVal, 123.45, isNull, false, err, nil, `error setting fixeddecimal:`)

	floatVal, isNull, err := reader.GetFloatValueFrom(`FloatField`, record)
	checkExpectedGetValueFrom(t, floatVal, float64(float32(654.321)), isNull, false, err, nil, `error setting float:`)

	doubleVal, isNull, err := reader.GetFloatValueFrom(`DoubleField`, record)
	checkExpectedGetValueFrom(t, doubleVal, 909.33, isNull, false, err, nil, `error setting double:`)

	stringVal, isNull, err := reader.GetStringValueFrom(`StringField`, record)
	checkExpectedGetValueFrom(t, stringVal, `ABCDEFG`, isNull, false, err, nil, `error setting string:`)

	wstringVal, isNull, err := reader.GetStringValueFrom(`WStringField`, record)
	checkExpectedGetValueFrom(t, wstringVal, `CXVY`, isNull, false, err, nil, `error setting wstring:`)

	dateVal, isNull, err := reader.GetDateValueFrom(`DateField`, record)
	expectedDate := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	checkExpectedGetValueFrom(t, dateVal, expectedDate, isNull, false, err, nil, `error setting date:`)

	dateTimeVal, isNull, err := reader.GetDateValueFrom(`DateTimeField`, record)
	expectedDate = time.Date(2021, 3, 4, 5, 6, 7, 0, time.UTC)
	checkExpectedGetValueFrom(t, dateTimeVal, expectedDate, isNull, false, err, nil, `error setting datetime:`)
}

func TestSetNullValuesAndGenerateRecord(t *testing.T) {
	recordInfo, reader := generateTestRecordInfo()
	setNullTestData(recordInfo)

	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	_, isNull, err := reader.GetIntValueFrom(`ByteField`, record)
	checkExpectedGetNullFrom(t, isNull, true, err, nil, `error setting byte:`)

	_, isNull, err = reader.GetBoolValueFrom(`BoolField`, record)
	checkExpectedGetNullFrom(t, isNull, true, err, nil, `error setting bool:`)

	_, isNull, err = reader.GetIntValueFrom(`Int16Field`, record)
	checkExpectedGetNullFrom(t, isNull, true, err, nil, `error setting int16:`)

	_, isNull, err = reader.GetIntValueFrom(`Int32Field`, record)
	checkExpectedGetNullFrom(t, isNull, true, err, nil, `error setting int32:`)

	_, isNull, err = reader.GetIntValueFrom(`Int64Field`, record)
	checkExpectedGetNullFrom(t, isNull, true, err, nil, `error setting int64:`)

	_, isNull, err = reader.GetFloatValueFrom(`FixedDecimalField`, record)
	checkExpectedGetNullFrom(t, isNull, true, err, nil, `error setting fixeddecimal:`)

	_, isNull, err = reader.GetFloatValueFrom(`FloatField`, record)
	checkExpectedGetNullFrom(t, isNull, true, err, nil, `error setting float:`)

	_, isNull, err = reader.GetFloatValueFrom(`DoubleField`, record)
	checkExpectedGetNullFrom(t, isNull, true, err, nil, `error setting double:`)

	_, isNull, err = reader.GetStringValueFrom(`StringField`, record)
	checkExpectedGetNullFrom(t, isNull, true, err, nil, `error setting string:`)

	_, isNull, err = reader.GetStringValueFrom(`WStringField`, record)
	checkExpectedGetNullFrom(t, isNull, true, err, nil, `error setting wstring:`)

	_, isNull, err = reader.GetDateValueFrom(`DateField`, record)
	checkExpectedGetNullFrom(t, isNull, true, err, nil, `error setting date:`)

	_, isNull, err = reader.GetDateValueFrom(`DateTimeField`, record)
	checkExpectedGetNullFrom(t, isNull, true, err, nil, `error setting datetime:`)
}

func TestCachedRecords(t *testing.T) {
	recordInfo, _ := generateTestRecordInfo()
	setRecordInfoTestData(recordInfo)

	record1, _ := recordInfo.GenerateRecord()
	record2, _ := recordInfo.GenerateRecord()
	if record1.Blob() != record2.Blob() {
		t.Fatalf(`record1 and record2 are 2 different pointers`)
	}
}

func TestSetLongVarDataFieldsAndGenerateRecord(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddByteField(`ByteField`, ``)
	generator.AddV_WStringField(`V_WStringField`, ``, 250)
	generator.AddV_StringField(`V_StringField`, ``, 250)
	recordInfo := generator.GenerateRecordInfo()

	_ = recordInfo.SetIntField(`ByteField`, 1)
	_ = recordInfo.SetStringField(`V_StringField`, strings.Repeat(`B`, 200))
	_ = recordInfo.SetStringField(`V_WStringField`, strings.Repeat(`A`, 100))

	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	reader := generator.GenerateRecordBlobReader()
	value, isNull, err := reader.GetStringValueFrom(`V_StringField`, record)
	checkExpectedGetValueFrom(t, value, strings.Repeat(`B`, 200), isNull, false, err, nil, `error setting long v_string:`)

	value, isNull, err = reader.GetStringValueFrom(`V_WStringField`, record)
	checkExpectedGetValueFrom(t, value, strings.Repeat(`A`, 100), isNull, false, err, nil, `error setting long v_wstring:`)
}

func TestSetShortVarDataFieldsAndGenerateRecord(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddByteField(`ByteField`, ``)
	generator.AddV_WStringField(`V_WStringField`, ``, 250)
	generator.AddV_StringField(`V_StringField`, ``, 250)
	recordInfo := generator.GenerateRecordInfo()

	_ = recordInfo.SetIntField(`ByteField`, 1)
	_ = recordInfo.SetStringField(`V_StringField`, strings.Repeat(`B`, 100))
	_ = recordInfo.SetStringField(`V_WStringField`, strings.Repeat(`A`, 50))

	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	reader := generator.GenerateRecordBlobReader()
	value, isNull, err := reader.GetStringValueFrom(`V_StringField`, record)
	checkExpectedGetValueFrom(t, value, strings.Repeat(`B`, 100), isNull, false, err, nil, `error setting short v_string:`)

	value, isNull, err = reader.GetStringValueFrom(`V_WStringField`, record)
	checkExpectedGetValueFrom(t, value, strings.Repeat(`A`, 50), isNull, false, err, nil, `error setting short v_wstring:`)
}

func TestSetTinyVarDataFieldsAndGenerateRecord(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddByteField(`ByteField`, ``)
	generator.AddV_WStringField(`V_WStringField`, ``, 250)
	generator.AddV_StringField(`V_StringField`, ``, 250)
	recordInfo := generator.GenerateRecordInfo()

	_ = recordInfo.SetIntField(`ByteField`, 1)
	_ = recordInfo.SetStringField(`V_StringField`, `B`)
	_ = recordInfo.SetStringField(`V_WStringField`, `A`)

	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	reader := generator.GenerateRecordBlobReader()
	value, isNull, err := reader.GetStringValueFrom(`V_StringField`, record)
	checkExpectedGetValueFrom(t, value, `B`, isNull, false, err, nil, `error setting tiny v_string:`)

	value, isNull, err = reader.GetStringValueFrom(`V_WStringField`, record)
	checkExpectedGetValueFrom(t, value, `A`, isNull, false, err, nil, `error setting tiny v_wstring:`)
}

func TestSetEmptyVarDataFieldsAndGenerateRecord(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddByteField(`ByteField`, ``)
	generator.AddV_WStringField(`V_WStringField`, ``, 250)
	generator.AddV_StringField(`V_StringField`, ``, 250)
	recordInfo := generator.GenerateRecordInfo()

	_ = recordInfo.SetIntField(`ByteField`, 1)
	_ = recordInfo.SetStringField(`V_StringField`, ``)
	_ = recordInfo.SetStringField(`V_WStringField`, ``)

	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	reader := generator.GenerateRecordBlobReader()
	value, isNull, err := reader.GetStringValueFrom(`V_StringField`, record)
	checkExpectedGetValueFrom(t, value, ``, isNull, false, err, nil, `error setting empty v_string:`)

	value, isNull, err = reader.GetStringValueFrom(`V_WStringField`, record)
	checkExpectedGetValueFrom(t, value, ``, isNull, false, err, nil, `error setting empty v_wstring:`)
}

func TestSetNullVarDataFieldsAndGenerateRecord(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddByteField(`ByteField`, ``)
	generator.AddV_WStringField(`V_WStringField`, ``, 250)
	generator.AddV_StringField(`V_StringField`, ``, 250)
	recordInfo := generator.GenerateRecordInfo()

	_ = recordInfo.SetIntField(`ByteField`, 1)
	_ = recordInfo.SetFieldNull(`V_StringField`)
	_ = recordInfo.SetFieldNull(`V_WStringField`)

	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	reader := generator.GenerateRecordBlobReader()
	value, isNull, err := reader.GetStringValueFrom(`V_StringField`, record)
	checkExpectedGetValueFrom(t, value, ``, isNull, true, err, nil, `error setting null v_string:`)

	value, isNull, err = reader.GetStringValueFrom(`V_WStringField`, record)
	checkExpectedGetValueFrom(t, value, ``, isNull, true, err, nil, `error setting null v_wstring:`)
}

func TestSetFixedLenFromRawBytes(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddByteField(`ByteField`, ``)
	recordInfo := generator.GenerateRecordInfo()

	err := recordInfo.SetFromRawBytes(`ByteField`, []byte{4, 0})
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	reader := generator.GenerateRecordBlobReader()
	value, isNull, err := reader.GetIntValueFrom(`ByteField`, record)
	checkExpectedGetValueFrom(t, value, 4, isNull, false, err, nil, `error setting raw bytes:`)
}

func TestSetVarLenFromRawBytes(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddV_StringField(`V_StringField`, ``, 250)
	recordInfo := generator.GenerateRecordInfo()

	err := recordInfo.SetFromRawBytes(`V_StringField`, []byte(`Hello world, how are you?`))
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	record, err := recordInfo.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	reader := generator.GenerateRecordBlobReader()
	value, isNull, err := reader.GetStringValueFrom(`V_StringField`, record)
	checkExpectedGetValueFrom(t, value, `Hello world, how are you?`, isNull, false, err, nil, `error setting raw bytes:`)
}

func TestSetStringOfSmallerLength(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddStringField(`StringField`, ``, 100)
	recordInfo := generator.GenerateRecordInfo()

	_ = recordInfo.SetStringField(`StringField`, `Start`)
	record, _ := recordInfo.GenerateRecord()

	reader := generator.GenerateRecordBlobReader()
	value, _, _ := reader.GetStringValueFrom(`StringField`, record)
	if value != `Start` {
		t.Fatalf(`expected 'Start' but got '%v'`, value)
	}

	_ = recordInfo.SetStringField(`StringField`, `End`)
	record, _ = recordInfo.GenerateRecord()

	value, _, _ = reader.GetStringValueFrom(`StringField`, record)
	if value != `End` {
		t.Fatalf(`expected 'End' but got '%v'`, value)
	}
}

func TestSetV_StringOfSmallerLength(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddV_StringField(`StringField`, ``, 100)
	recordInfo := generator.GenerateRecordInfo()

	_ = recordInfo.SetStringField(`StringField`, `Start`)
	record, _ := recordInfo.GenerateRecord()

	reader := generator.GenerateRecordBlobReader()
	value, _, _ := reader.GetStringValueFrom(`StringField`, record)
	if value != `Start` {
		t.Fatalf(`expected 'Start' but got '%v'`, value)
	}

	_ = recordInfo.SetStringField(`StringField`, `End`)
	record, _ = recordInfo.GenerateRecord()

	value, _, _ = reader.GetStringValueFrom(`StringField`, record)
	if value != `End` {
		t.Fatalf(`expected 'End' but got '%v'`, value)
	}
}

func TestSetWStringOfSmallerLength(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddWStringField(`StringField`, ``, 100)
	recordInfo := generator.GenerateRecordInfo()

	_ = recordInfo.SetStringField(`StringField`, `Start`)
	record, _ := recordInfo.GenerateRecord()

	reader := generator.GenerateRecordBlobReader()
	value, _, _ := reader.GetStringValueFrom(`StringField`, record)
	if value != `Start` {
		t.Fatalf(`expected 'Start' but got '%v'`, value)
	}

	_ = recordInfo.SetStringField(`StringField`, `End`)
	record, _ = recordInfo.GenerateRecord()

	value, _, _ = reader.GetStringValueFrom(`StringField`, record)
	if value != `End` {
		t.Fatalf(`expected 'End' but got '%v'`, value)
	}
}

func TestSetV_WStringOfSmallerLength(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddV_WStringField(`StringField`, ``, 100)
	recordInfo := generator.GenerateRecordInfo()

	_ = recordInfo.SetStringField(`StringField`, `Start`)
	record, _ := recordInfo.GenerateRecord()

	reader := generator.GenerateRecordBlobReader()
	value, _, _ := reader.GetStringValueFrom(`StringField`, record)
	if value != `Start` {
		t.Fatalf(`expected 'Start' but got '%v'`, value)
	}

	_ = recordInfo.SetStringField(`StringField`, `End`)
	record, _ = recordInfo.GenerateRecord()

	value, _, _ = reader.GetStringValueFrom(`StringField`, record)
	if value != `End` {
		t.Fatalf(`expected 'End' but got '%v'`, value)
	}
}

func TestSetTruncatedString(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddStringField(`StringField`, ``, 2)
	recordInfo := generator.GenerateRecordInfo()

	_ = recordInfo.SetStringField(`StringField`, `Start`)
	record, _ := recordInfo.GenerateRecord()

	reader := generator.GenerateRecordBlobReader()
	value, _, _ := reader.GetStringValueFrom(`StringField`, record)
	if value != `St` {
		t.Fatalf(`expected 'St' but got '%v'`, value)
	}
}

func TestSetTruncatedV_String(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddV_StringField(`StringField`, ``, 2)
	recordInfo := generator.GenerateRecordInfo()

	_ = recordInfo.SetStringField(`StringField`, `Start`)
	record, _ := recordInfo.GenerateRecord()

	reader := generator.GenerateRecordBlobReader()
	value, _, _ := reader.GetStringValueFrom(`StringField`, record)
	if value != `St` {
		t.Fatalf(`expected 'St' but got '%v'`, value)
	}
}

func TestSetTruncatedWString(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddWStringField(`StringField`, ``, 2)
	recordInfo := generator.GenerateRecordInfo()

	_ = recordInfo.SetStringField(`StringField`, `Start`)
	record, _ := recordInfo.GenerateRecord()

	reader := generator.GenerateRecordBlobReader()
	value, _, _ := reader.GetStringValueFrom(`StringField`, record)
	if value != `St` {
		t.Fatalf(`expected 'St' but got '%v'`, value)
	}
}

func TestSetTruncatedV_WString(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddV_WStringField(`StringField`, ``, 2)
	recordInfo := generator.GenerateRecordInfo()

	_ = recordInfo.SetStringField(`StringField`, `Start`)
	record, _ := recordInfo.GenerateRecord()

	reader := generator.GenerateRecordBlobReader()
	value, _, _ := reader.GetStringValueFrom(`StringField`, record)
	if value != `St` {
		t.Fatalf(`expected 'St' but got '%v'`, value)
	}
}

func TestSetValueNullValue(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddInt64Field(`Field`, ``)
	recordInfo := generator.GenerateRecordInfo()

	_ = recordInfo.SetIntField(`Field`, 10)
	record, _ := recordInfo.GenerateRecord()

	reader := generator.GenerateRecordBlobReader()
	value, isNull, err := reader.GetIntValueFrom(`Field`, record)
	checkExpectedGetValueFrom(t, value, 10, isNull, false, err, nil, `error setting 10`)

	_ = recordInfo.SetFieldNull(`Field`)
	record, _ = recordInfo.GenerateRecord()

	_, isNull, err = reader.GetIntValueFrom(`Field`, record)
	checkExpectedGetNullFrom(t, isNull, true, err, nil, `error setting null`)

	_ = recordInfo.SetIntField(`Field`, 20)
	record, _ = recordInfo.GenerateRecord()

	value, isNull, err = reader.GetIntValueFrom(`Field`, record)
	checkExpectedGetValueFrom(t, value, 20, isNull, false, err, nil, `error setting 20`)
}

func TestOverwriteFixedLenFieldWithNull(t *testing.T) {
	info, reader := generateTestRecordInfo()
	copierMap := make([]recordcopier.IndexMap, info.NumFields())
	for i := 0; i < info.NumFields(); i++ {
		copierMap[i] = recordcopier.IndexMap{
			DestinationIndex: i,
			SourceIndex:      i,
		}
	}
	copier, _ := recordcopier.New(info, info.GenerateRecordBlobReader(), copierMap)
	setNullTestData(info)
	blob, _ := info.GenerateRecord()
	_ = copier.Copy(blob)

	setRecordInfoTestData(info)
	blob, _ = info.GenerateRecord()

	_, isNull, _ := reader.GetIntValueFrom(`ByteField`, blob)
	checkExpectedGetNullFrom(t, isNull, false, nil, nil, `byte field`)

	_, isNull, _ = reader.GetBoolValueFrom(`BoolField`, blob)
	checkExpectedGetNullFrom(t, isNull, false, nil, nil, `bool field`)

	_, isNull, _ = reader.GetIntValueFrom(`Int16Field`, blob)
	checkExpectedGetNullFrom(t, isNull, false, nil, nil, `int16 field`)

	_, isNull, _ = reader.GetIntValueFrom(`Int32Field`, blob)
	checkExpectedGetNullFrom(t, isNull, false, nil, nil, `int32 field`)

	_, isNull, _ = reader.GetIntValueFrom(`Int64Field`, blob)
	checkExpectedGetNullFrom(t, isNull, false, nil, nil, `int64 field`)

	_, isNull, _ = reader.GetFloatValueFrom(`FixedDecimalField`, blob)
	checkExpectedGetNullFrom(t, isNull, false, nil, nil, `fixed decimal field`)

	_, isNull, _ = reader.GetFloatValueFrom(`FloatField`, blob)
	checkExpectedGetNullFrom(t, isNull, false, nil, nil, `float field`)

	_, isNull, _ = reader.GetFloatValueFrom(`DoubleField`, blob)
	checkExpectedGetNullFrom(t, isNull, false, nil, nil, `double field`)

	_, isNull, _ = reader.GetStringValueFrom(`StringField`, blob)
	checkExpectedGetNullFrom(t, isNull, false, nil, nil, `string field`)

	_, isNull, _ = reader.GetStringValueFrom(`WStringField`, blob)
	checkExpectedGetNullFrom(t, isNull, false, nil, nil, `wstring field`)

	_, isNull, _ = reader.GetDateValueFrom(`DateField`, blob)
	checkExpectedGetNullFrom(t, isNull, false, nil, nil, `date field`)

	_, isNull, _ = reader.GetDateValueFrom(`DateTimeField`, blob)
	checkExpectedGetNullFrom(t, isNull, false, nil, nil, `datetime field`)
}

func TestOverwriteVarStringFieldWithNull(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddV_WStringField(`V_WStringField`, ``, 100)
	generator.AddV_StringField(`V_StringField`, ``, 100)
	info := generator.GenerateRecordInfo()
	copierMap := make([]recordcopier.IndexMap, info.NumFields())
	for i := 0; i < info.NumFields(); i++ {
		copierMap[i] = recordcopier.IndexMap{
			DestinationIndex: i,
			SourceIndex:      i,
		}
	}
	copier, _ := recordcopier.New(info, info.GenerateRecordBlobReader(), copierMap)
	_ = info.SetFieldNull(`V_WStringField`)
	_ = info.SetFieldNull(`V_StringField`)
	blob, _ := info.GenerateRecord()
	_ = copier.Copy(blob)

	_ = info.SetStringField(`V_WStringField`, `hello world`)
	_ = info.SetStringField(`V_StringField`, `hello world again`)
	blob, _ = info.GenerateRecord()

	reader := generator.GenerateRecordBlobReader()

	_, isNull, _ := reader.GetStringValueFrom(`V_WStringField`, blob)
	checkExpectedGetNullFrom(t, isNull, false, nil, nil, `V_WStringField field`)

	_, isNull, _ = reader.GetStringValueFrom(`V_StringField`, blob)
	checkExpectedGetNullFrom(t, isNull, false, nil, nil, `V_StringField field`)
}

func generateTestRecordInfo() (recordinfo.RecordInfo, recordinfo.RecordBlobReader) {
	generator := recordinfo.NewGenerator()
	generator.AddByteField(`ByteField`, ``)
	generator.AddBoolField(`BoolField`, ``)
	generator.AddInt16Field(`Int16Field`, ``)
	generator.AddInt32Field(`Int32Field`, ``)
	generator.AddInt64Field(`Int64Field`, ``)
	generator.AddFixedDecimalField(`FixedDecimalField`, ``, 19, 2)
	generator.AddFloatField(`FloatField`, ``)
	generator.AddDoubleField(`DoubleField`, ``)
	generator.AddStringField(`StringField`, ``, 10)
	generator.AddWStringField(`WStringField`, ``, 5)
	generator.AddDateField(`DateField`, ``)
	generator.AddDateTimeField(`DateTimeField`, ``)
	return generator.GenerateRecordInfo(), generator.GenerateRecordBlobReader()
}

func setRecordInfoTestData(recordInfo recordinfo.RecordInfo) {
	_ = recordInfo.SetIntField(`ByteField`, 1)
	_ = recordInfo.SetBoolField(`BoolField`, true)
	_ = recordInfo.SetIntField(`Int16Field`, 2)
	_ = recordInfo.SetIntField(`Int32Field`, 3)
	_ = recordInfo.SetIntField(`Int64Field`, 4)
	_ = recordInfo.SetFloatField(`FixedDecimalField`, 123.45)
	_ = recordInfo.SetFloatField(`FloatField`, 654.321)
	_ = recordInfo.SetFloatField(`DoubleField`, 909.33)
	_ = recordInfo.SetStringField(`StringField`, `ABCDEFG`)
	_ = recordInfo.SetStringField(`WStringField`, `CXVY`)
	_ = recordInfo.SetDateField(`DateField`, time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC))
	_ = recordInfo.SetDateField(`DateTimeField`, time.Date(2021, 3, 4, 5, 6, 7, 0, time.UTC))
}

func setNullTestData(recordInfo recordinfo.RecordInfo) {
	_ = recordInfo.SetFieldNull(`ByteField`)
	_ = recordInfo.SetFieldNull(`BoolField`)
	_ = recordInfo.SetFieldNull(`Int16Field`)
	_ = recordInfo.SetFieldNull(`Int32Field`)
	_ = recordInfo.SetFieldNull(`Int64Field`)
	_ = recordInfo.SetFieldNull(`FixedDecimalField`)
	_ = recordInfo.SetFieldNull(`FloatField`)
	_ = recordInfo.SetFieldNull(`DoubleField`)
	_ = recordInfo.SetFieldNull(`StringField`)
	_ = recordInfo.SetFieldNull(`WStringField`)
	_ = recordInfo.SetFieldNull(`DateField`)
	_ = recordInfo.SetFieldNull(`DateTimeField`)
}

func checkExpectedGetNullFrom(t *testing.T, isNull bool, expectedIsNull bool, err error, expectedErr error, msg string) {
	if err != expectedErr {
		t.Fatalf("%v expected error: %v\ngot: %v", msg, expectedErr, err)
	}
	if isNull != expectedIsNull {
		t.Fatalf(`%v expected isNull=%v but got isNull=%v`, msg, expectedIsNull, isNull)
	}
}
