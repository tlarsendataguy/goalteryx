package sdk

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

func (i *EditingRecordInfo) AddTimeField(name string, source string, options ...AddFieldOptionSetter) string {
	return i.addField(name, `Time`, source, 8, 0, options...)
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
	info := &OutgoingRecordInfo{
		outgoingFields: nil,
		BlobFields:     make(map[string]OutgoingBlobField),
		BoolFields:     make(map[string]OutgoingBoolField),
		DateTimeFields: make(map[string]OutgoingDateTimeField),
		FloatFields:    make(map[string]OutgoingFloatField),
		IntFields:      make(map[string]OutgoingIntField),
		StringFields:   make(map[string]OutgoingStringField),
	}
	var outgoing *outgoingField

	for _, field := range i.fields {
		switch field.Type {
		case `Bool`:
			outgoing = NewBoolField(field.Name, field.Source)()
			info.BoolFields[field.Name] = outgoing
		case `Byte`:
			outgoing = NewByteField(field.Name, field.Source)()
			info.IntFields[field.Name] = outgoing
		case `Int16`:
			outgoing = NewInt16Field(field.Name, field.Source)()
			info.IntFields[field.Name] = outgoing
		case `Int32`:
			outgoing = NewInt32Field(field.Name, field.Source)()
			info.IntFields[field.Name] = outgoing
		case `Int64`:
			outgoing = NewInt64Field(field.Name, field.Source)()
			info.IntFields[field.Name] = outgoing
		case `Float`:
			outgoing = NewFloatField(field.Name, field.Source)()
			info.FloatFields[field.Name] = outgoing
		case `Double`:
			outgoing = NewDoubleField(field.Name, field.Source)()
			info.FloatFields[field.Name] = outgoing
		case `FixedDecimal`:
			outgoing = NewFixedDecimalField(field.Name, field.Source, field.Size, field.Scale)()
			info.FloatFields[field.Name] = outgoing
		case `Date`:
			outgoing = NewDateField(field.Name, field.Source)()
			info.DateTimeFields[field.Name] = outgoing
		case `DateTime`:
			outgoing = NewDateTimeField(field.Name, field.Source)()
			info.DateTimeFields[field.Name] = outgoing
		case `Time`:
			outgoing = NewTimeField(field.Name, field.Source)()
			info.DateTimeFields[field.Name] = outgoing
		case `String`:
			outgoing = NewStringField(field.Name, field.Source, field.Size)()
			info.StringFields[field.Name] = outgoing
		case `WString`:
			outgoing = NewWStringField(field.Name, field.Source, field.Size)()
			info.StringFields[field.Name] = outgoing
		case `V_String`:
			outgoing = NewV_StringField(field.Name, field.Source, field.Size)()
			info.StringFields[field.Name] = outgoing
		case `V_WString`:
			outgoing = NewV_WStringField(field.Name, field.Source, field.Size)()
			info.StringFields[field.Name] = outgoing
		case `Blob`:
			outgoing = NewBlobField(field.Name, field.Source, field.Size)()
			info.BlobFields[field.Name] = outgoing
		case `SpatialObj`:
			outgoing = NewSpatialObjField(field.Name, field.Source, field.Size)()
			info.BlobFields[field.Name] = outgoing
		default:
			panic(fmt.Sprintf(`field %v has an invalid field type (%v) for generating an OutgoingRecordInfo`, field.Name, field.Type))
		}
		outgoing.CopyFrom = field.GetBytes
		info.outgoingFields = append(info.outgoingFields, outgoing)
	}
	return info
}

func (i *EditingRecordInfo) RemoveFields(fieldNames ...string) {
	for index := i.NumFields() - 1; index >= 0; index-- {
		field := i.fields[index]
		for _, toDelete := range fieldNames {
			if field.Name == toDelete {
				i.fields = append(i.fields[:index], i.fields[index+1:]...)
				break
			}
		}
	}
}

func (i *EditingRecordInfo) MoveField(name string, newIndex int) error {
	if upperBound := i.NumFields() - 1; newIndex < 0 || newIndex > upperBound {
		return fmt.Errorf(`index out of range, must be between 0 and %v`, upperBound)
	}

	fields := i.Fields()
	for index, field := range fields {
		if field.Name == name {
			newFields := append(i.fields[:index], i.fields[index+1:]...)
			newFields = append(newFields[:newIndex], append([]IncomingField{field}, newFields[newIndex:]...)...)
			return nil
		}
	}
	return fmt.Errorf(`field '%v' does not exist in the record`, name)
}
