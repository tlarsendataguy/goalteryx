package main

import (
	"goalteryx/api"
	"goalteryx/output_connection"
	"goalteryx/presort"
	"goalteryx/recordinfo"
)

type PluginInput struct {
	ToolId int
	Output output_connection.OutputConnection
}

func (plugin *PluginInput) Init(toolId int, config string) bool {
	plugin.ToolId = toolId
	plugin.Output = output_connection.New(toolId, `Output`)
	return true
}

func (plugin *PluginInput) PushAllRecords(recordLimit int) bool {
	generator := recordinfo.NewGenerator()
	field := generator.AddInt64Field(`RecordCount`, `Go Input`)
	info := generator.GenerateRecordInfo()
	_ = plugin.Output.Init(info)

	if api.GetInitVar(plugin.ToolId, api.UpdateOnly) == `True` {
		return true
	}

	for i := 0; i < 100000; i++ {
		err := info.SetIntField(field, i)
		if err != nil {
			api.OutputMessage(plugin.ToolId, api.Error, err.Error())
			return false
		}
		record, _ := info.GenerateRecord()
		plugin.Output.PushRecord(record)
		if i%10000 == 0 {
			percent := float64(i) / 100000.0
			api.OutputToolProgress(plugin.ToolId, percent)
			plugin.Output.UpdateProgress(percent)
		}
	}

	api.OutputToolProgress(plugin.ToolId, 1.0)
	plugin.Output.UpdateProgress(1.0)
	api.OutputMessage(plugin.ToolId, api.Complete, ``)
	plugin.Output.Close()
	return true
}

func (plugin *PluginInput) Close(hasErrors bool) {
}

func (plugin *PluginInput) AddIncomingConnection(connectionType string, connectionName string) (api.IncomingInterface, *presort.PresortInfo) {
	return nil, nil
}

func (plugin *PluginInput) AddOutgoingConnection(connectionName string, connectionInterface *api.ConnectionInterfaceStruct) bool {
	plugin.Output.Add(connectionInterface)
	return true
}

func (plugin *PluginInput) GetToolId() int {
	return plugin.ToolId
}
