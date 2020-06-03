package presort_test

import (
	"goalteryx/presort"
	"testing"
)

func TestXmlGenerator(t *testing.T) {
	sortInfo := []presort.SortInfo{
		{Field: `Field1`, Order: presort.Asc},
		{Field: `Field2`, Order: presort.Desc},
	}
	fieldFilter := []string{
		`Field1`,
		`Field3`,
	}
	presortInfo := presort.PresortInfo{
		SortInfo:        sortInfo,
		FieldFilterList: fieldFilter,
	}
	generatedXml, err := presortInfo.ToXml()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	expectedXml := `<SortInfo><Field field="Field1" order="Asc"></Field><Field field="Field2" order="Desc"></Field></SortInfo><FieldFilterList><Field field="Field1"></Field><Field field="Field3"></Field></FieldFilterList>`
	if expectedXml != generatedXml {
		t.Fatalf("expected\n%v\nbut got\n%v", expectedXml, generatedXml)
	}
}

func TestXmlGeneratorWithoutFieldFilter(t *testing.T) {
	sortInfo := []presort.SortInfo{
		{Field: `Field1`, Order: presort.Asc},
		{Field: `Field2`, Order: presort.Desc},
	}
	presortInfo := presort.PresortInfo{
		SortInfo:        sortInfo,
		FieldFilterList: nil,
	}
	generatedXml, err := presortInfo.ToXml()
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	expectedXml := `<SortInfo><Field field="Field1" order="Asc"></Field><Field field="Field2" order="Desc"></Field></SortInfo>`
	if expectedXml != generatedXml {
		t.Fatalf("expected\n%v\nbut got\n%v", expectedXml, generatedXml)
	}
}
