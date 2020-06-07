# GoAlteryx

An unofficial SDK for building custom Alteryx tools with Go.

## Why a Go SDK?

With the announced deprecation of the .NET SDK, a gap formed between the C/C++ and Python SDKs.  C/C++ are low-level languages requiring great care and expertise to ensure proper memory management.  Python is very approachable but is slower.  I wanted to build tools with a middle-ground language having decent performance and simplified memory management.  Go fit the bill and is my favorite language to code in.

## Installation

Install goalteryx using Go modules: `go get github.com/tlarsen7572/goalteryx`

## Building your custom tools

You should specify the output DLL file and make sure `-buildmode` is set to `c-shared`.  For reference, the following command is used to build the included example tools:

```
go build -o "C:\Program Files\Alteryx\bin\Plugins\goalteryx.dll" -buildmode=c-shared goalteryx/implementation_example
```

I build directly to the Plugins folder in the Alteryx installation folder of my dev environment.  This allows me to rebuild my tools and run them directly in Alteryx without additional copying.  You do not need to close and restart Alteryx when you rebuild a DLL.  The next time you run a workflow with your custom tool, the new DLL will be used.  It should go without saying that you should not do this in production.

## Usage

Several examples are provided in the implementation_example folder.  Please refer to those examples for comprehensive, if simple, examples of how to use the goalteryx SDK.  If you have experience with the Python or C++ SDKs, goalteryx should be familiar to you.

The entry point of your custom tool must be defined in a C header file.  While I was hoping to avoid requiring the developer to have to touch the C layer of these custom tools, this was the only way I could find to allow multiple tool engines to be embedded in a single DLL.  The C layer you will have to manipulate is small and can be copied from this readme or the included examples.

#### implementation.h
```cgo
long __declspec(dllexport) PluginEntry(int nToolID,
   	void * pXmlProperties,
   	void *pEngineInterface,
   	void *r_pluginInterface);
```

The Go implementation of the entry point should be defined in a mirror Go file:

#### implementation.go
```go
package main

/*
#include "implementation.h"
*/
import "C"
import (
	"github.com/tlarsen7572/goalteryx/api"
	"unsafe"
)

func main() {}

//export PluginEntry
func PluginEntry(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &Plugin{}
	return C.long(api.ConfigurePlugin(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}
```

The Go implementation file must import the C package and include the header file defined earlier.  An empty `main()` function must also be present.  Finally, the entry point (PluginEntry, in this example) must be defined.  The `//export` comment preceding the entry point is required.  The entry point itself must instantiate a struct implementing the `api.Plugin` interface.  Call `api.ConfigurePlugin()` to connect the struct to the Alteryx lifecycle events.

#### Plugin Interface

The Plugin object must implement the following methods:

```go
Init(toolId int, config string) bool
```

Init is called when it is time to configure the plugin using the provided XML configuration string.

```go
PushAllRecords(recordLimit int) bool
```

PushAllRecords is called when our tool is an input tool and its time to start sending data to downstream tools.

```go
Close(hasErrors bool)
```

Close is called when the upstream tools are done sending our plugin data.

```go
AddIncomingConnection(connectionType string, connectionName string) (IncomingInterface, *presort.PresortInfo)
```

AddIncomingConnection is called when an upstream tool is being connected to our tool.  We should return a struct that implements the `api.IncomingInterface` interface.  We also have a chance to tell the engine that the incoming data should be presorted by returning a PresortInfo object.  If we do not wish to presort incoming data, return a nil PresortInfo.

```go
AddOutgoingConnection(connectionName string, connectionInterface *ConnectionInterfaceStruct) bool
```

AddOutgoingConnection is called when our tool is being connected to a downstream tool.  It is best to use the OutputConnection helper in the output_connection package when working with outbound connections.

```go
GetToolId() int
```

GetToolId should return the toolId provided by Init.

#### IncomingInterface Interface

The IncomingInterface object must implement the following methods:

```go
Init(recordInfoIn string) bool
```

Init is called when an upstream tool sends us its outgoing RecordInfo.  Use `recordinfo.RecordInfo` and `recordinfo.Generator` to manage record structures and generate outgoing record blobs.

```go
PushRecord(record unsafe.Pointer) bool
```

PushRecord is called when an upstream tool is pushing a record blob to our tool.

```go
UpdateProgress(percent float64)
```

UpdateProgress is called when an upstream tool is updating our tool with its progress.

```go
Close()
```

Close is called when an upstream tool is finished sending us data.

#### recordinfo.RecordInfo and recordinfo.Generator

Use `recordinfo.Generator` to build up recordinfo objects.  Generators can be created empty:

```go
generator := NewGenerator()
```

or by instantiating a Generator from XML:

```go
generator, err := GeneratorFromXml(recordInfoXml)
```

Once the required record structure is build, create a RecordInfo by calling `GenerateRecordInfo()`:

```go
recordInfo := generator.GenerateRecordInfo()
```

RecordInfo's can also be created directly from an XML string:

```go
recordInfo, err := recordinfo.FromXml(recordInfoXml)
```

RecordInfo's can be used to obtain values from incoming record blobs:

```go
value, isNull, err := recordInfo.GetIntValueFrom(`FieldName`, recordBlob)
```

RecordInfo's are also used to set values for new record blobs:

```go
err := recordInfo.SetIntField(`FieldName`, 16)
```

There are getters and setters for all of the different types of fields.  Attempting to set or get values from a field type that does not match the getter/setter field type results in an error.  For example, you cannot get an int value from a string field.  The only universal getters/setters are the RawBytes functions.  While high performance, these should be avoided unless you are familiar with Alteryx's record blob structure.

Once all field values are set, a record blob can be created by calling `GenerateRecord`:

```go
record, err := recordInfo.GenerateRecord()
```

#### OutputConnection

Use OutputConnection to manage the minutia of outgoing connections.  Your plugin struct should have an OutputConnection member for each outgoing anchor:

```go
type Plugin struct {
	ToolId int
	Output output_connection.OutputConnection
}
```

During the plugin's `Init` method, instantiate the connection using the tool ID and output connection name:

```go
plugin.Output = output_connection.New(toolId, `Output`)
```

Then, pass all outgoing connections to this object during the `AddOutgoingConnection` method:

```go
func (plugin *Plugin) AddOutgoingConnection(connectionName string, connectionInterface *api.ConnectionInterfaceStruct) bool {
	plugin.Output.Add(connectionInterface)
	return true
}
```

Once the outgoing RecordInfo structure is known (usually during `plugin.PushAllRecords` or `incomingInterface.Init`) initialize the outgoing connection:

```go
err = plugin.Output.Init(recordInfo)
```

Push records to downstream tools by passing record blobs to the OutputConnection:

```go
plugin.Output.PushRecord(outputRecord)
```

Finally, when it is time to close down the outgoing data streams, call the `Close` method:

```go
plugin.Output.Close()
```

#### RecordCopier

RecordCopier is a helper object for copying incoming record blobs to outgoing record blobs.  Create a RecordCopier using `recordcopier.New`:

```go
copier, err := recordcopier.New(destinationRecordInfo, sourceRecordInfo, indexMaps)
```

indexMaps is a list of `recordcopier.IndexMap`, which is a simple struct containing a destination index and a source index.

When you are ready to copy an incoming record (usually during `IncomingInterface.PushRecord`), do that by calling Copy:

```go
err := copier.Copy(record)
```
