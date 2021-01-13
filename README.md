<img src="https://github.com/tlarsen7572/goalteryx/blob/master/goalteryx_icon_whiteBackground.png?raw=true" width="200">

# GoAlteryx

An unofficial SDK for building custom Alteryx tools with Go.

## Why a Go SDK?

With the announced deprecation of the .NET SDK, a gap formed between the C/C++ and Python SDKs.  C/C++ are low-level languages requiring great care and expertise to ensure proper memory management.  Python is very approachable but is slower.  I wanted to build tools with a middle-ground language having decent performance and simplified memory management.  Go fit the bill and is my favorite language to code with.

## Table of contents

1. [Installation](#Installation)  
2. [Building your custom tools](#Building-your-custom-tools)  
3. [Sample tool](#Sample-tool)  
4. [Implementing the Plugin interface](#Implementing-the-Plugin-interface)  
5. [Registering your tool](#Registering-your-tool)  
6. [Using Provider](#Using-Provider)  
7. [Using OutputAnchor](#Using-OutputAnchor)  
8. [Using Io](#Using-Io)  
9. [Using Environment](#Using-Environment)  
10. [Using InputConnection](#Using-InputConnection)  
11. [RecordInfo](#RecordInfo)  
12. [Using RecordPacket](#Using-RecordPacket)  

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
}
```

The `Error` function sends an error to the Alteryx engine.  This shows up in Designer as an error.  When run in a unit test context, it prints an error message to stdout.

The `Warn` function sends a warning to the Alteryx engine.  This shows up in Designer as a warning.  When run in a unit test context, it prints a warning message to stdout.

The `Info` function sends a message to the Alteryx engine.  This shows up in Designer as an informational message.  When run in a unit test context, it prints the message to stdout.

The `UpdateProgress` function notifies the Alteryx engine of the current percentage completion of the custom tool.  This is the overall completion of the tool as opposed to the datastream completion percentage in the `OutputAnchor.UpdateProgress()` method.

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

`EditingRecordInfo` is used to edit incoming recordinfo's and then generate the final outgoing recordinfo once all edits are made.  It has the following interface:

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
func GenerateOutgoingRecordInfo() *OutgoingRecordInfo
```

[Back to table of contents](#Table-of-contents)

## Using RecordPacket

[Back to table of contents](#Table-of-contents)
