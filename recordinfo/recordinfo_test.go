package recordinfo_test

import (
	"github.com/tlarsen7572/goalteryx/recordinfo"
	"testing"
)

func TestTotalSizeVarField(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddV_StringField(`TestVarField`, ``, 100)
	info := generator.GenerateRecordInfo()
	_ = info.SetStringField(`TestVarField`, `Hello world`)
	record, _ := info.GenerateRecord()

	expectedSize := 4 + 4 + 12 // 4 bytes for var location, 4 bytes for var len, and 12 bytes for null-terminated string
	if actualSize := info.TotalSize(record); actualSize != expectedSize {
		t.Fatalf(`expected total size %v but got %v`, expectedSize, actualSize)
	}
}

func TestTotalSizeFixedField(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddStringField(`TestFixedField`, ``, 14)
	info := generator.GenerateRecordInfo()
	_ = info.SetStringField(`TestFixedField`, `Hello world`)
	record, _ := info.GenerateRecord()

	expectedSize := 15 + 4 // field size 14, plus 1 byte for null indicator, plus 4 bytes for var len
	if actualSize := info.TotalSize(record); actualSize != expectedSize {
		t.Fatalf(`expected total size %v but got %v`, expectedSize, actualSize)
	}
}
