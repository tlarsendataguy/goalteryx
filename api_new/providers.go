package api_new

type provider struct {
	sharedMemory *goPluginSharedMemory
	config       string
	io           Io
}

func (p *provider) ToolConfig() string {
	return p.config
}

func (p *provider) Io() Io {
	return p.io
}

func (p *provider) GetOutputAnchor(s string) OutputAnchor {
	panic("implement me")
}

func (p *provider) Environment() Environment {
	panic("implement me")
}
