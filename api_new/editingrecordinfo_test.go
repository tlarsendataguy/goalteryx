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

func TestAddByteField(t *testing.T) {
	editor := &api_new.EditingRecordInfo{}
	name := editor.AddByteField(`Field1`, ``)
	if name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, name)
	}
	if editor.NumFields() != 1 {
		t.Fatalf(`expected 1 but got %v`, editor.NumFields())
	}
}
