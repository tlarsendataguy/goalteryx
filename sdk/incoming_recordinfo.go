package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	b "github.com/tlarsen7572/goalteryx/sdk/field_base"
)

const dateFormat = `2006-01-02`
const dateTimeFormat = `2006-01-02 15:04:05`

type xmlMetaInfo struct {
	Connection string        `xml:"connection,attr"`
	RecordInfo xmlRecordInfo `xml:"RecordInfo"`
}

type xmlRecordInfo struct {
	Fields []IncomingField `xml:"Field"`
}

type IncomingRecordInfo struct {
	fields []IncomingField
}

func (i IncomingRecordInfo) NumFields() int {
	return len(i.fields)
}

func (i IncomingRecordInfo) Fields() []b.FieldBase {
	fields := make([]b.FieldBase, len(i.fields))
	for index, field := range i.fields {
		fields[index] = b.FieldBase{
			Name:   field.Name,
			Type:   field.Type,
			Source: field.Source,
			Size:   field.Size,
			Scale:  field.Scale,
		}
	}
	return fields
}

func (i IncomingRecordInfo) Clone() *EditingRecordInfo {
	return &EditingRecordInfo{fields: i.fields}
}

func (i IncomingRecordInfo) GetIntField(name string) (IncomingIntField, error) {
	for _, field := range i.fields {
		if field.Name != name {
			continue
		}
		switch field.Type {
		case `Byte`:
			return generateIncomingIntField(field, bytesToByte), nil
		case `Int16`:
			return generateIncomingIntField(field, bytesToInt16), nil
		case `Int32`:
			return generateIncomingIntField(field, bytesToInt32), nil
		case `Int64`:
			return generateIncomingIntField(field, bytesToInt64), nil
		default:
			return IncomingIntField{}, fmt.Errorf(`the '%v' field is not an integer field, it is '%v'`, name, field.Type)
		}
	}
	return IncomingIntField{}, fmt.Errorf(`there is no '%v' field in the record`, name)
}

func (i IncomingRecordInfo) GetFloatField(name string) (IncomingFloatField, error) {
	for _, field := range i.fields {
		if field.Name != name {
			continue
		}
		switch field.Type {
		case `Float`:
			return generateIncomingFloatField(field, bytesToFloat), nil
		case `Double`:
			return generateIncomingFloatField(field, bytesToDouble), nil
		case `FixedDecimal`:
			return generateFixedDecimalField(field), nil
		default:
			return IncomingFloatField{}, fmt.Errorf(`the '%v' field is not a float field, it is '%v'`, name, field.Type)
		}
	}
	return IncomingFloatField{}, fmt.Errorf(`there is no '%v' field in the record`, name)
}

func (i IncomingRecordInfo) GetBoolField(name string) (IncomingBoolField, error) {
	for _, field := range i.fields {
		if field.Name != name {
			continue
		}
		switch field.Type {
		case `Bool`:
			return generateBoolField(field), nil
		default:
			return IncomingBoolField{}, fmt.Errorf(`the '%v' field is not a bool field, it is '%v'`, name, field.Type)
		}
	}
	return IncomingBoolField{}, fmt.Errorf(`there is no '%v' field in the record`, name)
}

func (i IncomingRecordInfo) GetTimeField(name string) (IncomingTimeField, error) {
	for _, field := range i.fields {
		if field.Name != name {
			continue
		}
		switch field.Type {
		case `Date`:
			return generateTimeField(field, dateFormat, 10), nil
		case `DateTime`:
			return generateTimeField(field, dateTimeFormat, 19), nil
		default:
			return IncomingTimeField{}, fmt.Errorf(`the '%v' field is not a time field, it is '%v'`, name, field.Type)
		}
	}
	return IncomingTimeField{}, fmt.Errorf(`there is no '%v' field in the record`, name)
}

func (i IncomingRecordInfo) GetBlobField(name string) (IncomingBlobField, error) {
	for _, field := range i.fields {
		if field.Name != name {
			continue
		}
		switch field.Type {
		case `Blob`, `SpatialObj`:
			return generateBlobField(field), nil
		default:
			return IncomingBlobField{}, fmt.Errorf(`the '%v' field is not a blob field, it is '%v'`, name, field.Type)
		}
	}
	return IncomingBlobField{}, fmt.Errorf(`there is no '%v' field in the record`, name)
}

func (i IncomingRecordInfo) GetStringField(name string) (IncomingStringField, error) {
	for _, field := range i.fields {
		if field.Name != name {
			continue
		}
		switch field.Type {
		case `String`:
			return generateIncomingStringField(field, bytesToString), nil
		case `WString`:
			return generateIncomingStringField(field, bytesToWString), nil
		case `V_String`:
			return generateIncomingStringField(field, bytesToV_String), nil
		case `V_WString`:
			return generateIncomingStringField(field, bytesToV_WString), nil
		default:
			return IncomingStringField{}, fmt.Errorf(`the '%v' field is not a string field, it is '%v'`, name, field.Type)
		}
	}
	return IncomingStringField{}, fmt.Errorf(`there is no '%v' field in the record`, name)
}

func incomingRecordInfoFromString(config string) (IncomingRecordInfo, error) {
	if config[:9] != `<MetaInfo` {
		if config[:11] != `<RecordInfo` {
			return IncomingRecordInfo{}, errors.New(`config is not a valid IncomingRecordInfo xml string`)
		}
		config = `<MetaInfo>` + config + `</MetaInfo>`
	}
	metaInfo := xmlMetaInfo{}
	err := xml.Unmarshal([]byte(config), &metaInfo)
	if err != nil {
		return IncomingRecordInfo{}, err
	}
	startAt := 0
	for index, field := range metaInfo.RecordInfo.Fields {
		switch field.Type {
		case `V_String`, `V_WString`, `Blob`, `SpatialObj`:
			field.GetBytes = generateGetVarBytes(startAt)
			startAt += 4
		case `Bool`:
			field.GetBytes = generateGetFixedBytes(startAt, 1)
			startAt += 1
		case `Byte`:
			field.GetBytes = generateGetFixedBytes(startAt, 2)
			startAt += 2
		case `Int16`:
			field.GetBytes = generateGetFixedBytes(startAt, 3)
			startAt += 3
		case `Int32`, `Float`:
			field.GetBytes = generateGetFixedBytes(startAt, 5)
			startAt += 5
		case `Int64`, `Double`:
			field.GetBytes = generateGetFixedBytes(startAt, 9)
			startAt += 9
		case `String`, `FixedDecimal`:
			field.GetBytes = generateGetFixedBytes(startAt, field.Size+1)
			startAt += field.Size + 1
		case `WString`:
			size := field.Size * 2
			field.GetBytes = generateGetFixedBytes(startAt, size+1)
			startAt += size + 1
		case `Date`:
			field.GetBytes = generateGetFixedBytes(startAt, 11)
			startAt += 11
		case `DateTime`:
			field.GetBytes = generateGetFixedBytes(startAt, 20)
			startAt += 20
		default:
			return IncomingRecordInfo{}, fmt.Errorf(`field '%v' has invalid field type '%v'`, field.Name, field.Type)
		}
		metaInfo.RecordInfo.Fields[index] = field
	}
	return IncomingRecordInfo{fields: metaInfo.RecordInfo.Fields}, nil
}
