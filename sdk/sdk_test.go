package sdk_test

import (
	"bytes"
	"fmt"
	"github.com/tlarsendataguy/goalteryx/sdk"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

type TestImplementation struct {
	DidInit                    bool
	DidOnComplete              bool
	DidOnInputConnectionOpened bool
	DidOnRecordPacket          bool
	Config                     string
	Provider                   sdk.Provider
	Output                     sdk.OutputAnchor
}

func (t *TestImplementation) TestIo() {
	t.Provider.Io().Info(`test1`)
	t.Provider.Io().Warn(`test2`)
	t.Provider.Io().Error(`test3`)
	t.Provider.Io().UpdateProgress(0.10)
}

func (t *TestImplementation) Init(provider sdk.Provider) {
	t.DidInit = true
	t.Config = provider.ToolConfig()
	t.Provider = provider
	t.Output = provider.GetOutputAnchor(`Output`)
}

func (t *TestImplementation) OnInputConnectionOpened(_ sdk.InputConnection) {
	t.DidOnInputConnectionOpened = true
}

func (t *TestImplementation) OnRecordPacket(_ sdk.InputConnection) {
	t.DidOnRecordPacket = true
}

func (t *TestImplementation) OnComplete() {
	t.DidOnComplete = true
	t.Output.UpdateProgress(1)
}

type TestInputTool struct {
	Provider     sdk.Provider
	Output       sdk.OutputAnchor
	OutputConfig *sdk.OutgoingRecordInfo
}

func (i *TestInputTool) Init(provider sdk.Provider) {
	i.Provider = provider
	i.Output = provider.GetOutputAnchor(`Output`)
}

func (i *TestInputTool) OnInputConnectionOpened(_ sdk.InputConnection) {
	panic("This should never be called")
}

func (i *TestInputTool) OnRecordPacket(_ sdk.InputConnection) {
	panic("This should never be called")
}

func (i *TestInputTool) OnComplete() {
	source := `source`
	output, _ := sdk.NewOutgoingRecordInfo([]sdk.NewOutgoingField{
		sdk.NewBlobField(`Field1`, source, 100),
		sdk.NewBoolField(`Field2`, source),
		sdk.NewByteField(`Field3`, source),
		sdk.NewInt16Field(`Field4`, source),
		sdk.NewInt32Field(`Field5`, source),
		sdk.NewInt64Field(`Field6`, source),
		sdk.NewFloatField(`Field7`, source),
		sdk.NewDoubleField(`Field8`, source),
		sdk.NewFixedDecimalField(`Field9`, source, 19, 2),
		sdk.NewStringField(`Field10`, source, 100),
		sdk.NewWStringField(`Field11`, source, 100),
		sdk.NewV_StringField(`Field12`, source, 100000),
		sdk.NewV_WStringField(`Field13`, source, 100000),
		sdk.NewDateField(`Field14`, source),
		sdk.NewDateTimeField(`Field15`, source),
		sdk.NewSpatialObjField(`Field16`, source, 1000000),
		sdk.NewTimeField(`Field17`, source),
	})
	i.OutputConfig = output
	i.Output.Open(output)

	for index := 0; index < 10; index++ {
		output.BlobFields[`Field1`].SetBlob([]byte{byte(index)})
		output.BoolFields[`Field2`].SetBool(index%2 == 0)
		output.IntFields[`Field3`].SetInt(index)
		output.IntFields[`Field4`].SetInt(index)
		output.IntFields[`Field5`].SetInt(index)
		output.IntFields[`Field6`].SetInt(index)
		output.FloatFields[`Field7`].SetFloat(float64(index))
		output.FloatFields[`Field8`].SetFloat(float64(index))
		output.FloatFields[`Field9`].SetFloat(float64(index))
		output.StringFields[`Field10`].SetString(strconv.Itoa(index))
		output.StringFields[`Field11`].SetString(strconv.Itoa(index))
		output.StringFields[`Field12`].SetString(strconv.Itoa(index))
		output.StringFields[`Field13`].SetString(strconv.Itoa(index))
		output.DateTimeFields[`Field14`].SetDateTime(time.Date(2020, 1, index, 0, 0, 0, 0, time.UTC))
		output.DateTimeFields[`Field15`].SetDateTime(time.Date(2020, 1, index, 0, 0, 0, 0, time.UTC))
		output.BlobFields[`Field16`].SetBlob([]byte{byte(index)})
		output.DateTimeFields[`Field17`].SetDateTime(time.Date(0, 0, 0, 12, 31, index, 0, time.UTC))
		i.Output.Write()
	}
	i.Output.UpdateProgress(1)
}

type InputRecordLargerThanCache struct {
	Output sdk.OutputAnchor
}

func (i *InputRecordLargerThanCache) Init(provider sdk.Provider) {
	i.Output = provider.GetOutputAnchor(`Output`)
}

func (i *InputRecordLargerThanCache) OnInputConnectionOpened(_ sdk.InputConnection) {
	panic("this should never be called")
}

func (i *InputRecordLargerThanCache) OnRecordPacket(_ sdk.InputConnection) {
	panic("this should never be called")
}

func (i *InputRecordLargerThanCache) OnComplete() {
	info, _ := sdk.NewOutgoingRecordInfo([]sdk.NewOutgoingField{
		sdk.NewV_WStringField(`Field1`, `source`, 1000000000),
	})
	i.Output.Open(info)
	info.StringFields[`Field1`].SetString(`hello world`)
	i.Output.Write()
	info.StringFields[`Field1`].SetString(strings.Repeat(`ABCDEFGHIJKLMNOPQRSTUVWXYZ`, 200000))
	i.Output.Write()
	info.StringFields[`Field1`].SetString(`zyxwvutsrqponmlkjihgfedcba`)
	i.Output.Write()
}

type InputWithNulls struct {
	output sdk.OutputAnchor
}

func (i *InputWithNulls) Init(provider sdk.Provider) {
	i.output = provider.GetOutputAnchor(`Output`)
}

func (i *InputWithNulls) OnInputConnectionOpened(_ sdk.InputConnection) {
	panic("this should never be called")
}

func (i *InputWithNulls) OnRecordPacket(_ sdk.InputConnection) {
	panic("this should never be called")
}

func (i *InputWithNulls) OnComplete() {
	info, _ := sdk.NewOutgoingRecordInfo([]sdk.NewOutgoingField{
		sdk.NewBoolField(`Field1`, `source`),
		sdk.NewInt32Field(`Field2`, `source`),
		sdk.NewDoubleField(`Field3`, `source`),
		sdk.NewStringField(`Field4`, `source`, 100),
		sdk.NewDateField(`Field5`, `source`),
		sdk.NewBlobField(`Field6`, `source`, 1000000),
	})
	i.output.Open(info)
	info.BoolFields[`Field1`].SetBool(true)
	info.IntFields[`Field2`].SetInt(12)
	info.FloatFields[`Field3`].SetFloat(123.4)
	info.StringFields[`Field4`].SetString(`hello world`)
	info.DateTimeFields[`Field5`].SetDateTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
	info.BlobFields[`Field6`].SetBlob([]byte{1, 2, 3})
	i.output.Write()
	info.BoolFields[`Field1`].SetNull()
	info.IntFields[`Field2`].SetNull()
	info.FloatFields[`Field3`].SetNull()
	info.StringFields[`Field4`].SetNull()
	info.DateTimeFields[`Field5`].SetNull()
	info.BlobFields[`Field6`].SetNull()
	i.output.Write()
	info.BoolFields[`Field1`].SetBool(true)
	info.IntFields[`Field2`].SetInt(12)
	info.FloatFields[`Field3`].SetFloat(123.4)
	info.StringFields[`Field4`].SetString(`hello world`)
	info.DateTimeFields[`Field5`].SetDateTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
	info.BlobFields[`Field6`].SetBlob([]byte{1, 2, 3})
	i.output.Write()
}

type PassThroughTool struct {
	output sdk.OutputAnchor
	info   *sdk.OutgoingRecordInfo
}

func (p *PassThroughTool) Init(provider sdk.Provider) {
	p.output = provider.GetOutputAnchor(`Output`)
}

func (p *PassThroughTool) OnInputConnectionOpened(connection sdk.InputConnection) {
	p.info = connection.Metadata().Clone().GenerateOutgoingRecordInfo()
	p.output.Open(p.info)
}

func (p *PassThroughTool) OnRecordPacket(connection sdk.InputConnection) {
	packet := connection.Read()
	for packet.Next() {
		p.info.CopyFrom(packet.Record())
		p.output.Write()
	}
}

func (p *PassThroughTool) OnComplete() {}

func TestRegister(t *testing.T) {
	config := `<Configuration></Configuration>`
	implementation := &TestImplementation{}
	sdk.RegisterToolTest(implementation, 1, config)
	if !implementation.DidInit {
		t.Fatalf(`implementation did not init`)
	}
	if implementation.Config != config {
		t.Fatalf(`expected '%v' but got '%v'`, config, implementation.Config)
	}
}

func TestProviderIo(t *testing.T) {
	implementation := &TestImplementation{}
	sdk.RegisterToolTest(implementation, 1, ``)
	implementation.TestIo()
}

func TestDefaultTestProviderEnvironment(t *testing.T) {
	implementation := &TestImplementation{}
	sdk.RegisterToolTest(implementation, 5, ``)
	if id := implementation.Provider.Environment().ToolId(); id != 5 {
		t.Fatalf(`expected 5 but got %v`, id)
	}
	if updateOnly := implementation.Provider.Environment().UpdateOnly(); updateOnly {
		t.Fatalf(`expected false but got true`)
	}
	if installDir := implementation.Provider.Environment().AlteryxInstallDir(); installDir != `` {
		t.Fatalf(`expected '' but got '%v'`, installDir)
	}
	if locale := implementation.Provider.Environment().AlteryxLocale(); locale != `en` {
		t.Fatalf(`expected 'en' but got '%v'`, locale)
	}
	if version := implementation.Provider.Environment().DesignerVersion(); version != `TestHarness` {
		t.Fatalf(`expected 'TestHarness' but got '%v'`, version)
	}
	if updateMode := implementation.Provider.Environment().UpdateMode(); updateMode != `` {
		t.Fatalf(`expected '' but got '%v'`, updateMode)
	}
	if workflowDir := implementation.Provider.Environment().WorkflowDir(); workflowDir != `` {
		t.Fatalf(`expected '' but got '%v'`, workflowDir)
	}
}

func TestCustomTestProviderEnvironmentOptions(t *testing.T) {
	implementation := &TestImplementation{}
	sdk.RegisterToolTest(implementation, 5, ``,
		sdk.UpdateOnly(true),
		sdk.UpdateMode(`custom updateMode`),
		sdk.WorkflowDir(`custom workflowDir`),
		sdk.AlteryxLocale(`fr`))
	if updateOnly := implementation.Provider.Environment().UpdateOnly(); !updateOnly {
		t.Fatalf(`expected true but got false`)
	}
	if locale := implementation.Provider.Environment().AlteryxLocale(); locale != `fr` {
		t.Fatalf(`expected 'fr' but got '%v'`, locale)
	}
	if updateMode := implementation.Provider.Environment().UpdateMode(); updateMode != `custom updateMode` {
		t.Fatalf(`expected 'custom updateMode' but got '%v'`, updateMode)
	}
	if workflowDir := implementation.Provider.Environment().WorkflowDir(); workflowDir != `custom workflowDir` {
		t.Fatalf(`expected 'custom workflowDir' but got '%v'`, workflowDir)
	}
}

func TestUpdateConfig(t *testing.T) {
	implementation := &TestImplementation{}
	sdk.RegisterToolTest(implementation, 1, `<Configuration></Configuration>`)
	newConfig := `<Configuration><Something /></Configuration`
	implementation.Provider.Environment().UpdateToolConfig(newConfig)
	config := implementation.Provider.ToolConfig()
	if config != newConfig {
		t.Fatalf(`expected '%v' but got '%v'`, newConfig, config)
	}
}

func TestGettingOutputAnchorTwiceIsSameObject(t *testing.T) {
	implementation := &TestImplementation{}
	sdk.RegisterToolTest(implementation, 1, ``)
	output1 := implementation.Provider.GetOutputAnchor(`Output`)
	output2 := implementation.Provider.GetOutputAnchor(`Output`)
	if output1 != output2 {
		t.Fatalf(`expected the same outputAnchor object but got 2 different objects: %v and %v`, output1, output2)
	}
}

func TestSimulateInputTool(t *testing.T) {
	implementation := &TestImplementation{}
	runner := sdk.RegisterToolTest(implementation, 1, ``)
	if implementation.Output == nil {
		t.Fatalf(`expected an output anchor but got nil`)
	}
	runner.SimulateLifecycle()
	if !implementation.DidOnComplete {
		t.Fatalf(`did not run OnComplete but expected it to`)
	}
	if implementation.DidOnInputConnectionOpened {
		t.Fatalf(`OnInputConnectionOpened was called but it should not have been`)
	}
	if implementation.DidOnRecordPacket {
		t.Fatalf(`OnRecordPacket was called but it should not have been`)
	}
}

func TestOutputRecordsToTestRunner(t *testing.T) {
	implementation := &TestInputTool{}
	runner := sdk.RegisterToolTest(implementation, 1, ``)
	collector := runner.CaptureOutgoingAnchor(`Output`)
	runner.SimulateLifecycle()

	if collector.Name != `Output` {
		t.Fatalf(`expected 'Output' but got '%v'`, collector.Name)
	}
	if fields := collector.Config.NumFields(); fields != 17 {
		t.Fatalf("expected 17 fields but got %v", fields)
	}
	outputConfig := implementation.Output.Metadata()
	if outputConfig != implementation.OutputConfig {
		t.Fatalf(`expected same instance but got %v and %v`, outputConfig, implementation.OutputConfig)
	}
	if length := len(collector.Data); length != 17 {
		t.Fatalf(`expected 17 fields but got %v`, length)
	}
	if length := len(collector.Data[`Field3`]); length != 10 {
		t.Fatalf(`expected [0 1 2 3 4 5 6 7 8 9] but got %v`, collector.Data[`Field3`])
	}
	for i := 0; i < 10; i++ {
		if !bytes.Equal(collector.Data[`Field1`][i].([]byte), []byte{byte(i)}) {
			t.Fatalf(`expected [[0] [1] [2] [3] [4] [5] [6] [7] [8] [9]] but got %v`, collector.Data[`Field1`])
		}
		if collector.Data[`Field2`][i] != (i%2 == 0) {
			t.Fatalf(`expected [true false true false true false true false true false] but got %v`, collector.Data[`Field2`])
		}
		if collector.Data[`Field3`][i] != i {
			t.Fatalf(`expected [0 1 2 3 4 5 6 7 8 9] but got %v`, collector.Data[`Field3`])
		}
		if collector.Data[`Field4`][i] != i {
			t.Fatalf(`expected [0 1 2 3 4 5 6 7 8 9] but got %v`, collector.Data[`Field4`])
		}
		if collector.Data[`Field5`][i] != i {
			t.Fatalf(`expected [0 1 2 3 4 5 6 7 8 9] but got %v`, collector.Data[`Field5`])
		}
		if collector.Data[`Field6`][i] != i {
			t.Fatalf(`expected [0 1 2 3 4 5 6 7 8 9] but got %v`, collector.Data[`Field6`])
		}
		if collector.Data[`Field7`][i] != float64(i) {
			t.Fatalf(`expected [0 1 2 3 4 5 6 7 8 9] but got %v`, collector.Data[`Field7`])
		}
		if collector.Data[`Field8`][i] != float64(i) {
			t.Fatalf(`expected [0 1 2 3 4 5 6 7 8 9] but got %v`, collector.Data[`Field8`])
		}
		if collector.Data[`Field9`][i] != float64(i) {
			t.Fatalf(`expected [0 1 2 3 4 5 6 7 8 9] but got %v`, collector.Data[`Field9`])
		}
		if collector.Data[`Field10`][i] != strconv.Itoa(i) {
			t.Fatalf(`expected [0 1 2 3 4 5 6 7 8 9] but got %v`, collector.Data[`Field10`])
		}
		if collector.Data[`Field11`][i] != strconv.Itoa(i) {
			t.Fatalf(`expected [0 1 2 3 4 5 6 7 8 9] but got %v`, collector.Data[`Field11`])
		}
		if collector.Data[`Field12`][i] != strconv.Itoa(i) {
			t.Fatalf(`expected [0 1 2 3 4 5 6 7 8 9] but got %v`, collector.Data[`Field12`])
		}
		if collector.Data[`Field13`][i] != strconv.Itoa(i) {
			t.Fatalf(`expected [0 1 2 3 4 5 6 7 8 9] but got %v`, collector.Data[`Field13`])
		}
		if collector.Data[`Field14`][i] != time.Date(2020, 1, i, 0, 0, 0, 0, time.UTC) {
			t.Fatalf(`expected [2019-12-31 00:00:00 +0000 UTC 2020-01-01 00:00:00 +0000 UTC 2020-01-02 00:00:00 +0000 UTC 2020-01-03 00:00:00 +0000 UTC 2020-01-04 00:00:00 +0000 UTC 2020-01-05 00:00:00 +0000 UTC 2020-01-06 00:00:00 +0000 UTC 2020-01-07 00:00:00 +0000 UTC 2020-01-08 00:00:00 +0000 UTC 2020-01-09 00:00:00 +0000 UTC] but got %v`, collector.Data[`Field14`])
		}
		if collector.Data[`Field15`][i] != time.Date(2020, 1, i, 0, 0, 0, 0, time.UTC) {
			t.Fatalf(`expected [2019-12-31 00:00:00 +0000 UTC 2020-01-01 00:00:00 +0000 UTC 2020-01-02 00:00:00 +0000 UTC 2020-01-03 00:00:00 +0000 UTC 2020-01-04 00:00:00 +0000 UTC 2020-01-05 00:00:00 +0000 UTC 2020-01-06 00:00:00 +0000 UTC 2020-01-07 00:00:00 +0000 UTC 2020-01-08 00:00:00 +0000 UTC 2020-01-09 00:00:00 +0000 UTC] but got %v`, collector.Data[`Field15`])
		}
		if !bytes.Equal(collector.Data[`Field16`][i].([]byte), []byte{byte(i)}) {
			t.Fatalf(`expected [[0] [1] [2] [3] [4] [5] [6] [7] [8] [9]] but got %v`, collector.Data[`Field16`])
		}
		if collector.Data[`Field17`][i] != time.Date(0, 1, 1, 12, 31, i, 0, time.UTC) {
			t.Fatalf(`expected [0000-01-01 12:31:00 +0000 UTC 0000-01-01 12:31:01 +0000 UTC 0000-01-01 12:31:02 +0000 UTC 0000-01-01 12:31:03 +0000 UTC 0000-01-01 12:31:04 +0000 UTC  0000-01-01 12:31:05 +0000 UTC 0000-01-01 12:31:06 +0000 UTC 0000-01-01 12:31:07 +0000 UTC 0000-01-01 12:31:08 +0000 UTC 0000-01-01 12:31:09 +0000 UTC] but got %v`, collector.Data[`Field17`])
		}
	}
	if collector.Progress != 1.0 {
		t.Fatalf(`expected 1.0 but got %v`, collector.Progress)
	}
}

func TestMultipleOutputConnections(t *testing.T) {
	implementation := &TestInputTool{}
	runner := sdk.RegisterToolTest(implementation, 1, ``)
	collector1 := runner.CaptureOutgoingAnchor(`Output`)
	collector2 := runner.CaptureOutgoingAnchor(`Output`)
	runner.SimulateLifecycle()

	if length := len(collector1.Data[`Field1`]); length != 10 {
		t.Fatalf(`expected 10 records in collector1 but got %v`, collector1.Data[`Field1`])
	}
	if length := len(collector2.Data[`Field1`]); length != 10 {
		t.Fatalf(`expected 10 records in collector2 but got %v`, collector1.Data[`Field1`])
	}
	if connections := implementation.Output.NumConnections(); connections != 2 {
		t.Fatalf(`expected 2 output connections but got %v`, connections)
	}
}

func TestRecordLargerThanCache(t *testing.T) {
	if runtime.GOOS != `windows` {
		t.Skipf(`TestRecordLargerThanCache fails on Mac and I am not sure why`)
	}
	implementation := &InputRecordLargerThanCache{}
	runner := sdk.RegisterToolTest(implementation, 1, ``)
	collector := runner.CaptureOutgoingAnchor(`Output`)
	runner.SimulateLifecycle()
	if value := collector.Data[`Field1`][0]; value != `hello world` {
		t.Fatalf(`expected first record to be 'hello world' but got '%v'`, value)
	}
	if collector.Data[`Field1`][1] != strings.Repeat(`ABCDEFGHIJKLMNOPQRSTUVWXYZ`, 200000) {
		t.Fatalf(`The second record did not have the expected value`)
	}
	if value := collector.Data[`Field1`][2]; value != `zyxwvutsrqponmlkjihgfedcba` {
		t.Fatalf(`expected third record to be 'zyxwvutsrqponmlkjihgfedcba' but got '%v'`, value)
	}
}

func TestRecordsWithNulls(t *testing.T) {
	implementation := &InputWithNulls{}
	runner := sdk.RegisterToolTest(implementation, 1, ``)
	collector := runner.CaptureOutgoingAnchor(`Output`)
	runner.SimulateLifecycle()

	for row := 0; row < 3; row++ {
		for field := 1; field < 7; field++ {
			fieldName := fmt.Sprintf(`Field%v`, field)
			if row%2 == 0 {
				if collector.Data[fieldName][row] == nil {
					t.Fatalf(`expected non-nil in %v row %v but got nil`, fieldName, row)
				}
			} else {
				if value := collector.Data[fieldName][row]; value != nil {
					t.Fatalf(`expected nil in %v row %v but got %v`, fieldName, row, value)
				}
			}
		}
	}
}

func TestPassthroughSimulation(t *testing.T) {
	implementation := &PassThroughTool{}
	runner := sdk.RegisterToolTest(implementation, 1, ``)
	collector := runner.CaptureOutgoingAnchor(`Output`)
	runner.ConnectInput(`Input`, `sdk_test_passthrough_simulation.txt`)
	runner.SimulateLifecycle()
	if len(collector.Data) != 17 {
		t.Fatalf(`expected 17 fields but got %v`, len(collector.Data))
	}
	if recordCount := len(collector.Data[`Field1`]); recordCount != 4 {
		t.Fatalf(`expected 4 records but got %v`, recordCount)
	}
	if expectedValues := []interface{}{true, false, nil, true}; !reflect.DeepEqual(expectedValues, collector.Data[`Field1`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field1`])
	}
	if expectedValues := []interface{}{2, 2, nil, 42}; !reflect.DeepEqual(expectedValues, collector.Data[`Field2`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field2`])
	}
	if expectedValues := []interface{}{100, -100, nil, -110}; !reflect.DeepEqual(expectedValues, collector.Data[`Field3`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field3`])
	}
	if expectedValues := []interface{}{1000, -1000, nil, 392}; !reflect.DeepEqual(expectedValues, collector.Data[`Field4`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field4`])
	}
	if expectedValues := []interface{}{10000, -10000, nil, 2340}; !reflect.DeepEqual(expectedValues, collector.Data[`Field5`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field5`])
	}
	if expectedValues := []interface{}{float64(float32(12.34)), float64(float32(-12.34)), nil, float64(float32(12))}; !reflect.DeepEqual(expectedValues, collector.Data[`Field6`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field6`])
	}
	if expectedValues := []interface{}{1.23, -1.23, nil, 41.22}; !reflect.DeepEqual(expectedValues, collector.Data[`Field7`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field7`])
	}
	if expectedValues := []interface{}{234.56, -234.56, nil, 98.2}; !reflect.DeepEqual(expectedValues, collector.Data[`Field8`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field8`])
	}
	if expectedValues := []interface{}{`ABC`, `DE|"FG`, nil, ``}; !reflect.DeepEqual(expectedValues, collector.Data[`Field9`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field9`])
	}
	if expectedValues := []interface{}{`Hello `, `HIJK`, nil, ``}; !reflect.DeepEqual(expectedValues, collector.Data[`Field10`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field10`])
	}
	if expectedValues := []interface{}{` World`, `LMNOP`, nil, ``}; !reflect.DeepEqual(expectedValues, collector.Data[`Field11`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field11`])
	}
	if expectedValues := []interface{}{`abcdefg`, "QRSTU\r\nVWXYZ", nil, ``}; !reflect.DeepEqual(expectedValues, collector.Data[`Field12`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field12`])
	}
	if expectedValues := []interface{}{time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, 2, 3, 0, 0, 0, 0, time.UTC), nil, time.Date(2020, 2, 13, 0, 0, 0, 0, time.UTC)}; !reflect.DeepEqual(expectedValues, collector.Data[`Field13`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field13`])
	}
	if expectedValues := []interface{}{time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC), time.Date(2020, 1, 2, 13, 14, 15, 0, time.UTC), nil, time.Date(2020, 11, 2, 13, 14, 15, 0, time.UTC)}; !reflect.DeepEqual(expectedValues, collector.Data[`Field14`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field14`])
	}
	if expectedValues := []interface{}{nil, nil, nil, nil}; !reflect.DeepEqual(expectedValues, collector.Data[`Field15`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field15`])
	}
	if expectedValues := []interface{}{nil, nil, nil, nil}; !reflect.DeepEqual(expectedValues, collector.Data[`Field16`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field16`])
	}
	if expectedValues := []interface{}{time.Date(0, 1, 1, 10, 1, 1, 0, time.UTC), nil, time.Date(0, 1, 1, 17, 2, 1, 0, time.UTC), nil}; !reflect.DeepEqual(expectedValues, collector.Data[`Field17`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field17`])
	}
}

type WriteBeforeOpeningOutput struct {
	output sdk.OutputAnchor
}

func (p *WriteBeforeOpeningOutput) Init(provider sdk.Provider) {
	p.output = provider.GetOutputAnchor(`Output`)
}

func (p *WriteBeforeOpeningOutput) OnInputConnectionOpened(_ sdk.InputConnection) {
}

func (p *WriteBeforeOpeningOutput) OnRecordPacket(connection sdk.InputConnection) {
	packet := connection.Read()
	for packet.Next() {
		p.output.Write()
	}
}

func (p *WriteBeforeOpeningOutput) OnComplete() {
}

func TestWritingOutputBeforeOpenShouldPanic(t *testing.T) {
	defer func() {
		expected := `you are writing to output anchor 'Output' before it has been opened; call Open() before writing records`
		if r := recover(); r != expected {
			t.Fatalf("expected\n%v\nbut got\n%v", expected, r)
		} else {
			t.Logf(`recovered from: %v`, r)
		}
	}()

	implementation := &WriteBeforeOpeningOutput{}
	runner := sdk.RegisterToolTest(implementation, 1, ``)
	_ = runner.CaptureOutgoingAnchor(`Output`)
	runner.ConnectInput(`Input`, `sdk_test_passthrough_simulation.txt`)
	runner.SimulateLifecycle()
	t.Fatalf(`expected a panic but it did not happen`)
}

type OpenBeforeAddingConnectionPlugin struct {
	output     sdk.OutputAnchor
	recordInfo *sdk.OutgoingRecordInfo
	id         int
}

func (i *OpenBeforeAddingConnectionPlugin) Init(p sdk.Provider) {
	i.output = p.GetOutputAnchor(`Output`)
	i.recordInfo, _ = sdk.NewOutgoingRecordInfo([]sdk.NewOutgoingField{
		sdk.NewInt64Field(`ID`, ``),
	})
	i.output.Open(i.recordInfo)
}

func (i *OpenBeforeAddingConnectionPlugin) OnInputConnectionOpened(_ sdk.InputConnection) {}

func (i *OpenBeforeAddingConnectionPlugin) OnRecordPacket(connection sdk.InputConnection) {
	packet := connection.Read()
	for packet.Next() {
		i.recordInfo.IntFields[`ID`].SetInt(i.id)
		i.id++
		i.output.Write()
	}
}

func (i *OpenBeforeAddingConnectionPlugin) OnComplete() {}

func TestInitOutputBeforeAddingOutgoingConnection(t *testing.T) {
	plugin := &OpenBeforeAddingConnectionPlugin{}
	runner := sdk.RegisterToolTest(plugin, 1, ``)
	collector := runner.CaptureOutgoingAnchor(`Output`)
	runner.ConnectInput(`Input`, `sdk_test_passthrough_simulation.txt`)
	runner.SimulateLifecycle()
	if length := len(collector.Data); length != 1 {
		t.Fatalf(`expected length of 1 but got %v`, length)
	}
	if field, ok := collector.Data[`ID`]; !ok {
		t.Fatalf(`expected a field called 'ID' but it did not exist`)
	} else {
		t.Logf(`ID data: %v`, field)
	}
}

type statusTester struct {
	connection1 sdk.InputConnection
	connection2 sdk.InputConnection
	err1        error
	err2        error
}

func (t *statusTester) checkExpectedConn1Status(expected sdk.Status) {
	if t.err1 != nil {
		return
	}
	if status := t.connection1.Status(); status != expected {
		t.err1 = fmt.Errorf(`expected connection1 status of %v but got %v`, expected, status)
	}
}

func (t *statusTester) checkExpectedConn2Status(expected sdk.Status) {
	if t.err2 != nil {
		return
	}
	if status := t.connection2.Status(); status != expected {
		t.err2 = fmt.Errorf(`expected connection2 status of %v but got %v`, expected, status)
	}
}

func (t *statusTester) Init(_ sdk.Provider) {}

func (t *statusTester) OnInputConnectionOpened(connection sdk.InputConnection) {
	if connection.Name() == `Input1` {
		t.connection1 = connection
		t.checkExpectedConn1Status(sdk.Initialized)
	} else {
		t.connection2 = connection
		t.checkExpectedConn2Status(sdk.Initialized)
	}
}

func (t *statusTester) OnRecordPacket(connection sdk.InputConnection) {
	if connection.Name() == `Input1` {
		t.checkExpectedConn1Status(sdk.ReceivingRecords)
	} else {
		t.checkExpectedConn2Status(sdk.ReceivingRecords)
	}
}

func (t *statusTester) OnComplete() {
	t.checkExpectedConn1Status(sdk.Closed)
	if t.connection2 != nil {
		t.checkExpectedConn2Status(sdk.Closed)
	}
}

func TestStatus(t *testing.T) {
	plugin := &statusTester{}
	runner := sdk.RegisterToolTest(plugin, 1, ``)
	runner.ConnectInput(`Input1`, `sdk_test_passthrough_simulation.txt`)
	runner.SimulateLifecycle()
	if plugin.err1 != nil {
		t.Fatalf(`expected no error but got: %v`, plugin.err1)
	}
}

func TestStatusMultipleInputs(t *testing.T) {
	plugin := &statusTester{}
	runner := sdk.RegisterToolTest(plugin, 1, ``)
	runner.ConnectInput(`Input1`, `sdk_test_passthrough_simulation.txt`)
	runner.ConnectInput(`Input2`, `sdk_test_passthrough_simulation.txt`)
	runner.SimulateLifecycle()
	if plugin.err1 != nil {
		t.Fatalf(`expected no error for err1 but got: %v`, plugin.err1)
	}
	if plugin.err2 != nil {
		t.Fatalf(`expected no error for err2 but got: %v`, plugin.err2)
	}
}

type outputAnchorCloseTester struct {
	output1 sdk.OutputAnchor
	output2 sdk.OutputAnchor
	err     error
}

func (t *outputAnchorCloseTester) Init(provider sdk.Provider) {
	info, _ := sdk.NewOutgoingRecordInfo([]sdk.NewOutgoingField{
		sdk.NewByteField(`Id`, ``),
	})
	t.output1 = provider.GetOutputAnchor(`Output1`)
	t.output1.Open(info)
	t.output2 = provider.GetOutputAnchor(`Output2`)
	t.output2.Open(info)
	if !t.output1.IsOpen() {
		t.err = fmt.Errorf(`output1 was not open but should have been`)
		return
	}
	if !t.output2.IsOpen() {
		t.err = fmt.Errorf(`output2 was not open but should have been`)
	}
}

func (t *outputAnchorCloseTester) OnInputConnectionOpened(_ sdk.InputConnection) {
	panic("implement me")
}

func (t *outputAnchorCloseTester) OnRecordPacket(_ sdk.InputConnection) {
	panic("implement me")
}

func (t *outputAnchorCloseTester) OnComplete() {
	if t.err != nil {
		return
	}
	t.output1.Close()
	if t.output1.IsOpen() {
		t.err = fmt.Errorf(`output1 should have been closed but was open`)
		return
	}
}

func TestCloseOutputAnchor(t *testing.T) {
	plugin := &outputAnchorCloseTester{}
	runner := sdk.RegisterToolTest(plugin, 1, ``)
	runner.SimulateLifecycle()

	if plugin.err != nil {
		t.Fatalf(`expected no error but got: %v`, plugin.err.Error())
	}
	if plugin.output2.IsOpen() {
		t.Fatalf(`expected output2 to be closed but it was not`)
	}
}

type testCreateTempFile struct {
	filePath string
}

func (t *testCreateTempFile) Init(provider sdk.Provider) {
	t.filePath = provider.Io().CreateTempFile(`yxdb`)
}

func (t *testCreateTempFile) OnInputConnectionOpened(_ sdk.InputConnection) {}

func (t *testCreateTempFile) OnRecordPacket(_ sdk.InputConnection) {}

func (t *testCreateTempFile) OnComplete() {}

func TestCreateTempFile(t *testing.T) {
	plugin := &testCreateTempFile{}
	runner := sdk.RegisterToolTest(plugin, 1, ``)
	runner.SimulateLifecycle()

	if ext := filepath.Ext(plugin.filePath); ext != `.yxdb` {
		t.Fatalf(`expected '.yxdb' but got '%v'`, ext)
	}
	t.Logf(plugin.filePath)
}

type byteTester struct {
	output sdk.OutputAnchor
}

func (b *byteTester) Init(provider sdk.Provider) {
	b.output = provider.GetOutputAnchor(`Output`)
}

func (b *byteTester) OnInputConnectionOpened(_ sdk.InputConnection) {
}

func (b *byteTester) OnRecordPacket(_ sdk.InputConnection) {
}

func (b *byteTester) OnComplete() {
	info, _ := sdk.NewOutgoingRecordInfo([]sdk.NewOutgoingField{
		sdk.NewByteField(`bytes`, `Byte Tester`),
	})
	b.output.Open(info)
	info.IntFields[`bytes`].SetInt(250)
	b.output.Write()
	b.output.Close()
}

func TestWriteBytes(t *testing.T) {
	plugin := &byteTester{}
	runner := sdk.RegisterToolTest(plugin, 1, ``)
	output := runner.CaptureOutgoingAnchor(`Output`)
	runner.SimulateLifecycle()

	bytesField, ok := output.Data[`bytes`]
	if !ok {
		t.Fatalf(`[bytes] not outputted`)
	}
	if len(bytesField) != 1 {
		t.Fatalf(`expected 1 record but got %v`, len(bytesField))
	}
	if bytesField[0] != 250 {
		t.Fatalf(`expected 250 but got %v`, bytesField[0])
	}
}
