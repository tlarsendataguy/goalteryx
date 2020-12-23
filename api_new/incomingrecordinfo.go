package api_new

import (
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
	"unicode/utf16"
	"unsafe"
)

type IntGetter func(Record) (int, bool)
type FloatGetter func(Record) (float64, bool)
type BoolGetter func(Record) (bool, bool)
type TimeGetter func(Record) (time.Time, bool)
type StringGetter func(Record) (string, bool)

const dateFormat = `2006-01-02`
const dateTimeFormat = `2006-01-02 15:04:05`

type xmlMetaInfo struct {
	Connection string        `xml:"connection,attr"`
	RecordInfo xmlRecordInfo `xml:"RecordInfo"`
}

type xmlRecordInfo struct {
	Fields []IncomingField `xml:"Field"`
}

type IncomingField struct {
	Name     string `xml:"name,attr"`
	Type     string `xml:"type,attr"`
	Source   string `xml:"source,attr"`
	Size     int    `xml:"size,attr"`
	Scale    int    `xml:"scale,attr"`
	GetBytes BytesGetter
}

type IncomingIntField struct {
	Name     string
	Type     string
	Source   string
	GetValue IntGetter
}

type IncomingFloatField struct {
	Name     string
	Type     string
	Source   string
	GetValue FloatGetter
}

type IncomingBoolField struct {
	Name     string
	Type     string
	Source   string
	GetValue BoolGetter
}

type IncomingTimeField struct {
	Name     string
	Type     string
	Source   string
	GetValue TimeGetter
}

type IncomingBlobField struct {
	Name     string
	Type     string
	Source   string
	Size     int
	GetValue BytesGetter
}

type IncomingStringField struct {
	Name     string
	Type     string
	Source   string
	Size     int
	GetValue StringGetter
}

type IncomingRecordInfo struct {
	fields []IncomingField
}

func (i IncomingRecordInfo) NumFields() int {
	return len(i.fields)
}

func (i IncomingRecordInfo) Clone() EditingRecordInfo {
	return EditingRecordInfo{fields: i.fields}
}

func (i IncomingRecordInfo) GetIntField(name string) (IncomingIntField, error) {
	for _, field := range i.fields {
		if field.Name != name {
			continue
		}
		switch field.Type {
		case `Byte`:
			return generateIncomingIntField(field, bytesToByte), nil
		case `Int16`:
			return generateIncomingIntField(field, bytesToInt16), nil
		case `Int32`:
			return generateIncomingIntField(field, bytesToInt32), nil
		case `Int64`:
			return generateIncomingIntField(field, bytesToInt64), nil
		default:
			return IncomingIntField{}, fmt.Errorf(`the '%v' field is not an integer field, it is '%v'`, name, field.Type)
		}
	}
	return IncomingIntField{}, fmt.Errorf(`there is no '%v' field in the record`, name)
}

func (i IncomingRecordInfo) GetFloatField(name string) (IncomingFloatField, error) {
	for _, field := range i.fields {
		if field.Name != name {
			continue
		}
		switch field.Type {
		case `Float`:
			return generateIncomingFloatField(field, bytesToFloat), nil
		case `Double`:
			return generateIncomingFloatField(field, bytesToDouble), nil
		case `FixedDecimal`:
			return generateFixedDecimalField(field), nil
		default:
			return IncomingFloatField{}, fmt.Errorf(`the '%v' field is not a float field, it is '%v'`, name, field.Type)
		}
	}
	return IncomingFloatField{}, fmt.Errorf(`there is no '%v' field in the record`, name)
}

func (i IncomingRecordInfo) GetBoolField(name string) (IncomingBoolField, error) {
	for _, field := range i.fields {
		if field.Name != name {
			continue
		}
		switch field.Type {
		case `Bool`:
			return generateBoolField(field), nil
		default:
			return IncomingBoolField{}, fmt.Errorf(`the '%v' field is not a bool field, it is '%v'`, name, field.Type)
		}
	}
	return IncomingBoolField{}, fmt.Errorf(`there is no '%v' field in the record`, name)
}

func (i IncomingRecordInfo) GetTimeField(name string) (IncomingTimeField, error) {
	for _, field := range i.fields {
		if field.Name != name {
			continue
		}
		switch field.Type {
		case `Date`:
			return generateTimeField(field, dateFormat, 10), nil
		case `DateTime`:
			return generateTimeField(field, dateTimeFormat, 19), nil
		default:
			return IncomingTimeField{}, fmt.Errorf(`the '%v' field is not a time field, it is '%v'`, name, field.Type)
		}
	}
	return IncomingTimeField{}, fmt.Errorf(`there is no '%v' field in the record`, name)
}

func (i IncomingRecordInfo) GetBlobField(name string) (IncomingBlobField, error) {
	for _, field := range i.fields {
		if field.Name != name {
			continue
		}
		switch field.Type {
		case `Blob`, `SpatialObj`:
			return generateBlobField(field), nil
		default:
			return IncomingBlobField{}, fmt.Errorf(`the '%v' field is not a blob field, it is '%v'`, name, field.Type)
		}
	}
	return IncomingBlobField{}, fmt.Errorf(`there is no '%v' field in the record`, name)
}

func (i IncomingRecordInfo) GetStringField(name string) (IncomingStringField, error) {
	for _, field := range i.fields {
		if field.Name != name {
			continue
		}
		switch field.Type {
		case `String`:
			return generateIncomingStringField(field, bytesToString), nil
		case `WString`:
			return generateIncomingStringField(field, bytesToWString), nil
		case `V_String`:
			return generateIncomingStringField(field, bytesToV_String), nil
		case `V_WString`:
			return generateIncomingStringField(field, bytesToV_WString), nil
		default:
			return IncomingStringField{}, fmt.Errorf(`the '%v' field is not a string field, it is '%v'`, name, field.Type)
		}
	}
	return IncomingStringField{}, fmt.Errorf(`there is no '%v' field in the record`, name)
}

func incomingRecordInfoFromString(config string) (IncomingRecordInfo, error) {
	if config[:9] != `<MetaInfo` {
		if config[:11] != `<RecordInfo` {
			return IncomingRecordInfo{}, errors.New(`config is not a valid IncomingRecordInfo xml string`)
		}
		config = `<MetaInfo>` + config + `</MetaInfo>`
	}
	metaInfo := xmlMetaInfo{}
	err := xml.Unmarshal([]byte(config), &metaInfo)
	if err != nil {
		return IncomingRecordInfo{}, err
	}
	startAt := 0
	for index, field := range metaInfo.RecordInfo.Fields {
		switch field.Type {
		case `V_String`, `V_WString`, `Blob`, `SpatialObj`:
			field.GetBytes = generateGetVarBytes(startAt)
			startAt += 4
		case `Bool`:
			field.GetBytes = generateGetFixedBytes(startAt, 1)
			startAt += 1
		case `Byte`:
			field.GetBytes = generateGetFixedBytes(startAt, 2)
			startAt += 2
		case `Int16`:
			field.GetBytes = generateGetFixedBytes(startAt, 3)
			startAt += 3
		case `Int32`, `Float`:
			field.GetBytes = generateGetFixedBytes(startAt, 5)
			startAt += 5
		case `Int64`, `Double`:
			field.GetBytes = generateGetFixedBytes(startAt, 9)
			startAt += 9
		case `String`, `FixedDecimal`:
			field.GetBytes = generateGetFixedBytes(startAt, field.Size+1)
			startAt += field.Size + 1
		case `WString`:
			size := field.Size * 2
			field.GetBytes = generateGetFixedBytes(startAt, size+1)
			startAt += size + 1
		case `Date`:
			field.GetBytes = generateGetFixedBytes(startAt, 11)
			startAt += 11
		case `DateTime`:
			field.GetBytes = generateGetFixedBytes(startAt, 20)
			startAt += 20
		default:
			return IncomingRecordInfo{}, fmt.Errorf(`field '%v' has invalid field type '%v'`, field.Name, field.Type)
		}
		metaInfo.RecordInfo.Fields[index] = field
	}
	return IncomingRecordInfo{fields: metaInfo.RecordInfo.Fields}, nil
}

func generateIncomingIntField(field IncomingField, getter func(BytesGetter) IntGetter) IncomingIntField {
	return IncomingIntField{
		Name:     field.Name,
		Type:     field.Type,
		Source:   field.Source,
		GetValue: getter(field.GetBytes),
	}
}

func bytesToByte(getBytes BytesGetter) IntGetter {
	return func(record Record) (int, bool) {
		bytes := getBytes(record)
		if bytes[1] == 1 {
			return 0, true
		}
		return int(bytes[0]), false
	}
}

func bytesToInt16(getBytes BytesGetter) IntGetter {
	return func(record Record) (int, bool) {
		bytes := getBytes(record)
		if bytes[2] == 1 {
			return 0, true
		}
		return int(binary.LittleEndian.Uint16(bytes)), false
	}
}

func bytesToInt32(getBytes BytesGetter) IntGetter {
	return func(record Record) (int, bool) {
		bytes := getBytes(record)
		if bytes[4] == 1 {
			return 0, true
		}
		return int(binary.LittleEndian.Uint32(bytes)), false
	}
}

func bytesToInt64(getBytes BytesGetter) IntGetter {
	return func(record Record) (int, bool) {
		bytes := getBytes(record)
		if bytes[8] == 1 {
			return 0, true
		}
		return int(binary.LittleEndian.Uint64(bytes)), false
	}
}

func generateIncomingFloatField(field IncomingField, getter func(BytesGetter) FloatGetter) IncomingFloatField {
	return IncomingFloatField{
		Name:     field.Name,
		Type:     field.Type,
		Source:   field.Source,
		GetValue: getter(field.GetBytes),
	}
}

func bytesToFloat(getBytes BytesGetter) FloatGetter {
	return func(record Record) (float64, bool) {
		bytes := getBytes(record)
		if bytes[4] == 1 {
			return 0, true
		}
		return float64(math.Float32frombits(binary.LittleEndian.Uint32(bytes))), false
	}
}

func bytesToDouble(getBytes BytesGetter) FloatGetter {
	return func(record Record) (float64, bool) {
		bytes := getBytes(record)
		if bytes[8] == 1 {
			return 0, true
		}
		return math.Float64frombits(binary.LittleEndian.Uint64(bytes)), false
	}
}

func truncateAtNullByte(raw []byte) []byte {
	var dataLen int
	for dataLen = 0; dataLen < len(raw); dataLen++ {
		if raw[dataLen] == 0 {
			break
		}
	}
	return raw[:dataLen]
}

func generateFixedDecimalField(field IncomingField) IncomingFloatField {
	getter := func(record Record) (float64, bool) {
		bytes := field.GetBytes(record)
		if bytes[field.Size] == 1 {
			return 0, true
		}
		valueStr := string(truncateAtNullByte(bytes))
		value, _ := strconv.ParseFloat(valueStr, 64)
		return value, false
	}
	return IncomingFloatField{
		Name:     field.Name,
		Type:     field.Type,
		Source:   field.Source,
		GetValue: getter,
	}
}

func generateBoolField(field IncomingField) IncomingBoolField {
	getter := func(record Record) (bool, bool) {
		bytes := field.GetBytes(record)
		if bytes[0] == 2 {
			return false, true
		}
		return bytes[0] == 1, false
	}
	return IncomingBoolField{
		Name:     field.Name,
		Type:     field.Type,
		Source:   field.Source,
		GetValue: getter,
	}
}

func generateTimeField(field IncomingField, format string, size int) IncomingTimeField {
	getter := func(record Record) (time.Time, bool) {
		bytes := field.GetBytes(record)
		if bytes[size] == 1 {
			return time.Time{}, true
		}
		value, _ := time.Parse(format, string(bytes[0:size]))
		return value, false
	}
	return IncomingTimeField{
		Name:     field.Name,
		Type:     field.Type,
		Source:   field.Source,
		GetValue: getter,
	}
}

func generateBlobField(field IncomingField) IncomingBlobField {
	return IncomingBlobField{
		Name:     field.Name,
		Type:     field.Type,
		Source:   field.Source,
		Size:     field.Size,
		GetValue: field.GetBytes,
	}
}

func generateIncomingStringField(field IncomingField, getter func(BytesGetter, int) StringGetter) IncomingStringField {
	return IncomingStringField{
		Name:     field.Name,
		Type:     field.Type,
		Source:   field.Source,
		GetValue: getter(field.GetBytes, field.Size),
	}
}

func bytesToString(getBytes BytesGetter, size int) StringGetter {
	return func(record Record) (string, bool) {
		bytes := getBytes(record)
		if bytes[size] == 1 {
			return ``, true
		}
		return string(truncateAtNullByte(bytes)), false
	}
}

func truncateAtNullUtf16(raw []uint16) []uint16 {
	var dataLen int
	for dataLen = 0; dataLen < len(raw); dataLen++ {
		if raw[dataLen] == 0 {
			break
		}
	}
	return raw[:dataLen]
}

func bytesToWString(getBytes BytesGetter, size int) StringGetter {
	return func(record Record) (string, bool) {
		bytes := getBytes(record)
		if bytes[size*2] == 1 {
			return ``, true
		}
		var utf16Bytes []uint16
		rawHeader := (*reflect.SliceHeader)(unsafe.Pointer(&utf16Bytes))
		rawHeader.Data = uintptr(unsafe.Pointer(&bytes[0]))
		rawHeader.Len = size
		rawHeader.Cap = size
		utf16Bytes = truncateAtNullUtf16(utf16Bytes)
		value := string(utf16.Decode(utf16Bytes))
		return value, false
	}
}

func bytesToV_String(getBytes BytesGetter, _ int) StringGetter {
	return func(record Record) (string, bool) {
		bytes := getBytes(record)
		if bytes == nil {
			return ``, true
		}
		return string(bytes), false
	}
}

func bytesToV_WString(getBytes BytesGetter, size int) StringGetter {
	return func(record Record) (string, bool) {
		bytes := getBytes(record)
		if bytes == nil {
			return ``, true
		}
		var utf16Bytes []uint16
		rawHeader := (*reflect.SliceHeader)(unsafe.Pointer(&utf16Bytes))
		rawHeader.Data = uintptr(unsafe.Pointer(&bytes[0]))
		rawHeader.Len = len(bytes) / 2
		rawHeader.Cap = len(bytes) / 2
		value := string(utf16.Decode(utf16Bytes))
		return value, false
	}
}
