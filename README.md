<img src="https://github.com/tlarsen7572/goalteryx/blob/master/goalteryx_icon_whiteBackground.png?raw=true" width="200">

# GoAlteryx

An unofficial SDK for building custom Alteryx tools with Go.

## Why a Go SDK?

With the announced deprecation of the .NET SDK, a gap formed between the C/C++ and Python SDKs.  C/C++ are low-level languages requiring great care and expertise to ensure proper memory management.  Python is very approachable but is slower.  I wanted to build tools with a middle-ground language having decent performance and simplified memory management.  Go fit the bill and is my favorite language to code in.

## Table of contents

[Installation](https://github.com/tlarsen7572/goalteryx#Installation)
[Building your custom tools](https://github.com/tlarsen7572/goalteryx#Building-your-custom-tools)
[Sample tool](https://github.com/tlarsen7572/goalteryx#Sample-tool)
[Registering your tool](https://github.com/tlarsen7572/goalteryx#Registering-your-tool)
[Implementing the Plugin interface](https://github.com/tlarsen7572/goalteryx#Implementing-the-Plugin-interface)
[Using Provider](https://github.com/tlarsen7572/goalteryx#Using-Provider)
[Using OutputAnchor](https://github.com/tlarsen7572/goalteryx#Using-OutputAnchor)
[Using Io](https://github.com/tlarsen7572/goalteryx#Using-Io)
[Using Environment](https://github.com/tlarsen7572/goalteryx#Using-Environment)
[Using InputAnchor](https://github.com/tlarsen7572/goalteryx#Using-InputAnchor)

## Installation

Install goalteryx using Go modules: `go get github.com/tlarsen7572/goalteryx`

## Building your custom tools

You should specify the output DLL file and make sure `-buildmode` is set to `c-shared`.  For reference, the following command is used to build the included example tools:

```
go build -o "C:\Program Files\Alteryx\bin\Plugins\goalteryx.dll" -buildmode=c-shared goalteryx/implementation_example
```

I build directly to the Plugins folder in the Alteryx installation folder of my dev environment.  This allows me to rebuild my tools and run them directly in Alteryx without additional copying.  You do not need to close and restart Alteryx when you rebuild a DLL.  The next time you run a workflow with your custom tool, the new DLL will be used.  It should go without saying that you should not do this in production.

[Back to table of contents](https://github.com/tlarsen7572/goalteryx#Table-of-contents)

## Sample tool

[Back to table of contents](https://github.com/tlarsen7572/goalteryx#Table-of-contents)

## Registering your tool

[Back to table of contents](https://github.com/tlarsen7572/goalteryx#Table-of-contents)

## Implementing the Plugin interface

[Back to table of contents](https://github.com/tlarsen7572/goalteryx#Table-of-contents)

## Using Provider

[Back to table of contents](https://github.com/tlarsen7572/goalteryx#Table-of-contents)

## Using OutputAnchor

[Back to table of contents](https://github.com/tlarsen7572/goalteryx#Table-of-contents)

## Using Io

[Back to table of contents](https://github.com/tlarsen7572/goalteryx#Table-of-contents)

## Using Environment

[Back to table of contents](https://github.com/tlarsen7572/goalteryx#Table-of-contents)

## Using InputAnchor

[Back to table of contents](https://github.com/tlarsen7572/goalteryx#Table-of-contents)
