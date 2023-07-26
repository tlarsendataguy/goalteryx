package import_file_test

import (
	"bufio"
	"github.com/tlarsendataguy/goalteryx/sdk/field_base"
	"github.com/tlarsendataguy/goalteryx/sdk/import_file"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

const fieldNames string = "Field1\000Field2\000Field3\000Field4\000Field5\000Field6\000Field7\000Field8\000Field9\000Field10\000Field11\000Field12\000Field13\000Field14\000Field15\000Field16\000Field17"
const fieldTypes string = "Bool\000Byte\000Int16\000Int32\000Int64\000Float\000Double\000FixedDecimal;19;2\000String;100\000WString;100\000V_String;10000\000V_WString;100000\000Date\000DateTime\000Blob;10\000SpatialObj;100\000Time"
const record1 string = "true\0002\000100\0001000\00010000\00012.34\0001.23\000234.56\000\"ABC\"\000\"Hello \"\000\" World\"\000\"abcdefg\"\0002020-01-01\0002020-01-02 03:04:05\000\000\00010:01:01"
const record2 string = "false\0002\000-100\000-1000\000-10000\000-12.34\000-1.23\000-234.56\000\"DE|\"FG\"\000HIJK\000LMNOP\000\"QRSTU\r\nVWXYZ\"\0002020-02-03\0002020-01-02 13:14:15\000\000\000"
const record3 string = "\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\00017:02:01"
const record4 string = "true\00042\000-110\000392\0002340\00012\00041.22\00098.2\000\"\"\000\"\"\000\"\"\000\"\"\0002020-02-13\0002020-11-02 13:14:15\000\000\000"

func TestPreprocessTextFile(t *testing.T) {
	file, err := os.Open(filepath.Join(`..`, `sdk_test_passthrough_simulation.txt`))
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	lines := make([][]byte, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		value := scanner.Bytes()
		lines = append(lines, import_file.Preprocess(value))
	}

	if value := string(lines[0]); fieldNames != value {
		t.Fatalf("expected\n%v\nbut got\n%v", fieldNames, value)
	}

	if value := string(lines[1]); fieldTypes != value {
		t.Fatalf("expected\n%v\nbut got\n%v", fieldTypes, value)
	}

	if value := string(lines[2]); record1 != value {
		t.Fatalf("expected\n%v\nbut got\n%v", record1, value)
	}

	if value := string(lines[3]); record2 != value {
		t.Fatalf("expected\n%v\nbut got\n%v", record2, value)
	}

	if value := string(lines[4]); record3 != value {
		t.Fatalf("expected\n%v\nbut got\n%v", record3, value)
	}
}

func TestExtractData(t *testing.T) {
	extractor := import_file.NewExtractor([]byte(fieldNames), []byte(fieldTypes))
	fields := extractor.Fields()
	expectedFields := []field_base.FieldBase{
		{Name: `Field1`, Type: `Bool`},
		{Name: `Field2`, Type: `Byte`},
		{Name: `Field3`, Type: `Int16`},
		{Name: `Field4`, Type: `Int32`},
		{Name: `Field5`, Type: `Int64`},
		{Name: `Field6`, Type: `Float`},
		{Name: `Field7`, Type: `Double`},
		{Name: `Field8`, Type: `FixedDecimal`, Size: 19, Scale: 2},
		{Name: `Field9`, Type: `String`, Size: 100},
		{Name: `Field10`, Type: `WString`, Size: 100},
		{Name: `Field11`, Type: `V_String`, Size: 10000},
		{Name: `Field12`, Type: `V_WString`, Size: 100000},
		{Name: `Field13`, Type: `Date`},
		{Name: `Field14`, Type: `DateTime`},
		{Name: `Field15`, Type: `Blob`, Size: 10},
		{Name: `Field16`, Type: `SpatialObj`, Size: 100},
		{Name: `Field17`, Type: `Time`},
	}
	if !reflect.DeepEqual(fields, expectedFields) {
		t.Fatalf("expected\n%v\nbut got\n%v", expectedFields, fields)
	}

	data := extractor.Extract([]byte(record1))
	if value := data.BoolFields[`Field1`]; value != true {
		t.Fatalf(`expected true but got %v`, value)
	}
	if value := data.IntFields[`Field2`]; value != 2 {
		t.Fatalf(`expected 2 but got %v`, value)
	}
	if value := data.IntFields[`Field3`]; value != 100 {
		t.Fatalf(`expected 100 but got %v`, value)
	}
	if value := data.IntFields[`Field4`]; value != 1000 {
		t.Fatalf(`expected 1000 but got %v`, value)
	}
	if value := data.IntFields[`Field5`]; value != 10000 {
		t.Fatalf(`expected 10000 but got %v`, value)
	}
	if value := data.DecimalFields[`Field6`]; value != 12.34 {
		t.Fatalf(`expected 12.34 but got %v`, value)
	}
	if value := data.DecimalFields[`Field7`]; value != 1.23 {
		t.Fatalf(`expected 1.23 but got %v`, value)
	}
	if value := data.DecimalFields[`Field8`]; value != 234.56 {
		t.Fatalf(`expected 234.56 but got %v`, value)
	}
	if value := data.StringFields[`Field9`]; value != `ABC` {
		t.Fatalf(`expected 'ABC' but got '%v'`, value)
	}
	if value := data.StringFields[`Field10`]; value != `Hello ` {
		t.Fatalf(`expected 'Hello ' but got '%v'`, value)
	}
	if value := data.StringFields[`Field11`]; value != ` World` {
		t.Fatalf(`expected ' World' but got '%v'`, value)
	}
	if value := data.StringFields[`Field12`]; value != `abcdefg` {
		t.Fatalf(`expected 'abcdefg' but got '%v'`, value)
	}
	if value := data.DateTimeFields[`Field13`]; value != time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC) {
		t.Fatalf(`expected 2020-01-01 but got '%v'`, value)
	}
	if value := data.DateTimeFields[`Field14`]; value != time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC) {
		t.Fatalf(`expected 2020-01-02 03:04:05 but got '%v'`, value)
	}
	if value := data.BlobFields[`Field15`]; value != nil {
		t.Fatalf(`expected nil but got '%v'`, value)
	}

	data = extractor.Extract([]byte(record3))
	if value := data.BoolFields[`Field1`]; value != nil {
		t.Fatalf(`expected nil but got %v`, value)
	}
	if value := data.IntFields[`Field2`]; value != nil {
		t.Fatalf(`expected nil but got %v`, value)
	}
	if value := data.IntFields[`Field3`]; value != nil {
		t.Fatalf(`expected nil but got %v`, value)
	}
	if value := data.IntFields[`Field4`]; value != nil {
		t.Fatalf(`expected nil but got %v`, value)
	}
	if value := data.IntFields[`Field5`]; value != nil {
		t.Fatalf(`expected nil but got %v`, value)
	}
	if value := data.DecimalFields[`Field6`]; value != nil {
		t.Fatalf(`expected nil but got %v`, value)
	}
	if value := data.DecimalFields[`Field7`]; value != nil {
		t.Fatalf(`expected nil but got %v`, value)
	}
	if value := data.DecimalFields[`Field8`]; value != nil {
		t.Fatalf(`expected nil but got %v`, value)
	}
	if value := data.StringFields[`Field9`]; value != nil {
		t.Fatalf(`expected '' but got %v`, value)
	}
	if value := data.StringFields[`Field10`]; value != nil {
		t.Fatalf(`expected '' but got %v`, value)
	}
	if value := data.StringFields[`Field11`]; value != nil {
		t.Fatalf(`expected '' but got %v`, value)
	}
	if value := data.StringFields[`Field12`]; value != nil {
		t.Fatalf(`expected '' but got %v`, value)
	}
	if value := data.DateTimeFields[`Field13`]; value != nil {
		t.Fatalf(`expected nil but got %v`, value)
	}
	if value := data.DateTimeFields[`Field14`]; value != nil {
		t.Fatalf(`expected nil but got %v`, value)
	}
	if value := data.BlobFields[`Field15`]; value != nil {
		t.Fatalf(`expected nil but got %v`, value)
	}

	data = extractor.Extract([]byte(record4))
	if value := data.StringFields[`Field9`]; value != `` {
		t.Fatalf(`expected '' but got %v`, value)
	}
}

func TestPanicWhenFieldNameAndTypeDoNotEqual(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		expectedErr := `the number of field names and types did not match; got 2 field names but 1 field types`
		if expectedErr != errText {
			t.Fatalf("expected error message\n%v\nbut got\n%v", expectedErr, errText)
		}
		t.Logf(errText)
	}(t)

	import_file.NewExtractor([]byte("Field1\000Field2"), []byte(`Byte`))
}

func TestPanicOnInvalidFieldType(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		expectedErr := `'Invalid' is not a valid field type`
		if expectedErr != errText {
			t.Fatalf("expected error message\n%v\nbut got\n%v", expectedErr, errText)
		}
		t.Logf(errText)
	}(t)

	import_file.NewExtractor([]byte("Field1"), []byte("Invalid"))
}

func TestPanicWhenSizeIsMissing(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		t.Logf(errText)
	}(t)

	import_file.NewExtractor([]byte("Field1"), []byte("String"))
}

func TestPanicWhenSizeIsNotNumeric(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		t.Logf(errText)
	}(t)

	import_file.NewExtractor([]byte("Field1"), []byte("String;notNumeric"))
}

func TestPanicWhenScaleIsMissing(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		t.Logf(errText)
	}(t)

	import_file.NewExtractor([]byte("Field1"), []byte("FixedDecimal;10"))
}

func TestPanicWhenScaleIsNotNumeric(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		t.Logf(errText)
	}(t)

	import_file.NewExtractor([]byte("Field1"), []byte("FixedDecimal;10;NotNumeric"))
}

func TestPanicWhenDataHasIncorrectNumberOfFields(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		t.Logf(errText)
	}(t)

	extractor := import_file.NewExtractor([]byte("Field1\000Field2"), []byte("Bool\000Byte"))
	extractor.Extract([]byte("true"))
}

func TestPanicWhenBoolHasInvalidValue(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		t.Logf(errText)
	}(t)

	extractor := import_file.NewExtractor([]byte("Field1"), []byte("Bool"))
	extractor.Extract([]byte("not a bool"))
}

func TestPanicWhenBoolHasQuotedValue(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		t.Logf(errText)
	}(t)

	extractor := import_file.NewExtractor([]byte("Field1"), []byte("Bool"))
	extractor.Extract([]byte(`"true"`))
}

func TestPanicWhenIntHasInvalidValue(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		t.Logf(errText)
	}(t)

	extractor := import_file.NewExtractor([]byte("Field1"), []byte("Byte"))
	extractor.Extract([]byte("not an integer"))
}

func TestPanicWhenIntHasQuotedValue(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		t.Logf(errText)
	}(t)

	extractor := import_file.NewExtractor([]byte("Field1"), []byte("Int16"))
	extractor.Extract([]byte(`"123"`))
}

func TestPanicWhenDecimalHasInvalidValue(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		t.Logf(errText)
	}(t)

	extractor := import_file.NewExtractor([]byte("Field1"), []byte("Double"))
	extractor.Extract([]byte("not a decimal"))
}

func TestPanicWhenDecimalHasQuotedValue(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		t.Logf(errText)
	}(t)

	extractor := import_file.NewExtractor([]byte("Field1"), []byte("Double"))
	extractor.Extract([]byte(`"123.45"`))
}

func TestPanicWhenDateHasInvalidValue(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		t.Logf(errText)
	}(t)

	extractor := import_file.NewExtractor([]byte("Field1"), []byte("Date"))
	extractor.Extract([]byte("not a date"))
}

func TestPanicWhenDateTimeHasInvalidValue(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		t.Logf(errText)
	}(t)

	extractor := import_file.NewExtractor([]byte("Field1"), []byte("DateTime"))
	extractor.Extract([]byte("not a datetime"))
}

func TestPanicWhenDateHasQuotedValue(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		t.Logf(errText)
	}(t)

	extractor := import_file.NewExtractor([]byte("Field1"), []byte("Date"))
	extractor.Extract([]byte(`"2020-01-01"`))
}

func TestPanicWhenBlobHasInvalidValue(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		t.Logf(errText)
	}(t)

	extractor := import_file.NewExtractor([]byte("Field1"), []byte("Blob;100"))
	extractor.Extract([]byte("not a base64 string"))
}

func TestPanicWhenBlobHasQuotedValue(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		if err == nil {
			t.Fatalf(`no panic happened but one was expected`)
		}
		errText := err.(string)
		t.Logf(errText)
	}(t)

	extractor := import_file.NewExtractor([]byte("Field1"), []byte("Blob;100"))
	extractor.Extract([]byte(`""`))
}
