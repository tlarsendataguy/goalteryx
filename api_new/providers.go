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

func (p *provider) GetOutputAnchor(s string) OutputAnchor {
	panic("implement me")
}

func (p *provider) Environment() Environment {
	return p.environment
}
