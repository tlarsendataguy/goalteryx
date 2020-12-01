package api_new

/*
#include "sdk.h"
*/
import "C"
import (
	"unicode/utf16"
	"unsafe"
)

type goPluginSharedMemory struct {
	toolId                 uint32
	toolConfig             unsafe.Pointer
	engine                 unsafe.Pointer
	outputAnchors          *goOutputAnchorData
	totalInputConnections  uint32
	closedInputConnections uint32
	inputAnchors           *goInputAnchorData
}

type goOutputAnchorData struct {
	name                unsafe.Pointer
	metadata            unsafe.Pointer
	isOpen              byte
	firstChild          *goOutputConnectionData
	nextAnchor          *goOutputAnchorData
	recordCache         unsafe.Pointer
	recordCachePosition uint32
}

type goOutputConnectionData struct {
	isOpen         byte
	ii             unsafe.Pointer
	nextConnection *goOutputConnectionData
}

type goInputAnchorData struct {
	name       unsafe.Pointer
	firstChild *goInputConnectionData
	nextAnchor *goInputAnchorData
}

type goInputConnectionData struct {
	isOpen              byte
	metadata            unsafe.Pointer
	percent             float64
	nextConnection      *goInputConnectionData
	plugin              *goPluginSharedMemory
	fixedSize           uint32
	hasVarFields        byte
	recordCache         unsafe.Pointer
	recordCachePosition uint32
}

var tools = map[*goPluginSharedMemory]Plugin{} // = make(map[uint32]goPluginWrapper)

func RegisterTool(plugin Plugin, toolId int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) int {
	data := (*goPluginSharedMemory)(C.configurePlugin(C.uint32_t(toolId), (*C.wchar_t)(xmlProperties), (*C.struct_EngineInterface)(engineInterface), (*C.struct_PluginInterface)(pluginInterface)))
	tools[data] = plugin
	plugin.Init(nil)
	return 1
}

func RegisterToolTest(plugin Plugin, toolId int, xmlProperties string) int {
	xmlRunes := []rune(xmlProperties)
	xmlUtf16 := unsafe.Pointer(&utf16.Encode(xmlRunes)[0])
	engine := C.malloc(148)
	pluginInterface := C.malloc(44)
	return RegisterTool(plugin, toolId, xmlUtf16, engine, pluginInterface)
}

//export goOnInputConnectionOpened
func goOnInputConnectionOpened(handle unsafe.Pointer) {

}

//export goOnRecordPacket
func goOnRecordPacket(handle unsafe.Pointer) {

}

//export goOnSingleRecord
func goOnSingleRecord(handle unsafe.Pointer, record unsafe.Pointer) {

}

//export goOnComplete
func goOnComplete(handle unsafe.Pointer) {

}
