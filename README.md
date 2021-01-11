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

The following 3 code files represent a basic tool in Alteryx that copies incoming records and pushes them through its output.

#### entry.h

```
long __declspec(dllexport) PluginEntry(int nToolID,
	void * pXmlProperties,
	void *pEngineInterface,
	void *r_pluginInterface);
```

entry.h declares your plugin's entry function for the Alteryx engine and is one half of registering your plugin.  See [Registering your tool](#Registering-your-tool) for more info.

#### entry.go

```
package main

/*
#include "entry.h"
*/
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

entry.go is the second half of [plugin registration](#Registering-your-tool).

#### plugin.go

```
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

```
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

```
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

```
long __declspec(dllexport) NameOfPluginEntryPoint(int nToolID, void * pXmlProperties, void *pEngineInterface, void *r_pluginInterface);
```

For custom Go tools, the easiest way to do this is to create a file called entry.h with the declared entry points.  If you plan on packaging multiple tools into the DLL, you can specify all of them in entry.h.  Example:

```
long __declspec(dllexport) FirstPlugin(int nToolID, void * pXmlProperties, void *pEngineInterface, void *r_pluginInterface);
long __declspec(dllexport) SecondPlugin(int nToolID, void * pXmlProperties, void *pEngineInterface, void *r_pluginInterface);
```

Now that you have declared the plugin's entry point, you need to implement it.  The easiest way to do this is to create a file called entry.go that performs the necessary registration steps.  Example:

```
package main

/*
#include "entry.h"
*/
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

We start by importing the C package and including entry.h that we created earlier.  The next part of the file is an empty main function.  DLLs are expected to have a main function, but we do not make use of it so we can keep it empty.

The next section implements our plugin's entry point.  It starts with a comment, `//export PluginEntry`, which has to match the declared function name from entry.h.  Immediately after the comment is the function itself, also with the same name as that declared in entry.h.

The next line, `plugin := &Plugin{}`, creates a pointer of our plugin's struct.  We use that pointer in the `RegisterTool` function on the next line to actually register our tool and prepare it for use.

If you have multiple tools, you can provide all of their implementations in entry.go.  Example:

```
package main

/*
#include "entry.h"
*/
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

Registering your custom tools in this manner keeps all of the registration code neatly separated from your business logic and prevents your business logic from depending on the Unsafe and C packages.

[Back to table of contents](#Table-of-contents)

## Using Provider

Provider is used for obtaining information about your custom tool, sending messages to the Alteryx engine, and retrieving environmental information and variables from the Alteryx engine.  It has the following interface:

```
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

[Back to table of contents](#Table-of-contents)

## Using Io

[Back to table of contents](#Table-of-contents)

## Using Environment

[Back to table of contents](#Table-of-contents)

## Using InputConnection

[Back to table of contents](#Table-of-contents)

## RecordInfo

[Back to table of contents](#Table-of-contents)

## Using RecordPacket

[Back to table of contents](#Table-of-contents)
