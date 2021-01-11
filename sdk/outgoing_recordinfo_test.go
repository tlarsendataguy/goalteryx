package sdk_test

import (
	"github.com/tlarsen7572/goalteryx/sdk"
	"strings"
	"testing"
)

func TestDataSize(t *testing.T) {
	info, _ := sdk.NewOutgoingRecordInfo([]sdk.NewOutgoingField{
		sdk.NewBoolField(`Field1`, `source`),
	})
	size := info.DataSize()
	if size != 1 {
		t.Fatalf(`expected size 1 but got %v`, size)
	}
}

func TestVarDataSize(t *testing.T) {
	info, _ := sdk.NewOutgoingRecordInfo([]sdk.NewOutgoingField{
		sdk.NewV_StringField(`Field1`, `source`, 100000),
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
	info, fieldNames := sdk.NewOutgoingRecordInfo([]sdk.NewOutgoingField{
		sdk.NewBoolField(`Field1`, `source`),
		sdk.NewBoolField(`Field1`, `source`),
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

func TestSetStringsToEmptyString(t *testing.T) {
	info, _ := sdk.NewOutgoingRecordInfo([]sdk.NewOutgoingField{
		sdk.NewStringField(`Field1`, `source`, 100),
		sdk.NewWStringField(`Field2`, `source`, 100),
		sdk.NewV_StringField(`Field3`, `source`, 100),
		sdk.NewV_WStringField(`Field4`, `source`, 100),
	})
	info.StringFields[`Field1`].SetString(``)
	info.StringFields[`Field2`].SetString(``)
	info.StringFields[`Field3`].SetString(``)
	info.StringFields[`Field4`].SetString(``)
}
