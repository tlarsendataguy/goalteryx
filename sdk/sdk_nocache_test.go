package sdk_test

import (
	"github.com/tlarsendataguy/goalteryx/sdk"
	"testing"
)

type NoCache struct {
	RecordPackets int
	output        sdk.OutputAnchor
	info          *sdk.OutgoingRecordInfo
}

func (n *NoCache) Init(provider sdk.Provider) {
	n.output = provider.GetOutputAnchor(`Output`)
}

func (n *NoCache) OnInputConnectionOpened(connection sdk.InputConnection) {
	n.info = connection.Metadata().Clone().GenerateOutgoingRecordInfo()
	n.output.Open(n.info)
}

func (n *NoCache) OnRecordPacket(connection sdk.InputConnection) {
	n.RecordPackets++
	packet := connection.Read()
	for packet.Next() {
		n.info.CopyFrom(packet.Record())
		n.output.Write()
	}
}

func (n *NoCache) OnComplete(nRecordLimit int64) {
}

func TestNoCache(t *testing.T) {
	plugin := &NoCache{}
	runner := sdk.RegisterToolTest(plugin, 1, ``, sdk.NoCache(true))
	runner.ConnectInput(`Input`, `sdk_test_passthrough_simulation.txt`)
	collector := runner.CaptureOutgoingAnchor(`Output`)
	runner.SimulateLifecycle()

	if plugin.RecordPackets != 4 {
		t.Fatalf(`expected 4 packets but got %v`, plugin.RecordPackets)
	}
	if collector.PacketsReceived != 4 {
		t.Fatalf(`expected 4 packets but got %v`, collector.PacketsReceived)
	}
	if fields := len(collector.Data); fields != 16 {
		t.Fatalf(`expected 16 fields but got %v`, fields)
	}
	if records := len(collector.Data[`Field1`]); records != 4 {
		t.Fatalf(`expected 4 records but got %v`, records)
	}
	if value := collector.Data[`Field1`][0]; value != true {
		t.Fatalf(`expected true but got %v`, value)
	}
	field12, ok := collector.Data[`Field12`]
	if !ok {
		t.Fatalf(`expected Field12 but it did not exist`)
	}
	if field12[1] != "QRSTU\r\nVWXYZ" {
		t.Fatalf("expected 'QRSTU\r\nVWXYZ' but got '%v", field12[1])
	}
}
