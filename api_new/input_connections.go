package api_new

type ImpInputConnection struct {
	data *goInputConnectionData
}

func (i *ImpInputConnection) Name() string {
	nameLen := utf16PtrLen(i.data.anchor.name)
	name := utf16PtrToString(i.data.anchor.name, nameLen)
	return name
}

func (i *ImpInputConnection) Metadata() string {
	configLen := utf16PtrLen(i.data.metadata)
	config := utf16PtrToString(i.data.metadata, configLen)
	return config
}

func (i *ImpInputConnection) Read() Record {
	panic("implement me")
}

func (i *ImpInputConnection) Progress() float64 {
	panic("implement me")
}
