package main

import (
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/output_connection"
	"github.com/tlarsen7572/goalteryx/presort"
	"github.com/tlarsen7572/goalteryx/recordblob"
	"github.com/tlarsen7572/goalteryx/recordcopier"
	"github.com/tlarsen7572/goalteryx/recordinfo"
)

type PluginNoCache struct {
	ToolId int
	Output output_connection.OutputConnection
}

func (plugin *PluginNoCache) Init(toolId int, config string) bool {
	plugin.ToolId = toolId
	plugin.Output = output_connection.New(toolId, `Output`, 0)
	return true
}

func (plugin *PluginNoCache) PushAllRecords(recordLimit int) bool {
	return false
}

func (plugin *PluginNoCache) Close(hasErrors bool) {
}

func (plugin *PluginNoCache) AddIncomingConnection(connectionType string, connectionName string) (api.IncomingInterface, *presort.PresortInfo) {
	return &PluginNoCacheIncomingInterface{Parent: plugin}, nil
}

func (plugin *PluginNoCache) AddOutgoingConnection(connectionName string, connectionInterface *api.ConnectionInterfaceStruct) bool {
	plugin.Output.Add(connectionInterface)
	return true
}

func (plugin *PluginNoCache) GetToolId() int {
	return plugin.ToolId
}

type PluginNoCacheIncomingInterface struct {
	Parent *PluginNoCache
	inInfo recordinfo.RecordInfo
	copier *recordcopier.RecordCopier
}

func (ii *PluginNoCacheIncomingInterface) Init(recordInfoIn string) bool {
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
	ii.copier, _ = recordcopier.New(ii.inInfo, ii.inInfo, indexMaps)
	return true
}

func (ii *PluginNoCacheIncomingInterface) PushRecord(record recordblob.RecordBlob) bool {
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

func (ii *PluginNoCacheIncomingInterface) UpdateProgress(percent float64) {
	api.OutputToolProgress(ii.Parent.ToolId, percent)
	ii.Parent.Output.UpdateProgress(percent)
}

func (ii *PluginNoCacheIncomingInterface) Close() {
	ii.Parent.Output.Close()
}

func (ii *PluginNoCacheIncomingInterface) CacheSize() int {
	return 0
}
