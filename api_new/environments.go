package api_new

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

func (e *testEnvironment) UpdateToolConfig(s string) {
	panic("implement me")
}

type ayxEnvironment struct {
	sharedMemory *goPluginSharedMemory
}

func (e *ayxEnvironment) UpdateOnly() bool {
	panic("implement me")
}

func (e *ayxEnvironment) UpdateMode() string {
	panic("implement me")
}

func (e *ayxEnvironment) DesignerVersion() string {
	panic("implement me")
}

func (e *ayxEnvironment) WorkflowDir() string {
	panic("implement me")
}

func (e *ayxEnvironment) AlteryxInstallDir() string {
	panic("implement me")
}

func (e *ayxEnvironment) AlteryxLocale() string {
	panic("implement me")
}

func (e *ayxEnvironment) ToolId() int {
	return int(e.sharedMemory.toolId)
}

func (e *ayxEnvironment) UpdateToolConfig(s string) {
	panic("implement me")
}
