// Package recordinfo provides all of the functionality to read and generate Alteryx records
package recordinfo

import "C"
import (
	"fmt"
	"github.com/tlarsen7572/goalteryx/recordblob"
	"unsafe"
)

// recordInfo is the struct which implements the RecordInfo interface
type recordInfo struct {
	fixedLen   uintptr
	numFields  int
	fields     []*fieldInfoEditor
	fieldNames map[string]int
	blob       []byte
	blobHandle unsafe.Pointer
	blobLen    int
}

// FieldInfo is the struct used to provide field information back to calling code.
// No logic is performed by FieldInfo; it is purely a data structure
type FieldInfo struct {
	Name      string
	Source    string
	Size      int
	Precision int
	Type      FieldType
}

// fieldInfoEditor is used by recordInfo for persisting record values and reading record values from record blobs.
// All of the logic of setting and getting field values are performed by this struct.
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

func (info *recordInfo) GetFieldByName(name string) (FieldInfo, error) {
	index, ok := info.fieldNames[name]
	if !ok {
		return FieldInfo{}, fmt.Errorf(`field '%v' does not exist in the RecordInfo`, name)
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

// checkFieldName checks whether the provided field name already exists in the RecordInfo.  If it does,
// the function appends the fieldname with '2' and checks again.  This happens until an unused name is found,
// which is then returned to the caller.
func (info *recordInfo) checkFieldName(name string) string {
	_, exists := info.fieldNames[name]
	for exists {
		name = name + `2`
		_, exists = info.fieldNames[name]
	}
	return name
}

// Total size of a record blob is the fixed size plus 4 bytes for the variable length plus the variable length
func (info *recordInfo) TotalSize(record recordblob.RecordBlob) int {
	variable := int(*((*uint32)(unsafe.Pointer(uintptr(record.Blob()) + info.fixedLen))))
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
