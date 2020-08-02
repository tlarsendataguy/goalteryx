package recordcopier_test

import (
	"github.com/tlarsen7572/goalteryx/recordblob"
	"github.com/tlarsen7572/goalteryx/recordcopier"
	"github.com/tlarsen7572/goalteryx/recordinfo"
	"testing"
	"unsafe"
)

func TestRecordCopier(t *testing.T) {
	info1, info2 := generateRecordInfos()

	indexMap := []recordcopier.IndexMap{
		{0, 1},
		{1, 0},
	}

	copier, err := recordcopier.New(info2, info1, indexMap)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	_ = info1.SetStringField(`Some String`, `hello world`)
	_ = info1.SetFloatField(`Some Number`, 123.45)
	record, _ := info1.GenerateRecord()

	err = copier.Copy(record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	record, err = info2.GenerateRecord()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	expectedValue := `hello world`
	actualValue, isNull, err := info2.GetStringValueFrom(`String`, record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected not null but got null`)
	}
	if expectedValue != actualValue {
		t.Fatalf(`expected '%v' but got '%v'`, expectedValue, actualValue)
	}

	expectedNumber := 123.45
	actualNumber, isNull, err := info2.GetFloatValueFrom(`Number`, record)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected not null but got null`)
	}
	if expectedNumber != actualNumber {
		t.Fatalf(`expected %v but got %v`, expectedNumber, actualNumber)
	}
}

func TestRecordCopierInvalidIndices(t *testing.T) {
	info1, info2 := generateRecordInfos()

	indexMaps := []recordcopier.IndexMap{
		{-1, 0},
	}
	_, err := recordcopier.New(info2, info1, indexMaps)
	if err == nil {
		t.Fatalf(`expected an error but got none`)
	}

	indexMaps = []recordcopier.IndexMap{
		{0, 2},
	}
	_, err = recordcopier.New(info2, info1, indexMaps)
	if err == nil {
		t.Fatalf(`expected an error but got none`)
	}
}

func TestCopierShortWString(t *testing.T) {
	blob := []byte{100, 0, 0, 0, 0, 0, 0, 0, 0, 8, 0, 0, 0, 7, 0, 0, 0, 13, 49, 0, 48, 0, 48, 0, 33, 156, 241, 23, 49, 106, 0, 16, 96, 106, 68, 179, 193, 2, 0, 0, 192, 106, 68, 179, 193, 2, 0, 0, 10, 0}
	record := recordblob.NewRecordBlob(unsafe.Pointer(&blob[0]))

	generator := recordinfo.NewGenerator()
	generator.AddInt64Field(`Id`, ``)
	generator.AddV_WStringField(`IdStr`, ``, 1073741823)
	info := generator.GenerateRecordInfo()
	copier, err := recordcopier.New(info, info, []recordcopier.IndexMap{
		{0, 0},
		{1, 1},
	})
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	err = copier.Copy(record)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	value, isNull, err := info.GetCurrentString(`IdStr`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	if isNull {
		t.Fatalf(`expected a value but got Null`)
	}
	if value != `100` {
		t.Fatalf(`expected '100' but got '%v'`, value)
	}
}

func generateRecordInfos() (recordinfo.RecordInfo, recordinfo.RecordInfo) {
	info1 := recordinfo.NewGenerator()
	info1.AddV_WStringField(`Some String`, ``, 1000)
	info1.AddDoubleField(`Some Number`, ``)

	info2 := recordinfo.NewGenerator()
	info2.AddDoubleField(`Number`, ``)
	info2.AddV_WStringField(`String`, ``, 1000)

	return info1.GenerateRecordInfo(), info2.GenerateRecordInfo()
}
