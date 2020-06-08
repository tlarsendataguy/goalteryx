package main_test

import (
	"encoding/xml"
	implementation "github.com/tlarsen7572/goalteryx/implementation_example"
	"testing"
)

func TestXmlConfig(t *testing.T) {
	configXmlStr := `<Configuration><Field>FixedDecimalField</Field></Configuration><Annotation DisplayMode="0"><Name/><DefaultAnnotationText/><Left value="False"/></Annotation>`
	var c implementation.ConfigXml
	err := xml.Unmarshal([]byte(configXmlStr), &c)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if c.Field != `FixedDecimalField` {
		t.Fatalf(`expected 'FixedDecimalField' but got '%v'`, c.Field)
	}
}
