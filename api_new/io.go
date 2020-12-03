package api_new

import (
	"fmt"
)

type testIo struct{}

func (t *testIo) Error(message string) {
	println(fmt.Sprintf(`ERROR: %v`, message))
}

func (t *testIo) Warn(message string) {
	println(fmt.Sprintf(`WARNING: %v`, message))
}

func (t *testIo) Info(message string) {
	println(fmt.Sprintf(`INFO: %v`, message))
}

func (t *testIo) UpdateProgress(progress float64) {
	println(fmt.Sprintf(`Progress: %v`, progress))
}

type ayxIo struct {
	sharedMemory *goPluginSharedMemory
}

func (a *ayxIo) Error(message string) {
	sendMessageToEngine(a.sharedMemory, Error, message)
}

func (a *ayxIo) Warn(message string) {
	sendMessageToEngine(a.sharedMemory, Warning, message)
}

func (a *ayxIo) Info(message string) {
	sendMessageToEngine(a.sharedMemory, Info, message)
}

func (a *ayxIo) UpdateProgress(progress float64) {
	sendToolProgressToEngine(a.sharedMemory, progress)
}
