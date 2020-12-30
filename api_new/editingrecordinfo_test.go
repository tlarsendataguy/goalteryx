package api_new_test

import (
	"github.com/tlarsen7572/goalteryx/api_new"
	"testing"
	"time"
)

func TestAddBoolField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name := editor.AddBoolField(`Field1`, ``)
	if name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name)
	}
	if editor.NumFields() != 1 {
		t.Fatalf(`expected 1 but got %v`, editor.NumFields())
	}
	field := editor.Fields()[0]
	if field.Type != `Bool` || field.Size != 1 || field.Scale != 0 {
		t.Fatalf(`expected Bool size 1 scale 0 but got %v size %v scale %v`, field.Type, field.Size, field.Scale)
	}
}

func TestAddDuplicateBoolFields(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name1 := editor.AddBoolField(`Field1`, ``)
	name2 := editor.AddBoolField(`Field1`, ``)
	if name1 != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name1)
	}
	if name2 != `Field12` {
		t.Fatalf(`expected 'Field12' but got '%v'`, name2)
	}
	if editor.NumFields() != 2 {
		t.Fatalf(`expected 2 but got %v`, editor.NumFields())
	}
}

func TestInsertField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddBoolField(`Field1`, `blah blah`)
	editor.AddBoolField(`Field2`, `blah blah`, api_new.InsertAt(0))
	if editor.NumFields() != 2 {
		t.Fatalf(`expected 2 fields but got %v`, editor.NumFields())
	}
	for index, field := range editor.Fields() {
		if index == 0 && field.Name != `Field2` {
			t.Fatalf(`expected 'Field2' at index 0 but got '%v'`, field.Name)
		}
		if index == 1 && field.Name != `Field1` {
			t.Fatalf(`expected 'Field1' at index 1 but got '%v'`, field.Name)
		}
	}
}

func TestAddByteField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name := editor.AddByteField(`Field1`, ``)
	if name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name)
	}
	if editor.NumFields() != 1 {
		t.Fatalf(`expected 1 but got %v`, editor.NumFields())
	}
	field := editor.Fields()[0]
	if field.Type != `Byte` || field.Size != 1 || field.Scale != 0 {
		t.Fatalf(`expected Byte size 1 scale 0 but got %v size %v scale %v`, field.Type, field.Size, field.Scale)
	}
}

func TestAddInt16Field(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name := editor.AddInt16Field(`Field1`, ``)
	if name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name)
	}
	if editor.NumFields() != 1 {
		t.Fatalf(`expected 1 but got %v`, editor.NumFields())
	}
	field := editor.Fields()[0]
	if field.Type != `Int16` || field.Size != 2 || field.Scale != 0 {
		t.Fatalf(`expected Int16 size 2 scale 0 but got %v size %v scale %v`, field.Type, field.Size, field.Scale)
	}
}

func TestAddInt32Field(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name := editor.AddInt32Field(`Field1`, ``)
	if name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name)
	}
	if editor.NumFields() != 1 {
		t.Fatalf(`expected 1 but got %v`, editor.NumFields())
	}
	field := editor.Fields()[0]
	if field.Type != `Int32` || field.Size != 4 || field.Scale != 0 {
		t.Fatalf(`expected Int32 size 4 scale 0 but got %v size %v scale %v`, field.Type, field.Size, field.Scale)
	}
}

func TestAddInt64Field(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name := editor.AddInt64Field(`Field1`, ``)
	if name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name)
	}
	if editor.NumFields() != 1 {
		t.Fatalf(`expected 1 but got %v`, editor.NumFields())
	}
	field := editor.Fields()[0]
	if field.Type != `Int64` || field.Size != 8 || field.Scale != 0 {
		t.Fatalf(`expected Int64 size 8 scale 0 but got %v size %v scale %v`, field.Type, field.Size, field.Scale)
	}
}

func TestAddFloatField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name := editor.AddFloatField(`Field1`, ``)
	if name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name)
	}
	if editor.NumFields() != 1 {
		t.Fatalf(`expected 1 but got %v`, editor.NumFields())
	}
	field := editor.Fields()[0]
	if field.Type != `Float` || field.Size != 4 || field.Scale != 0 {
		t.Fatalf(`expected Float size 4 scale 0 but got %v size %v scale %v`, field.Type, field.Size, field.Scale)
	}
}

func TestAddDoubleField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name := editor.AddDoubleField(`Field1`, ``)
	if name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name)
	}
	if editor.NumFields() != 1 {
		t.Fatalf(`expected 1 but got %v`, editor.NumFields())
	}
	field := editor.Fields()[0]
	if field.Type != `Double` || field.Size != 8 || field.Scale != 0 {
		t.Fatalf(`expected Double size 8 scale 0 but got %v size %v scale %v`, field.Type, field.Size, field.Scale)
	}
}

func TestAddFixedDecimalField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name := editor.AddFixedDecimalField(`Field1`, ``, 19, 2)
	if name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name)
	}
	if editor.NumFields() != 1 {
		t.Fatalf(`expected 1 but got %v`, editor.NumFields())
	}
	field := editor.Fields()[0]
	if field.Type != `FixedDecimal` || field.Size != 19 || field.Scale != 2 {
		t.Fatalf(`expected FixedDecimal size 19 scale 2 but got %v size %v scale %v`, field.Type, field.Size, field.Scale)
	}
}

func TestAddStringField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name := editor.AddStringField(`Field1`, ``, 15)
	if name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name)
	}
	if editor.NumFields() != 1 {
		t.Fatalf(`expected 1 but got %v`, editor.NumFields())
	}
	field := editor.Fields()[0]
	if field.Type != `String` || field.Size != 15 || field.Scale != 0 {
		t.Fatalf(`expected String size 15 scale 0 but got %v size %v scale %v`, field.Type, field.Size, field.Scale)
	}
}

func TestAddWStringDecimalField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name := editor.AddWStringField(`Field1`, ``, 100)
	if name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name)
	}
	if editor.NumFields() != 1 {
		t.Fatalf(`expected 1 but got %v`, editor.NumFields())
	}
	field := editor.Fields()[0]
	if field.Type != `WString` || field.Size != 100 || field.Scale != 0 {
		t.Fatalf(`expected WString size 100 scale 0 but got %v size %v scale %v`, field.Type, field.Size, field.Scale)
	}
}

func TestAddV_StringField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name := editor.AddV_StringField(`Field1`, ``, 10000)
	if name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name)
	}
	if editor.NumFields() != 1 {
		t.Fatalf(`expected 1 but got %v`, editor.NumFields())
	}
	field := editor.Fields()[0]
	if field.Type != `V_String` || field.Size != 10000 || field.Scale != 0 {
		t.Fatalf(`expected V_String size 10000 scale 0 but got %v size %v scale %v`, field.Type, field.Size, field.Scale)
	}
}

func TestAddV_WStringField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name := editor.AddV_WStringField(`Field1`, ``, 256)
	if name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name)
	}
	if editor.NumFields() != 1 {
		t.Fatalf(`expected 1 but got %v`, editor.NumFields())
	}
	field := editor.Fields()[0]
	if field.Type != `V_WString` || field.Size != 256 || field.Scale != 0 {
		t.Fatalf(`expected V_WString size 256 scale 0 but got %v size %v scale %v`, field.Type, field.Size, field.Scale)
	}
}

func TestAddBlobField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name := editor.AddBlobField(`Field1`, ``, 100000)
	if name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name)
	}
	if editor.NumFields() != 1 {
		t.Fatalf(`expected 1 but got %v`, editor.NumFields())
	}
	field := editor.Fields()[0]
	if field.Type != `Blob` || field.Size != 100000 || field.Scale != 0 {
		t.Fatalf(`expected Blob size 100000 scale 0 but got %v size %v scale %v`, field.Type, field.Size, field.Scale)
	}
}

func TestSpatialObjField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name := editor.AddSpatialObjField(`Field1`, ``, 200000)
	if name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name)
	}
	if editor.NumFields() != 1 {
		t.Fatalf(`expected 1 but got %v`, editor.NumFields())
	}
	field := editor.Fields()[0]
	if field.Type != `SpatialObj` || field.Size != 200000 || field.Scale != 0 {
		t.Fatalf(`expected SpatialObj size 200000 scale 0 but got %v size %v scale %v`, field.Type, field.Size, field.Scale)
	}
}

func TestDateField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name := editor.AddDateField(`Field1`, ``)
	if name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name)
	}
	if editor.NumFields() != 1 {
		t.Fatalf(`expected 1 but got %v`, editor.NumFields())
	}
	field := editor.Fields()[0]
	if field.Type != `Date` || field.Size != 10 || field.Scale != 0 {
		t.Fatalf(`expected Date size 10 scale 0 but got %v size %v scale %v`, field.Type, field.Size, field.Scale)
	}
}

func TestDateTimeField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name := editor.AddDateTimeField(`Field1`, ``)
	if name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name)
	}
	if editor.NumFields() != 1 {
		t.Fatalf(`expected 1 but got %v`, editor.NumFields())
	}
	field := editor.Fields()[0]
	if field.Type != `DateTime` || field.Size != 19 || field.Scale != 0 {
		t.Fatalf(`expected DateTime size 19 scale 0 but got %v size %v scale %v`, field.Type, field.Size, field.Scale)
	}
}

func TestOutgoingBoolField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddBoolField(`Field1`, ``)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetBoolField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	field.SetBool(true)
	if currentValue, isNull := field.GetCurrentBool(); currentValue != true || isNull {
		t.Fatalf(`expected true and not null but got %v and %v`, currentValue, isNull)
	}
	field.SetNullBool()
	if currentValue, isNull := field.GetCurrentBool(); currentValue != false || isNull != true {
		t.Fatalf(`expected false and null but got %v and %v`, currentValue, isNull)
	}
}

func TestOutgoingByteField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddByteField(`Field1`, ``)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetIntField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := 45
	field.SetInt(expectedValue)
	if currentValue, isNull := field.GetCurrentInt(); currentValue != expectedValue || isNull {
		t.Fatalf(`expected %v and not null but got %v and %v`, expectedValue, currentValue, isNull)
	}
	field.SetNullInt()
	if currentValue, isNull := field.GetCurrentInt(); currentValue != 0 || isNull != true {
		t.Fatalf(`expected 0 and null but got %v and %v`, currentValue, isNull)
	}
	field.SetInt(10000)
	currentValue, isNull := field.GetCurrentInt()
	t.Logf(`value %v and null=%v`, currentValue, isNull)
}

func TestOutgoingInt16Field(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddInt16Field(`Field1`, ``)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetIntField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := 500
	field.SetInt(expectedValue)
	if currentValue, isNull := field.GetCurrentInt(); currentValue != expectedValue || isNull {
		t.Fatalf(`expected %v and not null but got %v and %v`, expectedValue, currentValue, isNull)
	}
	field.SetNullInt()
	if currentValue, isNull := field.GetCurrentInt(); currentValue != 0 || isNull != true {
		t.Fatalf(`expected 0 and null but got %v and %v`, currentValue, isNull)
	}
}

func TestOutgoingInt32Field(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddInt32Field(`Field1`, ``)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetIntField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := 50000
	field.SetInt(expectedValue)
	if currentValue, isNull := field.GetCurrentInt(); currentValue != expectedValue || isNull {
		t.Fatalf(`expected %v and not null but got %v and %v`, expectedValue, currentValue, isNull)
	}
	field.SetNullInt()
	if currentValue, isNull := field.GetCurrentInt(); currentValue != 0 || isNull != true {
		t.Fatalf(`expected 0 and null but got %v and %v`, currentValue, isNull)
	}
}

func TestOutgoingInt64Field(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddInt64Field(`Field1`, ``)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetIntField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := -500000000
	field.SetInt(expectedValue)
	if currentValue, isNull := field.GetCurrentInt(); currentValue != expectedValue || isNull {
		t.Fatalf(`expected %v and not null but got %v and %v`, expectedValue, currentValue, isNull)
	}
	field.SetNullInt()
	if currentValue, isNull := field.GetCurrentInt(); currentValue != 0 || isNull != true {
		t.Fatalf(`expected 0 and null but got %v and %v`, currentValue, isNull)
	}
}

func TestOutgoingFloatField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddFloatField(`Field1`, ``)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetFloatField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := 1.25
	field.SetFloat(expectedValue)
	if currentValue, isNull := field.GetCurrentFloat(); currentValue != expectedValue || isNull {
		t.Fatalf(`expected %v and not null but got %v and %v`, expectedValue, currentValue, isNull)
	}
	field.SetNullFloat()
	if currentValue, isNull := field.GetCurrentFloat(); currentValue != 0 || isNull != true {
		t.Fatalf(`expected 0 and null but got %v and %v`, currentValue, isNull)
	}
}

func TestOutgoingDoubleField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddDoubleField(`Field1`, ``)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetFloatField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := 10456.25
	field.SetFloat(expectedValue)
	if currentValue, isNull := field.GetCurrentFloat(); currentValue != expectedValue || isNull {
		t.Fatalf(`expected %v and not null but got %v and %v`, expectedValue, currentValue, isNull)
	}
	field.SetNullFloat()
	if currentValue, isNull := field.GetCurrentFloat(); currentValue != 0 || isNull != true {
		t.Fatalf(`expected 0 and null but got %v and %v`, currentValue, isNull)
	}
}

func TestOutgoingFixedDecimalField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddFixedDecimalField(`Field1`, ``, 19, 2)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetFloatField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := 123.4
	field.SetFloat(expectedValue)
	if currentValue, isNull := field.GetCurrentFloat(); currentValue != expectedValue || isNull {
		t.Fatalf(`expected %v and not null but got %v and %v`, expectedValue, currentValue, isNull)
	}
	field.SetNullFloat()
	if currentValue, isNull := field.GetCurrentFloat(); currentValue != 0 || isNull != true {
		t.Fatalf(`expected 0 and null but got %v and %v`, currentValue, isNull)
	}
}

func TestTruncateDecimals(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddFixedDecimalField(`Field1`, ``, 19, 2)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetFloatField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	field.SetFloat(123.456)
	if currentValue, isNull := field.GetCurrentFloat(); currentValue != 123.46 || isNull {
		t.Fatalf(`expected 123.45 and not null but got %v and %v`, currentValue, isNull)
	}
	field.SetNullFloat()
	if currentValue, isNull := field.GetCurrentFloat(); currentValue != 0 || isNull != true {
		t.Fatalf(`expected 0 and null but got %v and %v`, currentValue, isNull)
	}
}

func TestTruncateNumber(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddFixedDecimalField(`Field1`, ``, 5, 2)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetFloatField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	field.SetFloat(123.45)
	if currentValue, isNull := field.GetCurrentFloat(); currentValue != 123.4 || isNull {
		t.Fatalf(`expected 123.4 and not null but got %v and %v`, currentValue, isNull)
	}
	field.SetNullFloat()
	if currentValue, isNull := field.GetCurrentFloat(); currentValue != 0 || isNull != true {
		t.Fatalf(`expected 0 and null but got %v and %v`, currentValue, isNull)
	}
	field.SetFloat(-123.45)
	if currentValue, isNull := field.GetCurrentFloat(); currentValue != -123 || isNull {
		t.Fatalf(`expected -123 and not null but got %v and %v`, currentValue, isNull)
	}
	field.SetNullFloat()
	if currentValue, isNull := field.GetCurrentFloat(); currentValue != 0 || isNull != true {
		t.Fatalf(`expected 0 and null but got %v and %v`, currentValue, isNull)
	}
}

func TestOutgoingDateField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddDateField(`Field1`, ``)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetDatetimeField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	field.SetDateTime(expectedValue)
	if currentValue, isNull := field.GetCurrentDateTime(); currentValue != expectedValue || isNull {
		t.Fatalf(`expected %v and not null but got %v and %v`, expectedValue, currentValue, isNull)
	}
	field.SetNullDateTime()
	if currentValue, isNull := field.GetCurrentDateTime(); currentValue != time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC) || isNull != true {
		t.Fatalf(`expected 0000-00-00 and null but got %v and %v`, currentValue, isNull)
	}
}

func TestOutgoingDatetimeField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddDateTimeField(`Field1`, ``)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetDatetimeField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	field.SetDateTime(expectedValue)
	if currentValue, isNull := field.GetCurrentDateTime(); currentValue != expectedValue || isNull {
		t.Fatalf(`expected %v and not null but got %v and %v`, expectedValue, currentValue, isNull)
	}
	field.SetNullDateTime()
	if currentValue, isNull := field.GetCurrentDateTime(); currentValue != time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC) || isNull != true {
		t.Fatalf(`expected 0000-00-00 and null but got %v and %v`, currentValue, isNull)
	}
}

func TestOutgoingStringField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddStringField(`Field1`, ``, 20)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetStringField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := `hello world`
	field.SetString(expectedValue)
	if currentValue, isNull := field.GetCurrentString(); currentValue != expectedValue || isNull {
		t.Fatalf(`expected '%v' and not null but got '%v' and %v`, expectedValue, currentValue, isNull)
	}
	field.SetNullString()
	if currentValue, isNull := field.GetCurrentString(); currentValue != `` || isNull != true {
		t.Fatalf(`expected '' and null but got '%v' and %v`, currentValue, isNull)
	}
}

func TestStringSameSizeAsField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddStringField(`Field1`, ``, 10)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetStringField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := `0123456789`
	field.SetString(expectedValue)
	if currentValue, isNull := field.GetCurrentString(); currentValue != expectedValue || isNull {
		t.Fatalf(`expected '%v' and not null but got '%v' and %v`, expectedValue, currentValue, isNull)
	}
}

func TestStringLargerThanField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddStringField(`Field1`, ``, 10)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetStringField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := `blah blah blah`
	field.SetString(expectedValue)
	if currentValue, isNull := field.GetCurrentString(); currentValue != `blah blah ` || isNull {
		t.Fatalf(`expected 'blah blah ' and not null but got '%v' and %v`, currentValue, isNull)
	}
}

func TestOutgoingWStringField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddWStringField(`Field1`, ``, 20)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetStringField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := `hello world`
	field.SetString(expectedValue)
	if currentValue, isNull := field.GetCurrentString(); currentValue != expectedValue || isNull {
		t.Fatalf(`expected '%v' and not null but got '%v' and %v`, expectedValue, currentValue, isNull)
	}
	field.SetNullString()
	if currentValue, isNull := field.GetCurrentString(); currentValue != `` || isNull != true {
		t.Fatalf(`expected '' and null but got '%v' and %v`, currentValue, isNull)
	}
}

func TestWStringSameSizeAsField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddWStringField(`Field1`, ``, 10)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetStringField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := `0123456789`
	field.SetString(expectedValue)
	if currentValue, isNull := field.GetCurrentString(); currentValue != expectedValue || isNull {
		t.Fatalf(`expected '%v' and not null but got '%v' and %v`, expectedValue, currentValue, isNull)
	}
}

func TestWStringLargerThanField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddWStringField(`Field1`, ``, 10)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetStringField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := `blah blah blah`
	field.SetString(expectedValue)
	if currentValue, isNull := field.GetCurrentString(); currentValue != `blah blah ` || isNull {
		t.Fatalf(`expected 'blah blah ' and not null but got '%v' and %v`, currentValue, isNull)
	}
}

func TestOutgoingV_StringField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddV_StringField(`Field1`, ``, 20)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetStringField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := `hello world`
	field.SetString(expectedValue)
	if currentValue, isNull := field.GetCurrentString(); currentValue != expectedValue || isNull {
		t.Fatalf(`expected '%v' and not null but got '%v' and %v`, expectedValue, currentValue, isNull)
	}
	field.SetNullString()
	if currentValue, isNull := field.GetCurrentString(); currentValue != `` || isNull != true {
		t.Fatalf(`expected '' and null but got '%v' and %v`, currentValue, isNull)
	}
}

func TestV_StringSameSizeAsField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddV_StringField(`Field1`, ``, 10)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetStringField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := `0123456789`
	field.SetString(expectedValue)
	if currentValue, isNull := field.GetCurrentString(); currentValue != expectedValue || isNull {
		t.Fatalf(`expected '%v' and not null but got '%v' and %v`, expectedValue, currentValue, isNull)
	}
}

func TestV_StringLargerThanField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	editor.AddV_StringField(`Field1`, ``, 10)
	info := editor.GenerateOutgoingRecordInfo()
	field, err := info.GetStringField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	expectedValue := `blah blah blah`
	field.SetString(expectedValue)
	if currentValue, isNull := field.GetCurrentString(); currentValue != `blah blah ` || isNull {
		t.Fatalf(`expected 'blah blah ' and not null but got '%v' and %v`, currentValue, isNull)
	}
}
