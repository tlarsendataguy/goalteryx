package api_new

import "testing"

func TestNewIncomingRecordInfo(t *testing.T) {
	config := `<MetaInfo connection="Output">
<RecordInfo>
	<Field name="Field1" source="TextInput:" type="Byte"/>
	<Field name="Field2" size="1" source="TextInput:" type="String"/>
</RecordInfo>
</MetaInfo>`
	recordInfo, err := incomingRecordInfoFromString(config)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if fields := recordInfo.NumFields(); fields != 2 {
		t.Fatalf(`expected 2 fields but got %v`, fields)
	}
}

func TestNewIncomingRecordInfoWithoutMetaInfo(t *testing.T) {
	config := `<RecordInfo>
	<Field name="Field1" source="TextInput:" type="Byte"/>
	<Field name="Field2" size="1" source="TextInput:" type="String"/>
</RecordInfo>`
	recordInfo, err := incomingRecordInfoFromString(config)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if fields := recordInfo.NumFields(); fields != 2 {
		t.Fatalf(`expected 2 fields but got %v`, fields)
	}
}
