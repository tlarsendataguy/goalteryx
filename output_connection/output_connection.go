package output_connection

import "C"
import (
	"fmt"
	"goalteryx/api"
	"goalteryx/recordinfo"
	"unsafe"
)

type OutputConnection interface {
	Add(connection *api.ConnectionInterfaceStruct)
	Init(info recordinfo.RecordInfo) error
	PushRecord(record unsafe.Pointer)
	Close()
}

func New(toolId int, name string) OutputConnection {
	browseEverywhereAnchorId := uint(0) // TODO fix this with api.BrowseEverywhereReserveAnchor(toolId)
	//browseEverywhereAnchorId = api.BrowseEverywhereReserveAnchor(toolId)
	return &outputConnection{
		toolId:                   toolId,
		name:                     name,
		connections:              []*api.ConnectionInterfaceStruct{},
		finishedConnections:      []*api.ConnectionInterfaceStruct{},
		browseEverywhereAnchorId: browseEverywhereAnchorId,
	}
}

type outputConnection struct {
	toolId                   int
	name                     string
	connections              []*api.ConnectionInterfaceStruct
	finishedConnections      []*api.ConnectionInterfaceStruct
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

	errs := 0
	for index, connection := range output.connections {
		err := api.InitOutput(connection, output.name, info)
		if err != nil {
			output.connections = append(output.connections[:index], output.connections[index+1:]...)
			errs++
		}
	}
	if errs > 0 {
		return fmt.Errorf(`%v connection(s) failed to initialize`, errs)
	}
	return nil
}

func (output *outputConnection) PushRecord(record unsafe.Pointer) {
	for index, connection := range output.connections {
		err := api.PushRecord(connection, record)
		if err != nil {
			output.finishedConnections = append(output.finishedConnections, output.connections[index])
			output.connections = append(output.connections[:index], output.connections[index+1:]...)
		}
	}
}

func (output *outputConnection) Close() {
	for _, connection := range append(output.connections, output.finishedConnections...) {
		api.CloseOutput(connection)
	}
	api.OutputMessage(output.toolId, api.Complete, ``)
}
