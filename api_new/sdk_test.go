package api_new_test

import (
	"github.com/tlarsen7572/goalteryx/api_new"
	"testing"
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

func (t *TestImplementation) OnInputConnectionOpened(connection api_new.InputConnection) {
	t.DidOnInputConnectionOpened = true
}

func (t *TestImplementation) OnRecordPacket(connection api_new.InputConnection) {
	t.DidOnRecordPacket = true
}

func (t *TestImplementation) OnComplete() {
	t.DidOnComplete = true
}

type TestInputTool struct {
	Provider api_new.Provider
	Output   api_new.OutputAnchor
}

func (i *TestInputTool) Init(provider api_new.Provider) {
	i.Provider = provider
	i.Output = provider.GetOutputAnchor(`Output`)
}

func (i *TestInputTool) OnInputConnectionOpened(connection api_new.InputConnection) {
	panic("This should never be called")
}

func (i *TestInputTool) OnRecordPacket(connection api_new.InputConnection) {
	panic("This should never be called")
}

func (i *TestInputTool) OnComplete() {
	outputConfig := `<MetaInfo connection="Output">
<RecordInfo>
	<Field name="Field1" source="TextInput:" type="Byte"/>
	<Field name="Field2" size="1" source="TextInput:" type="String"/>
</RecordInfo>
</MetaInfo>`
	i.Output.Open(outputConfig)
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
	expectedConfig := `<MetaInfo connection="Output">
<RecordInfo>
	<Field name="Field1" source="TextInput:" type="Byte"/>
	<Field name="Field2" size="1" source="TextInput:" type="String"/>
</RecordInfo>
</MetaInfo>`
	if collector.Config != expectedConfig {
		t.Fatalf("expected\n'%v'\nbut got\n'%v'", expectedConfig, collector.Config)
	}
}
