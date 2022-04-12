package sdk_test

import (
	"bytes"
	"github.com/tlarsendataguy/goalteryx/sdk"
	"testing"
	"unsafe"
)

func TestRecordPacket(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6}
	collectedData := make([]byte, 0)
	recordPacket := sdk.NewRecordPacket(sdk.RecordCache(&data[0]), len(data), 1, false)
	rows := 0
	for recordPacket.Next() {
		collectedData = append(collectedData, *(*byte)(recordPacket.Record()))
		rows++
	}
	if rows != 6 {
		t.Fatalf(`expected 6 rows but got %v`, rows)
	}
	if !bytes.Equal(data, collectedData) {
		t.Fatalf(`expected %v but got %v`, data, collectedData)
	}
}

func TestVarDataRecordPacket(t *testing.T) {
	data := []byte{1, 0, 0, 0, 0, 2, 4, 0, 0, 0, 10, 0, 0, 0, 3, 0, 0, 0}
	collectedFixedData := make([]byte, 0)
	collectedVarLens := make([]int32, 0)
	recordPacket := sdk.NewRecordPacket(sdk.RecordCache(&data[0]), len(data), 1, true)
	rows := 0
	for recordPacket.Next() {
		collectedFixedData = append(collectedFixedData, *(*byte)(recordPacket.Record()))
		collectedVarLens = append(collectedVarLens, *(*int32)(unsafe.Pointer(uintptr(recordPacket.Record()) + 1)))
		rows++
	}
	if rows != 3 {
		t.Fatalf(`expected 3 rows but got %v`, rows)
	}
	if !bytes.Equal([]byte{1, 2, 3}, collectedFixedData) {
		t.Fatalf(`expected [1 2 3] but got %v`, collectedFixedData)
	}
	if collectedVarLens[0] != 0 || collectedVarLens[1] != 4 || collectedVarLens[2] != 0 {
		t.Fatalf(`expected [0 4 0] but got %v`, collectedFixedData)
	}
}
