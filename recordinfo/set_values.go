package recordinfo

import (
	"fmt"
	"time"
)

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

func (info *recordInfo) SetV_StringField(fieldName string, value string) error {
	return info.setValue(fieldName, value)
}

func (info *recordInfo) SetV_WStringField(fieldName string, value string) error {
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

func (info *recordInfo) SetFromInterface(fieldName string, value interface{}) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	var ok bool
	switch field.Type {
	case ByteType:
		field.value, ok = value.(byte)
	case BoolType:
		field.value, ok = value.(bool)
	case Int16Type:
		field.value, ok = value.(int16)
	case Int32Type:
		field.value, ok = value.(int32)
	case Int64Type:
		field.value, ok = value.(int64)
	case FloatType:
		field.value, ok = value.(float32)
	case FixedDecimalType, DoubleType:
		field.value, ok = value.(float64)
	case StringType, WStringType, V_StringType, V_WStringType:
		field.value, ok = value.(string)
	case DateType, DateTimeType:
		field.value, ok = value.(time.Time)
	default:
		return fmt.Errorf(`[%v] field type '%v' is not valid`, field.Name, field.Type)
	}
	if !ok {
		return fmt.Errorf(`tried to set [%v] with a %T but the field is a %v field`, field.Name, value, field.Type)
	}
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
