package api_new_test

import (
	"github.com/tlarsen7572/goalteryx/api_new"
	"strconv"
	"testing"
	"time"
)

type TestImplementation struct {
	DidInit                    bool
	DidOnComplete              bool
	DidOnInputConnectionOpened bool
	DidOnRecordPacket          bool
	Config                     string
	Provider                   api_new.Provider
	Output                     api_new.OutputAnchor
}

func (t *TestImplementation) TestIo() {
	t.Provider.Io().Info(`test1`)
	t.Provider.Io().Warn(`test2`)
	t.Provider.Io().Error(`test3`)
	t.Provider.Io().UpdateProgress(0.10)
}

func (t *TestImplementation) Init(provider api_new.Provider) {
	t.DidInit = true
	t.Config = provider.ToolConfig()
	t.Provider = provider
	t.Output = provider.GetOutputAnchor(`Output`)
}

func (t *TestImplementation) OnInputConnectionOpened(_ api_new.InputConnection) {
	t.DidOnInputConnectionOpened = true
}

func (t *TestImplementation) OnRecordPacket(_ api_new.InputConnection) {
	t.DidOnRecordPacket = true
}

func (t *TestImplementation) OnComplete() {
	t.DidOnComplete = true
}

type TestInputTool struct {
	Provider     api_new.Provider
	Output       api_new.OutputAnchor
	OutputConfig *api_new.OutgoingRecordInfo
}

func (i *TestInputTool) Init(provider api_new.Provider) {
	i.Provider = provider
	i.Output = provider.GetOutputAnchor(`Output`)
}

func (i *TestInputTool) OnInputConnectionOpened(_ api_new.InputConnection) {
	panic("This should never be called")
}

func (i *TestInputTool) OnRecordPacket(_ api_new.InputConnection) {
	panic("This should never be called")
}

func (i *TestInputTool) OnComplete() {
	source := `source`
	output := api_new.NewOutgoingRecordInfo([]api_new.NewOutgoingField{
		api_new.NewBlobField(`Field1`, source, 100),
		api_new.NewBoolField(`Field2`, source),
		api_new.NewByteField(`Field3`, source),
		api_new.NewInt16Field(`Field4`, source),
		api_new.NewInt32Field(`Field5`, source),
		api_new.NewInt64Field(`Field6`, source),
		api_new.NewFloatField(`Field7`, source),
		api_new.NewDoubleField(`Field8`, source),
		api_new.NewFixedDecimalField(`Field9`, source, 19, 2),
		api_new.NewStringField(`Field10`, source, 100),
		api_new.NewWStringField(`Field11`, source, 100),
		api_new.NewV_StringField(`Field12`, source, 100000),
		api_new.NewV_WStringField(`Field13`, source, 100000),
		api_new.NewDateField(`Field14`, source),
		api_new.NewDateTimeField(`Field15`, source),
		api_new.NewSpatialObjField(`Field16`, source, 1000000),
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
}

func TestRegister(t *testing.T) {
	config := `<Configuration></Configuration>`
	implementation := &TestImplementation{}
	api_new.RegisterToolTest(implementation, 1, config)
	if !implementation.DidInit {
		t.Fatalf(`implementation did not init`)
	}
	if implementation.Config != config {
		t.Fatalf(`expected '%v' but got '%v'`, config, implementation.Config)
	}
}

func TestProviderIo(t *testing.T) {
	implementation := &TestImplementation{}
	api_new.RegisterToolTest(implementation, 1, ``)
	implementation.TestIo()
}

func TestDefaultTestProviderEnvironment(t *testing.T) {
	implementation := &TestImplementation{}
	api_new.RegisterToolTest(implementation, 5, ``)
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
	api_new.RegisterToolTest(implementation, 5, ``,
		api_new.UpdateOnly(true),
		api_new.UpdateMode(`custom updateMode`),
		api_new.WorkflowDir(`custom workflowDir`),
		api_new.AlteryxLocale(`fr`))
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
	api_new.RegisterToolTest(implementation, 1, `<Configuration></Configuration>`)
	newConfig := `<Configuration><Something /></Configuration`
	implementation.Provider.Environment().UpdateToolConfig(newConfig)
	config := implementation.Provider.ToolConfig()
	if config != newConfig {
		t.Fatalf(`expected '%v' but got '%v'`, newConfig, config)
	}
}

func TestGettingOutputAnchorTwiceIsSameObject(t *testing.T) {
	implementation := &TestImplementation{}
	api_new.RegisterToolTest(implementation, 1, ``)
	output1 := implementation.Provider.GetOutputAnchor(`Output`)
	output2 := implementation.Provider.GetOutputAnchor(`Output`)
	if output1 != output2 {
		t.Fatalf(`expected the same outputAnchor object but got 2 different objects: %v and %v`, output1, output2)
	}
}

func TestSimulateInputTool(t *testing.T) {
	implementation := &TestImplementation{}
	runner := api_new.RegisterToolTest(implementation, 1, ``)
	if implementation.Output == nil {
		t.Fatalf(`expected an output anchor but got nil`)
	}
	runner.SimulateInputTool()
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
	runner := api_new.RegisterToolTest(implementation, 1, ``)
	collector := runner.CaptureOutgoingAnchor(`Output`)
	runner.SimulateInputTool()
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
}
