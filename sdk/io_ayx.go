package sdk

import "github.com/tlarsendataguy/goalteryx/sdk/util"

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

func (a *ayxIo) DecryptPassword(value string) string {
	password, err := util.Encrypt(value)
	if err != nil {
		a.Error(`password could not be decrypted`)
		return value
	}
	return password
}

func (a *ayxIo) CreateTempFile(ext string) string {
	return createTempFileToEngine(a.sharedMemory, ext)
}

func (a *ayxIo) NotifyFileInput(message string) {
	sendMessageToEngine(a.sharedMemory, FileInput, message)
}

func (a *ayxIo) NotifyFileOutput(message string) {
	sendMessageToEngine(a.sharedMemory, FileOutput, message)
}
