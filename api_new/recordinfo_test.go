package api_new

import (
	"math"
	"testing"
	"time"
	"unsafe"
)

func TestNewIncomingRecordInfo(t *testing.T) {
	config := `<MetaInfo connection="Output">
<RecordInfo>
	<Field name="Field1" source="TextInput:" type="Byte"/>
	<Field name="Field2" size="1" source="TextInput:" type="String"/>
</RecordInfo>
</MetaInfo>`
	recordInfo, err := incomingRecordInfoFromString(config)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if fields := recordInfo.NumFields(); fields != 2 {
		t.Fatalf(`expected 2 fields but got %v`, fields)
	}
}

func TestNewIncomingRecordInfoWithoutMetaInfo(t *testing.T) {
	config := `<RecordInfo>
	<Field name="Field1" source="TextInput:" type="Byte"/>
	<Field name="Field2" size="1" source="TextInput:" type="String"/>
</RecordInfo>`
	recordInfo, err := incomingRecordInfoFromString(config)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if fields := recordInfo.NumFields(); fields != 2 {
		t.Fatalf(`expected 2 fields but got %v`, fields)
	}
}

func TestIncomingFieldDoesNotExist(t *testing.T) {
	config := `<RecordInfo>
	<Field name="Field1" source="TextInput:" type="Byte"/>
	<Field name="Field2" size="1" source="TextInput:" type="String"/>
</RecordInfo>`
	recordInfo, _ := incomingRecordInfoFromString(config)
	_, err := recordInfo.GetIntField(`Hello world`)
	if err == nil {
		t.Fatalf(`expected an error but got none`)
	}
	t.Logf(err.Error())
}

func TestGetByteValue(t *testing.T) {
	config := `<RecordInfo>
	<Field name="Field1" type="Bool"/>
	<Field name="Field2" type="Byte"/>
</RecordInfo>`
	recordInfo, _ := incomingRecordInfoFromString(config)
	field, err := recordInfo.GetIntField(`Field2`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}

	record := unsafe.Pointer(&[]byte{2, 35, 0}[0])
	value, isNull := field.GetValue(record)
	if value != 35 {
		t.Fatalf(`expected 35 but got %v`, value)
	}
	if isNull {
		t.Fatalf(`expected not null but got null`)
	}

	record = unsafe.Pointer(&[]byte{2, 35, 1}[0])
	value, isNull = field.GetValue(record)
	if value != 0 {
		t.Fatalf(`expected 0 but got %v`, value)
	}
	if !isNull {
		t.Fatalf(`expected null but got not null`)
	}
}

func TestGetInt16Value(t *testing.T) {
	config := `<RecordInfo>
	<Field name="Field1" type="Bool"/>
	<Field name="Field2" type="Int16"/>
</RecordInfo>`
	recordInfo, _ := incomingRecordInfoFromString(config)
	field, err := recordInfo.GetIntField(`Field2`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}

	record := unsafe.Pointer(&[]byte{2, 255, 127, 0}[0])
	value, isNull := field.GetValue(record)
	if value != 32767 {
		t.Fatalf(`expected 32767 but got %v`, value)
	}
	if isNull {
		t.Fatalf(`expected not null but got null`)
	}

	record = unsafe.Pointer(&[]byte{2, 255, 127, 1}[0])
	value, isNull = field.GetValue(record)
	if value != 0 {
		t.Fatalf(`expected 0 but got %v`, value)
	}
	if !isNull {
		t.Fatalf(`expected null but got not null`)
	}
}

func TestGetInt32Value(t *testing.T) {
	config := `<RecordInfo>
	<Field name="Field1" type="Bool"/>
	<Field name="Field2" type="Int32"/>
</RecordInfo>`
	recordInfo, _ := incomingRecordInfoFromString(config)
	field, err := recordInfo.GetIntField(`Field2`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}

	record := unsafe.Pointer(&[]byte{2, 0, 16, 0, 0, 0}[0])
	value, isNull := field.GetValue(record)
	if value != 4096 {
		t.Fatalf(`expected 4096 but got %v`, value)
	}
	if isNull {
		t.Fatalf(`expected not null but got null`)
	}

	record = unsafe.Pointer(&[]byte{2, 0, 16, 0, 0, 1}[0])
	value, isNull = field.GetValue(record)
	if value != 0 {
		t.Fatalf(`expected 0 but got %v`, value)
	}
	if !isNull {
		t.Fatalf(`expected null but got not null`)
	}
}

func TestGetInt64Value(t *testing.T) {
	config := `<RecordInfo>
	<Field name="Field1" type="Bool"/>
	<Field name="Field2" type="Int64"/>
</RecordInfo>`
	recordInfo, _ := incomingRecordInfoFromString(config)
	field, err := recordInfo.GetIntField(`Field2`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}

	record := unsafe.Pointer(&[]byte{2, 86, 85, 85, 85, 85, 85, 255, 255, 0}[0])
	value, isNull := field.GetValue(record)
	if value != -187649984473770 {
		t.Fatalf(`expected -187649984473770 but got %v`, value)
	}
	if isNull {
		t.Fatalf(`expected not null but got null`)
	}

	record = unsafe.Pointer(&[]byte{2, 86, 85, 85, 85, 85, 85, 255, 255, 1}[0])
	value, isNull = field.GetValue(record)
	if value != 0 {
		t.Fatalf(`expected 0 but got %v`, value)
	}
	if !isNull {
		t.Fatalf(`expected null but got not null`)
	}
}

func TestGetFloatValue(t *testing.T) {
	config := `<RecordInfo>
	<Field name="Field1" type="Bool"/>
	<Field name="Field2" type="Float"/>
</RecordInfo>`
	recordInfo, _ := incomingRecordInfoFromString(config)
	field, err := recordInfo.GetFloatField(`Field2`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}

	record := unsafe.Pointer(&[]byte{2, 102, 230, 246, 66, 0}[0])
	value, isNull := field.GetValue(record)
	if math.Abs(value-123.45) > 0.00001 {
		t.Fatalf(`expected 123.45 but got %v`, value)
	}
	if isNull {
		t.Fatalf(`expected not null but got null`)
	}

	record = unsafe.Pointer(&[]byte{2, 102, 230, 246, 66, 1}[0])
	value, isNull = field.GetValue(record)
	if value != 0.0 {
		t.Fatalf(`expected 0 but got %v`, value)
	}
	if !isNull {
		t.Fatalf(`expected null but got not null`)
	}
}

func TestGetDoubleValue(t *testing.T) {
	config := `<RecordInfo>
	<Field name="Field1" type="Bool"/>
	<Field name="Field2" type="Double"/>
</RecordInfo>`
	recordInfo, _ := incomingRecordInfoFromString(config)
	field, err := recordInfo.GetFloatField(`Field2`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}

	record := unsafe.Pointer(&[]byte{2, 205, 204, 204, 204, 204, 220, 94, 64, 0}[0])
	value, isNull := field.GetValue(record)
	if math.Abs(value-123.45) > 0.00001 {
		t.Fatalf(`expected 123.45 but got %v`, value)
	}
	if isNull {
		t.Fatalf(`expected not null but got null`)
	}

	record = unsafe.Pointer(&[]byte{2, 205, 204, 204, 204, 204, 220, 94, 64, 1}[0])
	value, isNull = field.GetValue(record)
	if value != 0.0 {
		t.Fatalf(`expected 0 but got %v`, value)
	}
	if !isNull {
		t.Fatalf(`expected null but got not null`)
	}
}

func TestGetFixedDecimalValue(t *testing.T) {
	config := `<RecordInfo>
	<Field name="Field1" type="Bool"/>
	<Field name="Field2" type="FixedDecimal" size="19" scale="2" />
</RecordInfo>`
	recordInfo, _ := incomingRecordInfoFromString(config)
	field, err := recordInfo.GetFloatField(`Field2`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}

	record := unsafe.Pointer(&[]byte{2, 49, 50, 51, 46, 52, 53, 0, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 0}[0])
	value, isNull := field.GetValue(record)
	if value != 123.45 {
		t.Fatalf(`expected 123.45 but got %v`, value)
	}
	if isNull {
		t.Fatalf(`expected not null but got null`)
	}

	record = unsafe.Pointer(&[]byte{2, 49, 50, 51, 46, 52, 53, 0, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 1}[0])
	value, isNull = field.GetValue(record)
	if value != 0.0 {
		t.Fatalf(`expected 0 but got %v`, value)
	}
	if !isNull {
		t.Fatalf(`expected null but got not null`)
	}
}

func TestGetBoolValue(t *testing.T) {
	config := `<RecordInfo>
	<Field name="Field1" type="Bool"/>
	<Field name="Field2" type="Bool" />
</RecordInfo>`
	recordInfo, _ := incomingRecordInfoFromString(config)
	field, err := recordInfo.GetBoolField(`Field2`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}

	record := unsafe.Pointer(&[]byte{2, 1}[0])
	value, isNull := field.GetValue(record)
	if !value {
		t.Fatal(`expected true but got false`)
	}
	if isNull {
		t.Fatalf(`expected not null but got null`)
	}

	record = unsafe.Pointer(&[]byte{2, 2}[0])
	value, isNull = field.GetValue(record)
	if value {
		t.Fatal(`expected false but got true`)
	}
	if !isNull {
		t.Fatalf(`expected null but got not null`)
	}
}

func TestGetDateValue(t *testing.T) {
	config := `<RecordInfo>
	<Field name="Field1" type="Bool"/>
	<Field name="Field2" type="Date" />
</RecordInfo>`
	recordInfo, _ := incomingRecordInfoFromString(config)
	field, err := recordInfo.GetTimeField(`Field2`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}

	record := unsafe.Pointer(&[]byte{2, 50, 48, 50, 48, 45, 48, 49, 45, 48, 51, 0}[0])
	value, isNull := field.GetValue(record)
	if value != time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC) {
		t.Fatalf(`expected '2020-01-03' but got '%v'`, value)
	}
	if isNull {
		t.Fatalf(`expected not null but got null`)
	}

	record = unsafe.Pointer(&[]byte{2, 50, 48, 50, 48, 45, 48, 49, 45, 48, 51, 1}[0])
	value, isNull = field.GetValue(record)
	if empty := (time.Time{}); value != empty {
		t.Fatalf(`expected '%v' but got '%v'`, empty, value)
	}
	if !isNull {
		t.Fatalf(`expected null but got not null`)
	}
}
