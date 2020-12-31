package api_new_test

import (
	"github.com/tlarsen7572/goalteryx/api_new"
	"testing"
)

func TestDataSize(t *testing.T) {
	editor := api_new.EditingRecordInfo{}
	editor.AddBoolField(`Field1`, `source`)
	info := editor.GenerateOutgoingRecordInfo()
	field, _ := info.GetBoolField(`Field1`)
	field.SetBool(false)
	size := info.DataSize()
	if size != 1 {
		t.Fatalf(`expected size 1 but got %v`, size)
	}
}
