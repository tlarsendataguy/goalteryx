package sdk_test

import (
	"github.com/tlarsendataguy/goalteryx/sdk"
	"testing"
)

type NoCache struct {
	RecordPackets int
}

func (n *NoCache) Init(provider sdk.Provider) {
}

func (n *NoCache) OnInputConnectionOpened(connection sdk.InputConnection) {
}

func (n *NoCache) OnRecordPacket(connection sdk.InputConnection) {
	n.RecordPackets++
}

func (n *NoCache) OnComplete() {
}

func TestNoCache(t *testing.T) {
	plugin := &NoCache{}
	runner := sdk.RegisterToolTest(plugin, 1, ``, sdk.NoCache(true))
	runner.ConnectInput(`Input`, `sdk_test_passthrough_simulation.txt`)
	runner.SimulateLifecycle()

	if plugin.RecordPackets != 4 {
		t.Fatalf(`expected 4 packets but got %v`, plugin.RecordPackets)
	}
}
