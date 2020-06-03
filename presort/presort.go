package presort

import "encoding/xml"

type SortOrder string

const (
	Asc  SortOrder = `Asc`
	Desc SortOrder = `Desc`
)

type SortInfo struct {
	Field string    `xml:"field,attr"`
	Order SortOrder `xml:"order,attr"`
}

type PresortInfo struct {
	SortInfo        []SortInfo
	FieldFilterList []string
}

type xmlSortInfo struct {
	XMLName xml.Name `xml:"SortInfo"`
	Field   []SortInfo
}

type xmlFieldFilterList struct {
	XMLName xml.Name `xml:"FieldFilterList"`
	Field   []xmlField
}

type xmlField struct {
	Field string `xml:"field,attr"`
}

func (info *PresortInfo) ToXml() (string, error) {
	sortInfo := xmlSortInfo{Field: info.SortInfo}
	sortXml, err := xml.Marshal(sortInfo)
	if err != nil {
		return ``, err
	}
	if len(info.FieldFilterList) == 0 {
		return string(sortXml), nil
	}
	fieldFilterList := xmlFieldFilterList{}
	for _, field := range info.FieldFilterList {
		fieldFilterList.Field = append(fieldFilterList.Field, xmlField{Field: field})
	}
	filterXml, err := xml.Marshal(fieldFilterList)
	if err != nil {
		return ``, err
	}
	return string(append(sortXml, filterXml...)), nil
}
