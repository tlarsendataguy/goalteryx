// Package presort provides a struct that generates a compliant PreSort XML string for the Alteryx engine.
package presort

import "encoding/xml"

type SortOrder string

const (
	Asc  SortOrder = `Asc`
	Desc SortOrder = `Desc`
)

// SortInfo specifies the sort order of an incoming field
type SortInfo struct {
	Field string    `xml:"field,attr"`
	Order SortOrder `xml:"order,attr"`
}

// PresortInfo specifies the sort order for the incoming records, as well as any field filters.  Note that all
// fields specified in SortInfo must also be present in FieldFilterList if you wish to filter fields.  If you do not
// wish to filter fields, leave FieldFilterList nil.
type PresortInfo struct {
	SortInfo        []SortInfo
	FieldFilterList []string
}

// xmlSortInfo is a struct used for XML marshalling and is not exported.
type xmlSortInfo struct {
	XMLName xml.Name `xml:"SortInfo"`
	Field   []SortInfo
}

// xmlFieldFilterList is a struct used for XML marshalling and is not exported.
type xmlFieldFilterList struct {
	XMLName xml.Name `xml:"FieldFilterList"`
	Field   []xmlField
}

// xmlField is a struct used for XML marshalling and is not exported.
type xmlField struct {
	Field string `xml:"field,attr"`
}

// ToXml generates the PreSort XML string needed to tell the engine how to presort the incoming data.  The PresortInfo
// object is first converted to the temporary xml structs and then those structs are marshalled to XML.
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
