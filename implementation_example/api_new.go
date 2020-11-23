package main

/*
#include "implementation.h"
*/
import "C"
import (
	"github.com/tlarsen7572/goalteryx/api_new"
	"unsafe"
)

//export NewApiEntry
func NewApiEntry(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	return C.long(api_new.ConfigureTool(nil, int(toolId), xmlProperties, engineInterface, pluginInterface))
}
