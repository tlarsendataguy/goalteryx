package api_new

import "unsafe"

type MessageStatus int

const (
	Info                          MessageStatus = 1
	TransientInfo                 MessageStatus = 0x40000000 | 1
	Warning                       MessageStatus = 2
	TransientWarning              MessageStatus = 0x40000000 | 2
	Error                         MessageStatus = 3
	Complete                      MessageStatus = 4
	FieldConversionError          MessageStatus = 5
	TransientFieldConversionError MessageStatus = 0x40000000 | 5
	UpdateOutputMetaInfoXml       MessageStatus = 10
	RecordCountString             MessageStatus = 50
	BrowseEverywhereFileName      MessageStatus = 70
)

type Plugin interface {
	Init(Provider)
	OnInputConnectionOpened(InputConnection)
	OnRecordPacket(InputConnection)
	OnComplete()
}

type Io interface {
	Error(string)
	Warn(string)
	Info(string)
	UpdateProgress(float64)
}

type InputConnection interface {
	Name() string
	Metadata() IncomingRecordInfo
	Read() RecordPacket
	Progress() float64
}

type OutputAnchor interface {
	Name() string
	IsOpen() bool
	Metadata() *OutgoingRecordInfo
	Open(info *OutgoingRecordInfo)
	Write()
	UpdateProgress(float64)
}

type Environment interface {
	UpdateOnly() bool
	UpdateMode() string
	DesignerVersion() string
	WorkflowDir() string
	AlteryxInstallDir() string
	AlteryxLocale() string
	ToolId() int
	UpdateToolConfig(string)
}

type Provider interface {
	ToolConfig() string
	Io() Io
	GetOutputAnchor(string) OutputAnchor
	Environment() Environment
}

type RecordPacket interface {
	Next() bool
	Record() Record
}

type Record = unsafe.Pointer
