package recordinfo_test

import (
	"goalteryx/recordinfo"
	"testing"
)

func TestTotalSizeVarField(t *testing.T) {
	info := recordinfo.New()
	info.AddV_StringField(`TestVarField`, ``, 100)
	_ = info.SetStringField(`TestVarField`, `Hello world`)
	record, _ := info.GenerateRecord()

	expectedSize := 4 + 4 + 12 // 4 bytes for var location, 4 bytes for var len, and 12 bytes for null-terminated string
	if actualSize := info.TotalSize(record); actualSize != expectedSize {
		t.Fatalf(`expected total size %v but got %v`, expectedSize, actualSize)
	}
}

func TestTotalSizeFixedField(t *testing.T) {
	info := recordinfo.New()
	info.AddStringField(`TestFixedField`, ``, 14)
	_ = info.SetStringField(`TestFixedField`, `Hello world`)
	record, _ := info.GenerateRecord()

	expectedSize := 15 + 4 // field size 14, plus 1 byte for null indicator, plus 4 bytes for var len
	if actualSize := info.TotalSize(record); actualSize != expectedSize {
		t.Fatalf(`expected total size %v but got %v`, expectedSize, actualSize)
	}
}
