package api_new

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

type outputAnchor struct {
	data     *goOutputAnchorData
	metaData *OutgoingRecordInfo
}

func (a *outputAnchor) Name() string {
	name := utf16PtrToString(a.data.name, utf16PtrLen(a.data.name))
	return name
}

func (a *outputAnchor) IsOpen() bool {
	return a.data.isOpen == 1
}

func (a *outputAnchor) Metadata() *OutgoingRecordInfo {
	return a.metaData
}

func (a *outputAnchor) Open(info *OutgoingRecordInfo) {
	a.metaData = info
	xmlStr := info.toXml(a.Name())
	openOutgoingAnchor(a.data, xmlStr)
}

func (a *outputAnchor) Write() {
	panic("implement me")
}

func (a *outputAnchor) UpdateProgress(progress float64) {
	panic("implement me")
}
