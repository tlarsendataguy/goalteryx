package recordinfo_test

import (
	"github.com/tlarsen7572/goalteryx/recordinfo"
	"testing"
)

func TestTotalSizeVarField(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddV_StringField(`TestVarField`, ``, 100)
	info := generator.GenerateRecordInfo()
	reader := generator.GenerateRecordBlobReader()
	_ = info.SetStringField(`TestVarField`, `Hello world`)
	record, _ := info.GenerateRecord()

	expectedSize := 4 + 4 + 12 // 4 bytes for var location, 4 bytes for var len, and 12 bytes for null-terminated string
	if actualSize := reader.TotalSize(record); actualSize != expectedSize {
		t.Fatalf(`expected total size %v but got %v`, expectedSize, actualSize)
	}
}

func TestTotalSizeFixedField(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddStringField(`TestFixedField`, ``, 14)
	info := generator.GenerateRecordInfo()
	reader := generator.GenerateRecordBlobReader()
	_ = info.SetStringField(`TestFixedField`, `Hello world`)
	record, _ := info.GenerateRecord()

	expectedSize := 15 // field size 14, plus 1 byte for null indicator, plus 4 bytes for var len
	if actualSize := reader.TotalSize(record); actualSize != expectedSize {
		t.Fatalf(`expected total size %v but got %v`, expectedSize, actualSize)
	}
}

func TestAddFromFieldInfo(t *testing.T) {
	generator := recordinfo.NewGenerator()
	generator.AddInt64Field(`Field1`, ``)
	generator.AddStringField(`Field2`, ``, 10)
	generator.AddFixedDecimalField(`Field3`, ``, 19, 6)
	sourceInfo := generator.GenerateRecordInfo()

	generator = recordinfo.NewGenerator()
	for index := 0; index < sourceInfo.NumFields(); index++ {
		field, _ := sourceInfo.GetFieldByIndex(index)
		generator.AddField(field, ``)
	}
	newInfo := generator.GenerateRecordInfo()

	sourceXml, _ := sourceInfo.ToXml(`Test`)
	newXml, _ := newInfo.ToXml(`Test`)
	if sourceXml != newXml {
		t.Fatalf("expected\n%v\nbut got\n%v", sourceXml, newXml)
	}
}

func TestInstantiateFromXml(t *testing.T) {
	xml := `<MetaInfo connection="Output">
<RecordInfo>
	<Field name="Field1" source="TextInput:" type="Byte"/>
	<Field name="Field2" size="1" source="TextInput:" type="String"/>
</RecordInfo>
</MetaInfo>`

	_, err := recordinfo.FromXml(xml)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
}

func TestInstantiateFromXmlNoMetainfo(t *testing.T) {
	xml := `<RecordInfo>
	<Field name="Field1" source="TextInput:" type="Int32"/>
	<Field name="Field2" size="2147483647" source="TextInput:" type="V_String"/>
</RecordInfo>`

	_, err := recordinfo.FromXml(xml)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
}
