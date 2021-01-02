package api_new

import (
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
)

type NewOutgoingField func() *outgoingField

func NewBoolField(name string, source string) NewOutgoingField {
	return func() *outgoingField {
		return &outgoingField{
			Name:         name,
			Type:         `Bool`,
			Source:       source,
			Size:         1,
			Scale:        0,
			CurrentValue: make([]byte, 1),
			isFixedLen:   true,
		}
	}
}

func NewByteField(name string, source string) NewOutgoingField {
	return func() *outgoingField {
		return &outgoingField{
			Name:         name,
			Type:         `Byte`,
			Source:       source,
			Size:         1,
			CurrentValue: make([]byte, 2),
			nullSetter:   setNormalFieldNull,
			nullGetter:   getNormalFieldNull,
			intSetter:    setByte,
			intGetter:    getByte,
			isFixedLen:   true,
		}
	}
}

func NewInt16Field(name string, source string) NewOutgoingField {
	return func() *outgoingField {
		return &outgoingField{
			Name:         name,
			Type:         `Int16`,
			Source:       source,
			Size:         2,
			CurrentValue: make([]byte, 3),
			nullSetter:   setNormalFieldNull,
			nullGetter:   getNormalFieldNull,
			intSetter:    setInt16,
			intGetter:    getInt16,
			isFixedLen:   true,
		}
	}
}

func NewInt32Field(name string, source string) NewOutgoingField {
	return func() *outgoingField {
		return &outgoingField{
			Name:         name,
			Type:         `Int32`,
			Source:       source,
			Size:         4,
			CurrentValue: make([]byte, 5),
			nullSetter:   setNormalFieldNull,
			nullGetter:   getNormalFieldNull,
			intSetter:    setInt32,
			intGetter:    getInt32,
			isFixedLen:   true,
		}
	}
}

func NewInt64Field(name string, source string) NewOutgoingField {
	return func() *outgoingField {
		return &outgoingField{
			Name:         name,
			Type:         `Int64`,
			Source:       source,
			Size:         8,
			CurrentValue: make([]byte, 9),
			nullSetter:   setNormalFieldNull,
			nullGetter:   getNormalFieldNull,
			intSetter:    setInt64,
			intGetter:    getInt64,
			isFixedLen:   true,
		}
	}
}

func NewFloatField(name string, source string) NewOutgoingField {
	return func() *outgoingField {
		return &outgoingField{
			Name:         name,
			Type:         `Float`,
			Source:       source,
			Size:         4,
			CurrentValue: make([]byte, 5),
			nullSetter:   setNormalFieldNull,
			nullGetter:   getNormalFieldNull,
			floatSetter:  setFloat,
			floatGetter:  getFloat,
			isFixedLen:   true,
		}
	}
}

func NewDoubleField(name string, source string) NewOutgoingField {
	return func() *outgoingField {
		return &outgoingField{
			Name:         name,
			Type:         `Double`,
			Source:       source,
			Size:         8,
			CurrentValue: make([]byte, 9),
			nullSetter:   setNormalFieldNull,
			nullGetter:   getNormalFieldNull,
			floatSetter:  setDouble,
			floatGetter:  getDouble,
			isFixedLen:   true,
		}
	}
}

func NewFixedDecimalField(name string, source string, size int, scale int) NewOutgoingField {
	return func() *outgoingField {
		return &outgoingField{
			Name:            name,
			Type:            `FixedDecimal`,
			Source:          source,
			Size:            size,
			Scale:           scale,
			fixedDecimalFmt: fmt.Sprintf(`%%%d.%df`, size, scale),
			CurrentValue:    make([]byte, size+1),
			nullSetter:      setNormalFieldNull,
			nullGetter:      getNormalFieldNull,
			floatSetter:     setFixedDecimal,
			floatGetter:     getFixedDecimal,
			isFixedLen:      true,
		}
	}
}

func NewDateField(name string, source string) NewOutgoingField {
	return func() *outgoingField {
		return &outgoingField{
			Name:           name,
			Type:           `Date`,
			Source:         source,
			Size:           10,
			CurrentValue:   make([]byte, 11),
			nullSetter:     setNormalFieldNull,
			nullGetter:     getNormalFieldNull,
			dateTimeSetter: setDate,
			dateTimeGetter: getDate,
			isFixedLen:     true,
		}
	}
}

func NewDateTimeField(name string, source string) NewOutgoingField {
	return func() *outgoingField {
		return &outgoingField{
			Name:           name,
			Type:           `DateTime`,
			Source:         source,
			Size:           19,
			CurrentValue:   make([]byte, 20),
			nullSetter:     setNormalFieldNull,
			nullGetter:     getNormalFieldNull,
			dateTimeSetter: setDateTime,
			dateTimeGetter: getDateTime,
			isFixedLen:     true,
		}
	}
}

func NewStringField(name string, source string, size int) NewOutgoingField {
	return func() *outgoingField {
		return &outgoingField{
			Name:         name,
			Type:         `String`,
			Source:       source,
			Size:         size,
			CurrentValue: make([]byte, size+1),
			nullSetter:   setNormalFieldNull,
			nullGetter:   getNormalFieldNull,
			stringSetter: setString,
			stringGetter: getString,
			isFixedLen:   true,
		}
	}
}

func NewWStringField(name string, source string, size int) NewOutgoingField {
	return func() *outgoingField {
		return &outgoingField{
			Name:         name,
			Type:         `WString`,
			Source:       source,
			Size:         size,
			CurrentValue: make([]byte, (size*2)+1),
			nullSetter:   setWideFieldNull,
			nullGetter:   getWideFieldNull,
			stringSetter: setWString,
			stringGetter: getWString,
			isFixedLen:   true,
		}
	}
}

func NewV_StringField(name string, source string, size int) NewOutgoingField {
	return func() *outgoingField {
		return &outgoingField{
			Name:         name,
			Type:         `V_String`,
			Source:       source,
			Size:         size,
			CurrentValue: make([]byte, 1),
			nullSetter:   setVarFieldNull,
			nullGetter:   getVarFieldNull,
			stringSetter: setV_String,
			stringGetter: getV_String,
			isFixedLen:   false,
		}
	}
}

func NewV_WStringField(name string, source string, size int) NewOutgoingField {
	return func() *outgoingField {
		return &outgoingField{
			Name:         name,
			Type:         `V_WString`,
			Source:       source,
			Size:         size,
			CurrentValue: make([]byte, 1),
			nullSetter:   setVarFieldNull,
			nullGetter:   getVarFieldNull,
			stringSetter: setV_WString,
			stringGetter: getV_WString,
			isFixedLen:   false,
		}
	}
}

func NewBlobField(name string, source string, size int) NewOutgoingField {
	return func() *outgoingField {
		return &outgoingField{
			Name:         name,
			Type:         `Blob`,
			Source:       source,
			Size:         size,
			CurrentValue: make([]byte, 1),
			nullSetter:   setVarFieldNull,
			nullGetter:   getVarFieldNull,
			blobSetter:   setBlob,
			blobGetter:   getBlob,
			isFixedLen:   false,
		}
	}
}

func NewSpatialObjField(name string, source string, size int) NewOutgoingField {
	return func() *outgoingField {
		return &outgoingField{
			Name:         name,
			Type:         `SpatialObj`,
			Source:       source,
			Size:         size,
			CurrentValue: make([]byte, 1),
			nullSetter:   setVarFieldNull,
			nullGetter:   getVarFieldNull,
			blobSetter:   setBlob,
			blobGetter:   getBlob,
			isFixedLen:   false,
		}
	}
}

type outgoingField struct {
	XMLName         string                          `xml:"Field"`
	Name            string                          `xml:"name,attr"`
	Type            string                          `xml:"type,attr"`
	Source          string                          `xml:"source,attr"`
	Size            int                             `xml:"size,attr"`
	Scale           int                             `xml:"scale,attr"`
	CopyFrom        BytesGetter                     `xml:"-"`
	CurrentValue    []byte                          `xml:"-"`
	isFixedLen      bool                            `xml:"-"`
	nullSetter      func(byte, *outgoingField)      `xml:"-"`
	nullGetter      func(*outgoingField) bool       `xml:"-"`
	intSetter       func(int, *outgoingField)       `xml:"-"`
	intGetter       func(*outgoingField) int        `xml:"-"`
	floatSetter     func(float64, *outgoingField)   `xml:"-"`
	floatGetter     func(*outgoingField) float64    `xml:"-"`
	dateTimeSetter  func(time.Time, *outgoingField) `xml:"-"`
	dateTimeGetter  func(*outgoingField) time.Time  `xml:"-"`
	fixedDecimalFmt string                          `xml:"-"`
	stringSetter    func(string, *outgoingField)    `xml:"-"`
	stringGetter    func(*outgoingField) string     `xml:"-"`
	blobSetter      func([]byte, *outgoingField)    `xml:"-"`
	blobGetter      func(*outgoingField) []byte     `xml:"-"`
}

func (f *outgoingField) dataSize() int {
	return len(f.CurrentValue)
}

func setNormalFieldNull(isNull byte, f *outgoingField) {
	f.CurrentValue[f.Size] = isNull
}

func setWideFieldNull(isNull byte, f *outgoingField) {
	f.CurrentValue[f.Size*2] = isNull
}

func setVarFieldNull(isNull byte, f *outgoingField) {
	f.CurrentValue[0] = isNull
}

func getNormalFieldNull(f *outgoingField) bool {
	return f.CurrentValue[f.Size] == 1
}

func getWideFieldNull(f *outgoingField) bool {
	return f.CurrentValue[f.Size*2] == 1
}

func getVarFieldNull(f *outgoingField) bool {
	return f.CurrentValue[0] == 1
}

func (f *outgoingField) SetBool(value bool) {
	if value {
		f.CurrentValue[0] = 1
	} else {
		f.CurrentValue[0] = 0
	}
}

func (f *outgoingField) SetNullBool() {
	f.CurrentValue[0] = 2
}

func (f *outgoingField) GetCurrentBool() (bool, bool) {
	if f.CurrentValue[0] == 2 {
		return false, true
	}
	return f.CurrentValue[0] == 1, false
}

func getByte(f *outgoingField) int {
	return int(f.CurrentValue[0])
}

func setByte(value int, f *outgoingField) {
	f.CurrentValue[0] = byte(value)
}

func getInt16(f *outgoingField) int {
	return int(binary.LittleEndian.Uint16(f.CurrentValue[:2]))
}

func setInt16(value int, f *outgoingField) {
	binary.LittleEndian.PutUint16(f.CurrentValue[:2], uint16(value))
}

func getInt32(f *outgoingField) int {
	return int(binary.LittleEndian.Uint32(f.CurrentValue[:4]))
}

func setInt32(value int, f *outgoingField) {
	binary.LittleEndian.PutUint32(f.CurrentValue[:4], uint32(value))
}

func getInt64(f *outgoingField) int {
	return int(binary.LittleEndian.Uint64(f.CurrentValue[:8]))
}

func setInt64(value int, f *outgoingField) {
	binary.LittleEndian.PutUint64(f.CurrentValue[:8], uint64(value))
}

func (f *outgoingField) SetInt(value int) {
	f.intSetter(value, f)
	f.nullSetter(0, f)
}

func (f *outgoingField) SetNullInt() {
	f.nullSetter(1, f)
}

func (f *outgoingField) GetCurrentInt() (int, bool) {
	if f.nullGetter(f) {
		return 0, true
	}
	return f.intGetter(f), false
}

func getFloat(f *outgoingField) float64 {
	return float64(math.Float32frombits(binary.LittleEndian.Uint32(f.CurrentValue[:4])))
}

func setFloat(value float64, f *outgoingField) {
	binary.LittleEndian.PutUint32(f.CurrentValue[:4], math.Float32bits(float32(value)))
}

func getDouble(f *outgoingField) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(f.CurrentValue[:8]))
}

func setDouble(value float64, f *outgoingField) {
	binary.LittleEndian.PutUint64(f.CurrentValue[:8], math.Float64bits(value))
}

func getFixedDecimal(f *outgoingField) float64 {
	valueStr := string(truncateAtNullByte(f.CurrentValue[:f.Size]))
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		println(err.Error())
	}
	return value
}

func setFixedDecimal(value float64, f *outgoingField) {
	valueStr := strings.TrimLeft(fmt.Sprintf(f.fixedDecimalFmt, value), ` `)
	copy(f.CurrentValue[:f.Size], valueStr)
	if length := len(valueStr); length < f.Size {
		f.CurrentValue[length] = 0
	}
}

func (f *outgoingField) SetFloat(value float64) {
	f.floatSetter(value, f)
	f.nullSetter(0, f)
}

func (f *outgoingField) SetNullFloat() {
	f.nullSetter(1, f)
}

func (f *outgoingField) GetCurrentFloat() (float64, bool) {
	if f.nullGetter(f) {
		return 0, true
	}
	return f.floatGetter(f), false
}

func getDate(f *outgoingField) time.Time {
	value, _ := time.Parse(dateFormat, string(f.CurrentValue[:10]))
	return value
}

func setDate(value time.Time, f *outgoingField) {
	valueStr := value.Format(dateFormat)
	copy(f.CurrentValue[:10], valueStr)
}

func getDateTime(f *outgoingField) time.Time {
	value, _ := time.Parse(dateTimeFormat, string(f.CurrentValue[:19]))
	return value
}

func setDateTime(value time.Time, f *outgoingField) {
	valueStr := value.Format(dateTimeFormat)
	copy(f.CurrentValue[:19], valueStr)
}

func (f *outgoingField) SetDateTime(value time.Time) {
	f.dateTimeSetter(value, f)
	f.nullSetter(0, f)
}

func (f *outgoingField) SetNullDateTime() {
	f.nullSetter(1, f)
}

func (f *outgoingField) GetCurrentDateTime() (time.Time, bool) {
	if f.nullGetter(f) {
		return time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), true
	}
	return f.dateTimeGetter(f), false
}

func getString(f *outgoingField) string {
	value := truncateAtNullByte(f.CurrentValue[:f.Size])
	return string(value)
}

func setString(value string, f *outgoingField) {
	length := len(value)
	if length > f.Size {
		value = value[:f.Size]
		length = f.Size
	}
	if length < f.Size {
		f.CurrentValue[length] = 0
	}
	copy(f.CurrentValue[:length], value)
}

func getWString(f *outgoingField) string {
	utf16Bytes := bytesToUtf16(f.CurrentValue[:f.Size*2])
	utf16Bytes = truncateAtNullUtf16(utf16Bytes)
	value := string(utf16.Decode(utf16Bytes))
	return value
}

func setWString(value string, f *outgoingField) {
	utf16Bytes := utf16.Encode([]rune(value))
	length := len(utf16Bytes)
	if length > f.Size {
		utf16Bytes = utf16Bytes[:f.Size]
		length = f.Size
	}
	if length < f.Size {
		utf16Bytes = append(utf16Bytes, 0)
		length++
	}
	stringBytes := utf16ToBytes(utf16Bytes)
	copy(f.CurrentValue, stringBytes)
}

func getV_String(f *outgoingField) string {
	return string(f.CurrentValue[1:])
}

func setV_String(value string, f *outgoingField) {
	bytes := []byte(value)
	if length := len(bytes); length > f.Size {
		bytes = bytes[:f.Size]
	}
	requiredLen := len(bytes) + 1
	if requiredLen > cap(f.CurrentValue) {
		f.CurrentValue = make([]byte, requiredLen)
	}
	copy(f.CurrentValue[1:], value)
	f.CurrentValue = f.CurrentValue[:requiredLen]
}

func getV_WString(f *outgoingField) string {
	utf16Bytes := bytesToUtf16(f.CurrentValue[1:])
	return string(utf16.Decode(utf16Bytes))
}

func setV_WString(value string, f *outgoingField) {
	utf16Bytes := utf16.Encode([]rune(value))
	if length := len(utf16Bytes); length > f.Size {
		utf16Bytes = utf16Bytes[:f.Size]
	}
	bytes := utf16ToBytes(utf16Bytes)
	requiredLen := len(bytes) + 1
	if requiredLen > cap(f.CurrentValue) {
		f.CurrentValue = make([]byte, requiredLen)
	}
	copy(f.CurrentValue[1:], bytes)
	f.CurrentValue = f.CurrentValue[:requiredLen]
}

func (f *outgoingField) SetString(value string) {
	f.stringSetter(value, f)
	f.nullSetter(0, f)
}

func (f *outgoingField) SetNullString() {
	f.nullSetter(1, f)
}

func (f *outgoingField) GetCurrentString() (string, bool) {
	if f.nullGetter(f) {
		return ``, true
	}
	return f.stringGetter(f), false
}

func getBlob(f *outgoingField) []byte {
	return f.CurrentValue[1:]
}

func setBlob(value []byte, f *outgoingField) {
	requiredLen := len(value) + 1
	if requiredLen > cap(f.CurrentValue) {
		f.CurrentValue = make([]byte, requiredLen)
	}
	copy(f.CurrentValue[1:], value)
	f.CurrentValue = f.CurrentValue[:requiredLen]
}

func (f *outgoingField) SetBlob(value []byte) {
	f.blobSetter(value, f)
	f.nullSetter(0, f)
}

func (f *outgoingField) SetNullBlob() {
	f.nullSetter(1, f)
}

func (f *outgoingField) GetCurrentBlob() ([]byte, bool) {
	if f.nullGetter(f) {
		return nil, true
	}
	return f.blobGetter(f), false
}

type OutgoingBoolField interface {
	SetBool(bool)
	SetNullBool()
	GetCurrentBool() (bool, bool)
}

type OutgoingIntField interface {
	SetInt(int)
	SetNullInt()
	GetCurrentInt() (int, bool)
}

type OutgoingFloatField interface {
	SetFloat(float64)
	SetNullFloat()
	GetCurrentFloat() (float64, bool)
}

type OutgoingDateTimeField interface {
	SetDateTime(time.Time)
	SetNullDateTime()
	GetCurrentDateTime() (time.Time, bool)
}

type OutgoingStringField interface {
	SetString(string)
	SetNullString()
	GetCurrentString() (string, bool)
}

type OutgoingBlobField interface {
	SetBlob([]byte)
	SetNullBlob()
	GetCurrentBlob() ([]byte, bool)
}

func NewOutgoingRecordInfo(fields []NewOutgoingField) (*OutgoingRecordInfo, []string) {
	info := &OutgoingRecordInfo{
		outgoingFields: nil,
		BlobFields:     make(map[string]OutgoingBlobField),
		BoolFields:     make(map[string]OutgoingBoolField),
		DateTimeFields: make(map[string]OutgoingDateTimeField),
		FloatFields:    make(map[string]OutgoingFloatField),
		IntFields:      make(map[string]OutgoingIntField),
		StringFields:   make(map[string]OutgoingStringField),
	}
	var fieldNames []string

	for _, createField := range fields {
		field := createField()
		name := checkName(info, field.Name)
		field.Name = name
		info.outgoingFields = append(info.outgoingFields, field)
		fieldNames = append(fieldNames, name)
		switch field.Type {
		case `Bool`:
			info.BoolFields[field.Name] = field
		case `Byte`, `Int16`, `Int32`, `Int64`:
			info.IntFields[field.Name] = field
		case `Float`, `Double`, `FixedDecimal`:
			info.FloatFields[field.Name] = field
		case `Date`, `DateTime`:
			info.DateTimeFields[field.Name] = field
		case `String`, `WString`, `V_String`, `V_WString`:
			info.StringFields[field.Name] = field
		case `Blob`, `SpatialObj`:
			info.BlobFields[field.Name] = field
		}
	}
	return info, fieldNames
}

type OutgoingRecordInfo struct {
	outgoingFields []*outgoingField
	BlobFields     map[string]OutgoingBlobField
	BoolFields     map[string]OutgoingBoolField
	DateTimeFields map[string]OutgoingDateTimeField
	FloatFields    map[string]OutgoingFloatField
	IntFields      map[string]OutgoingIntField
	StringFields   map[string]OutgoingStringField
}

func (i *OutgoingRecordInfo) FixedSize() int {
	fixedSize := 0
	for _, field := range i.outgoingFields {
		if field.isFixedLen {
			fixedSize += field.dataSize()
			continue
		}
		fixedSize += 4
	}
	return fixedSize
}

func (i *OutgoingRecordInfo) HasVarFields() bool {
	for _, field := range i.outgoingFields {
		if field.isFixedLen {
			continue
		}
		return true
	}
	return false
}

func (i *OutgoingRecordInfo) DataSize() int {
	totalSize := 0
	varFields := 0
	for _, field := range i.outgoingFields {
		if field.isFixedLen {
			totalSize += field.dataSize()
			continue
		}
		varFields++
		size := field.dataSize() - 1
		if field.CurrentValue[0] == 1 || size < 4 {
			totalSize += 4 // everything fits into the fixed portion of record
			continue
		}
		if size < 127 {
			totalSize += 5 + size // 4 bytes in fixed portion of record and 1 byte for len
			continue
		}
		totalSize += 8 + size // 4 bytes in fixed portion of record and 4 bytes for len
	}
	if varFields > 0 {
		totalSize += 4 // 4 byte integer for the length of the variable portion of record
	}
	return totalSize
}

func (i *OutgoingRecordInfo) toXml(connName string) string {
	xmlBytes, _ := xml.Marshal(i.outgoingFields)
	return fmt.Sprintf(`<MetaInfo connection="%v"><RecordInfo>%v</RecordInfo></MetaInfo>`, connName, string(xmlBytes))
}

func checkName(info *OutgoingRecordInfo, name string) string {
	for _, field := range info.outgoingFields {
		if name == field.Name {
			name = fmt.Sprintf(`%v2`, name)
			return checkName(info, name)
		}
	}
	return name
}
