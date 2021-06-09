package main

import (
	"encoding/xml"
	"fmt"
	"github.com/tlarsen7572/goalteryx/sdk"
	"io/ioutil"
)

type Configuration struct {
	Password string
}

type Plugin struct {
	provider sdk.Provider
	output   sdk.OutputAnchor
	outInfo  *sdk.OutgoingRecordInfo
	config   Configuration
}

func (p *Plugin) Init(provider sdk.Provider) {
	provider.Io().Info(fmt.Sprintf(`Init tool %v`, provider.Environment().ToolId()))
	tempFilePath := provider.Io().CreateTempFile(`txt`)
	provider.Io().Info(fmt.Sprintf(`temp file: %v`, tempFilePath))
	err := ioutil.WriteFile(tempFilePath, []byte(`hello world`), 0600)
	if err != nil {
		provider.Io().Error(err.Error())
	}
	data, err := ioutil.ReadFile(tempFilePath)
	if err != nil {
		provider.Io().Error(err.Error())
	}
	provider.Io().Info(fmt.Sprintf(`temp file content: %v`, string(data)))
	configBytes := []byte(provider.ToolConfig())
	err = xml.Unmarshal(configBytes, &p.config)
	if err != nil {
		provider.Io().Error(err.Error())
	}
	password := provider.Io().DecryptPassword(p.config.Password)
	provider.Io().Info(fmt.Sprintf(`got password %v`, password))
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
