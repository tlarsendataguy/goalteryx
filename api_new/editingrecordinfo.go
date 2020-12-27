package api_new

import "fmt"

type AddFieldOptions struct {
	doInsert bool
	insertAt int
}

type AddFieldOptionSetter func(AddFieldOptions) AddFieldOptions

func InsertAt(position int) AddFieldOptionSetter {
	return func(options AddFieldOptions) AddFieldOptions {
		options.doInsert = true
		options.insertAt = position
		return options
	}
}

type EditingRecordInfo struct {
	fields []IncomingField
}

func (i *EditingRecordInfo) NumFields() int {
	return len(i.fields)
}

func (i *EditingRecordInfo) Fields() []IncomingField {
	value := make([]IncomingField, len(i.fields))
	copy(value, i.fields)
	return value
}

func (i *EditingRecordInfo) AddBoolField(name string, source string, options ...AddFieldOptionSetter) string {
	return i.addField(name, `Bool`, source, 1, 0, options...)
}

func (i *EditingRecordInfo) AddByteField(name string, source string, options ...AddFieldOptionSetter) string {
	return i.addField(name, `Byte`, source, 1, 0, options...)
}

func (i *EditingRecordInfo) addField(name string, typeName string, source string, size int, scale int, options ...AddFieldOptionSetter) string {
	addFieldOptions := AddFieldOptions{}
	for _, setter := range options {
		addFieldOptions = setter(addFieldOptions)
	}
	actualName := i.checkName(name)
	field := IncomingField{
		Name:     actualName,
		Type:     typeName,
		Source:   source,
		Size:     size,
		Scale:    scale,
		GetBytes: nil,
	}
	if addFieldOptions.doInsert && addFieldOptions.insertAt < len(i.fields) && addFieldOptions.insertAt >= 0 {
		newFields := make([]IncomingField, len(i.fields)+1)
		copy(newFields[:addFieldOptions.insertAt], i.fields[:addFieldOptions.insertAt])
		copy(newFields[addFieldOptions.insertAt+1:], i.fields[addFieldOptions.insertAt:])
		newFields[addFieldOptions.insertAt] = field
		i.fields = newFields
	} else {
		i.fields = append(i.fields, field)
	}
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
