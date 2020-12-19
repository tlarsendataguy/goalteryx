package api_new

import (
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
			return IncomingIntField{
				Name:     field.Name,
				Type:     field.Type,
				Source:   field.Source,
				GetValue: bytesToByte(field.GetBytes),
			}, nil
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
	return IncomingRecordInfo{fields: metaInfo.RecordInfo.Fields}, nil
}

func bytesToByte(getBytes BytesGetter) IntGetter {
	return func(record Record) (int, bool) {
		bytes := getBytes(record)
		if bytes[0] == 2 {
			return 0, true
		}
		return int(bytes[0]), false
	}
}
