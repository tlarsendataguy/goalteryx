package api_new

type FileTestRunner struct {
	io     *testIo
	plugin *goPluginSharedMemory
}

func (r *FileTestRunner) ConnectToOutgoingAnchor(name string) {

}
