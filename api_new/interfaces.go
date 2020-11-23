package api_new

import "unsafe"

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
	Metadata() string
	Read() Record
	Progress() float64
}

type OutputAnchor interface {
	Name() string
	IsOpen() bool
	Metadata() string
	Open(string)
	Write(Record)
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

type Record = unsafe.Pointer
