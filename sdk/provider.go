package sdk

type Provider interface {
	ToolConfig() string
	Io() Io
	GetOutputAnchor(string) OutputAnchor
	Environment() Environment
}

type provider struct {
	sharedMemory  *goPluginSharedMemory
	io            Io
	environment   Environment
	outputAnchors map[string]*outputAnchor
}

func (p *provider) ToolConfig() string {
	return utf16PtrToString(p.sharedMemory.toolConfig, int(p.sharedMemory.toolConfigLen))
}

func (p *provider) Io() Io {
	return p.io
}

func (p *provider) GetOutputAnchor(name string) OutputAnchor {
	anchor, ok := p.outputAnchors[name]
	if ok {
		return anchor
	}
	anchorData := getOrCreateOutputAnchor(p.sharedMemory, name)
	anchor = &outputAnchor{data: anchorData}
	p.outputAnchors[name] = anchor
	return anchor
}

func (p *provider) Environment() Environment {
	return p.environment
}
