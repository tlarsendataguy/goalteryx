package recordinfo

func (info *recordInfo) AddByteField(name string, source string) string {
	return info.addField(name, source, 1, 0, Byte, 1, 1)
}

func (info *recordInfo) AddBoolField(name string, source string) string {
	return info.addField(name, source, 1, 0, Bool, 1, 0)
}

func (info *recordInfo) AddInt16Field(name string, source string) string {
	return info.addField(name, source, 2, 0, Int16, 2, 1)
}

func (info *recordInfo) AddInt32Field(name string, source string) string {
	return info.addField(name, source, 4, 0, Int32, 4, 1)
}

func (info *recordInfo) AddInt64Field(name string, source string) string {
	return info.addField(name, source, 8, 0, Int64, 8, 1)
}

func (info *recordInfo) AddFixedDecimalField(name string, source string, size int, precision int) string {
	return info.addField(name, source, size, precision, FixedDecimal, uintptr(size), 1)
}

func (info *recordInfo) AddFloatField(name string, source string) string {
	return info.addField(name, source, 4, 0, Float, 4, 1)
}

func (info *recordInfo) AddDoubleField(name string, source string) string {
	return info.addField(name, source, 8, 0, Double, 8, 1)
}

func (info *recordInfo) AddStringField(name string, source string, size int) string {
	return info.addField(name, source, size, 0, String, uintptr(size), 1)
}

func (info *recordInfo) AddWStringField(name string, source string, size int) string {
	return info.addField(name, source, size, 0, WString, uintptr(size*2), 1)
}

func (info *recordInfo) AddV_StringField(name string, source string, size int) string {
	return info.addField(name, source, size, 0, V_String, 4, 0)
}

func (info *recordInfo) AddV_WStringField(name string, source string, size int) string {
	return info.addField(name, source, size, 0, V_WString, 4, 0)
}

func (info *recordInfo) AddDateField(name string, source string) string {
	return info.addField(name, source, 10, 0, Date, 10, 1)
}

func (info *recordInfo) AddDateTimeField(name string, source string) string {
	return info.addField(name, source, 19, 0, DateTime, 19, 1)
}

func (info *recordInfo) AddBlobField(name string, source string, size int) string {
	return info.addField(name, source, size, 0, Blob, 4, 0)
}

func (info *recordInfo) AddSpatialField(name string, source string, size int) string {
	return info.addField(name, source, size, 0, Spatial, 4, 0)
}

func (info *recordInfo) addField(name string, source string, size int, scale int, fieldType FieldType, fixedLen uintptr, nullByteLen uintptr) string {
	actualName := info.checkFieldName(name)
	info.fields = append(info.fields, &fieldInfoEditor{
		Name:        actualName,
		Source:      source,
		Size:        size,
		Precision:   scale,
		Type:        fieldType,
		location:    info.fixedLen,
		fixedLen:    fixedLen,
		nullByteLen: nullByteLen,
		value:       make([]byte, fixedLen+nullByteLen),
	})
	info.fieldNames[actualName] = info.numFields
	info.numFields++
	info.fixedLen += fixedLen + nullByteLen
	return actualName
}
