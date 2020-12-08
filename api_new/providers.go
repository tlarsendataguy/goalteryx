package api_new

type provider struct {
	sharedMemory *goPluginSharedMemory
	io           Io
	environment  Environment
}

func (p *provider) ToolConfig() string {
	return utf16PtrToString(p.sharedMemory.toolConfig, int(p.sharedMemory.toolConfigLen))
}

func (p *provider) Io() Io {
	return p.io
}

func (p *provider) GetOutputAnchor(name string) OutputAnchor {
	anchorData := getOrCreateOutputAnchor(p.sharedMemory, name)
	return &outputAnchor{data: anchorData}
}

func (p *provider) Environment() Environment {
	return p.environment
}

type outputAnchor struct {
	data *goOutputAnchorData
}

func (a *outputAnchor) Name() string {
	panic("implement me")
}

func (a *outputAnchor) IsOpen() bool {
	panic("implement me")
}

func (a *outputAnchor) Metadata() string {
	panic("implement me")
}

func (a *outputAnchor) Open(config string) {
	panic("implement me")
}

func (a *outputAnchor) Write(record Record) {
	panic("implement me")
}

func (a *outputAnchor) UpdateProgress(progress float64) {
	panic("implement me")
}
