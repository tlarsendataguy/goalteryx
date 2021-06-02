<img src="https://github.com/tlarsen7572/goalteryx/blob/main/goalteryx_icon_whiteBackground.png?raw=true" width="200">

# GoAlteryx

An unofficial SDK for building custom Alteryx tools with Go.

## Why a Go SDK?

With the announced deprecation of the .NET SDK, a gap formed between the C/C++ and Python SDKs.  C/C++ are low-level languages requiring great care and expertise to ensure proper memory management.  Python is very approachable but is slower.  I wanted to build tools with a middle-ground language having decent performance and simplified memory management.  Go fit the bill and is my favorite language to code with.

## Table of contents

1. [Prerequisites](#Prerequisites)
2. [Installation](#Installation)  
3. [Building your custom tools](#Building-your-custom-tools)  
4. [Sample tool](#Sample-tool)  
5. [Implementing the Plugin interface](#Implementing-the-Plugin-interface)  
6. [Registering your tool](#Registering-your-tool)  
7. [Using Provider](#Using-Provider)  
8. [Using OutputAnchor](#Using-OutputAnchor)  
9. [Using Io](#Using-Io)  
10. [Using Environment](#Using-Environment)  
11. [Using InputConnection](#Using-InputConnection)  
12. [RecordInfo](#RecordInfo)  
13. [Using RecordPacket](#Using-RecordPacket) 
14. [Testing your tools](#Testing-your-tools)
15. [Feature parity with the Python SDK](#Feature-parity-with-the-Python-SDK)

## Prerequisites

1. To use the Go SDK you must have Go installed on your machine.  You can download the latest version of Go [here](https://golang.org/dl/).
2. The Go SDK requires cgo, which means you must have a 64-bit C compiler on your system.  If you do not already have one, Mingw-w64 has been tested and works with the SDK.  Download from [here](http://mingw-w64.org/doku.php/download/mingw-builds) and install, making sure Mingw-w64 is added to PATH.
3. While not required, an IDE is highly recommended.  I prefer the [GoLand IDE](https://www.jetbrains.com/go/) from JetBrains.

## Installation

Install goalteryx using Go modules: `go get github.com/tlarsen7572/goalteryx`

## Building your custom tools

You should specify the output DLL file and make sure `-buildmode` is set to `c-shared`.  For reference, the following command is used to build the included example tools:

```
go build -o "C:\Program Files\Alteryx\bin\Plugins\goalteryx.dll" -buildmode=c-shared goalteryx/implementation_example
```

I build directly to the Plugins folder in the Alteryx installation folder of my dev environment.  This allows me to rebuild my tools and run them directly in Alteryx without additional copying.  You do not need to close and restart Alteryx when you rebuild a DLL.  The next time you run a workflow with your custom tool, the new DLL will be used.  It should go without saying that you should not do this in production.

[Back to table of contents](#Table-of-contents)

## Sample tool

The following 2 code files represent a basic tool in Alteryx that copies incoming records and pushes them through its output.

#### entry.go

```go
package main

import "C"
import (
	"github.com/tlarsen7572/goalteryx/sdk"
	"unsafe"
)

func main() {}

//export PluginEntry
func PluginEntry(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &Plugin{}
	return C.long(sdk.RegisterTool(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}
```

entry.go is used to register your plugin to the Alteryx engine.  See [plugin registration](#Registering-your-tool) for more info.

#### plugin.go

```go
package main

import (
	"fmt"
	"github.com/tlarsen7572/goalteryx/sdk"
)

type Plugin struct {
	provider sdk.Provider
	output   sdk.OutputAnchor
	outInfo  *sdk.OutgoingRecordInfo
}

func (p *Plugin) Init(provider sdk.Provider) {
	provider.Io().Info(fmt.Sprintf(`Init tool %v`, provider.Environment().ToolId()))
	p.provider = provider
	p.output = provider.GetOutputAnchor(`Output`)
}

func (p *Plugin) OnInputConnectionOpened(connection sdk.InputConnection) {
	p.provider.Io().Info(fmt.Sprintf(`got connection %v`, connection.Name()))
	p.outInfo = connection.Metadata().Clone().GenerateOutgoingRecordInfo()
	p.output.Open(p.outInfo)
}

func (p *Plugin) OnRecordPacket(connection sdk.InputConnection) {
	packet := connection.Read()
	for packet.Next() {
		p.outInfo.CopyFrom(packet.Record())
		p.output.Write()
	}
}

func (p *Plugin) OnComplete() {
	p.provider.Io().Info(`Done`)
}
```

plugin.go contains the implementation of your plugin.  Your implementation must satisfy the [Plugin interface](#Implementing-the-Plugin-interface).  In this example the tool simply copies incoming records and pushes them to its output.

[Back to table of contents](#Table-of-contents)

## Implementing the Plugin interface

Plugins must implement the Plugin interface:

```go
type Plugin interface {
	Init(Provider)
	OnInputConnectionOpened(InputConnection)
	OnRecordPacket(InputConnection)
	OnComplete()
}
```

The `Init` function is called immediately after the tool is registered and allows you to initialize your tool.  Your tool is given a [Provider](#Using-Provider), which allows you to retrieve your tool's configuration, interact with the Alteryx engine, retrieve environment information, and obtain output anchors for passing records to downstream tools.

The `OnInputConnectionOpened` function is called when an upstream tool is connected to your custom tool.  Your tool is given an [InputConnection](#Using-InputConnection), which allows you to retrieve the connection's name and metadata.  If your custom tool is an input tool this function will not be called.

The `OnRecordPacket` function is called when your custom tool recieves records from an upstream tool.  Your tool is given an [InputConnection](#Using-InputConnection), which allows you to check the incoming connection name, iterate through the incoming records, and retrieve the progress of the incoming datastream.  As with `OnInputConnectionOpened`, this function is not called if your custom tool is an input tool.

The `OnComplete` function is called at the end of your custom tool's lifecycle.  For tools which receive data from upstream tools, this happens after all incoming connections have been closed by the upstream tools.  For input tools, this happens when Alteryx is ready for your tool to start processing and sending data.

Below is an example of a struct that implements the Plugin interface:

```go
import (
	"github.com/tlarsen7572/goalteryx/sdk"
)

type Plugin struct {
	provider sdk.Provider
	output   sdk.OutputAnchor
	outInfo  *sdk.OutgoingRecordInfo
}

func (p *Plugin) Init(provider sdk.Provider) {
	p.provider = provider
	p.output = provider.GetOutputAnchor(`Output`)
}

func (p *Plugin) OnInputConnectionOpened(connection sdk.InputConnection) {
	p.outInfo = connection.Metadata().Clone().GenerateOutgoingRecordInfo()
	p.output.Open(p.outInfo)
}

func (p *Plugin) OnRecordPacket(connection sdk.InputConnection) {
	packet := connection.Read()
	for packet.Next() {
		p.outInfo.CopyFrom(packet.Record())
		p.output.Write()
	}
}

func (p *Plugin) OnComplete() {}
```

[Back to table of contents](#Table-of-contents)

## Registering your tool

Alteryx connects to custom tools through a C API function call.  All custom tools are expected to provide an entry point to the Alteryx engine that looks like the following:

```c
long NameOfPluginEntryPoint(int nToolID, void * pXmlProperties, void *pEngineInterface, void *r_pluginInterface);
```

For custom Go tools, the easiest way to do this is by creating an entry file that imports the C package and exports the declared entry points that perform the necessary registration steps.  Example:

```go
package main

import "C"
import (
	"github.com/tlarsen7572/goalteryx/sdk"
	"unsafe"
)

func main() {}

//export PluginEntry
func PluginEntry(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &Plugin{}
	return C.long(sdk.RegisterTool(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}
```

We start by importing the C and unsafe packages, as well as the SDK.  The next part of the file is an empty main function.  DLLs are expected to have a main function, but we do not make use of it, so we can keep it empty.

The next section implements our plugin's entry point.  It starts with a comment, `//export PluginEntry`, which has to match the declared entry point from the tool's Config.xml file.  Immediately after the comment is the function itself, with the signature the Alteryx engine is expecting.

The next line, `plugin := &Plugin{}`, creates a pointer of our plugin's struct.  We use that pointer in the `RegisterTool` function to actually register our tool and prepare it for use.

If you have multiple tools, you can register all of them in entry.go.  Example:

```go
package main

import "C"
import (
	"github.com/tlarsen7572/goalteryx/sdk"
	"unsafe"
)

func main() {}

//export FirstPlugin
func FirstPlugin(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &First{}
	return C.long(sdk.RegisterTool(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}

//export SecondPlugin
func SecondPlugin(toolId C.int, xmlProperties unsafe.Pointer, engineInterface unsafe.Pointer, pluginInterface unsafe.Pointer) C.long {
	plugin := &Second{}
	return C.long(sdk.RegisterTool(plugin, int(toolId), xmlProperties, engineInterface, pluginInterface))
}
```

Registering your custom tools in this manner keeps all the registration code neatly separated from your business logic and prevents your business logic from depending on the Unsafe and C packages.

[Back to table of contents](#Table-of-contents)

## Using Provider

`Provider` is used for obtaining information about your custom tool, sending messages to the Alteryx engine, and retrieving environmental information and variables from the Alteryx engine.  It has the following interface:

```go
type Provider interface {
	ToolConfig() string
	Io() Io
	GetOutputAnchor(string) OutputAnchor
	Environment() Environment
}
```

The `ToolConfig` function returns the current configuration for your custom tool.  It is provided as a raw XML string rather than attempting to provide a generic XML navigator object.  As tool configurations are unique to each tool, it is recommended to use Go's built-in parsing capabilities to unmarshal the XML into custom structs fit for purpose.

The `Io` function returns an [Io](#Using-Io), which is used primarily for sending messages to the Alteryx engine.

The `GetOutputAnchor` function returns an [OutgoingAnchor](#Using-OutputAnchor) which you can use to send records to downstream tools.

The `Environment` function returns an [Environment](#Environment), which you can use to obtain your custom tool's ID and retrieve environmental variables from the Alteryx engine.

[Back to table of contents](#Table-of-contents)

## Using OutputAnchor

`OutputAnchor` is the interface you use to send data to downstream tools.  It has the following interface:

```go
type OutputAnchor interface {
	Name() string
	IsOpen() bool
	Metadata() *OutgoingRecordInfo
	Open(info *OutgoingRecordInfo)
	Write()
	UpdateProgress(float64)
}
```

The `Name` function returns the name of the output anchor and should match the name provided in the tool's Config.xml file.

The `IsOpen` function tells whether the `Open` function has been called on the connection.

The `Metadata` function returns a pointer to the `OutgoingRecordInfo` that the anchor was opened with.  If `Open` has not been called yet, the return value is nil.  See the section on [RecordInfo](#RecordInfo) for more information about how to use and generate `OutgoingRecordInfo` structs.

The `Open` function opens the output anchor and sends metadata downstream to all connected tools.  The `OutgoingRecordInfo` you open the connection with is also where the `OutputAnchor` reads data from when Write() is called.

The `Write` function writes the current values in the `OutgoingRecordInfo` to downstream tools.

The `UpdateProgress` function notifies downstream tools on the percentage completion of the dataset being sent.  The value provided should be between 1 and 0, with 1 being 100% completed.

[Back to table of contents](#Table-of-contents)

## Using Io

`Io` is the interface you use to send messages to the Alteryx engine.  It has the following interface:

```go
type Io interface {
	Error(string)
	Warn(string)
	Info(string)
	UpdateProgress(float64)
	DecryptPassword(string) string
}
```

The `Error` function sends an error to the Alteryx engine.  This shows up in Designer as an error.  When run in a unit test context, it prints an error message to stdout.

The `Warn` function sends a warning to the Alteryx engine.  This shows up in Designer as a warning.  When run in a unit test context, it prints a warning message to stdout.

The `Info` function sends a message to the Alteryx engine.  This shows up in Designer as an informational message.  When run in a unit test context, it prints the message to stdout.

The `UpdateProgress` function notifies the Alteryx engine of the current percentage completion of the custom tool.  This is the overall completion of the tool as opposed to the datastream completion percentage in the `OutputAnchor.UpdateProgress()` method.

The `DecryptPassword` function decrypts a password encrypted by the front-end UI.

[Back to table of contents](#Table-of-contents)

## Using Environment

`Environment` is the interface you use to retrieve environment variables from the Alteryx engine.  It has the following interface:

```go
type Environment interface {
	UpdateOnly() bool
	UpdateMode() string
	DesignerVersion() string
	WorkflowDir() string
	AlteryxInstallDir() string
	AlteryxLocale() string
	ToolId() int
	UpdateToolConfig(string)
}
```

The `UpdateOnly` function identifies whether the Alteryx engine expects the tool to send data.  If the return value is `true`, the tool should not send records downstream.

The `UpdateMode` function returns one of a blank string, 'Quick', or 'Full'.

The `DesignerVersion` function returns the version of Designer being run.  If run in a unit test context, it returns the value 'TestHarness'.

The `WorkflowDir` function returns the folder of the workflow the tool is being run in.

The `AlteryxInstallDir` function returns the Alteryx installation folder.  If run in a unit test context, it returns an empty string.

The `AlteryxLocale` function returns the locale/language setting of the current user.

The `ToolId` function returns the ID of the custom tool in the current workflow.

The `UpdateToolConfig` function provides a way for the custom tool to update its own configuration and send it back to Designer for persistance.

[Back to table of contents](#Table-of-contents)

## Using InputConnection

`InputConnection` is provided to the custom tool by the SDK and is the interface by which you interact with incoming connections and data.  It has the following interface:

```go
type InputConnection interface {
	Name() string
	Metadata() IncomingRecordInfo
	Read() RecordPacket
	Progress() float64
}
```

The `Name` function returns the name of the incoming connection.  This name should match the name of one of the input connections defined in the tool's Config.xml file.

The `Metadata` function returns the structure of the incoming data.  See [RecordInfo](#RecordInfo) for more information about using this interface.

The `Read` function returns a `RecordPacket` containing a cache of records that have been pushed to your custom tool.  If you have multiple input connections, it is important to always first read the name of the input connection so you know how to process the incoming data.  Input connections are not guaranteed to arrive in any specific order, nor is it guaranteed that all of an input connection's records will arrive before another input connection starts sending its data.  The `Read` function should only be used during the `OnRecordPacket` function of the [Plugin](#Implementing-the-Plugin-interface).

The `Progress` function returns the percentage of records that have been passed through the `InputConnection`.

[Back to table of contents](#Table-of-contents)

## RecordInfo

There are 3 different `RecordInfo` structs that you may use during the lifecycle of your custom tool:

[IncomingRecordInfo](#IncomingRecordInfo)  
[EditingRecordInfo](#EditingRecordInfo)  
[OutgoingRecordInfo](#OutgoingRecordInfo)  

#### IncomingRecordInfo

`IncomingRecordInfo` is provided during your [custom tool's](#Implementing-the-Plugin-interface) `OnInputConnectionOpened` and `OnRecordPacket` functions.  It provides for a way to inspect the structure of your incoming data and generate outgoing record information that can copy data from incoming datastreams.  `IncomingRecordInfo` has the following interface:

```go
func NumFields() int
func Fields() []b.FieldBase
func Clone() *EditingRecordInfo
func GetBlobField(name string) (IncomingBlobField, error)
func GetBoolField(name string) (IncomingBoolField, error)
func GetIntField(name string) (IncomingIntField, error)
func GetFloatField(name string) (IncomingFloatField, error)
func GetStringField(name string) (IncomingStringField, error)
func GetTimeField(name string) (IncomingTimeField, error)
```

The `NumFields` function returns the number of fields in the `IncomingRecordInfo`.

The `Fields` function returns the list of fields.  Each field provides the name, type, source, size, and scale of the field.

The `Clone` function clones the `IncomingRecordInfo` into an [EditingRecordInfo](#EditingRecordInfo).  Using the `Clone` function to build your outgoing recordinfo allows you to easily copy data from incoming records to your outgoing records.

The `GetBlobField` function returns a struct that lets you extract blob values (slice of bytes) from an incoming record.  This function only returns correctly if the field type of the named field is 'Blob' or 'SpatialObj'.  If the field does not exist or is the incorrect type, an error is returned.

The `GetBoolField` function returns a struct that lets you extract boolean values from an incoming record.  This function only returns correctly if the field type of the named field is 'Bool'.  If the field does not exist or is the incorrect type, an error is returned.

The `GetIntField` function returns a struct that lets you extract integers from an incoming record.  This function only returns correctly if the field type of the named field is 'Byte', 'Int16', 'Int32', or 'Int64'.  If the field does not exist or is the incorrect type, an error is returned.

The `GetFloatField` function returns a struct that lets you extract decimal numbers from an incoming record.  This function only returns correctly if the field type of the named field is 'Float', 'Double', or 'FixedDecimal'.  If the field does not exist or is the incorrect type, an error is returned.

The `GetStringField` function returns a struct that lets you extract text values from an incoming record.  This function only returns correctly if the field type of the named field is 'String', 'WString', 'V_String', or 'V_WString'.  If the field does not exist or is the incorrect type, an error is returned.

The `GetTimeField` function returns a struct that lets you extract temporal values from an incoming record.  This function only returns correctly if the field type of the named field is 'Date' or 'DateTime'.  If the field does not exist or is the incorrect type, an error is returned.

Each of the GetXxxField functions returns a field struct that provides the name, type source, size (if applicable), and scale (if FixedDecimal).  The field struct also provides a `GetValue` function that allows you to retrieve the field's value.  The `GetValue` function signatures for the various incoming fields are as follows:

IncomingBlobField: GetValue(Record) (value []byte, isNull bool)  
IncomingBoolField: GetValue(Record) (value bool, isNull bool)  
IncomingIntField: GetValue(Record) (value int, isNull bool)  
IncomingFloatField: GetValue(Record) (value float64, isNull bool)  
IncomingStringField: GetValue(Record) (value stirng, isNull bool)  
IncomingTimeField: GetValue(Record) (value time.Time, isNull bool)  

An example of a tool that uses GetXxxField to extract values from specific fields is below:

```go
type Plugin struct {
	field    sdk.IncomingStringField
}

func (p *Plugin) Init(provider sdk.Provider) {}

func (p *Plugin) OnInputConnectionOpened(connection sdk.InputConnection) {
	var err error
	p.field, err = connection.Metadata().GetStringField(`MyField`)
	if err != nil {
		panic(`field not found or is of the wrong type`)
	}
}

func (p *Plugin) OnRecordPacket(connection sdk.InputConnection) {
	packet := connection.Read()
	for packet.Next() {
		value, isNull := p.field.GetValue(packet.Record())
	}
}

func (p *Plugin) OnComplete() {}
```

#### EditingRecordInfo

`EditingRecordInfo` is used to edit an incoming recordinfo and then generate the final outgoing recordinfo once all edits are made.  It has the following interface:

```go
func NumFields() int
func Fields() []IncomingField
func AddBoolField(name string, source string, options ...AddFieldOptionSetter) string
func AddByteField(name string, source string, options ...AddFieldOptionSetter) string
func AddInt16Field(name string, source string, options ...AddFieldOptionSetter) string
func AddInt32Field(name string, source string, options ...AddFieldOptionSetter) string
func AddInt64Field(name string, source string, options ...AddFieldOptionSetter) string
func AddFloatField(name string, source string, options ...AddFieldOptionSetter) string
func AddDoubleField(name string, source string, options ...AddFieldOptionSetter) string
func AddFixedDecimalField(name string, source string, size int, scale int, options ...AddFieldOptionSetter) string
func AddStringField(name string, source string, size int, options ...AddFieldOptionSetter) string
func AddWStringField(name string, source string, size int, options ...AddFieldOptionSetter) string
func AddV_StringField(name string, source string, size int, options ...AddFieldOptionSetter) string
func AddV_WStringField(name string, source string, size int, options ...AddFieldOptionSetter) string
func AddDateField(name string, source string, options ...AddFieldOptionSetter) string
func AddDateTimeField(name string, source string, options ...AddFieldOptionSetter) string
func AddBlobField(name string, source string, size int, options ...AddFieldOptionSetter) string
func AddSpatialObjField(name string, source string, size int, options ...AddFieldOptionSetter) string
func RemoveFields(fieldNames ...string)
func MoveField(name string, newIndex int) error
func GenerateOutgoingRecordInfo() *OutgoingRecordInfo
```

The `NumFields` function returns the number of fields currently in the recordinfo.

The `Fields` function returns a list of basic field information of all of the fields currently in the recordinfo.

The `AddXxxField` functions adds a new field to the recordinfo.  Each function represents a different storage type for the underlying data.  All functions require a name and source, with size and scale being required on specific field types such as strings and fixed decimal fields.  You may also provide a list of options when creating the field.  The currently supported options are:

* `InsertAt(position int)`: Use this option to insert the field in the beginning or middle of the record.  For example, to insert a new Int32 field at the beginning of the recordinfo, use:
	```go
	editor.AddInt32Field(`FieldName`, `some source`, sdk.InsertAt(0))
	```

The `RemoveFields` function removes the provided list of fields from the record, if they exist.

The `MoveField` function moves a field to a different position in the record.  An error is returned if `newIndex` is out of bounds or if the name provided does not exist in the record.

The `GenerateOutgoingRecordInfo` function returns a pointer to an [OutgoingRecordInfo](#OutgoingRecordInfo) struct, which is used to open [OutputAnchors](#Using-OutputAnchor) and set values for writing to downstream tools.

#### OutgoingRecordInfo

`OutgoingRecordInfo` is used to send metadata to downstream tools and store values that will be written to the custom tool's output anchors.  You can create an `OutgoingRecordInfo` from an `EditingRecordInfo` or by using the `NewOutgoingRecordInfo` function in the SDK.  The following example creates an `OutgoingRecordInfo` with a Bool field, an Int64 field, and a V_WString field:

```go
recordInfo, fieldNames := sdk.NewOutgoingRecordInfo([]sdk.NewOutgoingField{
	sdk.NewBoolField(`Field 1`, `source`),
	sdk.NewInt64Field(`Field 2`, `source`),
	sdk.NewV_WStringField(`Field 3`, `source`, 1000),
})
```

If duplicate field names are specified in the list of `NewOutgoingField`, then `NewOutgoingRecordInfo()` will rename the duplicate fields.  The second return value from `NewOutgoingRecordInfo` contains the actual field names in the `OutgoingRecordInfo`.

`OutgoingRecordInfo` has the following interface:

```go
func FixedSize() int
func HasVarFields() int
func DataSize() uint32
func CopyFrom(Record)
```

The `FixedSize` function returns the size of the fixed portion of the `RecordInfo` data structure.

The `HasVarFields` function identifies whether the recordinfo contains variable-length fields (V_String, V_WString, Blob, or SpatialObj).

The `DataSize` functions returns the record size of the current values in the `OutgoingRecordInfo` struct.

The `CopyFrom` function copies values from the incoming record into its current values.  This function only copies those fields which originated from an `IncomingRecordInfo` via the `Clone` method.

The following code shows an end-to-end example of how to use the various recordinfo structs by implementing a custom tool that adds a record ID to the beginning of the record.

```go
package awesomeProject

import (
	"github.com/tlarsen7572/goalteryx/sdk"
)

type Plugin struct {
	outputAnchor      sdk.OutputAnchor
	outputInfo        *sdk.OutgoingRecordInfo
	recordIdFieldName string
	recordId          int
}

func (p *Plugin) Init(provider sdk.Provider) {
	p.outputAnchor = provider.GetOutputAnchor(`Output`)
	p.recordId = 0
}

func (p *Plugin) OnInputConnectionOpened(connection sdk.InputConnection) {
	// convert the incoming recordinfo into an editor
	editor := connection.Metadata().Clone()
	
	// add the record ID field
	p.recordIdFieldName = editor.AddInt32Field(`RecordId`, `my custom tool`, sdk.InsertAt(0))
	
	// generate the outgoing recordinfo
	p.outputInfo = editor.GenerateOutgoingRecordInfo()
	
	// open the output anchor with the metadata from the outgoing recordinfo
	p.outputAnchor.Open(p.outputInfo)
}

func (p *Plugin) OnRecordPacket(connection sdk.InputConnection) {
	packet := connection.Read()
	for packet.Next() {
		// copy data from the incoming record to the current values of the outgoing recordinfo
		p.outputInfo.CopyFrom(packet.Record())
		
		// set the record ID field
		p.outputInfo.IntFields[p.recordIdFieldName].SetInt(p.recordId)
		
		// write the current outgoing recordinfo values to downstream tools
		p.outputAnchor.Write()
		p.recordId++
	}
}

func (p *Plugin) OnComplete() {}
```

[Back to table of contents](#Table-of-contents)

## Using RecordPacket

`RecordPacket` is an abstraction used to iterate through the packet of records sent to your custom tool by upstream tools.  Records are recieved (and sent) in 4mb chunks.  This is done to minimize the number of calls between the Alteryx engine and the Go runtime, each of which has bookkeeping overhead.

`RecordPacket` has the following interface:

```go
func Next() bool
func Record() Record
```

The `Next` function tries to retrieve the next record in the packet.  If there are no more records, it returns false; otherwise, it returns true.

The `Record` function returns the record retrieved during the call to `Next`.

The easiest way to interact with `RecordPacket` is to iterate through it using a for loop:

```go
func iteratePacket(packet RecordPacket) {
	for packet.Next() {
		record := packet.Record()
		
		// do something with the record
	}
}
```

[Back to table of contents](#Table-of-contents)

## Testing your tools

GoAlteryx includes testing facilities to assist your development of custom tools.  They are designed to mimic the lifecycle events your tool will experience during the run of a workflow.  As a result, you can develop and test your tools without running them in Alteryx and still be confident that they will work.  This also frees the developer to choose non-Windows development environments such as macOS.

A basic example of unit testing input and passthrough tools is below:

```go
package awesomeProject_test

import (
	"awesomeProject"
	"github.com/tlarsen7572/goalteryx/sdk"
	"testing"
)

func TestInputTool(t *testing.T) {
	plugin := &awesomeProject.InputPlugin{}
	runner := sdk.RegisterToolTest(plugin, 1, `<Configuration></Configuration>`)
	collector := runner.CaptureOutgoingAnchor(`Output`)
	runner.SimulateLifecycle()
	t.Logf(`%v`, collector.Data)
}

func TestPassthroughTool(t *testing.T) {
	plugin := &awesomeProject.PassthroughPlugin{}
	runner := sdk.RegisterToolTest(plugin, 1, `<Configuration></Configuration>`)
	collector := runner.CaptureOutgoingAnchor(`Output`)
	runner.ConnectInput(`Input`, `testfile.txt`)
	runner.SimulateLifecycle()
	t.Logf(`%v`, collector.Data)
}
```

In both cases, the unit test begins by creating a pointer to your plugin and then registering it with the test harness by calling `sdk.RegisterToolTest()`.  You provide the pointer, tool ID, and configuration XML as a string in the registration call.  The registration function returns a test runner.  Using the test runner, you can capture outgoing records and connect input testing files to your custom tools.  Calling `SimulateLifecycle()` on the runner will then execute the lifecycle events and test your tool.  You can inspect and verify your custom tools' outputs by inspecting the `Data` member on the captured outgoing anchor.

A detailed review of the test harness features are below.  We start with the signature of the `RegisterToolTest` function:

```go
func RegisterToolTest(plugin Plugin, toolId int, xmlProperties string, optionSetters ...OptionSetter) *FileTestRunner
```

`plugin` is a struct that fulfills the Plugin interface specified by the SDK.

`toolId` is an arbitrary integer that represents the tool's ID when it is placed on a workflow's canvas.

`xmlProperties` is a string containing the XML configuration you want your tool to receive for the test.

`optionSetters` is a list of options to pass to the test harness.  The following options are currently available:

* `func UpdateOnly(bool)`: Sets the engine's UpdateOnly environment variable
* `func UpdateMode(string)`: Sets the engine's UpdateMode environment variable
* `func WorkflowDir(string)`: Sets a custom workflow directory for the test
* `func AlteryxLocal(string)`: Sets the locale for the test

Any, all, or no options may be specified.  An example of registering a tool with the test harness that specifies the UpdateOnly and AlteryxLocale options is below:

```go
runner := sdk.RegisterToolTest(plugin, 1, `<Configuration></Configuration>`, sdk.UpdateOnly(true), sdk.AlteryxLocale(`en-us`))
```

The `RegisterToolTest` function returns a pointer to a `FileTestRunner`.  The interface for `FileTestRunner` is:

```go
func CaptureOutgoingAnchor(name string) *RecordCollector
func ConnectInput(name string, dataFile string)
func SimulateLifecycle()
```

The `CaptureOutgoingAnchor` function adds an outgoing connection to the specified anchor of your tool.  It returns a pointer to a `RecordCollector`, which you can use to inspect the data output from your tool.  Retrieving `RecordCollector.Data` will return a `map[string][]interface{}` containing the output data.  The map key is the output field name and the map value is a list of `interface{}` containing the values that were output for that field.

The `ConnectInput` function connects input data to the specified anchor of your tool.  You specify the path to a data file in the second argument.  Data files can be best thought of as pipe-delimited files with a few special rules.  The rules to follow are:

1. The first row must contain the field names
2. The second row must contain the field types
3. Bool, integer, decimal, date, and binary fields should not be quoted
4. Strings should be double-quoted if leading/trailing whitespace is desired or the field value contains a pipe
5. If a string field needs a double quote in the value, escape it with a backslash (\\")
6. If a string field needs a backslash in the value, escape it with a backslash (\\\\)
7. You may use `\r` and `\n` to specify carriage return and newline characters
8. Leading or trailing spaces outside of double quotes is ignored
9. Empty fields are interpreted as nulls
10. String fields with a value of 2 double quotes ("") are interpreted as empty strings rather than null
11. Dates should be entered with a `YYYY-mm-dd` format
12. DateTimes should be entered with a `YYY-mm-dd HH:MM:SS` format
13. Size and scale, for fields that require them, are specified after the field type and separated by semi-colons (;)

An example data file is below that illustrates how to set up each different type of field and how the data should be formatted.

```
Field1|Field2|Field3|Field4|Field5|Field6|Field7|Field8           |Field9    |Field10    |Field11       |Field12         |Field13   |Field14            |Field15|Field16
Bool  |Byte  |Int16 |Int32 |Int64 |Float |Double|FixedDecimal;19;2|String;100|WString;100|V_String;10000|V_WString;100000|Date      |DateTime           |Blob;10|SpatialObj;100
true  |2     |100   |1000  |10000 |12.34 |1.23  |     234.56      |"ABC"     |"Hello "   |" World"      |"abcdefg"       |2020-01-01|2020-01-02 03:04:05|       |
false |-2    |-100  |-1000 |-10000|-12.34|-1.23 |    -234.56      |"DE|\"FG" |HIJK       |  LMNOP       |"QRSTU\r\nVWXYZ"|2020-02-03|2020-01-02 13:14:15|       |
      |      |      |      |      |      |      |                 |          |           |              |                |          |                   |       |
true  |42    |-110  |392   |2340  |12    |41.22 |  98.2           |""        |"HIJK"     |  LMN         |"qrstuvwxyz"    |2020-02-13|2020-11-02 13:14:15|       |
```

[Back to table of contents](#Table-of-contents)

## Feature parity with the Python SDK

The graph below identifies elements of the Python SDK API that are implemented, or not implemented, in goalteryx.

🟢 &nbsp;= Implemented, 🟡 &nbsp;= Not implemented, but planned, ⚪ &nbsp;= Not planned for implementation

* 🟢 &nbsp;Plugin
    * 🟢 &nbsp;Init
    * 🟢 &nbsp;OnInputConnectionOpened
    * 🟢 &nbsp;OnRecordPacket
    * 🟢 &nbsp;OnComplete

* 🟢 &nbsp;Provider
    * 🟢 &nbsp;ToolConfig
    * ⚪ &nbsp;Logger
    * 🟢 &nbsp;IO
    * 🟢 &nbsp;Environment
    * ⚪ &nbsp;GetInputAnchor
    * 🟢 &nbsp;GetOutputAnchor

* 🟡 &nbsp;IO
    * 🟢 &nbsp;Error
    * 🟢 &nbsp;Warn
    * 🟢 &nbsp;Info
    * 🟢 &nbsp;UpdateProgress
    * 🟡 &nbsp;CreateTempFile
    * 🟢 &nbsp;DecryptPassword

* 🟢 &nbsp;Environment
    * 🟢 &nbsp;UpdateOnly
    * 🟢 &nbsp;UpdateMode
    * 🟢 &nbsp;DesignerVersion
    * 🟢 &nbsp;WorkflowDir
    * 🟢 &nbsp;AlteryxInstallDir
    * 🟢 &nbsp;Locale
    * 🟢 &nbsp;ToolId
    * 🟢 &nbsp;UpdateToolConfig

* 🟡 &nbsp;OutputAnchor
    * 🟢 &nbsp;Name
    * ⚪ &nbsp;AllowMultiple
    * ⚪ &nbsp;Optional
    * ⚪ &nbsp;NumConnections
    * 🟡 &nbsp;IsOpen
    * 🟢 &nbsp;Metadata
    * 🟢 &nbsp;Open
    * 🟢 &nbsp;Write
    * ⚪ &nbsp;Flush
    * 🟡 &nbsp;Close
    * 🟢 &nbsp;UpdateProgress

* ⚪ &nbsp;InputAnchor
    * ⚪ &nbsp;Name
    * ⚪ &nbsp;AllowMultiple
    * ⚪ &nbsp;Optional
    * ⚪ &nbsp;Connections

* 🟡 &nbsp;InputConnection
    * 🟢 &nbsp;Name
    * 🟢 &nbsp;Metadata
    * ⚪ &nbsp;Anchor
    * 🟢 &nbsp;Read
    * ⚪ &nbsp;MaxPacketSize
    * 🟢 &nbsp;Progress
    * 🟡 &nbsp;Status

* RecordPacket
    * RecordPacket is intentionally different than the Python implementation. Python translates record packets to and from data frames. This makes sense for Python tools, but not for Go. The Go implementation of RecordPacket mimics the behavior of the Go SQL package. Records in a record packet are accessed through an iterator and field-specific extractors.



[Back to table of contents](#Table-of-contents)
