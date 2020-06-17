package main

import (
	"encoding/xml"
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/output_connection"
	"github.com/tlarsen7572/goalteryx/presort"
	"github.com/tlarsen7572/goalteryx/recordblob"
	"github.com/tlarsen7572/goalteryx/recordcopier"
	"github.com/tlarsen7572/goalteryx/recordinfo"
)

type PluginPresort struct {
	ToolId int
	Field  string
	Output output_connection.OutputConnection
}

func (plugin *PluginPresort) Init(toolId int, config string) bool {
	plugin.ToolId = toolId
	var c ConfigXml
	err := xml.Unmarshal([]byte(config), &c)
	if err != nil {
		api.OutputMessage(toolId, api.Error, err.Error())
		return false
	}
	if c.Field == `` {
		api.OutputMessage(toolId, api.Error, `'Field' was not provided.`)
		return false
	}
	plugin.Field = c.Field
	plugin.Output = output_connection.New(toolId, `Output`)
	return true
}

func (plugin *PluginPresort) PushAllRecords(recordLimit int) bool {
	return false
}

func (plugin *PluginPresort) Close(hasErrors bool) {
}

func (plugin *PluginPresort) AddIncomingConnection(connectionType string, connectionName string) (api.IncomingInterface, *presort.PresortInfo) {
	presortInfo := &presort.PresortInfo{
		SortInfo: []presort.SortInfo{
			{
				Field: plugin.Field,
				Order: presort.Desc,
			},
		},
		FieldFilterList: nil,
	}
	return &PluginPresortIncomingInterface{Parent: plugin}, presortInfo
}

func (plugin *PluginPresort) AddOutgoingConnection(connectionName string, connectionInterface *api.ConnectionInterfaceStruct) bool {
	plugin.Output.Add(connectionInterface)
	return true
}

func (plugin *PluginPresort) GetToolId() int {
	return plugin.ToolId
}

type PluginPresortIncomingInterface struct {
	Parent *PluginPresort
	inInfo recordinfo.RecordInfo
	copier *recordcopier.RecordCopier
}

func (ii *PluginPresortIncomingInterface) Init(recordInfoIn string) bool {
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

func (ii *PluginPresortIncomingInterface) PushRecord(record recordblob.RecordBlob) bool {
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

func (ii *PluginPresortIncomingInterface) UpdateProgress(percent float64) {
	api.OutputToolProgress(ii.Parent.ToolId, percent)
	ii.Parent.Output.UpdateProgress(percent)
}

func (ii *PluginPresortIncomingInterface) Close() {
	ii.Parent.Output.Close()
}

func (ii *PluginPresortIncomingInterface) CacheSize() int {
	return 10
}
