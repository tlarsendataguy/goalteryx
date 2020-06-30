// Package output_connection provides a helper interface that manages output connections.
package output_connection

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"fmt"
	"github.com/tlarsen7572/goalteryx/api"
	"github.com/tlarsen7572/goalteryx/recordblob"
	"github.com/tlarsen7572/goalteryx/recordinfo"
	"reflect"
	"time"
	"unsafe"
)

// OutputConnection defines the lifecycle methods needed to manage output connections.
type OutputConnection interface {
	Add(connection *api.ConnectionInterfaceStruct)
	Init(info recordinfo.RecordInfo) error
	PushRecord(record recordblob.RecordBlob)
	UpdateProgress(percent float64)
	Close()
}

// New generates a new OutputConnection for the specified tool ID and connection name.  It also reserves a
// BrowseEverywhere anchor to allow for that capability in Designer.
func New(toolId int, name string, bufferSize int) OutputConnection {
	browseEverywhereAnchorId := api.BrowseEverywhereReserveAnchor(toolId)

	output := &outputConnection{
		toolId:                   toolId,
		name:                     name,
		connections:              []*api.ConnectionInterfaceStruct{},
		finishedConnections:      []*api.ConnectionInterfaceStruct{},
		browseEverywhereAnchorId: browseEverywhereAnchorId,
		bufferSize:               bufferSize,
		buffer:                   make([]unsafe.Pointer, bufferSize),
		blobSizes:                make([]int, bufferSize),
	}
	if bufferSize <= 0 {
		output.pushRecordCallback = output.pushSingleRecord
	} else {
		output.pushRecordCallback = output.pushRecordBuffer
	}

	return output
}

// outputConnection is the struct which implements OutputConnection.
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
	pushRecordCallback       func(record recordblob.RecordBlob)
	bufferSize               int
	currentBufferIndex       int
	buffer                   []unsafe.Pointer
	blobSizes                []int
}

// Add adds a connection to the list of connections.
func (output *outputConnection) Add(connection *api.ConnectionInterfaceStruct) {
	output.connections = append(output.connections, connection)
}

// Init initializes all output connection with the specified RecordInfo.  Any connections that fail to initialize
// are removed and no longer managed.  BrowseEverywhere connections are also added using the anchor ID obtained in Init.
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

// PushRecord pushes a record blob to all output connections.  Any output connections that return an error
// are removed from the connections list and added to the finished connections list.
func (output *outputConnection) PushRecord(record recordblob.RecordBlob) {
	output.pushRecordCallback(record)
}

func (output *outputConnection) pushSingleRecord(record recordblob.RecordBlob) {
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

func (output *outputConnection) pushRecordBuffer(record recordblob.RecordBlob) {
	output.recordCount++
	size := output.recordInfo.TotalSize(record)
	output.recordSize += size
	output.OutputRecordCount(false)

	if size > output.blobSizes[output.currentBufferIndex] {
		if output.recordCount > 1 {
			C.free(output.buffer[output.currentBufferIndex])
		}
		output.buffer[output.currentBufferIndex] = C.malloc(C.ulonglong(size))
	}

	output.blobSizes[output.currentBufferIndex] = size

	var from []byte
	fromHeader := (*reflect.SliceHeader)(unsafe.Pointer(&from))
	fromHeader.Data = uintptr(record.Blob())
	fromHeader.Len = size
	fromHeader.Cap = size

	var to []byte
	toHeader := (*reflect.SliceHeader)(unsafe.Pointer(&to))
	toHeader.Data = uintptr(output.buffer[output.currentBufferIndex])
	toHeader.Len = size
	toHeader.Cap = size

	copy(to, from)

	output.currentBufferIndex += 1
	if output.currentBufferIndex == output.bufferSize {
		errs := api.OutputPushBuffer(output.connections, output.buffer, output.bufferSize)
		for index, err := range errs {
			if err != nil {
				output.finishedConnections = append(output.finishedConnections, output.connections[index])
				output.connections = append(output.connections[:index], output.connections[index+1:]...)
			}
		}
		output.currentBufferIndex = 0
	}
}

// OutputRecordCount updates the engine with the number of records sent downstream.
func (output *outputConnection) OutputRecordCount(final bool) {
	if output.recordCount < 256 || output.recordCount%256 == 0 || final {
		now := time.Now()
		if now.Sub(output.lastCountOutput).Seconds() > 10 || final {
			api.OutputMessage(output.toolId, api.RecordCountString, fmt.Sprintf(`%v|%v|%v`, output.name, output.recordCount, output.recordSize))
			output.lastCountOutput = now
		}
	}
}

// UpdateProgress calls UpdateProgress on all connections.
func (output *outputConnection) UpdateProgress(percent float64) {
	for _, connection := range output.connections {
		api.OutputUpdateProgress(connection, percent)
	}
}

// Close closes all open and finished connections, outputs the final record count, and tells the engine that this
// tool is complete.
func (output *outputConnection) Close() {
	if output.bufferSize > 0 {
		api.OutputPushBuffer(output.connections, output.buffer, output.currentBufferIndex)

		freeSize := output.bufferSize
		if output.recordCount < freeSize {
			freeSize = output.recordCount
		}
		for index := range output.buffer {
			C.free(output.buffer[index])
		}
	}

	output.OutputRecordCount(true)
	for _, connection := range append(output.connections, output.finishedConnections...) {
		api.OutputUpdateProgress(connection, 1)
		api.OutputClose(connection)
	}

	api.OutputMessage(output.toolId, api.Complete, ``)
}
