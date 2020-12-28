package api_new_test

import (
	"github.com/tlarsen7572/goalteryx/api_new"
	"testing"
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
	field.SetInt(45)
	if currentValue, isNull := field.GetCurrentInt(); currentValue != 45 || isNull {
		t.Fatalf(`expected 45 and not null but got %v and %v`, currentValue, isNull)
	}
	field.SetNullInt()
	if currentValue, isNull := field.GetCurrentInt(); currentValue != 0 || isNull != true {
		t.Fatalf(`expected 0 and null but got %v and %v`, currentValue, isNull)
	}
	field.SetInt(10000)
	currentValue, isNull := field.GetCurrentInt()
	t.Logf(`value %v and null=%v`, currentValue, isNull)
}
