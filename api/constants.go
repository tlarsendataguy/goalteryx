package api

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

type InitVar string

const (
	RunMode                 InitVar = `RunMode`
	ActionApplies           InitVar = `ActionApplies`
	RunningAsMacro          InitVar = `RunningAsMacro`
	RuntimeDataPath         InitVar = `RuntimeDataPath`
	SettingsPath            InitVar = `SettingsPath`
	NumThreads              InitVar = `NumThreads`
	UpdateOnly              InitVar = `UpdateOnly`
	UpdateMode              InitVar = `UpdateMode`
	AllowDesktopInteraction InitVar = `AllowDesktopInteraction`
	DefaultDir              InitVar = `DefaultDir`
	Version                 InitVar = `Version`
	SerialNumber            InitVar = `SerialNumber`
	OutputRecordCounts      InitVar = `OutputRecordCounts`
)
