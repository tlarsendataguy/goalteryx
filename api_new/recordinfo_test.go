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

func TestGetIncomingField(t *testing.T) {
	config := `<RecordInfo>
	<Field name="Field1" source="TextInput:" type="Byte"/>
	<Field name="Field2" size="1" source="TextInput:" type="String"/>
</RecordInfo>`
	recordInfo, _ := incomingRecordInfoFromString(config)
	field, err := recordInfo.GetIntField(`Field1`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	if field.Name != `Field1` {
		t.Fatalf(`expected 'Field1' but got '%v'`, field.Name)
	}
}

func TestIncomingFieldDoesNotExist(t *testing.T) {
	config := `<RecordInfo>
	<Field name="Field1" source="TextInput:" type="Byte"/>
	<Field name="Field2" size="1" source="TextInput:" type="String"/>
</RecordInfo>`
	recordInfo, _ := incomingRecordInfoFromString(config)
	_, err := recordInfo.GetIntField(`Hello world`)
	if err == nil {
		t.Fatalf(`expected an error but got none`)
	}
	t.Logf(err.Error())
}
