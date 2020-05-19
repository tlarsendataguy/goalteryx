package recordinfo

import (
	"fmt"
	"time"
	"unsafe"
)

type RecordInfo interface {
	NumFields() int
	GetFieldByIndex(index int) (FieldInfo, error)
	AddByteField(name string, source string) string
	AddBoolField(name string, source string) string
	AddInt16Field(name string, source string) string
	AddInt32Field(name string, source string) string
	AddInt64Field(name string, source string) string
	AddFixedDecimalField(name string, source string, size int, precision int) string
	AddFloatField(name string, source string) string
	AddDoubleField(name string, source string) string
	AddStringField(name string, source string, size int) string
	AddWStringField(name string, source string, size int) string
	AddV_StringField(name string, source string, size int) string
	AddV_WStringField(name string, source string, size int) string
	AddDateField(name string, source string) string
	AddDateTimeField(name string, source string) string

	GetByteValueFrom(fieldName string, record unsafe.Pointer) (value byte, isNull bool, err error)
	GetBoolValueFrom(fieldName string, record unsafe.Pointer) (value bool, isNull bool, err error)
	GetInt16ValueFrom(fieldName string, record unsafe.Pointer) (value int16, isNull bool, err error)
	GetInt32ValueFrom(fieldName string, record unsafe.Pointer) (value int32, isNull bool, err error)
	GetInt64ValueFrom(fieldName string, record unsafe.Pointer) (value int64, isNull bool, err error)
	GetFixedDecimalValueFrom(fieldName string, record unsafe.Pointer) (value float64, isNull bool, err error)
	GetFloatValueFrom(fieldName string, record unsafe.Pointer) (value float32, isNull bool, err error)
	GetDoubleValueFrom(fieldName string, record unsafe.Pointer) (value float64, isNull bool, err error)
	GetStringValueFrom(fieldName string, record unsafe.Pointer) (value string, isNull bool, err error)
	GetWStringValueFrom(fieldName string, record unsafe.Pointer) (value string, isNull bool, err error)
	GetDateValueFrom(fieldName string, record unsafe.Pointer) (value time.Time, isNull bool, err error)
	GetDateTimeValueFrom(fieldName string, record unsafe.Pointer) (value time.Time, isNull bool, err error)
	GetInterfaceValueFrom(fieldName string, record unsafe.Pointer) (value interface{}, isNull bool, err error)

	SetByteField(fieldName string, value byte) error
	SetBoolField(fieldName string, value bool) error
	SetInt16Field(fieldName string, value int16) error
	SetInt32Field(fieldName string, value int32) error
	SetInt64Field(fieldName string, value int64) error
	SetFixedDecimalField(fieldName string, value float64) error
	SetFloatField(fieldName string, value float32) error
	SetDoubleField(fieldName string, value float64) error
	SetStringField(fieldName string, value string) error
	SetWStringField(fieldName string, value string) error
	SetDateField(fieldName string, value time.Time) error
	SetDateTimeField(fieldName string, value time.Time) error
	SetFieldNull(fieldName string) error

	GenerateRecord() (unsafe.Pointer, error)
	ToXml(connection string) (string, error)
}

type recordInfo struct {
	currentLen uintptr
	numFields  int
	fields     []*fieldInfoEditor
	fieldNames map[string]int
	blob       []byte
}

type FieldInfo struct {
	Name      string
	Source    string
	Size      int
	Precision int
	Type      string
}

type fieldInfoEditor struct {
	Name        string
	Source      string
	Size        int
	Precision   int
	Type        string
	location    uintptr
	fixedLen    uintptr
	nullByteLen uintptr
	value       interface{}
	generator   generateBytes
}

func New() RecordInfo {
	return &recordInfo{fieldNames: map[string]int{}}
}

func (info *recordInfo) NumFields() int {
	return info.numFields
}

func (info *recordInfo) GetFieldByIndex(index int) (FieldInfo, error) {
	if count := len(info.fields); index < 0 || index >= count {
		return FieldInfo{}, fmt.Errorf(`index was not between 0 and %v`, count)
	}
	field := info.fields[index]
	return FieldInfo{
		Name:      field.Name,
		Source:    field.Source,
		Size:      field.Size,
		Precision: field.Precision,
		Type:      field.Type,
	}, nil
}

func (info *recordInfo) getFieldInfo(fieldName string) (*fieldInfoEditor, error) {
	index, ok := info.fieldNames[fieldName]
	if !ok {
		return nil, fmt.Errorf(`field '%v' does not exist`, fieldName)
	}
	return info.fields[index], nil
}

func (info *recordInfo) checkFieldName(name string) string {
	_, exists := info.fieldNames[name]
	for exists {
		name = name + `2`
		_, exists = info.fieldNames[name]
	}
	return name
}

func (info *recordInfo) getRecordSize() int {
	size := 0
	for _, field := range info.fields {
		size += int(field.fixedLen + field.nullByteLen)
	}
	return size
}
