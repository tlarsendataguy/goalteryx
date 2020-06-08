// Package api provides all of the glue to join Alteryx's C API with Go.
package api

/*
#include "plugins.h"
*/
import "C"
import (
	"fmt"
	"github.com/mattn/go-pointer"
	"github.com/tlarsen7572/goalteryx/convert_strings"
	"github.com/tlarsen7572/goalteryx/presort"
	"github.com/tlarsen7572/goalteryx/recordinfo"
	"os"
	"time"
	"unsafe"
)

// engine is our pointer to the Alteryx engine given to our tool.  We obtain this reference when a tool calls
// ConfigurePlugin.
var engine *C.struct_EngineInterface

// incomingInterfaces stores up to 1,000 IncomingInterface objects.  This means each DLL is limited to 1,000 input
// anchors.  We do this because we need a handle that can cross between Go and C.  We use a pointer to an integer
// that defines the index in incomingInterfaces where we can find out IncomingInterface.  This means our interface
// will never be garbage collected and our handle can cross back and forth from C.
var incomingInterfaces = make([]IncomingInterface, 1000)

// Plugin defines the methods a plugin must implement.
type Plugin interface {
	// Init is called when it is time to configure the plugin using the provided XML configuration string.
	Init(toolId int, config string) bool

	// PushAllRecords is called when our tool is an input tool and its time to start sending data to downstream tools.
	PushAllRecords(recordLimit int) bool

	// Close is called when the upstream tools are done sending our plugin data.
	Close(hasErrors bool)

	// AddIncomingConnection is called when an upstream tool is being connected to our tool.
	AddIncomingConnection(connectionType string, connectionName string) (IncomingInterface, *presort.PresortInfo)

	// AddOutgoingConnection is called when our tool is being connected to a downstream tool.
	AddOutgoingConnection(connectionName string, connectionInterface *ConnectionInterfaceStruct) bool

	// GetToolId should return the toolId provided by Init.
	GetToolId() int
}

// IncomingInterface handles incoming data.
type IncomingInterface interface {
	// Init is called when an upstream tool sends us its outgoing RecordInfo.
	Init(recordInfoIn string) bool

	// PushRecord is called when an upstream tool is pushing a record blob to our tool.
	PushRecord(record unsafe.Pointer) bool

	// UpdateProgress is called when an upstream tool is updating our tool with its progress.
	UpdateProgress(percent float64)

	// Close is called when an upstream tool is finished sending us data.
	Close()
}

// ConnectionInterfaceStruct is a wrapper around the C function pointer struct for Incoming Interfaces.
type ConnectionInterfaceStruct struct {
	connection *C.struct_IncomingConnectionInterface
}

// NewConnectionInterfaceStruct creates a new ConnectionInterfaceStruct.  Normally these are provided to us by the
// Alteryx engine.  However, to perform certain tests we need to generate them ourselves.  This function provides
// that facility.
func NewConnectionInterfaceStruct(incomingInterface IncomingInterface) *ConnectionInterfaceStruct {
	iiIndexHandle := C.getIiIndex()
	iiIndex := int(*(*C.int)(iiIndexHandle))
	incomingInterfaces[iiIndex] = incomingInterface
	var ii *C.struct_IncomingConnectionInterface = C.newIi(iiIndexHandle)
	return &ConnectionInterfaceStruct{connection: ii}
}

// Plugin methods

// ConfigurePlugin wires up a Go plugin to the C struct provided by the Alteryx engine.  ConfigurePlugin should
// be called by a plugin as soon as its entry point is called by the Alteryx engine.  Once ConfigurePlugin is called,
// all lifecycle events will be called by the Alteryx engine and this cgo layer.
func ConfigurePlugin(plugin Plugin, toolId int, pXmlProperties unsafe.Pointer, pEngineInterface unsafe.Pointer, r_pluginInterface unsafe.Pointer) int {
	engine = (*C.struct_EngineInterface)(pEngineInterface)
	C.c_setEngine(engine)

	config := convert_strings.WideCToString(pXmlProperties)
	if !plugin.Init(toolId, config) {
		return 0
	}

	pluginInterface := (*C.struct_PluginInterface)(r_pluginInterface)
	handle := getPlugin(plugin)
	C.c_configurePlugin(handle, pluginInterface)
	return 1
}

// getPlugin generates a pointer we can use as a plugin's handle.  We cannot use Go pointers in C, so this function
// is needed to traverse the barrier.
func getPlugin(plugin Plugin) unsafe.Pointer {
	return pointer.Save(plugin)
}

//export go_piPushAllRecords
// go_piPushAllRecords calls the PushAllRecords method on the plugin.  It is called from c_piPushAllRecords.
func go_piPushAllRecords(handle unsafe.Pointer, recordLimit C.__int64) C.long {
	alteryxPlugin := pointer.Restore(handle).(Plugin)
	if alteryxPlugin.PushAllRecords(int(recordLimit)) {
		return C.long(1)
	}
	return C.long(0)
}

//export go_piClose
// go_piClose calls the Close method on the plugin.  It is called from c_piClose.
func go_piClose(handle unsafe.Pointer, hasErrors C.bool) {
	alteryxPlugin := pointer.Restore(handle).(Plugin)
	alteryxPlugin.Close(bool(hasErrors))
}

//export go_piAddIncomingConnection
// go_piAddIncomingConnection calls the AddIncomingConnection method on the plugin.  The plugin should return an
// IncomingInterface and an optional PreSort configuration.  go_piAddIncomingConnection generates a handle for the
// returned IncomingInterface and returns is back to the C layer.  The optional PreSort configuration is converted to
// an XML string and also returned back to C.
//
// It is called from c_piAddIncomingConnection.
func go_piAddIncomingConnection(handle unsafe.Pointer, connectionType unsafe.Pointer, connectionName unsafe.Pointer) *C.struct_IncomingConnectionInfo {
	alteryxPlugin := pointer.Restore(handle).(Plugin)
	goName := convert_strings.WideCToString(connectionName)
	goType := convert_strings.WideCToString(connectionType)
	goIncomingInterface, presortInfo := alteryxPlugin.AddIncomingConnection(goType, goName)
	if goIncomingInterface == nil {
		return nil
	}
	iiIndexHandle := C.getIiIndex()
	iiIndex := int(*(*C.int)(iiIndexHandle))
	incomingInterfaces[iiIndex] = goIncomingInterface
	if presortInfo == nil {
		return C.newUnsortedIncomingConnectionInfo(iiIndexHandle)
	}
	presortInfoXml, _ := presortInfo.ToXml()
	cPresortInfoXml, _ := convert_strings.StringToWideC(presortInfoXml)
	return C.newSortedIncomingConnectionInfo(iiIndexHandle, cPresortInfoXml)
}

//export go_piAddOutgoingConnection
// go_piAddOutgoingConnection calls the Close method on the plugin.  It is called from c_piAddOutgoingConnection.
func go_piAddOutgoingConnection(handle unsafe.Pointer, connectionName unsafe.Pointer, incomingConnection *C.struct_IncomingConnectionInterface) C.long {
	alteryxPlugin := pointer.Restore(handle).(Plugin)
	goName := convert_strings.WideCToString(connectionName)
	connectionInterface := &ConnectionInterfaceStruct{connection: incomingConnection}
	if alteryxPlugin.AddOutgoingConnection(goName, connectionInterface) {
		return C.long(1)
	}
	return C.long(0)
}

// Incoming interface methods

//export go_iiInit
// go_iiInit calls the Init method on the IncomingInterface.  It also saves the RecordInfo's fixed size back
// to the C layer to allow C to cache records.  It is called from c_iiInit.
func go_iiInit(handle unsafe.Pointer, recordInfoIn unsafe.Pointer) C.long {
	incomingInterface := getIncomingInterfaceFromHandle(handle)
	goRecordInfoIn := convert_strings.WideCToString(recordInfoIn)
	if incomingInterface.Init(goRecordInfoIn) {
		recordInfo, _ := recordinfo.FromXml(goRecordInfoIn)
		fixedSize := recordInfo.FixedSize()
		C.saveIncomingInterfaceFixedSize(handle, C.int(fixedSize))
		return C.long(1)
	}
	return C.long(0)
}

//export go_iiPushRecordCache
// go_iiPushRecordCache iterates through a cache of records and calls PushRecord on the IncomingInterface for each
// one.  For performance reasons, calls from C to Go must be minimized, so we cache records in C and then batch
// them to Go.  This process is transparent to the IncomingInterface, which only deals with 1 record at a time.  This
// keeps the Go SDK consistent with the other SDKs.
func go_iiPushRecordCache(handle unsafe.Pointer, cache unsafe.Pointer, cacheSize C.int) C.long {
	incomingInterface := getIncomingInterfaceFromHandle(handle)
	cacheArray := *((*[10]unsafe.Pointer)(cache))
	for i := 0; i < int(cacheSize); i++ {
		ok := incomingInterface.PushRecord(cacheArray[i])
		if !ok {
			return C.long(0)
		}
	}
	return C.long(1)
}

//export go_iiUpdateProgress
// go_iiUpdateProgress calls the Close method on the plugin.  It is called from c_iiUpdateProgress.
func go_iiUpdateProgress(handle unsafe.Pointer, percent C.double) {
	incomingInterface := getIncomingInterfaceFromHandle(handle)
	incomingInterface.UpdateProgress(float64(percent))
}

//export go_iiClose
// go_iiClose calls the Close method on the plugin.  It is called from c_iiClose.
func go_iiClose(handle unsafe.Pointer) {
	incomingInterface := getIncomingInterfaceFromHandle(handle)
	incomingInterface.Close()
}

// getIncomingInterfaceFromHandle converts our handle from C back to a Go IncomingInterface.
func getIncomingInterfaceFromHandle(handle unsafe.Pointer) IncomingInterface {
	iiIndex := int(*(*C.int)(handle))
	return incomingInterfaces[iiIndex]
}

// Output methods

// OutputInit initializes an output connection.  Usually you would use an OutputConnection rather than call this
// function directly.
func OutputInit(connection *ConnectionInterfaceStruct, name string, recordInfo recordinfo.RecordInfo) error {
	recordInfoXml, err := recordInfo.ToXml(name)
	if err != nil {
		return fmt.Errorf(`error intializing output connection '%v': %v`, name, err.Error())
	}
	cRecordInfoXml, err := convert_strings.StringToWideC(recordInfoXml)
	if err != nil {
		return fmt.Errorf(`error initializing output connection '%v': %v`, name, err.Error())
	}
	result := C.c_outputInit(connection.connection, cRecordInfoXml)
	if result == C.long(0) {
		return fmt.Errorf(`error calling pII_InitOutput on output '%v'`, name)
	}
	return nil
}

// OutputPushRecord pushes a record to an output connection.  Usually you would use an OutputConnection rather
// than call this function directly.
func OutputPushRecord(connection *ConnectionInterfaceStruct, record unsafe.Pointer) error {
	result := C.c_outputPushRecord(connection.connection, record)
	if result == C.long(0) {
		return fmt.Errorf(`error calling pII_PushRecord`)
	}
	return nil
}

// OutputClose closes an output connection.  Usually you would use an OutputConnection rather than call this
// function directly.
func OutputClose(connection *ConnectionInterfaceStruct) {
	C.c_outputClose(connection.connection)
}

// OutputUpdateProgress updates the progress of an output connection.  Usually you would use an OutputConnection
// rather than call this function directly.
func OutputUpdateProgress(connection *ConnectionInterfaceStruct, percent float64) {
	C.c_outputUpdateProgress(connection.connection, C.double(percent))
}

// Engine methods

// OutputMessage sends a message to the Alteryx engine.  You would typically use this to send Info, Warnings, or Errors
// back to Designer.
func OutputMessage(toolId int, status MessageStatus, message string) {
	if engine == nil {
		return
	}
	cMessage, err := convert_strings.StringToWideC(message)
	if err != nil {
		return
	}
	if cMessage == nil {
		return
	}

	C.callEngineOutputMessage(C.int(toolId), C.int(status), cMessage)
}

// OutputToolProgress notifies the engine of this tool's progress.  Calling this function also updates the progress
// in Designer when a workflow is running.
func OutputToolProgress(toolId int, percent float64) bool {
	if engine == nil {
		return true
	}
	if C.callEngineOutputToolProgress(C.int(toolId), C.double(percent)) == C.long(1) {
		return true
	}
	return false
}

// BrowseEverywhereReserveAnchor reserves a BrowseEverywhere anchor.  Rather than calling this yourself, use an
// OutputConnection which handles this for you.
func BrowseEverywhereReserveAnchor(toolId int) uint {
	if engine == nil {
		return 0
	}
	anchorId := C.callEngineBrowseEverywhereReserveAnchor(C.int(toolId))
	return uint(anchorId)
}

// BrowseEverywhereGetII returns an IncomingInterface for a downstream BrowseEverywhere anchor.  Rather than calling
// this yourself, use an OutputConnection which handles this for you.
func BrowseEverywhereGetII(browseEverywhereReservationId uint, toolId int, name string) *ConnectionInterfaceStruct {
	cName, _ := convert_strings.StringToWideC(name)
	ii := C.callEngineBrowseEverywhereGetII(C.unsigned(browseEverywhereReservationId), C.int(toolId), cName)
	return &ConnectionInterfaceStruct{connection: ii}
}

// CreateTempFileName will return a path to a temporary file that will be cleaned up by Alteryx when the workflow finishes.
func CreateTempFileName(ext string) string {
	cExt, _ := convert_strings.StringToWideC(ext)
	cFileName := C.callEngineCreateTempFileName(cExt)
	return convert_strings.WideCToString(cFileName)
}

// GetInitVar returns the value of the specified InitVar.
func GetInitVar(toolId int, initVar InitVar) string {
	cInitVar, _ := convert_strings.StringToWideC(string(initVar))

	if initVar == RunMode || initVar == ActionApplies {
		cInitVarValue := C.callEngineGetInitVar2(C.int(toolId), cInitVar)
		return convert_strings.WideCToString(cInitVarValue)
	}
	cInitVarValue := C.callEngineGetInitVar(cInitVar)
	return convert_strings.WideCToString(cInitVarValue)
}

// This section will get removed, eventually

func printLogf(message string, args ...interface{}) {
	file, _ := os.OpenFile("C:\\temp\\output.txt", os.O_WRONLY|os.O_APPEND, 0644)
	defer file.Close()
	_, _ = file.WriteString(fmt.Sprintf(time.Now().String()+": "+message+"\n", args...))
}
