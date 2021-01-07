package import_file_test

import (
	"bufio"
	"github.com/tlarsen7572/goalteryx/api_new/import_file"
	"os"
	"testing"
)

func TestPreprocessTextFile(t *testing.T) {
	file, err := os.Open(`..\sdk_test_passthrough_simulation.txt`)
	if err != nil {
		t.Fatalf(`expected no error but got %v`, err.Error())
	}
	lines := make([][]byte, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		value := scanner.Bytes()
		lines = append(lines, import_file.Preprocess(value))
	}
	expected := "Field1\000Field2\000Field3\000Field4\000Field5\000Field6\000Field7\000Field8\000Field9\000Field10\000Field11\000Field12\000Field13\000Field14\000Field15\000Field16"
	if value := string(lines[0]); expected != value {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, value)
	}

	expected = "true\0002\000100\0001000\00010000\00012.34\0001.23\000234.56\000ABC\000Hello \000 World\000abcdefg\0002020-01-01\0002020-01-02 03:04:05\000\000"
	if value := string(lines[2]); expected != value {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, value)
	}

	expected = "false\000-2\000-100\000-1000\000-10000\000-12.34\000-1.23\000-234.56\000DE|\"FG\000HIJK\000LMNOP\000QRSTU\r\nVWXYZ\0002020-02-03\0002020-01-02 13:14:15\000\000"
	if value := string(lines[3]); expected != value {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, value)
	}
}
