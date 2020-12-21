package api_new

import (
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
)

type IntGetter func(Record) (int, bool)

type xmlMetaInfo struct {
	Connection string        `xml:"connection,attr"`
	RecordInfo xmlRecordInfo `xml:"RecordInfo"`
}

type xmlRecordInfo struct {
	Fields []IncomingField `xml:"Field"`
}

type IncomingField struct {
	Name     string `xml:"name,attr"`
	Type     string `xml:"type,attr"`
	Source   string `xml:"source,attr"`
	Size     int    `xml:"size,attr"`
	Scale    int    `xml:"scale,attr"`
	GetBytes BytesGetter
}

type IncomingIntField struct {
	Name     string
	Type     string
	Source   string
	GetValue IntGetter
}

type IncomingRecordInfo struct {
	fields []IncomingField
}

func (i IncomingRecordInfo) NumFields() int {
	return len(i.fields)
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
		default:
			return IncomingIntField{}, fmt.Errorf(`the '%v' field is not an integer field, it is '%v'`, name, field.Type)
		}
	}
	return IncomingIntField{}, fmt.Errorf(`there is no '%v' field in the record`, name)
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
		case `String`, `WString`, `FixedDecimal`:
			field.GetBytes = generateGetFixedBytes(startAt, field.Size+1)
			startAt += field.Size + 1
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

func generateIncomingIntField(field IncomingField, getter func(BytesGetter) IntGetter) IncomingIntField {
	return IncomingIntField{
		Name:     field.Name,
		Type:     field.Type,
		Source:   field.Source,
		GetValue: getter(field.GetBytes),
	}
}

func bytesToByte(getBytes BytesGetter) IntGetter {
	return func(record Record) (int, bool) {
		bytes := getBytes(record)
		if bytes[1] == 1 {
			return 0, true
		}
		return int(bytes[0]), false
	}
}

func bytesToInt16(getBytes BytesGetter) IntGetter {
	return func(record Record) (int, bool) {
		bytes := getBytes(record)
		if bytes[2] == 1 {
			return 0, true
		}
		return int(binary.LittleEndian.Uint16(bytes)), false
	}
}

func bytesToInt32(getBytes BytesGetter) IntGetter {
	return func(record Record) (int, bool) {
		bytes := getBytes(record)
		if bytes[4] == 1 {
			return 0, true
		}
		return int(binary.LittleEndian.Uint32(bytes)), false
	}
}
