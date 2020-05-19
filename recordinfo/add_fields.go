package recordinfo

func (info *recordInfo) AddByteField(name string, source string) string {
	return info.addField(name, source, 1, 0, ByteType, 1, 1, generateByte)
}

func (info *recordInfo) AddBoolField(name string, source string) string {
	return info.addField(name, source, 1, 0, BoolType, 1, 0, generateBool)
}

func (info *recordInfo) AddInt16Field(name string, source string) string {
	return info.addField(name, source, 2, 0, Int16Type, 2, 1, generateInt16)
}

func (info *recordInfo) AddInt32Field(name string, source string) string {
	return info.addField(name, source, 4, 0, Int32Type, 4, 1, generateInt32)
}

func (info *recordInfo) AddInt64Field(name string, source string) string {
	return info.addField(name, source, 8, 0, Int64Type, 8, 1, generateInt64)
}

func (info *recordInfo) AddFixedDecimalField(name string, source string, size int, precision int) string {
	return info.addField(name, source, size, precision, FixedDecimalType, uintptr(size), 1, generateFixedDecimal)
}

func (info *recordInfo) AddFloatField(name string, source string) string {
	return info.addField(name, source, 4, 0, FloatType, 4, 1, generateFloat32)
}

func (info *recordInfo) AddDoubleField(name string, source string) string {
	return info.addField(name, source, 8, 0, DoubleType, 8, 1, generateFloat64)
}

func (info *recordInfo) AddStringField(name string, source string, size int) string {
	return info.addField(name, source, size, 0, StringType, uintptr(size), 1, generateString)
}

func (info *recordInfo) AddWStringField(name string, source string, size int) string {
	return info.addField(name, source, size, 0, WStringType, uintptr(size)*2, 1, generateWString)
}

func (info *recordInfo) AddV_StringField(name string, source string, size int) string {
	return info.addField(name, source, size, 0, V_StringType, 4, 0, nil)
}

func (info *recordInfo) AddV_WStringField(name string, source string, size int) string {
	return info.addField(name, source, size, 0, V_WStringType, 4, 0, nil)
}

func (info *recordInfo) AddDateField(name string, source string) string {
	return info.addField(name, source, 10, 0, DateType, 10, 1, generateDate)
}

func (info *recordInfo) AddDateTimeField(name string, source string) string {
	return info.addField(name, source, 19, 0, DateTimeType, 19, 1, generateDateTime)
}

func (info *recordInfo) addField(name string, source string, size int, scale int, fieldType string, fixedLen uintptr, nullByteLen uintptr, generator generateBytes) string {
	actualName := info.checkFieldName(name)
	info.fields = append(info.fields, &fieldInfoEditor{
		Name:        actualName,
		Source:      source,
		Size:        size,
		Precision:   scale,
		Type:        fieldType,
		location:    info.currentLen,
		fixedLen:    fixedLen,
		nullByteLen: nullByteLen,
		generator:   generator,
	})
	info.fieldNames[actualName] = info.numFields
	info.numFields++
	info.currentLen += fixedLen + nullByteLen
	return actualName
}
