package api_new

import (
	"encoding/xml"
	"errors"
)

type xmlMetaInfo struct {
	Connection string        `xml:"connection,attr"`
	RecordInfo xmlRecordInfo `xml:"RecordInfo"`
}

type xmlRecordInfo struct {
	Fields []xmlField `xml:"Field"`
}

type xmlField struct {
	Name   string `xml:"name,attr"`
	Type   string `xml:"type,attr"`
	Source string `xml:"source,attr"`
	Size   string `xml:"size,attr"`
	Scale  string `xml:"scale,attr"`
}

type IncomingRecordInfo struct {
	fields int
}

func (i IncomingRecordInfo) NumFields() int {
	return i.fields
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
	return IncomingRecordInfo{fields: len(metaInfo.RecordInfo.Fields)}, nil
}
