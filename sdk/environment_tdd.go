package sdk

type testEnvironment struct {
	sharedMemory *goPluginSharedMemory
	updateOnly   bool
	updateMode   string
	workflowDir  string
	locale       string
}

func (e *testEnvironment) UpdateOnly() bool {
	return e.updateOnly
}

func (e *testEnvironment) UpdateMode() string {
	return e.updateMode
}

func (e *testEnvironment) DesignerVersion() string {
	return `TestHarness`
}

func (e *testEnvironment) WorkflowDir() string {
	return e.workflowDir
}

func (e *testEnvironment) AlteryxInstallDir() string {
	return ``
}

func (e *testEnvironment) AlteryxLocale() string {
	return e.locale
}

func (e *testEnvironment) ToolId() int {
	return int(e.sharedMemory.toolId)
}

func (e *testEnvironment) UpdateToolConfig(newConfig string) {
	updateConfig(e.sharedMemory, newConfig)
}
