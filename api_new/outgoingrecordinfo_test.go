package api_new_test

import (
	"github.com/tlarsen7572/goalteryx/api_new"
	"testing"
)

func TestDataSize(t *testing.T) {
	info := api_new.NewOutgoingRecordInfo([]api_new.NewOutgoingField{
		api_new.NewBoolField(`Field1`, `source`),
	})
	size := info.DataSize()
	if size != 1 {
		t.Fatalf(`expected size 1 but got %v`, size)
	}
}
