package api

/*
#include "plugins.h"
*/
import "C"
import (
	"fmt"
	"github.com/mattn/go-pointer"
	"goalteryx/convert_strings"
	"goalteryx/recordinfo"
	"os"
	"time"
	"unsafe"
)

var engine *C.struct_EngineInterface

type MessageStatus int

var incomingInterfaces = make([]IncomingInterface, 100)

const (
	Info                          MessageStatus = 1
	TransientInfo                 MessageStatus = 0x40000000 | 1
	Warning                       MessageStatus = 2
	TransientWarning              MessageStatus = 0x40000000 | 2
	Error                         MessageStatus = 3
	Complete                      MessageStatus = 4
	FieldConversionError          MessageStatus = 5
	TransientFieldConversionError MessageStatus = 0x40000000 | 5
	UpdateOutputMetaInfoXml       MessageStatus = 10
	RecordCountString             MessageStatus = 50
	BrowseEverywhereFileName      MessageStatus = 70
)

type Plugin interface {
	Init(toolId int, config string) bool
	PushAllRecords(recordLimit int) bool
	Close(hasErrors bool)
	AddIncomingConnection(connectionType string, connectionName string) IncomingInterface
	AddOutgoingConnection(connectionName string, connectionInterface *ConnectionInterfaceStruct) bool
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
	var ii *C.struct_IncomingConnectionInterface = C.newIi()
	iiIndexHandle := C.getIiIndex()
	iiIndex := int(*(*C.int)(iiIndexHandle))
	incomingInterfaces[iiIndex] = incomingInterface
	ii.handle = iiIndexHandle
	ii.pII_Init = C.T_II_Init(C.iiInit)
	ii.pII_PushRecord = C.T_II_PushRecord(C.iiPushRecord)
	ii.pII_UpdateProgress = C.T_II_UpdateProgress(C.iiUpdateProgress)
	ii.pII_Close = C.T_II_Close(C.iiClose)
	ii.pII_Free = C.T_II_Free(C.iiFree)
	return &ConnectionInterfaceStruct{connection: ii}
}

func ConfigurePlugin(plugin Plugin, toolId int, pXmlProperties unsafe.Pointer, pEngineInterface unsafe.Pointer, r_pluginInterface unsafe.Pointer) int {
	config := convert_strings.WideCToString(pXmlProperties)
	engine = (*C.struct_EngineInterface)(pEngineInterface)
	if !plugin.Init(toolId, config) {
		return 0
	}

	pluginInterface := (*C.struct_PluginInterface)(r_pluginInterface)
	pluginInterface.handle = getPlugin(plugin)
	pluginInterface.pPI_PushAllRecords = C.T_PI_PushAllRecords(C.piPushAllRecords)
	pluginInterface.pPI_Close = C.T_PI_Close(C.piClose)
	pluginInterface.pPI_AddIncomingConnection = C.T_PI_AddIncomingConnection(C.piAddIncomingConnection)
	pluginInterface.pPI_AddOutgoingConnection = C.T_PI_AddOutgoingConnection(C.piAddOutgoingConnection)
	return 1
}

//export piPushAllRecords
func piPushAllRecords(handle unsafe.Pointer, recordLimit C.__int64) C.long {
	alteryxPlugin := pointer.Restore(handle).(Plugin)
	if alteryxPlugin.PushAllRecords(int(recordLimit)) {
		return C.long(1)
	}
	return C.long(0)
}

//export piClose
func piClose(handle unsafe.Pointer, hasErrors C.bool) {
	alteryxPlugin := pointer.Restore(handle).(Plugin)
	alteryxPlugin.Close(bool(hasErrors))
}

//export piAddIncomingConnection
func piAddIncomingConnection(handle unsafe.Pointer, connectionType unsafe.Pointer, connectionName unsafe.Pointer, incomingInterface *C.struct_IncomingConnectionInterface) C.long {
	alteryxPlugin := pointer.Restore(handle).(Plugin)
	goName := convert_strings.WideCToString(connectionName)
	goType := convert_strings.WideCToString(connectionType)
	goIncomingInterface := alteryxPlugin.AddIncomingConnection(goType, goName)
	iiIndexHandle := C.getIiIndex()
	iiIndex := int(*(*C.int)(iiIndexHandle))
	incomingInterfaces[iiIndex] = goIncomingInterface
	incomingInterface.handle = iiIndexHandle
	incomingInterface.pII_Init = C.T_II_Init(C.iiInit)
	incomingInterface.pII_PushRecord = C.T_II_PushRecord(C.iiPushRecord)
	incomingInterface.pII_UpdateProgress = C.T_II_UpdateProgress(C.iiUpdateProgress)
	incomingInterface.pII_Close = C.T_II_Close(C.iiClose)
	incomingInterface.pII_Free = C.T_II_Free(C.iiFree)
	return C.long(1)
}

//export piAddOutgoingConnection
func piAddOutgoingConnection(handle unsafe.Pointer, connectionName unsafe.Pointer, incomingConnection *C.struct_IncomingConnectionInterface) C.long {
	alteryxPlugin := pointer.Restore(handle).(Plugin)
	goName := convert_strings.WideCToString(connectionName)
	connectionInterface := &ConnectionInterfaceStruct{connection: incomingConnection}
	if alteryxPlugin.AddOutgoingConnection(goName, connectionInterface) {
		return C.long(1)
	}
	return C.long(0)
}

//export iiInit
func iiInit(handle unsafe.Pointer, recordInfoIn unsafe.Pointer) C.long {
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

//export pushRecordCache
func pushRecordCache(handle unsafe.Pointer, cache unsafe.Pointer, cacheSize C.int) C.long {
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

//export iiUpdateProgress
func iiUpdateProgress(handle unsafe.Pointer, percent C.double) {
	incomingInterface := getIncomingInterfaceFromHandle(handle)
	incomingInterface.UpdateProgress(float64(percent))
}

//export iiClose
func iiClose(handle unsafe.Pointer) {
	C.closeRecordCache(handle)
	incomingInterface := getIncomingInterfaceFromHandle(handle)
	incomingInterface.Close()
}

//export iiFree
func iiFree(handle unsafe.Pointer) {
	C.freeRecordCache(handle)
}

//export getPlugin
func getPlugin(plugin Plugin) unsafe.Pointer {
	return pointer.Save(plugin)
}

func getIncomingInterfaceFromHandle(handle unsafe.Pointer) IncomingInterface {
	iiIndex := int(*(*C.int)(handle))
	return incomingInterfaces[iiIndex]
}

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

	C.callEngineOutputMessage(engine, C.int(toolId), C.int(status), cMessage)
}

func OutputToolProgress(toolId int, percent float64) bool {
	if C.callEngineOutputToolProgress(engine, C.int(toolId), C.double(percent)) == C.long(1) {
		return true
	}
	return false
}

func BrowseEverywhereReserveAnchor(toolId int) uint {
	//printLogf(`start reserving browse everywhere anchor ID`)
	anchorId := C.callEngineBrowseEverywhereReserveAnchor(engine, C.int(toolId))
	//printLogf(`done reserving browse everywhere anchor ID`)
	//printLogf(`returned browse everywhere anchor ID: %v`, anchorId)
	return uint(anchorId)
}

func BrowseEverywhereGetII(browseEverywhereReservationId uint, toolId int, name string) *ConnectionInterfaceStruct {
	//printLogf(`start getting browse everywhere II`)
	cName, _ := convert_strings.StringToWideC(name)
	ii := C.callEngineBrowseEverywhereGetII(engine, C.unsigned(browseEverywhereReservationId), C.int(toolId), cName)
	//printLogf(`done getting browse everywhere II`)
	return &ConnectionInterfaceStruct{connection: ii}
}

func InitOutput(connection *ConnectionInterfaceStruct, name string, recordInfo recordinfo.RecordInfo) error {
	recordInfoXml, err := recordInfo.ToXml(name)
	if err != nil {
		return fmt.Errorf(`error intializing output connection '%v': %v`, name, err.Error())
	}
	cRecordInfoXml, err := convert_strings.StringToWideC(recordInfoXml)
	if err != nil {
		return fmt.Errorf(`error initializing output connection '%v': %v`, name, err.Error())
	}
	result := C.callInitOutput(connection.connection, cRecordInfoXml)
	if result == C.long(0) {
		return fmt.Errorf(`error calling pII_InitOutput on output '%v'`, name)
	}
	return nil
}

func PushRecord(connection *ConnectionInterfaceStruct, record unsafe.Pointer) error {
	result := C.callPushRecord(connection.connection, record)
	if result == C.long(0) {
		return fmt.Errorf(`error calling pII_PushRecord`)
	}
	return nil
}

func CloseOutputConnection(connection *ConnectionInterfaceStruct) {
	C.callCloseOutput(connection.connection)
}

func UpdateOutputConnectionProgress(connection *ConnectionInterfaceStruct, percent float64) {
	C.updateProgress(connection.connection, C.double(percent))
}

func CreateTempFileName(ext string) string {
	cExt, _ := convert_strings.StringToWideC(ext)
	cFileName := C.callEngineCreateTempFileName(engine, cExt)
	return convert_strings.WideCToString(cFileName)
}

func printLogf(message string, args ...interface{}) {
	file, _ := os.OpenFile("C:\\temp\\output.txt", os.O_WRONLY|os.O_APPEND, 0644)
	defer file.Close()
	_, _ = file.WriteString(fmt.Sprintf(time.Now().String()+": "+message+"\n", args...))
}
