package recordinfo_test

import (
	"goalteryx/recordinfo"
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
		FieldType string
	}
	expectedFields := []expectedStruct{
		{`ByteField`, 1, 0, recordinfo.ByteType},
		{`BoolField`, 1, 0, recordinfo.BoolType},
		{`Int16Field`, 2, 0, recordinfo.Int16Type},
		{`Int32Field`, 4, 0, recordinfo.Int32Type},
		{`Int64Field`, 8, 0, recordinfo.Int64Type},
		{`FixedDecimalField`, 19, 6, recordinfo.FixedDecimalType},
		{`FloatField`, 4, 0, recordinfo.FloatType},
		{`DoubleField`, 8, 0, recordinfo.DoubleType},
		{`StringField`, 64, 0, recordinfo.StringType},
		{`WStringField`, 64, 0, recordinfo.WStringType},
		{`V_StringShortField`, 1000, 0, recordinfo.V_StringType},
		{`V_StringLongField`, 2147483647, 0, recordinfo.V_StringType},
		{`V_WStringShortField`, 10, 0, recordinfo.V_WStringType},
		{`V_WStringLongField`, 1073741823, 0, recordinfo.V_WStringType},
		{`DateField`, 10, 0, recordinfo.DateType},
		{`DateTimeField`, 19, 0, recordinfo.DateTimeType},
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

	value, isNull, err := recordInfo.GetByteValueFrom(`ByteField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, byte(1), isNull, false, err, nil)

	value, isNull, err = recordInfo.GetByteValueFrom(`ByteField`, nullRecord)
	checkExpectedGetValueFrom(t, value, byte(0), isNull, true, err, nil)
}

func TestCorrectlyRetrieveBoolValue(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetBoolValueFrom(`BoolField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, true, isNull, false, err, nil)

	value, isNull, err = recordInfo.GetBoolValueFrom(`BoolField`, nullRecord)
	checkExpectedGetValueFrom(t, value, false, isNull, true, err, nil)
}

func TestCorrectlyRetrieveInt16Value(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetInt16ValueFrom(`Int16Field`, sampleRecord)
	checkExpectedGetValueFrom(t, value, int16(16), isNull, false, err, nil)

	value, isNull, err = recordInfo.GetInt16ValueFrom(`Int16Field`, nullRecord)
	checkExpectedGetValueFrom(t, value, int16(0), isNull, true, err, nil)
}

func TestCorrectlyRetrieveInt32Value(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetInt32ValueFrom(`Int32Field`, sampleRecord)
	checkExpectedGetValueFrom(t, value, int32(32), isNull, false, err, nil)

	value, isNull, err = recordInfo.GetInt32ValueFrom(`Int32Field`, nullRecord)
	checkExpectedGetValueFrom(t, value, int32(0), isNull, true, err, nil)
}

func TestCorrectlyRetrieveInt64Value(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetInt64ValueFrom(`Int64Field`, sampleRecord)
	checkExpectedGetValueFrom(t, value, int64(64), isNull, false, err, nil)

	value, isNull, err = recordInfo.GetInt64ValueFrom(`Int64Field`, nullRecord)
	checkExpectedGetValueFrom(t, value, int64(0), isNull, true, err, nil)
}

func TestCorrectlyRetrieveFixedDecimalValue(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetFixedDecimalValueFrom(`FixedDecimalField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, 123.450000, isNull, false, err, nil)

	value, isNull, err = recordInfo.GetFixedDecimalValueFrom(`FixedDecimalField`, nullRecord)
	checkExpectedGetValueFrom(t, value, 0.0, isNull, true, err, nil)
}

func TestCorrectlyRetrieveFloatValue(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetFloatValueFrom(`FloatField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, float32(678.9), isNull, false, err, nil)

	value, isNull, err = recordInfo.GetFloatValueFrom(`FloatField`, nullRecord)
	checkExpectedGetValueFrom(t, value, float32(0.0), isNull, true, err, nil)
}

func TestCorrectlyRetrieveDoubleValue(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetDoubleValueFrom(`DoubleField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, 0.12345, isNull, false, err, nil)

	value, isNull, err = recordInfo.GetDoubleValueFrom(`DoubleField`, nullRecord)
	checkExpectedGetValueFrom(t, value, 0.0, isNull, true, err, nil)
}

func TestCorrectlyRetrieveStringValue(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetStringValueFrom(`StringField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, `A`, isNull, false, err, nil)

	value, isNull, err = recordInfo.GetStringValueFrom(`StringField`, nullRecord)
	checkExpectedGetValueFrom(t, value, ``, isNull, true, err, nil)
}

func TestCorrectlyRetrieveWStringValue(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetWStringValueFrom(`WStringField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, `AB`, isNull, false, err, nil)

	value, isNull, err = recordInfo.GetWStringValueFrom(`WStringField`, nullRecord)
	checkExpectedGetValueFrom(t, value, ``, isNull, true, err, nil)
}

var zeroDate = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

func TestCorrectlyRetrieveDateValue(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetDateValueFrom(`DateField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), isNull, false, err, nil)

	value, isNull, err = recordInfo.GetDateValueFrom(`DateField`, nullRecord)
	checkExpectedGetValueFrom(t, value, zeroDate, isNull, true, err, nil)
}

func TestCorrectlyRetrieveDateTimeValue(t *testing.T) {
	recordInfo, err := recordinfo.FromXml(recordInfoXml)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, isNull, err := recordInfo.GetDateTimeValueFrom(`DateTimeField`, sampleRecord)
	checkExpectedGetValueFrom(t, value, time.Date(2020, 2, 3, 4, 5, 6, 0, time.UTC), isNull, false, err, nil)

	value, isNull, err = recordInfo.GetDateTimeValueFrom(`DateTimeField`, nullRecord)
	checkExpectedGetValueFrom(t, value, zeroDate, isNull, true, err, nil)
}

func checkExpectedGetValueFrom(t *testing.T, value interface{}, expectedValue interface{}, isNull bool, expectedIsNull bool, err error, expectedErr error) {
	if err != expectedErr {
		t.Fatalf("expected error: %v\ngot: %v", expectedErr, err)
	}
	if value != expectedValue {
		t.Fatalf(`expected '%v' but got '%v'`, expectedValue, value)
	}
	if isNull != expectedIsNull {
		t.Fatalf(`expected isNull=%v but got isNull=%v`, expectedIsNull, isNull)
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

var fixedFieldRecord = unsafe.Pointer(&[]byte{
	1, 0, 1, 16, 0, 0, 32, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 49, 50, 51, 46, 52, 53, 48, 48, 48, 48, 0, 0, 116,
	0, 108, 0, 97, 0, 114, 0, 154, 185, 41, 68, 0, 124, 242, 176, 80, 107, 154, 191, 63, 0, 65, 0, 97, 0, 116, 0, 97,
	0, 92, 0, 76, 0, 111, 0, 99, 0, 97, 0, 108, 0, 92, 0, 84, 0, 101, 0, 109, 0, 112, 0, 92, 0, 69, 0, 110, 0, 103, 0,
	105, 0, 110, 0, 101, 0, 95, 0, 54, 0, 50, 0, 52, 0, 56, 0, 95, 0, 50, 0, 100, 0, 57, 0, 98, 0, 0, 65, 0, 66, 0, 0,
	0, 0, 101, 0, 55, 0, 101, 0, 52, 0, 52, 0, 53, 0, 55, 0, 102, 0, 57, 0, 98, 0, 98, 0, 57, 0, 51, 0, 99, 0, 56, 0,
	102, 0, 102, 0, 49, 0, 54, 0, 48, 0, 57, 0, 100, 0, 48, 0, 56, 0, 95, 0, 92, 0, 69, 0, 110, 0, 103, 0, 105, 0, 110,
	0, 101, 0, 95, 0, 50, 0, 56, 0, 49, 0, 54, 0, 95, 0, 100, 0, 50, 0, 49, 0, 56, 0, 55, 0, 56, 0, 99, 0, 56, 0, 52,
	0, 98, 0, 55, 0, 49, 0, 52, 0, 57, 0, 100, 0, 49, 0, 56, 0, 49, 0, 98, 0, 48, 0, 56, 0, 99, 0, 97, 0, 50, 48, 50,
	48, 45, 48, 49, 45, 48, 49, 0, 50, 48, 50, 48, 45, 48, 50, 45, 48, 51, 32, 48, 52, 58, 48, 53, 58, 48, 54, 0,
})

var sampleRecord = unsafe.Pointer(&[]byte{
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
}[0])

var nullRecord = unsafe.Pointer(&[]byte{
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
}[0])
