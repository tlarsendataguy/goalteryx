package output_connection

import "C"
import (
	"fmt"
	"goalteryx/api"
	"goalteryx/recordinfo"
	"time"
	"unsafe"
)

type OutputConnection interface {
	Add(connection *api.ConnectionInterfaceStruct)
	Init(info recordinfo.RecordInfo) error
	PushRecord(record unsafe.Pointer)
	UpdateProgress(percent float64)
	Close()
}

func New(toolId int, name string) OutputConnection {
	browseEverywhereAnchorId := api.BrowseEverywhereReserveAnchor(toolId)
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
	recordCount              int
	recordSize               int
	recordInfo               recordinfo.RecordInfo
	lastCountOutput          time.Time
}

func (output *outputConnection) Add(connection *api.ConnectionInterfaceStruct) {
	output.connections = append(output.connections, connection)
}

func (output *outputConnection) Init(info recordinfo.RecordInfo) error {
	output.recordInfo = info
	output.lastCountOutput = time.Now()
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
		err := api.OutputInit(connection, output.name, info)
		if err != nil {
			output.connections = append(output.connections[:index], output.connections[index+1:]...)
			errs++
			continue
		}
		api.OutputUpdateProgress(connection, 0)
	}
	if errs > 0 {
		return fmt.Errorf(`%v connection(s) failed to initialize`, errs)
	}
	return nil
}

func (output *outputConnection) PushRecord(record unsafe.Pointer) {
	output.recordCount++
	output.recordSize += output.recordInfo.TotalSize(record)
	output.OutputRecordCount(false)

	for index, connection := range output.connections {
		err := api.OutputPushRecord(connection, record)
		if err != nil {
			output.finishedConnections = append(output.finishedConnections, output.connections[index])
			output.connections = append(output.connections[:index], output.connections[index+1:]...)
		}
	}
}

func (output *outputConnection) OutputRecordCount(final bool) {
	if output.recordCount < 256 || output.recordCount%256 == 0 || final {
		now := time.Now()
		if now.Sub(output.lastCountOutput).Seconds() > 10 || final {
			api.OutputMessage(output.toolId, api.RecordCountString, fmt.Sprintf(`%v|%v|%v`, output.name, output.recordCount, output.recordSize))
			output.lastCountOutput = now
		}
	}
}

func (output *outputConnection) UpdateProgress(percent float64) {
	for _, connection := range output.connections {
		api.OutputUpdateProgress(connection, percent)
	}
}

func (output *outputConnection) Close() {
	output.OutputRecordCount(true)
	for _, connection := range append(output.connections, output.finishedConnections...) {
		api.OutputUpdateProgress(connection, 1)
		api.OutputClose(connection)
	}
	api.OutputMessage(output.toolId, api.Complete, ``)
}
