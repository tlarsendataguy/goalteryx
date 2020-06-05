package main

/*
#include "implementation.h"
*/
import "C"
import (
	"goalteryx/api"
	"goalteryx/output_connection"
	"goalteryx/presort"
	"goalteryx/recordcopier"
	"goalteryx/recordinfo"
	"unsafe"
)

func main() {}

//export AlteryxGoPlugin
func AlteryxGoPlugin(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	myPlugin := &MyNewPlugin{}
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
	plugin.Output1 = output_connection.New(toolId, `Output1`)
	return true
}

func (plugin *MyNewPlugin) PushAllRecords(recordLimit int) bool {
	return true
}

func (plugin *MyNewPlugin) Close(hasErrors bool) {
}

func (plugin *MyNewPlugin) AddIncomingConnection(connectionType string, connectionName string) (api.IncomingInterface, *presort.PresortInfo) {
	return &MyPluginIncomingInterface{Parent: plugin}, nil
}

func (plugin *MyNewPlugin) AddOutgoingConnection(connectionName string, connectionInterface *api.ConnectionInterfaceStruct) bool {
	plugin.Output1.Add(connectionInterface)
	return true
}

func (plugin *MyNewPlugin) GetToolId() int {
	return plugin.ToolId
}

type MyPluginIncomingInterface struct {
	Parent *MyNewPlugin
	inInfo recordinfo.RecordInfo
	copier *recordcopier.RecordCopier
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

	indexMaps := make([]recordcopier.IndexMap, ii.inInfo.NumFields())
	for index := range indexMaps {
		indexMaps[index] = recordcopier.IndexMap{
			DestinationIndex: index,
			SourceIndex:      index,
		}
	}
	ii.copier, _ = recordcopier.New(ii.inInfo, ii.inInfo, indexMaps)
	return true
}

func (ii *MyPluginIncomingInterface) PushRecord(record unsafe.Pointer) bool {
	err := ii.copier.Copy(record)
	if err != nil {
		api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
		return false
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
	api.OutputToolProgress(ii.Parent.ToolId, percent)
	ii.Parent.Output1.UpdateProgress(percent)
}

func (ii *MyPluginIncomingInterface) Close() {
	ii.Parent.Output1.Close()
}
