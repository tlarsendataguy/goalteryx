package api_new

import (
	"bytes"
	"testing"
	"unsafe"
)

func TestExtractFixedBytes(t *testing.T) {
	getFixed := generateGetFixedBytes(2, 4)
	data := unsafe.Pointer(&[]byte{12, 6, 24, 122, 86, 4, 200, 0, 73, 15, 3}[0])
	result := getFixed(data)
	if expected := []byte{24, 122, 86, 4}; !bytes.Equal(result, expected) {
		t.Fatalf(`expected '%v' but got '%v'`, expected, result)
	}
}

func TestExtractVarBytesNullValue(t *testing.T) {
	getVar := generateGetVarBytes(3)
	data := unsafe.Pointer(&[]byte{0, 122, 65, 1, 0, 0, 0, 3, 0, 0, 0}[0])
	result := getVar(data)
	if result != nil {
		t.Fatalf(`expected nil but got %v`, result)
	}
}

func TestExtractVarBytesEmptyValue(t *testing.T) {
	getVar := generateGetVarBytes(3)
	data := unsafe.Pointer(&[]byte{0, 122, 65, 0, 0, 0, 0, 3, 0, 0, 0}[0])
	result := getVar(data)
	if !bytes.Equal(result, []byte{}) {
		t.Fatalf(`expected empty slice but got %v`, result)
	}
}

func TestExtractVarBytesTinyValue(t *testing.T) {
	// byte with 1, v_wstring with 'A', v_string with 'B'
	var tinyValue = unsafe.Pointer(&[]byte{1, 0, 65, 0, 0, 32, 66, 0, 0, 16, 0, 0, 0, 0}[0])

	getWVar := generateGetVarBytes(2)
	result := getWVar(tinyValue)
	if expected := []byte{65, 0}; !bytes.Equal(result, expected) {
		t.Fatalf(`expected '%v' but got '%v'`, expected, result)
	}

	getVar := generateGetVarBytes(6)
	result = getVar(tinyValue)
	if expected := []byte{66}; !bytes.Equal(result, expected) {
		t.Fatalf(`expected '%v' but got '%v'`, expected, result)
	}
}

func TestExtractVarBytesShortValue(t *testing.T) {
	// byte with 1, v_wstring with 50 A's, v_string with 100 B's
	shortValue := unsafe.Pointer(&[]byte{
		1, 0, 12, 0, 0, 0, 109, 0, 0, 0, 202, 0, 0, 0, 201, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 201, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	}[0])

	getWVar := generateGetVarBytes(2)
	result := getWVar(shortValue)
	expected := []byte{65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0}
	if !bytes.Equal(expected, result) {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, result)
	}
}

func TestExtractVaryBytesLongValue(t *testing.T) {
	// byte with 1, v_wstring with 100 A's, v_string with 200 B's
	longValue := unsafe.Pointer(&[]byte{
		1, 0, 12, 0, 0, 0, 212, 0, 0, 0, 152, 1, 0, 0, 144, 1, 0, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 144, 1, 0, 0, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66,
	}[0])

	getWVar := generateGetVarBytes(2)
	result := getWVar(longValue)
	expected := []byte{65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0, 65, 0}
	if !bytes.Equal(expected, result) {
		t.Fatalf("expected\n%v\nbut got\n%v", expected, result)
	}
}

func TestArrayAndSlice(t *testing.T) {
	value := make([]int, 10)
	for i := 0; i < 10; i++ {
		value[i] = i
	}
	t.Logf(`array size %v and cap %v: %v`, len(value), cap(value), value)
	value = value[0:5]
	t.Logf(`array size %v and cap %v: %v`, len(value), cap(value), value)
	value = value[0:10]
	t.Logf(`array size %v and cap %v: %v`, len(value), cap(value), value)
}
