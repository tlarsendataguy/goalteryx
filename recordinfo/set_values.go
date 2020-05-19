package recordinfo

import "time"

func (info *recordInfo) SetByteField(fieldName string, value byte) error {
	return info.setValue(fieldName, value)
}

func (info *recordInfo) SetBoolField(fieldName string, value bool) error {
	return info.setValue(fieldName, value)
}

func (info *recordInfo) SetInt16Field(fieldName string, value int16) error {
	return info.setValue(fieldName, value)
}

func (info *recordInfo) SetInt32Field(fieldName string, value int32) error {
	return info.setValue(fieldName, value)
}

func (info *recordInfo) SetInt64Field(fieldName string, value int64) error {
	return info.setValue(fieldName, value)
}

func (info *recordInfo) SetFixedDecimalField(fieldName string, value float64) error {
	return info.setValue(fieldName, value)
}

func (info *recordInfo) SetFloatField(fieldName string, value float32) error {
	return info.setValue(fieldName, value)
}

func (info *recordInfo) SetDoubleField(fieldName string, value float64) error {
	return info.setValue(fieldName, value)
}

func (info *recordInfo) SetStringField(fieldName string, value string) error {
	return info.setValue(fieldName, value)
}

func (info *recordInfo) SetWStringField(fieldName string, value string) error {
	return info.setValue(fieldName, value)
}

func (info *recordInfo) SetDateField(fieldName string, value time.Time) error {
	return info.setValue(fieldName, value)
}

func (info *recordInfo) SetDateTimeField(fieldName string, value time.Time) error {
	return info.setValue(fieldName, value)
}

func (info *recordInfo) SetFieldNull(fieldName string) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return nil
	}
	field.value = nil
	return nil
}

func (info *recordInfo) setValue(fieldName string, value interface{}) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}
	field.value = value
	return nil
}
