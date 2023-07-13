package sdk

import "testing"

type InternalTest struct{}

func (t *InternalTest) Init(_ Provider) {}

func (t *InternalTest) OnInputConnectionOpened(_ InputConnection) {}

func (t *InternalTest) OnRecordPacket(_ InputConnection) {}

func (t *InternalTest) OnComplete(nRecordLimit int64) {}

func TestPluginsAreRemovedAfterOnComplete(t *testing.T) {
	if len(tools) != 0 {
		t.Fatalf(`expected 0 tools but got %v`, len(tools))
	}
	implementation := &InternalTest{}
	runner := RegisterToolTest(implementation, 1, ``)

	if len(tools) != 1 {
		t.Fatalf(`expected 1 tool but got %v`, len(tools))
	}

	runner.SimulateLifecycle()

	if len(tools) != 0 {
		t.Fatalf(`expected 0 tools but got %v`, len(tools))
	}

}
