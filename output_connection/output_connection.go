package output_connection

import "C"
import (
	"goalteryx/api"
	"goalteryx/recordinfo"
	"unsafe"
)

type OutputConnection interface {
	Add(connection *api.ConnectionInterfaceStruct)
	Init(info recordinfo.RecordInfo) error
	PushRecord(record unsafe.Pointer) error
	Close() error
}

func New(toolId int, name string) OutputConnection {
	browseEverywhereAnchorId := uint(0) // TODO fix this with api.BrowseEverywhereReserveAnchor(toolId)
	//browseEverywhereAnchorId = api.BrowseEverywhereReserveAnchor(toolId)
	return &outputConnection{
		toolId:                   toolId,
		name:                     name,
		connections:              []*api.ConnectionInterfaceStruct{},
		browseEverywhereAnchorId: browseEverywhereAnchorId,
	}
}

type outputConnection struct {
	toolId                   int
	name                     string
	connections              []*api.ConnectionInterfaceStruct
	browseEverywhereAnchorId uint
}

func (output *outputConnection) Add(connection *api.ConnectionInterfaceStruct) {
	output.connections = append(output.connections, connection)
}

func (output *outputConnection) Init(info recordinfo.RecordInfo) error {
	infoXml, err := info.ToXml(output.name)
	if err != nil {
		return err
	}
	api.OutputMessage(output.toolId, api.UpdateOutputMetaInfoXml, infoXml)
	if output.browseEverywhereAnchorId != 0 {
		ii := api.BrowseEverywhereGetII(output.browseEverywhereAnchorId, output.toolId, output.name)
		output.connections = append(output.connections, ii)
	}

	for _, connection := range output.connections {
		err := api.InitOutput(connection, output.name, info)
		if err != nil {
			return err
		}
	}
	return nil
}

func (output *outputConnection) PushRecord(record unsafe.Pointer) error {
	for _, connection := range output.connections {
		err := api.PushRecord(connection, record)
		if err != nil {
			return err
		}
	}
	return nil
}

func (output *outputConnection) Close() error {
	return nil
}
