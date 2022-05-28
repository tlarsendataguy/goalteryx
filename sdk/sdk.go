package sdk

/*
#include "sdk.h"
*/
import "C"
import (
	"reflect"
	"unicode/utf16"
	"unsafe"
)

const cacheSize uint32 = 4194304

type goPluginSharedMemory struct {
	toolId                 uint32
	toolConfig             unsafe.Pointer
	toolConfigLen          uint32
	engine                 unsafe.Pointer
	ayxInterface           unsafe.Pointer
	outputAnchors          *goOutputAnchorData
	totalInputConnections  uint32
	closedInputConnections uint32
	inputAnchors           *goInputAnchorData
}

type goOutputAnchorData struct {
	name                unsafe.Pointer
	metadata            unsafe.Pointer
	browseEverywhereId  uint32
	isOpen              byte
	plugin              *goPluginSharedMemory
	firstChild          *goOutputConnectionData
	nextAnchor          *goOutputAnchorData
	fixedSize           uint32
	hasVarFields        byte
	recordCache         unsafe.Pointer
	recordCachePosition uint32
	recordCacheSize     uint32
	recordCount         uint64
	totalDataSize       uint64
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

type Status byte

const (
	Created          = 1
	Initialized      = 2
	ReceivingRecords = 3
	Closed           = 4
)

type goInputConnectionData struct {
	anchor              *goInputAnchorData
	isOpen              byte
	status              Status
	metadata            unsafe.Pointer
	percent             float64
	nextConnection      *goInputConnectionData
	plugin              *goPluginSharedMemory
	fixedSize           uint32
	hasVarFields        byte
	recordCache         unsafe.Pointer
	recordCachePosition uint32
	recordCacheSize     uint32
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

func stringToUtf16Ptr(value string) *C.utf16char {
	utf16Bytes := append(utf16.Encode([]rune(value)), 0)

	length := len(utf16Bytes)
	byteData := allocateCache(uint32(length * 2))

	//We need to copy the UTF16 bytes to a pointer allocated from C so no Go pointers end up in C space
	var byteSlice []uint16
	rawHeader := (*reflect.SliceHeader)(unsafe.Pointer(&byteSlice))
	rawHeader.Data = uintptr(byteData)
	rawHeader.Len = length
	rawHeader.Cap = length
	copy(byteSlice, utf16Bytes)

	return (*C.utf16char)(byteData)
}

func simulateInputLifecycle(pluginInterface unsafe.Pointer) {
	C.simulateInputLifecycle((*C.struct_PluginInterface)(pluginInterface))
}

func sendMessageToEngine(data *goPluginSharedMemory, status MessageStatus, message string) {
	C.sendMessage((*C.struct_EngineInterface)(data.engine), (C.int)(data.toolId), (C.int)(status), (*C.utf16char)(stringToUtf16Ptr(message)))
}

func sendToolProgressToEngine(data *goPluginSharedMemory, progress float64) {
	C.outputToolProgress((*C.struct_EngineInterface)(data.engine), (C.int)(data.toolId), (C.double)(progress))
}

func sendProgressToAnchor(anchor *goOutputAnchorData, progress float64) {
	C.sendProgressToAnchor((*C.struct_OutputAnchor)(unsafe.Pointer(anchor)), (C.double)(progress))
}

func getInitVarToEngine(data *goPluginSharedMemory, initVar string) string {
	initVarPtr := stringToUtf16Ptr(initVar)
	resultPtr := C.getInitVar((*C.struct_EngineInterface)(data.engine), initVarPtr)
	length := utf16PtrLen(resultPtr)
	return utf16PtrToString(resultPtr, length)
}

func createTempFileToEngine(data *goPluginSharedMemory, ext string) string {
	filePathPtr := C.createTempFile((*C.struct_EngineInterface)(data.engine), (*C.utf16char)(stringToUtf16Ptr(ext)))
	length := utf16PtrLen(filePathPtr)
	return utf16PtrToString(filePathPtr, length)
}

func getOrCreateOutputAnchor(sharedMemory *goPluginSharedMemory, name string) *goOutputAnchorData {
	anchor := sharedMemory.outputAnchors

	for {
		if anchor == nil {
			nameUtf16 := stringToUtf16Ptr(name)
			cAnchor := unsafe.Pointer(C.appendOutgoingAnchor((*C.struct_PluginSharedMemory)(unsafe.Pointer(sharedMemory)), nameUtf16))
			anchor = (*goOutputAnchorData)(cAnchor)
			return anchor
		}
		if anchorName := utf16PtrToString(anchor.name, utf16PtrLen(anchor.name)); anchorName == name {
			return anchor
		}
		anchor = anchor.nextAnchor
	}
}

func registerAndInit(plugin Plugin, data *goPluginSharedMemory, provider Provider) {
	tools[data] = plugin
	plugin.Init(provider)
}

func generateIncomingConnectionInterface() unsafe.Pointer {
	return unsafe.Pointer(C.generateIncomingConnectionInterface())
}

func callPiAddIncomingConnection(plugin *goPluginSharedMemory, name string, ii unsafe.Pointer) {
	namePtr := stringToUtf16Ptr(name)
	C.callPiAddIncomingConnection((*C.struct_PluginSharedMemory)(unsafe.Pointer(plugin)), namePtr, (*C.struct_IncomingConnectionInterface)(ii))
}

func callPiAddOutgoingConnection(plugin *goPluginSharedMemory, name string, ii unsafe.Pointer) {
	namePtr := stringToUtf16Ptr(name)
	C.callPiAddOutgoingConnection((*C.struct_PluginSharedMemory)(unsafe.Pointer(plugin)), namePtr, (*C.struct_IncomingConnectionInterface)(ii))
}

func RegisterTool(plugin Plugin, toolId int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer, optionSetters ...ToolOptionSetter) int {
	options := toolOptions{}
	for _, setter := range optionSetters {
		options = setter(options)
	}
	var data *goPluginSharedMemory
	if options.noCache {
		data = (*goPluginSharedMemory)(C.configurePluginNoCache(C.uint32_t(toolId), (*C.utf16char)(xmlProperties), (*C.struct_EngineInterface)(engineInterface), (*C.struct_PluginInterface)(pluginInterface)))
	} else {
		data = (*goPluginSharedMemory)(C.configurePlugin(C.uint32_t(toolId), (*C.utf16char)(xmlProperties), (*C.struct_EngineInterface)(engineInterface), (*C.struct_PluginInterface)(pluginInterface)))
	}
	io := &ayxIo{sharedMemory: data}
	environment := &ayxEnvironment{sharedMemory: data}
	toolProvider := &provider{
		sharedMemory:  data,
		io:            io,
		environment:   environment,
		outputAnchors: make(map[string]*outputAnchor),
	}

	registerAndInit(plugin, data, toolProvider)
	return 1
}

func RegisterToolTest(plugin Plugin, toolId int, xmlProperties string, optionSetters ...OptionSetter) *FileTestRunner {
	options := testOptions{
		updateOnly:  false,
		updateMode:  "",
		workflowDir: "",
		locale:      "en",
		noCache:     false,
	}
	for _, optionSetter := range optionSetters {
		options = optionSetter(options)
	}
	xmlRunes := []rune(xmlProperties)
	xmlUtf16 := append(utf16.Encode(xmlRunes), 0)
	xmlPtr := unsafe.Pointer(&xmlUtf16[0])
	pluginInterface := unsafe.Pointer(C.generatePluginInterface())
	var data *goPluginSharedMemory
	if options.noCache {
		data = (*goPluginSharedMemory)(C.configurePluginNoCache(C.uint32_t(toolId), (*C.utf16char)(xmlPtr), nil, (*C.struct_PluginInterface)(pluginInterface)))
	} else {
		data = (*goPluginSharedMemory)(C.configurePlugin(C.uint32_t(toolId), (*C.utf16char)(xmlPtr), nil, (*C.struct_PluginInterface)(pluginInterface)))
	}
	io := &testIo{}
	environment := &testEnvironment{
		sharedMemory: data,
		updateOnly:   options.updateOnly,
		updateMode:   options.updateMode,
		workflowDir:  options.workflowDir,
		locale:       options.locale,
	}
	toolProvider := &provider{
		sharedMemory:  data,
		io:            io,
		environment:   environment,
		outputAnchors: make(map[string]*outputAnchor),
	}
	registerAndInit(plugin, data, toolProvider)
	return &FileTestRunner{
		io:     io,
		plugin: data,
		inputs: make(map[string]*FilePusher),
	}
}

func registerTestHarness(plugin Plugin) *goPluginSharedMemory {
	var toolId uint32 = 1
	for {
		found := false
		for key := range tools {
			if key.toolId == toolId {
				found = true
				break
			}
		}
		if found {
			toolId++
			continue
		}
		break
	}

	pluginInterface := unsafe.Pointer(C.generatePluginInterface())
	config := stringToUtf16Ptr("<Configuration></Configuration>")
	data := (*goPluginSharedMemory)(C.configurePlugin(C.uint32_t(toolId), (*C.utf16char)(config), nil, (*C.struct_PluginInterface)(pluginInterface)))
	io := &testIo{}
	environment := &testEnvironment{
		sharedMemory: data,
	}
	toolProvider := &provider{
		sharedMemory:  data,
		io:            io,
		environment:   environment,
		outputAnchors: make(map[string]*outputAnchor),
	}
	registerAndInit(plugin, data, toolProvider)
	return data
}

func openOutgoingAnchor(anchor *goOutputAnchorData, config string) {
	configPtr := stringToUtf16Ptr(config)
	C.openOutgoingAnchor((*C.struct_OutputAnchor)(unsafe.Pointer(anchor)), configPtr)
}

//export goOnInputConnectionOpened
func goOnInputConnectionOpened(handle unsafe.Pointer) {
	var data = (*goInputConnectionData)(handle)
	plugin := tools[data.plugin]
	inputConnection := &ImpInputConnection{
		data: data,
	}
	var hasVarFields byte = 0
	var fixedSize uint32 = 0
	fields := inputConnection.Metadata().fields
	for _, field := range fields {
		switch field.Type {
		case `Bool`:
			fixedSize += 1
		case `Byte`:
			fixedSize += 2
		case `Int16`:
			fixedSize += 3
		case `Int32`, `Float`:
			fixedSize += 5
		case `Int64`, `Double`:
			fixedSize += 9
		case `FixedDecimal`:
			fixedSize += uint32(field.Size) + 1
		case `Date`:
			fixedSize += 11
		case `DateTime`:
			fixedSize += 20
		case `String`:
			fixedSize += uint32(field.Size) + 1
		case `WString`:
			fixedSize += uint32(field.Size*2) + 1
		case `V_String`, `V_WString`, `Blob`, `SpatialObj`:
			fixedSize += 4
			hasVarFields = 1
		}
	}
	data.fixedSize = fixedSize
	data.hasVarFields = hasVarFields
	plugin.OnInputConnectionOpened(inputConnection)
}

//export goOnRecordPacket
func goOnRecordPacket(handle unsafe.Pointer) {
	data := (*goInputConnectionData)(handle)
	connection := &ImpInputConnection{data: data}
	implementation := tools[data.plugin]
	implementation.OnRecordPacket(connection)
}

//export goOnComplete
func goOnComplete(handle unsafe.Pointer) {
	data := (*goPluginSharedMemory)(handle)
	implementation := tools[data]
	implementation.OnComplete()
	for anchor := data.outputAnchors; anchor != nil; anchor = anchor.nextAnchor {
		if anchor.recordCachePosition > 0 {
			callWriteRecords(unsafe.Pointer(anchor))
		}
	}
	delete(tools, data)
}

func callWriteRecords(handle unsafe.Pointer) {
	C.callWriteRecords((*C.struct_OutputAnchor)(handle))
}

func callCloseOutputAnchor(anchor *goOutputAnchorData) {
	if anchor.recordCachePosition > 0 {
		callWriteRecords(unsafe.Pointer(anchor))
	}
	C.closeOutputAnchor((*C.struct_OutputAnchor)(unsafe.Pointer(anchor)))
}

func allocateCache(size uint32) unsafe.Pointer {
	return C.allocateCache(C.int(size))
}

func freeCache(cache unsafe.Pointer) {
	C.free(cache)
}
