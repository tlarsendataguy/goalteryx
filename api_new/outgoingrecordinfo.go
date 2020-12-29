package api_new

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type OutgoingRecordInfo struct {
	outgoingFields []*outgoingField
}

type outgoingField struct {
	Name            string
	Type            string
	Source          string
	Size            int
	Scale           int
	CopyFrom        BytesGetter
	CurrentValue    []byte
	intSetter       func(int, *outgoingField)
	intGetter       func(*outgoingField) int
	floatSetter     func(float64, *outgoingField)
	floatGetter     func(*outgoingField) float64
	dateTimeSetter  func(time.Time, *outgoingField)
	dateTimeGetter  func(*outgoingField) time.Time
	fixedDecimalFmt string
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
	f.CurrentValue[f.Size] = 0
}

func (f *outgoingField) SetNullInt() {
	f.CurrentValue[f.Size] = 1
}

func (f *outgoingField) GetCurrentInt() (int, bool) {
	if f.CurrentValue[f.Size] == 1 {
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
	f.CurrentValue[f.Size] = 0
}

func (f *outgoingField) SetNullFloat() {
	f.CurrentValue[f.Size] = 1
}

func (f *outgoingField) GetCurrentFloat() (float64, bool) {
	if f.CurrentValue[f.Size] == 1 {
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
	f.CurrentValue[f.Size] = 0
}

func (f *outgoingField) SetNullDateTime() {
	f.CurrentValue[f.Size] = 1
}

func (f *outgoingField) GetCurrentDateTime() (time.Time, bool) {
	if f.CurrentValue[f.Size] == 1 {
		return time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), true
	}
	return f.dateTimeGetter(f), false
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

func (i *OutgoingRecordInfo) GetBoolField(name string) (OutgoingBoolField, error) {
	return i.getField(name, []string{`Bool`}, `Bool`)
}

func (i *OutgoingRecordInfo) GetIntField(name string) (OutgoingIntField, error) {
	return i.getField(name, []string{`Byte`, `Int16`, `Int32`, `Int64`}, `Int`)
}

func (i *OutgoingRecordInfo) GetFloatField(name string) (OutgoingFloatField, error) {
	return i.getField(name, []string{`Float`, `Double`, `FixedDecimal`}, `Decimal`)
}

func (i *OutgoingRecordInfo) GetDatetimeField(name string) (OutgoingDateTimeField, error) {
	return i.getField(name, []string{`Date`, `DateTime`}, `DateTime`)
}

func (i *OutgoingRecordInfo) getField(name string, types []string, label string) (*outgoingField, error) {
	for _, field := range i.outgoingFields {
		if field.Name == name {
			for _, ofType := range types {
				if field.Type == ofType {
					return field, nil
				}
			}
			return nil, fmt.Errorf(`the '%v' field is not a %v field, it is '%v'`, name, label, field.Type)
		}
	}
	return nil, fmt.Errorf(`there is no '%v' field in the record`, name)
}
