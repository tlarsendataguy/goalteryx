package sdk

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
