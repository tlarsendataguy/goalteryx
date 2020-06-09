package main

/*
#include "implementation.h"
*/
import "C"
import (
	"github.com/tlarsen7572/goalteryx/api"
	"unsafe"
)

func main() {}

//export PluginEntry
func PluginEntry(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &Plugin{}
	return C.long(api.ConfigurePlugin(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}

//export PluginPresortEntry
func PluginPresortEntry(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &PluginPresort{}
	return C.long(api.ConfigurePlugin(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}

//export PluginInputEntry
func PluginInputEntry(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &PluginInput{}
	return C.long(api.ConfigurePlugin(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}

//export PluginNoCacheEntry
func PluginNoCacheEntry(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &PluginNoCache{}
	return C.long(api.ConfigurePlugin(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}

type ConfigXml struct {
	Field string `xml:"Field"`
}
