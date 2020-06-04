package api

/*
#include "plugins.h"
*/
import "C"
import (
	"fmt"
	"github.com/mattn/go-pointer"
	"goalteryx/convert_strings"
	"goalteryx/presort"
	"goalteryx/recordinfo"
	"os"
	"time"
	"unsafe"
)

var engine *C.struct_EngineInterface

var incomingInterfaces = make([]IncomingInterface, 1000)

type Plugin interface {
	Init(toolId int, config string) bool
	PushAllRecords(recordLimit int) bool
	Close(hasErrors bool)
	AddIncomingConnection(connectionType string, connectionName string) (IncomingInterface, *presort.PresortInfo)
	AddOutgoingConnection(connectionName string, connectionInterface *ConnectionInterfaceStruct) bool
	GetToolId() int
}

type IncomingInterface interface {
	Init(recordInfoIn string) bool
	PushRecord(record unsafe.Pointer) bool
	UpdateProgress(percent float64)
	Close()
}

type ConnectionInterfaceStruct struct {
	connection *C.struct_IncomingConnectionInterface
}

func NewConnectionInterfaceStruct(incomingInterface IncomingInterface) *ConnectionInterfaceStruct {
	iiIndexHandle := C.getIiIndex()
	iiIndex := int(*(*C.int)(iiIndexHandle))
	incomingInterfaces[iiIndex] = incomingInterface
	var ii *C.struct_IncomingConnectionInterface = C.newIi(iiIndexHandle)
	return &ConnectionInterfaceStruct{connection: ii}
}

// Plugin methods

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

func getPlugin(plugin Plugin) unsafe.Pointer {
	return pointer.Save(plugin)
}

//export go_piPushAllRecords
func go_piPushAllRecords(handle unsafe.Pointer, recordLimit C.__int64) C.long {
	alteryxPlugin := pointer.Restore(handle).(Plugin)
	if alteryxPlugin.PushAllRecords(int(recordLimit)) {
		return C.long(1)
	}
	return C.long(0)
}

//export go_piClose
func go_piClose(handle unsafe.Pointer, hasErrors C.bool) {
	alteryxPlugin := pointer.Restore(handle).(Plugin)
	alteryxPlugin.Close(bool(hasErrors))
}

//export go_piAddIncomingConnection
func go_piAddIncomingConnection(handle unsafe.Pointer, connectionType unsafe.Pointer, connectionName unsafe.Pointer) unsafe.Pointer {
	alteryxPlugin := pointer.Restore(handle).(Plugin)
	goName := convert_strings.WideCToString(connectionName)
	goType := convert_strings.WideCToString(connectionType)
	goIncomingInterface, _ := alteryxPlugin.AddIncomingConnection(goType, goName)
	iiIndexHandle := C.getIiIndex()
	iiIndex := int(*(*C.int)(iiIndexHandle))
	incomingInterfaces[iiIndex] = goIncomingInterface
	return iiIndexHandle
}

//export go_piAddOutgoingConnection
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
func go_iiUpdateProgress(handle unsafe.Pointer, percent C.double) {
	incomingInterface := getIncomingInterfaceFromHandle(handle)
	incomingInterface.UpdateProgress(float64(percent))
}

//export go_iiClose
func go_iiClose(handle unsafe.Pointer) {
	incomingInterface := getIncomingInterfaceFromHandle(handle)
	incomingInterface.Close()
}

func getIncomingInterfaceFromHandle(handle unsafe.Pointer) IncomingInterface {
	iiIndex := int(*(*C.int)(handle))
	return incomingInterfaces[iiIndex]
}

// Output methods

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

func OutputPushRecord(connection *ConnectionInterfaceStruct, record unsafe.Pointer) error {
	result := C.c_outputPushRecord(connection.connection, record)
	if result == C.long(0) {
		return fmt.Errorf(`error calling pII_PushRecord`)
	}
	return nil
}

func OutputClose(connection *ConnectionInterfaceStruct) {
	C.c_outputClose(connection.connection)
}

func OutputUpdateProgress(connection *ConnectionInterfaceStruct, percent float64) {
	C.c_outputUpdateProgress(connection.connection, C.double(percent))
}

// Engine methods

func OutputMessage(toolId int, status MessageStatus, message string) {
	cMessage, err := convert_strings.StringToWideC(message)
	if err != nil {
		return
	}
	if cMessage == nil {
		return
	}

	C.callEngineOutputMessage(C.int(toolId), C.int(status), cMessage)
}

func OutputToolProgress(toolId int, percent float64) bool {
	if C.callEngineOutputToolProgress(C.int(toolId), C.double(percent)) == C.long(1) {
		return true
	}
	return false
}

func BrowseEverywhereReserveAnchor(toolId int) uint {
	//printLogf(`start reserving browse everywhere anchor ID`)
	anchorId := C.callEngineBrowseEverywhereReserveAnchor(C.int(toolId))
	//printLogf(`done reserving browse everywhere anchor ID`)
	//printLogf(`returned browse everywhere anchor ID: %v`, anchorId)
	return uint(anchorId)
}

func BrowseEverywhereGetII(browseEverywhereReservationId uint, toolId int, name string) *ConnectionInterfaceStruct {
	//printLogf(`start getting browse everywhere II`)
	cName, _ := convert_strings.StringToWideC(name)
	ii := C.callEngineBrowseEverywhereGetII(C.unsigned(browseEverywhereReservationId), C.int(toolId), cName)
	//printLogf(`done getting browse everywhere II`)
	return &ConnectionInterfaceStruct{connection: ii}
}

func CreateTempFileName(ext string) string {
	cExt, _ := convert_strings.StringToWideC(ext)
	cFileName := C.callEngineCreateTempFileName(cExt)
	return convert_strings.WideCToString(cFileName)
}

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
