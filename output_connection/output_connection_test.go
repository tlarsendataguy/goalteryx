package output_connection_test

import (
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/output_connection"
	"github.com/tlarsen7572/goalteryx/recordinfo"
	"testing"
	"unsafe"
)

/*
NOTE: Running output_connection_test as a Directory test results in failures.  I don't know why.  Run
as a File test for correct results.
*/

func TestPassAndFailInit(t *testing.T) {
	iiInitOk := &IiTestStruct{InitReturnValue: true, PushRecordReturnValue: true}
	iiInitFail := &IiTestStruct{InitReturnValue: false, PushRecordReturnValue: true}
	connection := output_connection.New(1, `Test`)
	connection.Add(api.NewConnectionInterfaceStruct(iiInitOk))
	connection.Add(api.NewConnectionInterfaceStruct(iiInitFail))

	generator := recordinfo.NewGenerator()
	generator.AddByteField(`SomeByte`, ``)
	info := generator.GenerateRecordInfo()
	err := connection.Init(info)
	if err == nil {
		t.Fatalf(`expected error but got none`)
	}
	t.Logf(err.Error())
	if !iiInitOk.IsInitialized {
		t.Fatalf(`iiInitOk did not initialize`)
	}
	record, _ := info.GenerateRecord()

	connection.PushRecord(record)
	connection.Close()

	if iiInitOk.PushRecordCalls != 1 {
		t.Fatalf(`expected 1 push record call for iiInitOk but got %v`, iiInitOk.PushRecordCalls)
	}
	if iiInitFail.PushRecordCalls != 0 {
		t.Fatalf(`expected 0 push record calls for iiInitFail but got %v`, iiInitFail.PushRecordCalls)
	}
	if iiInitOk.CloseCalls != 1 {
		t.Fatalf(`expected 1 close call for iiInitOk but got %v`, iiInitOk.CloseCalls)
	}
	if iiInitFail.CloseCalls != 0 {
		t.Fatalf(`expected 0 close calls for iiInitFail but got %v`, iiInitFail.CloseCalls)
	}
}

func TestPassAndFailPushRecord(t *testing.T) {
	iiPushOk := &IiTestStruct{InitReturnValue: true, PushRecordReturnValue: true}
	iiPushFail := &IiTestStruct{InitReturnValue: true, PushRecordReturnValue: false}
	connection := output_connection.New(1, `Test`)
	connection.Add(api.NewConnectionInterfaceStruct(iiPushOk))
	connection.Add(api.NewConnectionInterfaceStruct(iiPushFail))

	generator := recordinfo.NewGenerator()
	generator.AddByteField(`SomeByte`, ``)
	info := generator.GenerateRecordInfo()
	err := connection.Init(info)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	record, _ := info.GenerateRecord()
	connection.PushRecord(record)
	connection.PushRecord(record)
	connection.Close()

	if iiPushOk.PushRecordCalls != 2 {
		t.Fatalf(`expected 2 push record calls but got %v`, iiPushOk.PushRecordCalls)
	}
	if iiPushFail.PushRecordCalls != 1 {
		t.Fatalf(`expected 1 push record call but got %v`, iiPushFail.PushRecordCalls)
	}
	if iiPushOk.CloseCalls != 1 {
		t.Fatalf(`expected 1 close call but got %v`, iiPushOk.CloseCalls)
	}
	if iiPushFail.CloseCalls != 1 {
		t.Fatalf(`expected 1 close call but got %v`, iiPushFail.CloseCalls)
	}
}

func TestUpdateProgress(t *testing.T) {
	iiPush1 := &IiTestStruct{InitReturnValue: true, PushRecordReturnValue: true}
	iiPush2 := &IiTestStruct{InitReturnValue: true, PushRecordReturnValue: true}
	connection := output_connection.New(1, `Test`)
	connection.Add(api.NewConnectionInterfaceStruct(iiPush1))
	connection.Add(api.NewConnectionInterfaceStruct(iiPush2))

	generator := recordinfo.NewGenerator()
	generator.AddByteField(`SomeByte`, ``)
	info := generator.GenerateRecordInfo()
	err := connection.Init(info)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}

	connection.UpdateProgress(0.66)
	if iiPush1.UpdateProgressResult != 0.66 {
		t.Fatalf(`expected progress for 1 to be 0.66 but got %v`, iiPush1.UpdateProgressResult)
	}

	if iiPush2.UpdateProgressResult != 0.66 {
		t.Fatalf(`expected progress for 2 to be 0.66 but got %v`, iiPush2.UpdateProgressResult)
	}

	connection.Close()
	if iiPush1.UpdateProgressResult != 1.0 {
		t.Fatalf(`expected progress for 1 to be 1.0 but got %v`, iiPush1.UpdateProgressResult)
	}

	if iiPush2.UpdateProgressResult != 1.0 {
		t.Fatalf(`expected progress for 2 to be 1.0 but got %v`, iiPush2.UpdateProgressResult)
	}
}

type IiTestStruct struct {
	InitReturnValue       bool
	PushRecordReturnValue bool
	UpdateProgressResult  float64
	IsClosed              bool
	IsInitialized         bool
	PushRecordCalls       int
	CloseCalls            int
}

func (i *IiTestStruct) Init(recordInfoIn string) bool {
	i.IsInitialized = i.InitReturnValue
	return i.InitReturnValue
}

func (i *IiTestStruct) PushRecord(record unsafe.Pointer) bool {
	i.PushRecordCalls++
	return i.PushRecordReturnValue
}

func (i *IiTestStruct) UpdateProgress(percent float64) {
	i.UpdateProgressResult = percent
}

func (i *IiTestStruct) Close() {
	i.CloseCalls++
	i.IsClosed = true
}

func (i *IiTestStruct) Free() {
	return
}
