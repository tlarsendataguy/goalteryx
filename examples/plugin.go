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
