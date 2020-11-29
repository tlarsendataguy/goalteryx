package api_new

/*
#include "sdk.h"
*/
import "C"
import "unsafe"

type PluginSharedMemory struct {
	toolId                 uint32
	toolConfig             unsafe.Pointer
	engine                 unsafe.Pointer
	outputAnchors          *OutputAnchorData
	totalInputConnections  uint32
	closedInputConnections uint32
	inputAnchors           *InputAnchorData
}

type OutputAnchorData struct {
	name        unsafe.Pointer
	metadata    unsafe.Pointer
	isOpen      uint32
	firstChild  *OutputConnectionData
	nextAnchor  *OutputAnchorData
	recordCache unsafe.Pointer
}

type OutputConnectionData struct {
	isOpen         uint32
	ii             unsafe.Pointer
	nextConnection *OutputConnectionData
}

type InputAnchorData struct {
	name       unsafe.Pointer
	firstChild *InputConnectionData
	nextAnchor *InputAnchorData
}

type InputConnectionData struct {
	isOpen         uint32
	metadata       unsafe.Pointer
	percent        float64
	nextConnection *InputConnectionData
	plugin         *PluginSharedMemory
	recordCache    unsafe.Pointer
}

var tools = make(map[*PluginSharedMemory]Plugin)

func RegisterTool(plugin Plugin, toolId int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) int {
	data := (*PluginSharedMemory)(C.configurePlugin(C.uint32_t(toolId), xmlProperties, (*C.struct_EngineInterface)(engineInterface), (*C.struct_PluginInterface)(pluginInterface)))
	tools[data] = plugin
	return 1
}

//export Init
func Init(handle unsafe.Pointer) {

}

//export OnInputConnectionOpened
func OnInputConnectionOpened(handle unsafe.Pointer) {

}

//export OnRecordPacket
func OnRecordPacket(handle unsafe.Pointer) {

}

//export OnComplete
func OnComplete(handle unsafe.Pointer) {

}
