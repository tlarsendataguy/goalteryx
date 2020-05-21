package recordinfo

import (
	"time"
)

func (info *recordInfo) SetIntField(fieldName string, value int) error {
	return nil
}

func (info *recordInfo) SetBoolField(fieldName string, value bool) error {
	return nil
}

func (info *recordInfo) SetFloatField(fieldName string, value float64) error {
	return nil
}

func (info *recordInfo) SetStringField(fieldName string, value string) error {
	return nil
}

func (info *recordInfo) SetDateField(fieldName string, value time.Time) error {
	return nil
}

func (info *recordInfo) SetFieldNull(fieldName string) error {
	field, err := info.getFieldInfo(fieldName)
	if err != nil {
		return nil
	}
	field.fixedValue = nil
	return nil
}

func (info *recordInfo) SetFromRawBytes(fieldName string, value []byte) error {
	_, err := info.getFieldInfo(fieldName)
	if err != nil {
		return err
	}

	return nil
}
