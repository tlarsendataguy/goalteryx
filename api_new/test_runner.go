package api_new

import "unsafe"

type FileReceiver struct {
}

type FileTestRunner struct {
	io           *testIo
	plugin       *goPluginSharedMemory
	ayxInterface unsafe.Pointer
}

func (r *FileTestRunner) SimulateInputTool() {
	simulateInputLifecycle(r.ayxInterface)
}

func (r *FileTestRunner) CaptureOutgoingAnchor(name string) *FileReceiver {

	return &FileReceiver{}
}
