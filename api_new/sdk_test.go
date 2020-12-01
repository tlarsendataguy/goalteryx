package api_new_test

import (
	"github.com/tlarsen7572/goalteryx/api_new"
	"testing"
)

type TestImplementation struct {
	DidInit bool
}

func (t *TestImplementation) Init(provider api_new.Provider) {
	t.DidInit = true
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
	implementation := &TestImplementation{DidInit: false}
	result := api_new.RegisterToolTest(implementation, 1, `<Configuration></Configuration>`)
	if result != 1 {
		t.Fatalf(`expected 1 but got %v`, result)
	}
	if !implementation.DidInit {
		t.Fatalf(`implementation did not init`)
	}
}
