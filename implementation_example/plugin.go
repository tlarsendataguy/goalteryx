package main

import (
	"fmt"
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/output_connection"
	"github.com/tlarsen7572/goalteryx/presort"
	"github.com/tlarsen7572/goalteryx/recordblob"
	"github.com/tlarsen7572/goalteryx/recordcopier"
	"github.com/tlarsen7572/goalteryx/recordinfo"
)

type Plugin struct {
	ToolId int
	Output output_connection.OutputConnection
}

func (plugin *Plugin) Init(toolId int, config string) bool {
	plugin.ToolId = toolId
	plugin.Output = output_connection.New(toolId, `Output`, 10)
	return true
}

func (plugin *Plugin) PushAllRecords(recordLimit int) bool {
	return false
}

func (plugin *Plugin) Close(hasErrors bool) {
}

func (plugin *Plugin) AddIncomingConnection(connectionType string, connectionName string) (api.IncomingInterface, *presort.PresortInfo) {
	return &PluginIncomingInterface{Parent: plugin, records: 0}, nil
}

func (plugin *Plugin) AddOutgoingConnection(connectionName string, connectionInterface *api.ConnectionInterfaceStruct) bool {
	plugin.Output.Add(connectionInterface)
	return true
}

func (plugin *Plugin) GetToolId() int {
	return plugin.ToolId
}

type PluginIncomingInterface struct {
	Parent  *Plugin
	inInfo  recordinfo.RecordInfo
	copier  *recordcopier.RecordCopier
	records int
}

func (ii *PluginIncomingInterface) Init(recordInfoIn string) bool {
	api.OutputMessage(ii.Parent.ToolId, api.Info, recordInfoIn)
	var err error
	ii.inInfo, err = recordinfo.FromXml(recordInfoIn)
	if err != nil {
		api.OutputMessage(ii.Parent.ToolId, api.Error, err.Error())
		return false
	}

	err = ii.Parent.Output.Init(ii.inInfo)
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
	ii.copier, _ = recordcopier.New(ii.inInfo, ii.inInfo.GenerateRecordBlobReader(), indexMaps)
	return true
}

func (ii *PluginIncomingInterface) PushRecord(record recordblob.RecordBlob) bool {
	if ii.records == 0 {
		api.OutputMessage(ii.Parent.ToolId, api.Info, fmt.Sprintf("first record bytes:\r\n%v", *((*[7]byte)(record.Blob()))))
	}
	ii.records++
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
	ii.Parent.Output.PushRecord(outputRecord)
	return true
}

func (ii *PluginIncomingInterface) UpdateProgress(percent float64) {
	api.OutputToolProgress(ii.Parent.ToolId, percent)
	ii.Parent.Output.UpdateProgress(percent)
}

func (ii *PluginIncomingInterface) Close() {
	ii.Parent.Output.Close()
}

func (ii *PluginIncomingInterface) CacheSize() int {
	return 10
}
