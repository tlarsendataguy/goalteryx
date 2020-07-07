package recordinfo_test

import (
	"github.com/tlarsen7572/goalteryx/recordblob"
	"github.com/tlarsen7572/goalteryx/recordcopier"
	"github.com/tlarsen7572/goalteryx/recordinfo"
	"strings"
	"testing"
	"time"
	"unsafe"
)

func TestInstantiateRecordInfoFromXml(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if count := recordInfo.NumFields(); count != 16 {
		t.Fatalf(`expecpted 16 fields but got %v`, count)
	}
	type expectedStruct struct {
		Name      string
		Size      int
		Scale     int
		FieldType recordinfo.FieldType
	}
	expectedFields := []expectedStruct{
		{`ByteField`, 1, 0, recordinfo.Byte},
		{`BoolField`, 1, 0, recordinfo.Bool},
		{`Int16Field`, 2, 0, recordinfo.Int16},
		{`Int32Field`, 4, 0, recordinfo.Int32},
		{`Int64Field`, 8, 0, recordinfo.Int64},
		{`FixedDecimalField`, 19, 6, recordinfo.FixedDecimal},
		{`FloatField`, 4, 0, recordinfo.Float},
		{`DoubleField`, 8, 0, recordinfo.Double},
		{`StringField`, 64, 0, recordinfo.String},
		{`WStringField`, 64, 0, recordinfo.WString},
		{`V_StringShortField`, 1000, 0, recordinfo.V_String},
		{`V_StringLongField`, 2147483647, 0, recordinfo.V_String},
		{`V_WStringShortField`, 10, 0, recordinfo.V_WString},
		{`V_WStringLongField`, 1073741823, 0, recordinfo.V_WString},
		{`DateField`, 10, 0, recordinfo.Date},
		{`DateTimeField`, 19, 0, recordinfo.DateTime},
	}
	for index, expectedField := range expectedFields {
		field, _ := recordInfo.GetFieldByIndex(index)
		if field.Name != expectedField.Name {
			t.Fatalf(`expected name '%v' but got '%v' at field %v`, expectedField.Name, field.Name, index)
		}
		if field.Size != expectedField.Size {
			t.Fatalf(`expected size %v but got %v at field %v`, expectedField.Size, field.Size, index)
		}
		if field.Precision != expectedField.Scale {
			t.Fatalf(`expected scale %v but got %v at field %v`, expectedField.Scale, field.Precision, index)
		}
		if field.Type != expectedField.FieldType {
			t.Fatalf(`expected '%v' but got '%v' at field %v`, expectedField.FieldType, field.Type, index)
		}
	}
}

func TestSaveRecordInfoToXml(t *testing.T) {
	recordInfo, _ := recordinfo.FromXml(recordInfoXml)
	xmlConfig, err := recordInfo.ToXml(`Output`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	expectedXml := `<MetaInfo connection="Output"><RecordInfo><Field name="ByteField" source="TextInput:" size="1" scale="0" type="Byte"></Field><Field name="BoolField" source="Formula: 1" size="1" scale="0" type="Bool"></Field><Field name="Int16Field" source="Formula: 16" size="2" scale="0" type="Int16"></Field><Field name="Int32Field" source="Formula: 32" size="4" scale="0" type="Int32"></Field><Field name="Int64Field" source="Formula: 64" size="8" scale="0" type="Int64"></Field><Field name="FixedDecimalField" source="Formula: 123.45" size="19" scale="6" type="FixedDecimal"></Field><Field name="FloatField" source="Formula: 678.9" size="4" scale="0" type="Float"></Field><Field name="DoubleField" source="Formula: 0.12345" size="8" scale="0" type="Double"></Field><Field name="StringField" source="Formula: &#34;A&#34;" size="64" scale="0" type="String"></Field><Field name="WStringField" source="Formula: &#34;AB&#34;" size="64" scale="0" type="WString"></Field><Field name="V_StringShortField" source="Formula: &#34;ABC&#34;" size="1000" scale="0" type="V_String"></Field><Field name="V_StringLongField" source="Formula: PadLeft(&#34;&#34;, 500, &#39;B&#39;)" size="2147483647" scale="0" type="V_String"></Field><Field name="V_WStringShortField" source="Formula: &#34;XZY&#34;" size="10" scale="0" type="V_WString"></Field><Field name="V_WStringLongField" source="Formula: PadLeft(&#34;&#34;, 500, &#39;W&#39;)" size="1073741823" scale="0" type="V_WString"></Field><Field name="DateField" source="Formula: &#39;2020-01-01&#39;" size="10" scale="0" type="Date"></Field><Field name="DateTimeField" source="Formula: &#39;2020-02-03 04:05:06&#39;" size="19" scale="0" type="DateTime"></Field></RecordInfo></MetaInfo>`
	if xmlConfig != expectedXml {
		t.Fatalf("expected:\n%v\nbut got:\n%v", expectedXml, xmlConfig)
	}
	t.Logf(xmlConfig)
}

func TestCorrectlyRetrieveByteValue(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetIntValueFrom(`ByteField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, 1, isNull, false, err, nil, `error retrieving byte:`)

	value, isNull, err = recordInfo.GetIntValueFrom(`ByteField`, nullRecord)
	checkExpectedGetValueFrom(t, value, 0, isNull, true, err, nil, `error retrieving null byte:`)
}

func TestCorrectlyRetrieveBoolValue(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetBoolValueFrom(`BoolField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, true, isNull, false, err, nil, `error retrieving bool:`)

	value, isNull, err = recordInfo.GetBoolValueFrom(`BoolField`, nullRecord)
	checkExpectedGetValueFrom(t, value, false, isNull, true, err, nil, `error retrieving null bool:`)
}

func TestCorrectlyRetrieveInt16Value(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetIntValueFrom(`Int16Field`, sampleRecord)
	checkExpectedGetValueFrom(t, value, 16, isNull, false, err, nil, `error retrieving int16:`)

	value, isNull, err = recordInfo.GetIntValueFrom(`Int16Field`, nullRecord)
	checkExpectedGetValueFrom(t, value, 0, isNull, true, err, nil, `error retrieving null int16:`)
}

func TestCorrectlyRetrieveInt32Value(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetIntValueFrom(`Int32Field`, sampleRecord)
	checkExpectedGetValueFrom(t, value, 32, isNull, false, err, nil, `error retrieving int32:`)

	value, isNull, err = recordInfo.GetIntValueFrom(`Int32Field`, nullRecord)
	checkExpectedGetValueFrom(t, value, 0, isNull, true, err, nil, `error retrieving null int32:`)
}

func TestCorrectlyRetrieveInt64Value(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetIntValueFrom(`Int64Field`, sampleRecord)
	checkExpectedGetValueFrom(t, value, 64, isNull, false, err, nil, `error retrieving int64:`)

	value, isNull, err = recordInfo.GetIntValueFrom(`Int64Field`, nullRecord)
	checkExpectedGetValueFrom(t, value, 0, isNull, true, err, nil, `error retrieving null int64:`)
}

func TestCorrectlyRetrieveFixedDecimalValue(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetFloatValueFrom(`FixedDecimalField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, 123.450000, isNull, false, err, nil, `error retrieving fixed decimal:`)

	value, isNull, err = recordInfo.GetFloatValueFrom(`FixedDecimalField`, nullRecord)
	checkExpectedGetValueFrom(t, value, 0.0, isNull, true, err, nil, `error retrieving null fixed decimal:`)
}

func TestCorrectlyRetrieveFloatValue(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetFloatValueFrom(`FloatField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, float64(float32(678.9)), isNull, false, err, nil, `error retrieving float:`)

	value, isNull, err = recordInfo.GetFloatValueFrom(`FloatField`, nullRecord)
	checkExpectedGetValueFrom(t, value, 0.0, isNull, true, err, nil, `error retrieving null float:`)
}

func TestCorrectlyRetrieveDoubleValue(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetFloatValueFrom(`DoubleField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, 0.12345, isNull, false, err, nil, `error retrieving double:`)

	value, isNull, err = recordInfo.GetFloatValueFrom(`DoubleField`, nullRecord)
	checkExpectedGetValueFrom(t, value, 0.0, isNull, true, err, nil, `error retrieving null double:`)
}

func TestCorrectlyRetrieveStringValue(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetStringValueFrom(`StringField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, `A`, isNull, false, err, nil, `error retrieving string:`)

	value, isNull, err = recordInfo.GetStringValueFrom(`StringField`, nullRecord)
	checkExpectedGetValueFrom(t, value, ``, isNull, true, err, nil, `error retrieving null string:`)
}

func TestCorrectlyRetrieveWStringValue(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetStringValueFrom(`WStringField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, `AB`, isNull, false, err, nil, `error retrieving wstring:`)

	value, isNull, err = recordInfo.GetStringValueFrom(`WStringField`, nullRecord)
	checkExpectedGetValueFrom(t, value, ``, isNull, true, err, nil, `error retrieving null wstring:`)
}

var zeroDate = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

func TestCorrectlyRetrieveDateValue(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetDateValueFrom(`DateField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), isNull, false, err, nil, `error retrieving date:`)

	value, isNull, err = recordInfo.GetDateValueFrom(`DateField`, nullRecord)
	checkExpectedGetValueFrom(t, value, zeroDate, isNull, true, err, nil, `error retrieving null date:`)
}

func TestCorrectlyRetrieveDateTimeValue(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetDateValueFrom(`DateTimeField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, time.Date(2020, 2, 3, 4, 5, 6, 0, time.UTC), isNull, false, err, nil, `error retrieving datetime:`)

	value, isNull, err = recordInfo.GetDateValueFrom(`DateTimeField`, nullRecord)
	checkExpectedGetValueFrom(t, value, zeroDate, isNull, true, err, nil, `error retrieving null datetime:`)
}

func TestCorrectlyRetrieveV_StringLongValue(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddByteField(`ByteField`, ``)
	generator.AddV_WStringField(`V_WStringField`, ``, 250)
	generator.AddV_StringField(`V_StringField`, ``, 250)
	recordInfo := generator.GenerateRecordInfo()

	value, isNull, err := recordInfo.GetStringValueFrom(`V_StringField`, varFieldLongRecord)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected is not null but got null`)
	}
	if expected := strings.Repeat(`B`, 200); value != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, value)
	}
}

func TestCorrectlyRetrieveV_WStringLongValue(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddByteField(`ByteField`, ``)
	generator.AddV_WStringField(`V_WStringField`, ``, 250)
	generator.AddV_StringField(`V_StringField`, ``, 250)
	recordInfo := generator.GenerateRecordInfo()
	value, isNull, err := recordInfo.GetStringValueFrom(`V_WStringField`, varFieldLongRecord)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected is not null but got null`)
	}
	if expected := strings.Repeat(`A`, 100); value != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, value)
	}
}

func TestCorrectlyRetrieveVarStringsNull(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddByteField(`ByteField`, ``)
	generator.AddV_WStringField(`V_WStringField`, ``, 250)
	generator.AddV_StringField(`V_StringField`, ``, 250)
	recordInfo := generator.GenerateRecordInfo()

	value, isNull, err := recordInfo.GetStringValueFrom(`V_WStringField`, varFieldNullRecord)
	checkExpectedGetValueFrom(t, value, ``, isNull, true, err, nil, `error retrieving null v_wstring:`)

	value, isNull, err = recordInfo.GetStringValueFrom(`V_StringField`, varFieldNullRecord)
	checkExpectedGetValueFrom(t, value, ``, isNull, true, err, nil, `error retrieving null v_string:`)
}

func TestCorrectlyRetrieveVarStringsEmpty(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddByteField(`ByteField`, ``)
	generator.AddV_WStringField(`V_WStringField`, ``, 250)
	generator.AddV_StringField(`V_StringField`, ``, 250)
	recordInfo := generator.GenerateRecordInfo()

	value, isNull, err := recordInfo.GetStringValueFrom(`V_WStringField`, varFieldEmptyStrings)
	checkExpectedGetValueFrom(t, value, ``, isNull, false, err, nil, `error retrieving empty v_wstring:`)

	value, isNull, err = recordInfo.GetStringValueFrom(`V_StringField`, varFieldEmptyStrings)
	checkExpectedGetValueFrom(t, value, ``, isNull, false, err, nil, `error retrieving empty v_string:`)
}

func TestCorrectlyRetrieveV_StringShortValue(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddByteField(`ByteField`, ``)
	generator.AddV_WStringField(`V_WStringField`, ``, 250)
	generator.AddV_StringField(`V_StringField`, ``, 250)
	recordInfo := generator.GenerateRecordInfo()
	value, isNull, err := recordInfo.GetStringValueFrom(`V_StringField`, varFieldShortRecord)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected is not null but got null`)
	}
	if expected := strings.Repeat(`B`, 100); value != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, value)
	}
}

func TestCorrectlyRetrieveV_WStringShortValue(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddByteField(`ByteField`, ``)
	generator.AddV_WStringField(`V_WStringField`, ``, 250)
	generator.AddV_StringField(`V_StringField`, ``, 250)
	recordInfo := generator.GenerateRecordInfo()
	value, isNull, err := recordInfo.GetStringValueFrom(`V_WStringField`, varFieldShortRecord)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected is not null but got null`)
	}
	if expected := strings.Repeat(`A`, 50); value != expected {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, value)
	}
}

func TestCorrectlyRetrieveVarTinyValue(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddByteField(`ByteField`, ``)
	generator.AddV_WStringField(`V_WStringField`, ``, 250)
	generator.AddV_StringField(`V_StringField`, ``, 250)
	recordInfo := generator.GenerateRecordInfo()

	value, isNull, err := recordInfo.GetStringValueFrom(`V_StringField`, varFieldTinyRecord)
	checkExpectedGetValueFrom(t, value, `B`, isNull, false, err, nil, `error retrieving tiny v_string:`)

	value, isNull, err = recordInfo.GetStringValueFrom(`V_WStringField`, varFieldTinyRecord)
	checkExpectedGetValueFrom(t, value, `A`, isNull, false, err, nil, `error retrieving tine v_wstring:`)
}

func TestGetCurrentByte(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddByteField(`MyField`, ``)
	recordInfo := generator.GenerateRecordInfo()
	_ = recordInfo.SetIntField(`MyField`, 12)

	value, isNull, err := recordInfo.GetCurrentInt(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 12 {
		t.Fatalf(`expected 12 but got %v`, value)
	}

	_ = recordInfo.SetFieldNull(`MyField`)
	value, isNull, err = recordInfo.GetCurrentInt(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestGetCurrentBool(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddBoolField(`MyField`, ``)
	recordInfo := generator.GenerateRecordInfo()
	_ = recordInfo.SetBoolField(`MyField`, true)

	value, isNull, err := recordInfo.GetCurrentBool(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if !value {
		t.Fatalf(`expected true but got %v`, value)
	}

	_ = recordInfo.SetFieldNull(`MyField`)
	value, isNull, err = recordInfo.GetCurrentBool(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestGetCurrentInt16(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddInt16Field(`MyField`, ``)
	recordInfo := generator.GenerateRecordInfo()
	_ = recordInfo.SetIntField(`MyField`, 12)

	value, isNull, err := recordInfo.GetCurrentInt(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 12 {
		t.Fatalf(`expected 12 but got %v`, value)
	}

	_ = recordInfo.SetFieldNull(`MyField`)
	value, isNull, err = recordInfo.GetCurrentInt(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestGetCurrentInt32(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddInt32Field(`MyField`, ``)
	recordInfo := generator.GenerateRecordInfo()
	_ = recordInfo.SetIntField(`MyField`, 12)

	value, isNull, err := recordInfo.GetCurrentInt(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 12 {
		t.Fatalf(`expected 12 but got %v`, value)
	}

	_ = recordInfo.SetFieldNull(`MyField`)
	value, isNull, err = recordInfo.GetCurrentInt(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestGetCurrentInt64(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddInt64Field(`MyField`, ``)
	recordInfo := generator.GenerateRecordInfo()
	_ = recordInfo.SetIntField(`MyField`, 12)

	value, isNull, err := recordInfo.GetCurrentInt(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 12 {
		t.Fatalf(`expected 12 but got %v`, value)
	}

	_ = recordInfo.SetFieldNull(`MyField`)
	value, isNull, err = recordInfo.GetCurrentInt(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestGetNullInt64FromRecordblob(t *testing.T) {
	generator1 := recordinfo.NewGenerator()
	generator1.AddInt64Field(`MyField`, ``)
	recordInfo1 := generator1.GenerateRecordInfo()

	generator2 := recordinfo.NewGenerator()
	generator2.AddInt64Field(`MyField`, ``)
	recordInfo2 := generator2.GenerateRecordInfo()

	_ = recordInfo1.SetFieldNull(`MyField`)
	record, _ := recordInfo1.GenerateRecord()

	copier, _ := recordcopier.New(recordInfo2, recordInfo1, []recordcopier.IndexMap{{
		DestinationIndex: 0,
		SourceIndex:      0,
	}})

	_ = copier.Copy(record)

	_, isNull, err := recordInfo2.GetCurrentInt(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestGetCurrentFixedDecimal(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddFixedDecimalField(`MyField`, ``, 16, 4)
	recordInfo := generator.GenerateRecordInfo()
	_ = recordInfo.SetFloatField(`MyField`, 123.45)

	value, isNull, err := recordInfo.GetCurrentFloat(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 123.45 {
		t.Fatalf(`expected 123.45 but got %v`, value)
	}

	_ = recordInfo.SetFieldNull(`MyField`)
	value, isNull, err = recordInfo.GetCurrentFloat(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestGetCurrentFloat(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddFloatField(`MyField`, ``)
	recordInfo := generator.GenerateRecordInfo()
	_ = recordInfo.SetFloatField(`MyField`, 123.45)

	value, isNull, err := recordInfo.GetCurrentFloat(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != float64(float32(123.45)) {
		t.Fatalf(`expected 123.45 but got %v`, value)
	}

	_ = recordInfo.SetFieldNull(`MyField`)
	value, isNull, err = recordInfo.GetCurrentFloat(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestGetCurrentDouble(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddDoubleField(`MyField`, ``)
	recordInfo := generator.GenerateRecordInfo()
	_ = recordInfo.SetFloatField(`MyField`, 123.45)

	value, isNull, err := recordInfo.GetCurrentFloat(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != 123.45 {
		t.Fatalf(`expected 123.45 but got %v`, value)
	}

	_ = recordInfo.SetFieldNull(`MyField`)
	value, isNull, err = recordInfo.GetCurrentFloat(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestGetCurrentString(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddStringField(`MyField`, ``, 50)
	recordInfo := generator.GenerateRecordInfo()
	_ = recordInfo.SetStringField(`MyField`, `hello world`)

	value, isNull, err := recordInfo.GetCurrentString(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != `hello world` {
		t.Fatalf(`expected 'hello world' but got '%v'`, value)
	}

	_ = recordInfo.SetFieldNull(`MyField`)
	value, isNull, err = recordInfo.GetCurrentString(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestGetCurrentV_String(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddV_StringField(`MyField`, ``, 50)
	recordInfo := generator.GenerateRecordInfo()
	_ = recordInfo.SetStringField(`MyField`, `hello world`)

	value, isNull, err := recordInfo.GetCurrentString(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != `hello world` {
		t.Fatalf(`expected 'hello world' but got '%v'`, value)
	}

	_ = recordInfo.SetFieldNull(`MyField`)
	value, isNull, err = recordInfo.GetCurrentString(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestGetCurrentWString(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddWStringField(`MyField`, ``, 50)
	recordInfo := generator.GenerateRecordInfo()
	_ = recordInfo.SetStringField(`MyField`, `hello world`)

	value, isNull, err := recordInfo.GetCurrentString(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != `hello world` {
		t.Fatalf(`expected 'hello world' but got '%v'`, value)
	}

	_ = recordInfo.SetFieldNull(`MyField`)
	value, isNull, err = recordInfo.GetCurrentString(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestGetCurrentV_WString(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddV_WStringField(`MyField`, ``, 50)
	recordInfo := generator.GenerateRecordInfo()
	_ = recordInfo.SetStringField(`MyField`, `hello world`)

	value, isNull, err := recordInfo.GetCurrentString(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != `hello world` {
		t.Fatalf(`expected 'hello world' but got '%v'`, value)
	}

	_ = recordInfo.SetFieldNull(`MyField`)
	value, isNull, err = recordInfo.GetCurrentString(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestGetCurrentDate(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddDateField(`MyField`, ``)
	recordInfo := generator.GenerateRecordInfo()
	date := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	_ = recordInfo.SetDateField(`MyField`, date)

	value, isNull, err := recordInfo.GetCurrentDate(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != date {
		t.Fatalf(`expected '2020-01-02' but got '%v'`, value)
	}

	_ = recordInfo.SetFieldNull(`MyField`)
	value, isNull, err = recordInfo.GetCurrentDate(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func TestGetCurrentDateTime(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddDateTimeField(`MyField`, ``)
	recordInfo := generator.GenerateRecordInfo()
	date := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	_ = recordInfo.SetDateField(`MyField`, date)

	value, isNull, err := recordInfo.GetCurrentDate(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected non-null but got null`)
	}
	if value != date {
		t.Fatalf(`expected '2020-01-02 03:04:05' but got '%v'`, value)
	}

	_ = recordInfo.SetFieldNull(`MyField`)
	value, isNull, err = recordInfo.GetCurrentDate(`MyField`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if !isNull {
		t.Fatalf(`expected null but got non-null`)
	}
}

func checkExpectedGetValueFrom(t *testing.T, value interface{}, expectedValue interface{}, isNull bool, expectedIsNull bool, err error, expectedErr error, msg string) {
	if err != expectedErr {
		t.Fatalf("%v expected error: %v\ngot: %v", msg, expectedErr, err)
	}
	if value != expectedValue {
		t.Fatalf("%v expected\n%v\nbut got\n%v\n", msg, expectedValue, value)
	}
	if isNull != expectedIsNull {
		t.Fatalf(`%v expected isNull=%v but got isNull=%v`, msg, expectedIsNull, isNull)
	}
}

var recordInfoXml = `<MetaInfo connection="Output">
	<RecordInfo>
		<Field name="ByteField" source="TextInput:" type="Byte"/>
		<Field name="BoolField" source="Formula: 1" type="Bool"/>
		<Field name="Int16Field" source="Formula: 16" type="Int16"/>
		<Field name="Int32Field" source="Formula: 32" type="Int32"/>
		<Field name="Int64Field" source="Formula: 64" type="Int64"/>
		<Field name="FixedDecimalField" scale="6" size="19" source="Formula: 123.45" type="FixedDecimal"/>
		<Field name="FloatField" source="Formula: 678.9" type="Float"/>
		<Field name="DoubleField" source="Formula: 0.12345" type="Double"/>
		<Field name="StringField" size="64" source="Formula: &quot;A&quot;" type="String"/>
		<Field name="WStringField" size="64" source="Formula: &quot;AB&quot;" type="WString"/>
		<Field name="V_StringShortField" size="1000" source="Formula: &quot;ABC&quot;" type="V_String"/>
		<Field name="V_StringLongField" size="2147483647" source="Formula: PadLeft(&quot;&quot;, 500, &apos;B&apos;)" type="V_String"/>
		<Field name="V_WStringShortField" size="10" source="Formula: &quot;XZY&quot;" type="V_WString"/>
		<Field name="V_WStringLongField" size="1073741823" source="Formula: PadLeft(&quot;&quot;, 500, &apos;W&apos;)" type="V_WString"/>
		<Field name="DateField" source="Formula: &apos;2020-01-01&apos;" type="Date"/>
		<Field name="DateTimeField" source="Formula: &apos;2020-02-03 04:05:06&apos;" type="DateTime"/>
	</RecordInfo>
</MetaInfo>
`

// byte with 1, v_wstring with '', v_string with ''
var varFieldEmptyStrings = recordblob.NewRecordBlob(unsafe.Pointer(&[]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}[0]))

// byte with 1, v_wstring with 'A', v_string with 'B'
var varFieldTinyRecord = recordblob.NewRecordBlob(unsafe.Pointer(&[]byte{1, 0, 65, 0, 0, 32, 66, 0, 0, 16, 0, 0, 0, 0}[0]))

// byte with 1, v_wstring with 50 A's, v_string with 100 B's
var varFieldShortRecord = recordblob.NewRecordBlob(unsafe.Pointer(&[]byte{
	1, 0, 12, 0, 0, 0, 109, 0, 0, 0, 202, 0, 0, 0, 201, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 201, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
}[0]))

// byte with 1, v_wstring with 100 A's, v_string with 200 B's
var varFieldLongRecord = recordblob.NewRecordBlob(unsafe.Pointer(&[]byte{
	1, 0, 12, 0, 0, 0, 212, 0, 0, 0, 152, 1, 0, 0, 144, 1, 0, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 144, 1, 0, 0, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
}[0]))

var varFieldNullRecord = recordblob.NewRecordBlob(unsafe.Pointer(&[]byte{
	1, 0, 1, 0, 0, 0, 1, 0, 0, 0,
}[0]))

var sampleRecord = recordblob.NewRecordBlob(unsafe.Pointer(&[]byte{
	1, 0, 1, 16, 0, 0, 32, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 49, 50, 51, 46, 52, 53, 48, 48, 48, 48, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 154, 185, 41, 68, 0, 124, 242, 176, 80, 107, 154, 191, 63, 0, 65, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 65, 0, 66, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 65, 66, 67, 48, 47, 0,
	0, 0, 35, 2, 0, 0, 38, 2, 0, 0, 50, 48, 50, 48, 45, 48, 49, 45, 48, 49, 0, 50, 48, 50, 48, 45, 48, 50, 45, 48, 51,
	32, 48, 52, 58, 48, 53, 58, 48, 54, 0, 235, 5, 0, 0, 232, 3, 0, 0, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 13, 88, 0, 90, 0,
	89, 0, 208, 7, 0, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87,
	0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0,
	87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87,
	0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0,
	87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87,
	0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0,
	87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87,
	0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0,
	87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87,
	0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0,
	87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87,
	0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0,
	87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87,
	0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0,
	87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87,
	0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0,
	87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87,
	0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0,
	87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87,
	0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0,
	87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87,
	0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0,
	87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87,
	0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0,
	87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87,
	0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0,
	87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87,
	0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0,
	87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87,
	0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0,
	87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0, 87, 0,
}[0]))

var nullRecord = recordblob.NewRecordBlob(unsafe.Pointer(&[]byte{
	0, 1, 2, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 115, 0, 92, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 32, 117, 20, 141, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 1, 0, 45, 0, 0, 45, 16, 136, 2, 0, 144, 102, 21, 141, 211, 1, 0, 0, 80, 121, 26, 141, 211, 1,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0,
	1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 158, 15, 221, 63, 180, 136, 2, 0, 96, 163, 147, 137, 211, 1, 0, 0, 16, 206, 146, 137, 211, 1, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 208, 15, 212, 120, 157, 136, 2, 32, 128, 33,
	25, 141, 211, 1, 0, 0, 192, 137, 20, 141, 211, 1, 0, 0, 16, 0, 0, 0, 211, 1, 0, 0, 120, 144, 161, 137, 211, 1, 0,
	0, 224, 22, 122, 137, 211, 1, 0, 0, 208, 123, 26, 141, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 67, 0, 6, 0, 4, 0, 0, 0, 254, 0, 3, 0, 7, 0, 0, 0, 136, 147, 161, 137, 211, 1, 0, 0, 200, 145, 161, 137, 211, 1,
	0, 0, 96, 21, 122, 137, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 7, 0, 0, 0, 72, 0, 6, 0, 2, 0, 0, 0, 232, 158, 161, 137, 211, 1, 0, 0, 88, 145, 161, 137, 211, 1, 0, 0, 224,
	16, 122, 137, 211, 1, 0, 0, 64, 132, 21, 141, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 165, 0,
	95, 0, 3, 0, 0, 0, 168, 0, 2, 0, 7, 0, 0, 0, 104, 122, 24, 141, 211, 1, 0, 0, 216, 122, 24, 141, 211, 1, 0, 0, 32,
	22, 122, 137, 211, 1, 0, 0, 208, 159, 161, 137, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	240, 0, 2, 0, 0, 0, 252, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 232, 97, 24, 141, 211, 1, 0, 0, 224, 22, 122,
	137, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0,
	254, 0, 3, 0, 2, 0, 0, 0, 168, 99, 24, 141, 211, 1, 0, 0, 232, 104, 24, 141, 211, 1, 0, 0, 224, 16, 122, 137, 211,
	1, 0, 0, 96, 228, 21, 141, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 80, 0, 3, 0, 0, 0,
	168, 0, 2, 0, 3, 0, 0, 0, 40, 103, 24, 141, 211, 1, 0, 0, 248, 93, 24, 141, 211, 1, 0, 0, 32, 22, 122, 137, 211, 1,
	0, 0, 32, 89, 24, 141, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 218, 0, 2, 0, 0, 0, 252,
	0, 1, 0, 1, 0, 0, 0, 40, 110, 24, 141, 211, 1, 0, 0, 184, 95, 24, 141, 211, 1, 0, 0, 32, 22, 122, 137, 211, 1, 0,
	0, 48, 121, 24, 141, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 218, 0, 2, 0, 0, 0, 252,
	0, 1, 0, 1, 0, 0, 0, 168, 113, 24, 141, 211, 1, 0, 0, 72, 102, 24, 141, 211, 1, 0, 0, 32, 22, 122, 137, 211, 1, 0,
	0, 80, 164, 21, 141, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 249, 1, 2, 0, 0, 0, 251,
	1, 1, 0, 1, 0, 0, 0, 200, 119, 24, 141, 211, 1, 0, 0, 136, 114, 24, 141, 211, 1, 0, 0, 32, 22, 122, 137, 211, 1, 0,
	0, 144, 45, 25, 141, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 23, 0, 2, 0, 0, 0, 251, 1,
	1, 0, 1, 0, 0, 0, 104, 108, 24, 141, 211, 1, 0, 0, 168, 120, 24, 141, 211, 1, 0, 0, 32, 22, 122, 137, 211, 1, 0, 0,
	112, 4, 22, 141, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 217, 1, 2, 0, 0, 0, 251, 1, 1,
	0, 1, 0, 0, 0, 40, 117, 24, 141, 211, 1, 0, 0, 104, 115, 24, 141, 211, 1, 0, 0, 32, 22, 122, 137, 211, 1, 0, 0,
	128, 68, 22, 141, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 194, 0, 208, 1, 2, 0, 0, 0, 251, 1,
	1, 0, 3, 0, 0, 0, 232, 182, 161, 137, 211, 1, 0, 0, 152, 160, 161, 137, 211, 1, 0, 0, 96, 27, 122, 137, 211, 1, 0,
	0, 64, 44, 159, 137, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 15, 0, 8, 0, 1, 1, 0, 0, 15, 0,
	96, 0, 7, 0, 0, 0, 216, 186, 161, 137, 211, 1, 0, 0, 24, 171, 161, 137, 211, 1, 0, 0, 96, 30, 122, 137, 211, 1, 0,
	0, 240, 195, 26, 141, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 15, 0, 0, 0, 1, 2, 0, 0, 15, 0,
	112, 0, 7, 0, 0, 0, 72, 173, 161, 137, 211, 1, 0, 0, 200, 169, 161, 137, 211, 1, 0, 0, 224, 19, 122, 137, 211, 1,
	0, 0, 64, 153, 24, 141, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 33, 0, 5, 0, 13, 0, 0, 0, 39,
	0, 12, 0, 7, 0, 0, 0, 24, 164, 161, 137, 211, 1, 0, 0, 104, 165, 161, 137, 211, 1, 0, 0, 216, 172, 161, 137, 211, 1,
	0, 0, 40, 174,
}[0]))
