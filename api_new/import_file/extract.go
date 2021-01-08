package import_file

import (
	"encoding/base64"
	"fmt"
	b "github.com/tlarsen7572/goalteryx/api_new/field_base"
	"strconv"
	"strings"
	"time"
)

const dateFormat = `2006-01-02`
const dateTimeFormat = `2006-01-02 15:04:05`

type Extractor struct {
	fields []b.FieldBase
}

func NewExtractor(fieldNameBytes []byte, fieldTypeBytes []byte) *Extractor {
	fieldNames := strings.Split(string(fieldNameBytes), "\000")
	fieldTypes := strings.Split(string(fieldTypeBytes), "\000")
	fields := make([]b.FieldBase, len(fieldNames))
	for index, fieldName := range fieldNames {
		fieldType := strings.Split(fieldTypes[index], `;`)
		size := 0
		scale := 0
		var err error

		switch fieldType[0] {
		case `String`, `WString`, `V_String`, `V_WString`, `Blob`, `SpatialObj`:
			if len(fieldType) < 2 {
				panic(fmt.Sprintf(`field %v does not have a size specifier in its field type (%v)`, fieldName, fieldType))
			}
			size, err = strconv.Atoi(fieldType[1])
			if err != nil {
				panic(err.Error())
			}
		case `FixedDecimal`:
			if len(fieldType) < 3 {
				panic(fmt.Sprintf(`field %v does not have a size or scale in its field type (%v)`, fieldName, fieldType))
			}
			size, err = strconv.Atoi(fieldType[1])
			if err != nil {
				panic(err.Error())
			}
			scale, err = strconv.Atoi(fieldType[2])
			if err != nil {
				panic(err.Error())
			}
		}

		fields[index] = b.FieldBase{
			Name:  fieldName,
			Type:  fieldType[0],
			Size:  size,
			Scale: scale,
		}
	}
	return &Extractor{fields: fields}
}

func (e *Extractor) Fields() []b.FieldBase {
	return e.fields
}

func (e *Extractor) Extract(data []byte) FileData {
	boolFields := make(map[string]interface{})
	intFields := make(map[string]interface{})
	decimalFields := make(map[string]interface{})
	stringFields := make(map[string]interface{})
	dateTimeFields := make(map[string]interface{})
	blobFields := make(map[string]interface{})

	dataStrings := strings.Split(string(data), "\000")
	for index, field := range e.fields {
		value := dataStrings[index]
		switch field.Type {
		case `Bool`:
			if value == `` {
				boolFields[field.Name] = nil
				continue
			}
			if value == `true` {
				boolFields[field.Name] = true
				continue
			}
			if value == `false` {
				boolFields[field.Name] = false
			}
		case `Byte`, `Int16`, `Int32`, `Int64`:
			if value == `` {
				intFields[field.Name] = nil
				continue
			}
			intValue, _ := strconv.Atoi(value)
			intFields[field.Name] = intValue
		case `Float`, `Double`, `FixedDecimal`:
			if value == `` {
				decimalFields[field.Name] = nil
				continue
			}
			floatValue, _ := strconv.ParseFloat(value, 64)
			decimalFields[field.Name] = floatValue
		case `String`, `WString`, `V_String`, `V_WString`:
			stringFields[field.Name] = value
		case `Date`:
			if value == `` {
				dateTimeFields[field.Name] = nil
				continue
			}
			dateValue, err := time.Parse(dateFormat, value)
			if err != nil {
				panic(err.Error())
			}
			dateTimeFields[field.Name] = dateValue
		case `DateTime`:
			if value == `` {
				dateTimeFields[field.Name] = nil
				continue
			}
			dateValue, err := time.Parse(dateTimeFormat, value)
			if err != nil {
				panic(err.Error())
			}
			dateTimeFields[field.Name] = dateValue
		case `Blob`, `SpatialObj`:
			if value == `` {
				blobFields[field.Name] = nil
				continue
			}
			blobValue, err := base64.StdEncoding.DecodeString(value)
			if err != nil {
				panic(err.Error())
			}
			blobFields[field.Name] = blobValue
		default:
			panic(fmt.Sprintf(`field %v has invalid type %v`, field.Name, field.Type))
		}
	}
	return FileData{
		BoolFields:     boolFields,
		IntFields:      intFields,
		DecimalFields:  decimalFields,
		StringFields:   stringFields,
		DateTimeFields: dateTimeFields,
		BlobFields:     blobFields,
	}
}

type FileData struct {
	BoolFields     map[string]interface{}
	IntFields      map[string]interface{}
	DecimalFields  map[string]interface{}
	StringFields   map[string]interface{}
	DateTimeFields map[string]interface{}
	BlobFields     map[string]interface{}
}
