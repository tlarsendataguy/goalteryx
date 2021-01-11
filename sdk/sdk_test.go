package sdk_test

import (
	"bytes"
	"fmt"
	"github.com/tlarsen7572/goalteryx/sdk"
	"reflect"
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
	info.BoolFields[`Field1`].SetNullBool()
	info.IntFields[`Field2`].SetNullInt()
	info.FloatFields[`Field3`].SetNullFloat()
	info.StringFields[`Field4`].SetNullString()
	info.DateTimeFields[`Field5`].SetNullDateTime()
	info.BlobFields[`Field6`].SetNullBlob()
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
	if fields := collector.Config.NumFields(); fields != 16 {
		t.Fatalf("expected 16 fields but got %v", fields)
	}
	outputConfig := implementation.Output.Metadata()
	if outputConfig != implementation.OutputConfig {
		t.Fatalf(`expected same instance but got %v and %v`, outputConfig, implementation.OutputConfig)
	}
	if length := len(collector.Data); length != 16 {
		t.Fatalf(`expected 16 fields but got %v`, length)
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
	}
	if progress := collector.Input.Progress(); progress != 1.0 {
		t.Fatalf(`expected 1.0 but got %v`, progress)
	}
}

func TestRecordLargerThanCache(t *testing.T) {
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
	if len(collector.Data) != 16 {
		t.Fatalf(`expected 16 fields but got %v`, len(collector.Data))
	}
	if recordCount := len(collector.Data[`Field1`]); recordCount != 4 {
		t.Fatalf(`expected 3 records but got %v`, recordCount)
	}
	if expectedValues := []interface{}{true, false, nil, true}; !reflect.DeepEqual(expectedValues, collector.Data[`Field1`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field1`])
	}
	if expectedValues := []interface{}{2, -2, nil, 42}; !reflect.DeepEqual(expectedValues, collector.Data[`Field2`]) {
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
	if expectedValues := []interface{}{`ABC`, `DE|"FG`, ``, ``}; !reflect.DeepEqual(expectedValues, collector.Data[`Field9`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field9`])
	}
	if expectedValues := []interface{}{`Hello `, `HIJK`, ``, `HIJK`}; !reflect.DeepEqual(expectedValues, collector.Data[`Field10`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field10`])
	}
	if expectedValues := []interface{}{` World`, `LMNOP`, ``, `LMN`}; !reflect.DeepEqual(expectedValues, collector.Data[`Field11`]) {
		t.Fatalf(`expected %v but got %v`, expectedValues, collector.Data[`Field11`])
	}
	if expectedValues := []interface{}{`abcdefg`, "QRSTU\r\nVWXYZ", ``, `qrstuvwxyz`}; !reflect.DeepEqual(expectedValues, collector.Data[`Field12`]) {
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
}
