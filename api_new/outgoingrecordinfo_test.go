package api_new_test

import (
	"github.com/tlarsen7572/goalteryx/api_new"
	"strings"
	"testing"
)

func TestDataSize(t *testing.T) {
	info, _ := api_new.NewOutgoingRecordInfo([]api_new.NewOutgoingField{
		api_new.NewBoolField(`Field1`, `source`),
	})
	size := info.DataSize()
	if size != 1 {
		t.Fatalf(`expected size 1 but got %v`, size)
	}
}

func TestVarDataSize(t *testing.T) {
	info, _ := api_new.NewOutgoingRecordInfo([]api_new.NewOutgoingField{
		api_new.NewV_StringField(`Field1`, `source`, 100000),
	})
	field := info.StringFields[`Field1`]

	size := info.DataSize()
	if size != 8 {
		t.Fatalf(`expected size 8 but got %v`, size)
	}

	field.SetString(`Hello world`)
	size = info.DataSize()
	if size != 20 {
		t.Fatalf(`expected size 20 but got %v`, size)
	}

	field.SetString(strings.Repeat(`A`, 200))
	size = info.DataSize()
	if size != 212 {
		t.Fatalf(`expected size 212 but got %v`, size)
	}
}

func TestDuplicateFieldNames(t *testing.T) {
	info, fieldNames := api_new.NewOutgoingRecordInfo([]api_new.NewOutgoingField{
		api_new.NewBoolField(`Field1`, `source`),
		api_new.NewBoolField(`Field1`, `source`),
	})
	_, ok := info.BoolFields[`Field1`]
	if !ok {
		t.Fatalf(`expected Field1 but got no field`)
	}
	_, ok = info.BoolFields[`Field12`]
	if !ok {
		t.Fatalf(`expected Field12 but got no field`)
	}
	if fieldNames[0] != `Field1` || fieldNames[1] != `Field12` {
		t.Fatalf(`expected [Field1 Field12] but got %v`, fieldNames)
	}
}
