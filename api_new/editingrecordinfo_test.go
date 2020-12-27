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

func TestAddFixedDecinalField(t *testing.T) {
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
