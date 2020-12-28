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

func (i *EditingRecordInfo) AddInt16Field(name string, source string, options ...AddFieldOptionSetter) string {
	return i.addField(name, `Int16`, source, 2, 0, options...)
}

func (i *EditingRecordInfo) AddInt32Field(name string, source string, options ...AddFieldOptionSetter) string {
	return i.addField(name, `Int32`, source, 4, 0, options...)
}

func (i *EditingRecordInfo) AddInt64Field(name string, source string, options ...AddFieldOptionSetter) string {
	return i.addField(name, `Int64`, source, 8, 0, options...)
}

func (i *EditingRecordInfo) AddFloatField(name string, source string, options ...AddFieldOptionSetter) string {
	return i.addField(name, `Float`, source, 4, 0, options...)
}

func (i *EditingRecordInfo) AddDoubleField(name string, source string, options ...AddFieldOptionSetter) string {
	return i.addField(name, `Double`, source, 8, 0, options...)
}

func (i *EditingRecordInfo) AddFixedDecimalField(name string, source string, size int, scale int, options ...AddFieldOptionSetter) string {
	return i.addField(name, `FixedDecimal`, source, size, scale, options...)
}

func (i *EditingRecordInfo) AddStringField(name string, source string, size int, options ...AddFieldOptionSetter) string {
	return i.addField(name, `String`, source, size, 0, options...)
}

func (i *EditingRecordInfo) AddWStringField(name string, source string, size int, options ...AddFieldOptionSetter) string {
	return i.addField(name, `WString`, source, size, 0, options...)
}

func (i *EditingRecordInfo) AddV_StringField(name string, source string, size int, options ...AddFieldOptionSetter) string {
	return i.addField(name, `V_String`, source, size, 0, options...)
}

func (i *EditingRecordInfo) AddV_WStringField(name string, source string, size int, options ...AddFieldOptionSetter) string {
	return i.addField(name, `V_WString`, source, size, 0, options...)
}

func (i *EditingRecordInfo) AddBlobField(name string, source string, size int, options ...AddFieldOptionSetter) string {
	return i.addField(name, `Blob`, source, size, 0, options...)
}

func (i *EditingRecordInfo) AddSpatialObjField(name string, source string, size int, options ...AddFieldOptionSetter) string {
	return i.addField(name, `SpatialObj`, source, size, 0, options...)
}

func (i *EditingRecordInfo) AddDateField(name string, source string, options ...AddFieldOptionSetter) string {
	return i.addField(name, `Date`, source, 10, 0, options...)
}

func (i *EditingRecordInfo) AddDateTimeField(name string, source string, options ...AddFieldOptionSetter) string {
	return i.addField(name, `DateTime`, source, 19, 0, options...)
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

func (i *EditingRecordInfo) GenerateOutgoingRecordInfo() *OutgoingRecordInfo {
	var valueCache []byte
	fields := make([]*outgoingField, i.NumFields())

	for index, field := range i.fields {
		switch field.Type {
		case `Bool`:
			valueCache = make([]byte, 1)
		case `Byte`:
			valueCache = make([]byte, 2)
		default:
			panic(fmt.Sprintf(`field %v has an invalid field type (%v) for generating an OutgoingRecordInfo`, field.Name, field.Type))
		}
		fields[index] = &outgoingField{
			Name:         field.Name,
			Type:         field.Type,
			Source:       field.Source,
			Size:         field.Size,
			Scale:        field.Scale,
			CopyFrom:     field.GetBytes,
			CurrentValue: valueCache,
		}
	}
	info := &OutgoingRecordInfo{outgoingFields: fields}
	return info
}
