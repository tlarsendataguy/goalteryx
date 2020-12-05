package api_new_test

import (
	"github.com/tlarsen7572/goalteryx/api_new"
	"testing"
)

type TestImplementation struct {
	DidInit  bool
	Config   string
	Provider api_new.Provider
}

func (t *TestImplementation) TestIo() {
	t.Provider.Io().Info(`test1`)
	t.Provider.Io().Warn(`test2`)
	t.Provider.Io().Error(`test3`)
	t.Provider.Io().UpdateProgress(0.10)
}

func (t *TestImplementation) TestEnvironmentToolId() int {
	return t.Provider.Environment().ToolId()
}

func (t *TestImplementation) TestEnvironment() {
	t.Provider.Environment().AlteryxInstallDir()
	t.Provider.Environment().AlteryxLocale()
	t.Provider.Environment().DesignerVersion()
	t.Provider.Environment().UpdateMode()
	t.Provider.Environment().UpdateOnly()
	t.Provider.Environment().WorkflowDir()
}

func (t *TestImplementation) Init(provider api_new.Provider) {
	t.DidInit = true
	t.Config = provider.ToolConfig()
	t.Provider = provider
}

func (t *TestImplementation) OnInputConnectionOpened(connection api_new.InputConnection) {
	panic("implement me")
}

func (t *TestImplementation) OnRecordPacket(connection api_new.InputConnection) {
	panic("implement me")
}

func (t *TestImplementation) OnComplete() {
	panic("implement me")
}

func TestRegister(t *testing.T) {
	config := `<Configuration></Configuration>`
	implementation := &TestImplementation{}
	result := api_new.RegisterToolTest(implementation, 1, config)
	if result != 1 {
		t.Fatalf(`expected 1 but got %v`, result)
	}
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

func TestProviderEnvironment(t *testing.T) {
	implementation := &TestImplementation{}
	api_new.RegisterToolTest(implementation, 5, ``)
	id := implementation.TestEnvironmentToolId()
	if id != 5 {
		t.Fatalf(`expected 5 but got %v`, id)
	}
	implementation.TestEnvironment()
}
