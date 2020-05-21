package recordinfo

func (info *recordInfo) AddByteField(name string, source string) string {
	return info.addField(name, source, 1, 0, ByteType, 1, 1)
}

func (info *recordInfo) AddBoolField(name string, source string) string {
	return info.addField(name, source, 1, 0, BoolType, 1, 0)
}

func (info *recordInfo) AddInt16Field(name string, source string) string {
	return info.addField(name, source, 2, 0, Int16Type, 2, 1)
}

func (info *recordInfo) AddInt32Field(name string, source string) string {
	return info.addField(name, source, 4, 0, Int32Type, 4, 1)
}

func (info *recordInfo) AddInt64Field(name string, source string) string {
	return info.addField(name, source, 8, 0, Int64Type, 8, 1)
}

func (info *recordInfo) AddFixedDecimalField(name string, source string, size int, precision int) string {
	return info.addField(name, source, size, precision, FixedDecimalType, uintptr(size), 1)
}

func (info *recordInfo) AddFloatField(name string, source string) string {
	return info.addField(name, source, 4, 0, FloatType, 4, 1)
}

func (info *recordInfo) AddDoubleField(name string, source string) string {
	return info.addField(name, source, 8, 0, DoubleType, 8, 1)
}

func (info *recordInfo) AddStringField(name string, source string, size int) string {
	return info.addField(name, source, size, 0, StringType, uintptr(size), 1)
}

func (info *recordInfo) AddWStringField(name string, source string, size int) string {
	return info.addField(name, source, size, 0, WStringType, uintptr(size)*2, 1)
}

func (info *recordInfo) AddV_StringField(name string, source string, size int) string {
	return info.addField(name, source, size, 0, V_StringType, 4, 0)
}

func (info *recordInfo) AddV_WStringField(name string, source string, size int) string {
	return info.addField(name, source, size, 0, V_WStringType, 4, 0)
}

func (info *recordInfo) AddDateField(name string, source string) string {
	return info.addField(name, source, 10, 0, DateType, 10, 1)
}

func (info *recordInfo) AddDateTimeField(name string, source string) string {
	return info.addField(name, source, 19, 0, DateTimeType, 19, 1)
}

func (info *recordInfo) addField(name string, source string, size int, scale int, fieldType string, fixedLen uintptr, nullByteLen uintptr) string {
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
		fixedValue:  make([]byte, fixedLen+nullByteLen),
	})
	info.fieldNames[actualName] = info.numFields
	info.numFields++
	info.fixedLen += fixedLen + nullByteLen
	return actualName
}
