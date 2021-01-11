package api_new

import (
	"strconv"
	"time"
)

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

func generateIncomingIntField(field IncomingField, getter func(BytesGetter) IntGetter) IncomingIntField {
	return IncomingIntField{
		Name:     field.Name,
		Type:     field.Type,
		Source:   field.Source,
		GetValue: getter(field.GetBytes),
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
