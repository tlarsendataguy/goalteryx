package main

/*
#include "implementation.h"
*/
import "C"
import (
	"goalteryx/api"
	"goalteryx/output_connection"
	"goalteryx/recordinfo"
	"unsafe"
)

func main() {}

//export AlteryxGoPlugin
func AlteryxGoPlugin(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	myPlugin := &MyNewPlugin{
		Output1: output_connection.New(int(toolId), `Output1`),
	}
	return C.long(api.ConfigurePlugin(myPlugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}

type MyNewPlugin struct {
	ToolId  int
	Field   string
	Output1 output_connection.OutputConnection
}

type ConfigXml struct {
	Field string `xml:"Field"`
}

func (plugin *MyNewPlugin) Init(toolId int, config string) bool {
	plugin.ToolId = toolId
	return true
}

func (plugin *MyNewPlugin) PushAllRecords(recordLimit int) bool {
	return true
}

func (plugin *MyNewPlugin) Close(hasErrors bool) {

}

func (plugin *MyNewPlugin) AddIncomingConnection(connectionType string, connectionName string) api.IncomingInterface {
	return &MyPluginIncomingInterface{Parent: plugin}
}

func (plugin *MyNewPlugin) AddOutgoingConnection(connectionName string, connectionInterface *api.ConnectionInterfaceStruct) bool {
	plugin.Output1.Add(connectionInterface)
	return true
}

type MyPluginIncomingInterface struct {
	Parent *MyNewPlugin
	inInfo recordinfo.RecordInfo
}

func (ii *MyPluginIncomingInterface) Init(recordInfoIn string) bool {
	var err error
	ii.inInfo, err = recordinfo.FromXml(recordInfoIn)
	if err != nil {
		api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
		return false
	}

	err = ii.Parent.Output1.Init(ii.inInfo)
	if err != nil {
		api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
		return false
	}
	return true
}

func (ii *MyPluginIncomingInterface) PushRecord(record unsafe.Pointer) bool {
	for index := 0; index < ii.inInfo.NumFields(); index++ {
		field, err := ii.inInfo.GetFieldByIndex(index)
		if err != nil {
			api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
			return false
		}
		value, err := ii.inInfo.GetRawBytesFrom(field.Name, record)
		if err != nil {
			api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
			return false
		}
		err = ii.inInfo.SetFromRawBytes(field.Name, value)
		if err != nil {
			api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
			return false
		}
	}

	outputRecord, err := ii.inInfo.GenerateRecord()
	if err != nil {
		api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
		return false
	}
	ii.Parent.Output1.PushRecord(outputRecord)
	return true
}

func (ii *MyPluginIncomingInterface) UpdateProgress(percent float64) {

}

func (ii *MyPluginIncomingInterface) Close() {
	ii.Parent.Output1.Close()
}

func (ii *MyPluginIncomingInterface) Free() {

}
