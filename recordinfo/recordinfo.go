// Package recordinfo provides all of the functionality to read and generate Alteryx records
package recordinfo

import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

// RecordInfo is the interface which defines all of the behaviors needed to read and generate Alteryx records.
type RecordInfo interface {
	// NumFields return the number of fields contained in the RecordInfo.
	NumFields() int

	// GetFieldByIndex returns field information at the specified index.  If an out-of-range index is specified,
	// an error is returned with an empty FieldInfo struct.
	GetFieldByIndex(index int) (FieldInfo, error)

	// AddByteField appends a Byte field to the end of the RecordInfo.  It returns a string with the actual field
	// name that was added.  If a field with the specified name already exists in the RecordInfo, this function will
	// add '2' to the end of the name until it finds a unique name.  This prevents duplicate fields from occurring
	// in the data and let's the calling code know the actual name of the field that was added.
	AddByteField(name string, source string) string

	// AddBoolField appends a Bool field to the end of the RecordInfo.  It returns a string with the actual field
	// name that was added.  If a field with the specified name already exists in the RecordInfo, this function will
	// add '2' to the end of the name until it finds a unique name.  This prevents duplicate fields from occurring
	// in the data and let's the calling code know the actual name of the field that was added.
	AddBoolField(name string, source string) string

	// AddInt16Field appends an Int16 field to the end of the RecordInfo.  It returns a string with the actual field
	// name that was added.  If a field with the specified name already exists in the RecordInfo, this function will
	// add '2' to the end of the name until it finds a unique name.  This prevents duplicate fields from occurring
	// in the data and let's the calling code know the actual name of the field that was added.
	AddInt16Field(name string, source string) string

	// AddInt32Field appends an Int32 field to the end of the RecordInfo.  It returns a string with the actual field
	// name that was added.  If a field with the specified name already exists in the RecordInfo, this function will
	// add '2' to the end of the name until it finds a unique name.  This prevents duplicate fields from occurring
	// in the data and let's the calling code know the actual name of the field that was added.
	AddInt32Field(name string, source string) string

	// AddInt64Field appends an Int64 field to the end of the RecordInfo.  It returns a string with the actual field
	// name that was added.  If a field with the specified name already exists in the RecordInfo, this function will
	// add '2' to the end of the name until it finds a unique name.  This prevents duplicate fields from occurring
	// in the data and let's the calling code know the actual name of the field that was added.
	AddInt64Field(name string, source string) string

	// AddFixedDecimalField appends a FixedDecimal field to the end of the RecordInfo.  It returns a string with the actual field
	// name that was added.  If a field with the specified name already exists in the RecordInfo, this function will
	// add '2' to the end of the name until it finds a unique name.  This prevents duplicate fields from occurring
	// in the data and let's the calling code know the actual name of the field that was added.
	AddFixedDecimalField(name string, source string, size int, precision int) string

	// AddFloatField appends a Float field to the end of the RecordInfo.  It returns a string with the actual field
	// name that was added.  If a field with the specified name already exists in the RecordInfo, this function will
	// add '2' to the end of the name until it finds a unique name.  This prevents duplicate fields from occurring
	// in the data and let's the calling code know the actual name of the field that was added.
	AddFloatField(name string, source string) string

	// AddDoubleField appends a Double field to the end of the RecordInfo.  It returns a string with the actual field
	// name that was added.  If a field with the specified name already exists in the RecordInfo, this function will
	// add '2' to the end of the name until it finds a unique name.  This prevents duplicate fields from occurring
	// in the data and let's the calling code know the actual name of the field that was added.
	AddDoubleField(name string, source string) string

	// AddStringField appends a String field to the end of the RecordInfo.  It returns a string with the actual field
	// name that was added.  If a field with the specified name already exists in the RecordInfo, this function will
	// add '2' to the end of the name until it finds a unique name.  This prevents duplicate fields from occurring
	// in the data and let's the calling code know the actual name of the field that was added.
	AddStringField(name string, source string, size int) string

	// AddWStringField appends a WString field to the end of the RecordInfo.  It returns a string with the actual field
	// name that was added.  If a field with the specified name already exists in the RecordInfo, this function will
	// add '2' to the end of the name until it finds a unique name.  This prevents duplicate fields from occurring
	// in the data and let's the calling code know the actual name of the field that was added.
	AddWStringField(name string, source string, size int) string

	// AddV_StringField appends a V_String field to the end of the RecordInfo.  It returns a string with the actual field
	// name that was added.  If a field with the specified name already exists in the RecordInfo, this function will
	// add '2' to the end of the name until it finds a unique name.  This prevents duplicate fields from occurring
	// in the data and let's the calling code know the actual name of the field that was added.
	AddV_StringField(name string, source string, size int) string

	// AddV_WStringField appends a V_WString field to the end of the RecordInfo.  It returns a string with the actual field
	// name that was added.  If a field with the specified name already exists in the RecordInfo, this function will
	// add '2' to the end of the name until it finds a unique name.  This prevents duplicate fields from occurring
	// in the data and let's the calling code know the actual name of the field that was added.
	AddV_WStringField(name string, source string, size int) string

	// AddDateField appends a Date field to the end of the RecordInfo.  It returns a string with the actual field
	// name that was added.  If a field with the specified name already exists in the RecordInfo, this function will
	// add '2' to the end of the name until it finds a unique name.  This prevents duplicate fields from occurring
	// in the data and let's the calling code know the actual name of the field that was added.
	AddDateField(name string, source string) string

	// AddDateTimeField appends a DateTime field to the end of the RecordInfo.  It returns a string with the actual field
	// name that was added.  If a field with the specified name already exists in the RecordInfo, this function will
	// add '2' to the end of the name until it finds a unique name.  This prevents duplicate fields from occurring
	// in the data and let's the calling code know the actual name of the field that was added.
	AddDateTimeField(name string, source string) string

	// AddBlobField appends a Blob field to the end of the RecordInfo.  It returns a string with the actual field
	// name that was added.  If a field with the specified name already exists in the RecordInfo, this function will
	// add '2' to the end of the name until it finds a unique name.  This prevents duplicate fields from occurring
	// in the data and let's the calling code know the actual name of the field that was added.
	AddBlobField(name string, source string, size int) string

	// AddSpatialField appends a SpatialObj field to the end of the RecordInfo.  It returns a string with the actual field
	// name that was added.  If a field with the specified name already exists in the RecordInfo, this function will
	// add '2' to the end of the name until it finds a unique name.  This prevents duplicate fields from occurring
	// in the data and let's the calling code know the actual name of the field that was added.
	AddSpatialField(name string, source string, size int) string

	// GetIntValueFrom obtains a value from the specific integer field in the record.  It can only be called on
	// Byte, Int16, Int32, and Int64 fields.  All other fields will return an error.
	GetIntValueFrom(fieldName string, record unsafe.Pointer) (value int, isNull bool, err error)

	// GetBoolValueFrom obtains a value from the specific bool field in the record.  It can only be called on
	// Bool fields.  All other fields will return an error.
	GetBoolValueFrom(fieldName string, record unsafe.Pointer) (value bool, isNull bool, err error)

	// GetFloatValueFrom obtains a value from the specific decimal field in the record.  It can only be called on
	// Float, Double, and FixedDecimal fields.  All other fields will return an error.
	GetFloatValueFrom(fieldName string, record unsafe.Pointer) (value float64, isNull bool, err error)

	// GetStringValueFrom obtains a value from the specific text field in the record.  It can only be called on
	// String, WString, V_String, and V_WString fields.  All other fields will return an error.
	GetStringValueFrom(fieldName string, record unsafe.Pointer) (value string, isNull bool, err error)

	// GetDateValueFrom obtains a value from the specific date/datetime field in the record.  It can only be called on
	// Date and DateTime fields.  All other fields will return an error.
	GetDateValueFrom(fieldName string, record unsafe.Pointer) (value time.Time, isNull bool, err error)

	// GetRawBytesFrom obtains a value from the specific field in the record.  It can
	// be called on any field type.
	GetRawBytesFrom(fieldName string, record unsafe.Pointer) (value []byte, err error)

	// GetRawBytesFromIndex obtains a value from the specific position in the record.  It can
	// be called on any field type and is the fastest way to obtain a value from a field.
	GetRawBytesFromIndex(index int, record unsafe.Pointer) (value []byte, err error)

	// SetIntField sets the specified integer field with a value.  It can only be called on
	// Byte, Int16, Int32, and Int64 fields.  All other fields will return an error.
	SetIntField(fieldName string, value int) error

	// SetBoolField sets the specified bool field with a value.  It can only be called on
	// Bool fields.  All other fields will return an error.
	SetBoolField(fieldName string, value bool) error

	// SetFloatField sets the specified decimal field with a value.  It can only be called on
	// Float, Double, and FixedDecimal fields.  All other fields will return an error.
	SetFloatField(fieldName string, value float64) error

	// SetStringField sets the specified text field with a value.  It can only be called on
	// String, WString, V_String, and V_WString fields.  All other fields will return an error.
	SetStringField(fieldName string, value string) error

	// SetDateField sets the specified date/datetime field with a value.  It can only be called on
	// Date and DateTime fields.  All other fields will return an error.
	SetDateField(fieldName string, value time.Time) error

	// SetFieldNull sets the specified field to NULL.  It can can be called on any field.
	SetFieldNull(fieldName string) error

	// SetFromRawBytes sets the specified field with a value.  It can be called on any field, but it is
	// up to the caller to ensure the bytes provided are appropriate for the field.
	SetFromRawBytes(fieldName string, value []byte) error

	// SetIndexFromRawBytes sets the specified field position with a value.  It can be called on any field, but it is
	// up to the caller to ensure the bytes provided are appropriate for the field.  This is the fastest way to
	// set a record's value.
	SetIndexFromRawBytes(index int, value []byte) error

	// GenerateRecord creates a record blob from the current field values that can be passed to downstream tools.
	GenerateRecord() (unsafe.Pointer, error)

	// ToXml outputs a string XML representation of the RecordInfo object.  This allows passing the RecordInfo
	// metadata to downstream tools.
	ToXml(connection string) (string, error)

	// FixedSize returns the total length of the fixed-size portion of the RecordInfo object.  This is used to
	// identify the fixed portion of the generated record blobs.
	FixedSize() int

	// TotalSize returns the total length of the specified record blob.
	TotalSize(unsafe.Pointer) int
}

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

// New returns a RecordInfo object.  The recordInfo struct is not exported, so this is the only way to obtain
// a RecordInfo object.
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
