package api_new

/*
#include "sdk.h"
*/
import "C"
import "unsafe"

func ConfigureTool(plugin Plugin, toolId int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	return C.configurePlugin(C.uint32_t(toolId), xmlProperties, (*C.struct_EngineInterface)(engineInterface), (*C.struct_PluginInterface)(pluginInterface))
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
