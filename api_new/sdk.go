package api_new

/*
#include "sdk.h"
*/
import "C"
import (
	"reflect"
	"unicode/utf16"
	"unsafe"
)

type goPluginSharedMemory struct {
	toolId                 uint32
	toolConfig             unsafe.Pointer
	toolConfigLen          uint32
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

var tools = map[*goPluginSharedMemory]Plugin{}

func utf16PtrToString(utf16Ptr unsafe.Pointer, len int) string {
	var utf16Slice []uint16
	rawHeader := (*reflect.SliceHeader)(unsafe.Pointer(&utf16Slice))
	rawHeader.Data = uintptr(utf16Ptr)
	rawHeader.Len = len
	rawHeader.Cap = len
	return string(utf16.Decode(utf16Slice))
}

func utf16PtrLen(utf16Ptr unsafe.Pointer) int {
	length := uintptr(0)
	for {
		currentChar := *(*uint16)(unsafe.Pointer(uintptr(utf16Ptr) + (length * 2)))
		if currentChar == 0 {
			break
		}
		length++
	}
	return int(length)
}

func stringToUtf16Ptr(value string) *C.wchar_t {
	utf16Bytes := append(utf16.Encode([]rune(value)), 0)
	return (*C.wchar_t)(&utf16Bytes[0])
}

func sendMessageToEngine(data *goPluginSharedMemory, status MessageStatus, message string) {
	C.sendMessage((*C.struct_EngineInterface)(data.engine), (C.int)(data.toolId), (C.int)(status), (*C.wchar_t)(stringToUtf16Ptr(message)))
}

func sendToolProgressToEngine(data *goPluginSharedMemory, progress float64) {
	C.outputToolProgress((*C.struct_EngineInterface)(data.engine), (C.int)(data.toolId), (C.double)(progress))
}

func getInitVarToEngine(data *goPluginSharedMemory, initVar string) string {
	initVarPtr := stringToUtf16Ptr(initVar)
	resultPtr := C.getInitVar((*C.struct_EngineInterface)(data.engine), initVarPtr)
	length := utf16PtrLen(resultPtr)
	return utf16PtrToString(resultPtr, length)
}

func registerAndInit(plugin Plugin, data *goPluginSharedMemory, provider Provider) {
	tools[data] = plugin
	plugin.Init(provider)
}

func RegisterTool(plugin Plugin, toolId int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) int {
	data := (*goPluginSharedMemory)(C.configurePlugin(C.uint32_t(toolId), (*C.wchar_t)(xmlProperties), (*C.struct_EngineInterface)(engineInterface), (*C.struct_PluginInterface)(pluginInterface)))
	io := &ayxIo{sharedMemory: data}
	environment := &ayxEnvironment{sharedMemory: data}
	config := utf16PtrToString(data.toolConfig, int(data.toolConfigLen))
	provider := &provider{
		sharedMemory: data,
		config:       config,
		io:           io,
		environment:  environment,
	}

	registerAndInit(plugin, data, provider)
	return 1
}

func RegisterToolTest(plugin Plugin, toolId int, xmlProperties string, optionSetters ...OptionSetter) *FileTestRunner {
	options := testOptions{
		updateOnly:  false,
		updateMode:  "",
		workflowDir: "",
		locale:      "en",
	}
	for _, optionSetter := range optionSetters {
		options = optionSetter(options)
	}
	xmlRunes := []rune(xmlProperties)
	xmlUtf16 := append(utf16.Encode(xmlRunes), 0)
	xmlPtr := unsafe.Pointer(&xmlUtf16[0])
	pluginInterface := C.malloc(44)
	data := (*goPluginSharedMemory)(C.configurePlugin(C.uint32_t(toolId), (*C.wchar_t)(xmlPtr), nil, (*C.struct_PluginInterface)(pluginInterface)))
	io := &testIo{}
	environment := &testEnvironment{
		sharedMemory: data,
		updateOnly:   options.updateOnly,
		updateMode:   options.updateMode,
		workflowDir:  options.workflowDir,
		locale:       options.locale,
	}
	provider := &provider{
		sharedMemory: data,
		config:       xmlProperties,
		io:           io,
		environment:  environment,
	}
	registerAndInit(plugin, data, provider)
	return &FileTestRunner{
		io: io,
	}
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
