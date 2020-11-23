package api_new

/*
#include "sdk.h"
*/
import "C"
import "unsafe"

func ConfigureTool(plugin Plugin, toolId int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	return C.configurePlugin(C.int(toolId), xmlProperties, (*C.struct_EngineInterface)(engineInterface), (*C.struct_PluginInterface)(pluginInterface))
}
