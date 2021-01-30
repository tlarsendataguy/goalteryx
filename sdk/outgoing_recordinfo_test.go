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

func TestGetNull(t *testing.T) {
	info, _ := sdk.NewOutgoingRecordInfo([]sdk.NewOutgoingField{
		sdk.NewBoolField(`Field1`, `source`),
		sdk.NewByteField(`Field2`, `source`),
		sdk.NewInt16Field(`Field3`, `source`),
		sdk.NewInt32Field(`Field4`, `source`),
		sdk.NewInt64Field(`Field5`, `source`),
		sdk.NewStringField(`Field6`, `source`, 10),
		sdk.NewWStringField(`Field7`, `source`, 10),
		sdk.NewV_StringField(`Field8`, `source`, 1000),
		sdk.NewV_WStringField(`Field9`, `source`, 1000),
		sdk.NewFloatField(`Field10`, `source`),
		sdk.NewDoubleField(`Field11`, `source`),
		sdk.NewFixedDecimalField(`Field12`, `source`, 18, 2),
		sdk.NewDateField(`Field13`, `source`),
		sdk.NewDateTimeField(`Field14`, `source`),
		sdk.NewBlobField(`Field15`, `source`, 1000000),
		sdk.NewSpatialObjField(`Field16`, `source`, 1000000),
	})

	for _, field := range info.FloatFields {
		field.SetNull()
		isNull := field.GetNull()
		if !isNull {
			t.Fatalf(`expected null but got not null`)
		}
	}
	for _, field := range info.StringFields {
		field.SetNull()
		isNull := field.GetNull()
		if !isNull {
			t.Fatalf(`expected null but got not null`)
		}
	}
	for _, field := range info.BoolFields {
		field.SetNull()
		isNull := field.GetNull()
		if !isNull {
			t.Fatalf(`expected null but got not null`)
		}
	}
	for _, field := range info.DateTimeFields {
		field.SetNull()
		isNull := field.GetNull()
		if !isNull {
			t.Fatalf(`expected null but got not null`)
		}
	}
	for _, field := range info.IntFields {
		field.SetNull()
		isNull := field.GetNull()
		if !isNull {
			t.Fatalf(`expected null but got not null`)
		}
	}
	for _, field := range info.BlobFields {
		field.SetNull()
		isNull := field.GetNull()
		if !isNull {
			t.Fatalf(`expected null but got not null`)
		}
	}
}
