package api_new

import (
	"unsafe"
)

type FileTestRunner struct {
	io           *testIo
	plugin       *goPluginSharedMemory
	ayxInterface unsafe.Pointer
}

func (r *FileTestRunner) SimulateInputTool() {
	simulateInputLifecycle(r.ayxInterface)
}

func (r *FileTestRunner) CaptureOutgoingAnchor(name string) *RecordCollector {
	collector := &RecordCollector{}
	sharedMemory := registerTestHarness(collector)

	ii := generateIncomingConnectionInterface()
	callPiAddIncomingConnection(sharedMemory, name, ii)
	callPiAddOutgoingConnection(r.plugin, name, ii)

	return collector
}

type RecordCollector struct {
	Config IncomingRecordInfo
	Name   string
}

func (r *RecordCollector) Init(provider Provider) {}

func (r *RecordCollector) OnInputConnectionOpened(connection InputConnection) {
	r.Name = connection.Name()
	r.Config = connection.Metadata()
}

func (r *RecordCollector) OnRecordPacket(connection InputConnection) {
	panic("implement me")
}

func (r *RecordCollector) OnComplete() {
	panic("implement me")
}
