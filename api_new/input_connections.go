package api_new

type InputConnection interface {
	Name() string
	Metadata() IncomingRecordInfo
	Read() RecordPacket
	Progress() float64
}

type ImpInputConnection struct {
	data *goInputConnectionData
}

func (i *ImpInputConnection) Name() string {
	nameLen := utf16PtrLen(i.data.anchor.name)
	name := utf16PtrToString(i.data.anchor.name, nameLen)
	return name
}

func (i *ImpInputConnection) Metadata() IncomingRecordInfo {
	configLen := utf16PtrLen(i.data.metadata)
	configStr := utf16PtrToString(i.data.metadata, configLen)
	config, _ := incomingRecordInfoFromString(configStr)
	return config
}

func (i *ImpInputConnection) Read() RecordPacket {
	return NewRecordPacket(RecordCache(i.data.recordCache), int(i.data.recordCachePosition), int(i.data.fixedSize), i.data.hasVarFields == 1)
}

func (i *ImpInputConnection) Progress() float64 {
	return i.data.percent
}
