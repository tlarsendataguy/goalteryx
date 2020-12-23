package api_new

import "fmt"

type EditingRecordInfo struct {
	fields []IncomingField
}

func (i *EditingRecordInfo) NumFields() int {
	return len(i.fields)
}

func (i *EditingRecordInfo) AddBoolField(name string, source string) string {
	return i.addField(name, `Bool`, source, 1, 0)
}

func (i *EditingRecordInfo) AddByteField(name string, source string) string {
	return i.addField(name, `Byte`, source, 1, 0)
}

func (i *EditingRecordInfo) addField(name string, typeName string, source string, size int, scale int) string {
	actualName := i.checkName(name)
	i.fields = append(i.fields, IncomingField{
		Name:     actualName,
		Type:     typeName,
		Source:   source,
		Size:     size,
		Scale:    scale,
		GetBytes: nil,
	})
	return actualName
}

func (i *EditingRecordInfo) checkName(name string) string {
	for _, field := range i.fields {
		if name == field.Name {
			name = fmt.Sprintf(`%v2`, name)
			return i.checkName(name)
		}
	}
	return name
}
