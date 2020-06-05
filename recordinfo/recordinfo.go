package recordinfo

import "C"
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

	GetIntValueFrom(fieldName string, record unsafe.Pointer) (value int, isNull bool, err error)
	GetBoolValueFrom(fieldName string, record unsafe.Pointer) (value bool, isNull bool, err error)
	GetFloatValueFrom(fieldName string, record unsafe.Pointer) (value float64, isNull bool, err error)
	GetStringValueFrom(fieldName string, record unsafe.Pointer) (value string, isNull bool, err error)
	GetDateValueFrom(fieldName string, record unsafe.Pointer) (value time.Time, isNull bool, err error)
	GetRawBytesFrom(fieldName string, record unsafe.Pointer) (value []byte, err error)
	GetRawBytesFromIndex(index int, record unsafe.Pointer) (value []byte, err error)

	SetIntField(fieldName string, value int) error
	SetBoolField(fieldName string, value bool) error
	SetFloatField(fieldName string, value float64) error
	SetStringField(fieldName string, value string) error
	SetDateField(fieldName string, value time.Time) error
	SetFieldNull(fieldName string) error
	SetFromRawBytes(fieldName string, value []byte) error
	SetIndexFromRawBytes(index int, value []byte) error

	GenerateRecord() (unsafe.Pointer, error)
	ToXml(connection string) (string, error)
	FixedSize() int
	TotalSize(unsafe.Pointer) int
}

type recordInfo struct {
	fixedLen   uintptr
	numFields  int
	fields     []*fieldInfoEditor
	fieldNames map[string]int
	blob       []byte
	blobHandle unsafe.Pointer
	blobLen    int
}

type FieldInfo struct {
	Name      string
	Source    string
	Size      int
	Precision int
	Type      FieldType
}

type fieldInfoEditor struct {
	Name        string
	Source      string
	Size        int
	Precision   int
	Type        FieldType
	location    uintptr
	fixedLen    uintptr
	nullByteLen uintptr
	value       []byte
	varLen      int
	isNull      bool
}

func New() RecordInfo {
	return &recordInfo{
		fieldNames: map[string]int{},
		blobLen:    0,
	}
}

func (info *recordInfo) NumFields() int {
	return info.numFields
}

func (info *recordInfo) GetFieldByIndex(index int) (FieldInfo, error) {
	if index < 0 || index >= info.numFields {
		return FieldInfo{}, fmt.Errorf(`index was not between 0 and %v`, info.numFields)
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

func (info *recordInfo) FixedSize() int {
	return int(info.fixedLen)
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

func (info *recordInfo) TotalSize(record unsafe.Pointer) int {
	variable := int(*((*uint32)(unsafe.Pointer(uintptr(record) + info.fixedLen))))
	return int(info.fixedLen) + 4 + variable
}

func (info *recordInfo) getRecordSizes() (fixed int, variable int) {
	fixed = int(info.fixedLen)
	variable = 0
	for _, field := range info.fields {
		variable += calcVarSizeFromLen(field.varLen)
	}
	return
}

func calcVarSizeFromLen(valueLen int) int {
	if valueLen < 4 {
		return 0
	}
	if valueLen < 128 {
		return valueLen + 1
	}
	return valueLen + 4
}
