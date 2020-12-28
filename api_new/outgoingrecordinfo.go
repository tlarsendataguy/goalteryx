package api_new

import (
	"fmt"
)

type OutgoingRecordInfo struct {
	outgoingFields []*outgoingField
}

type outgoingField struct {
	Name         string
	Type         string
	Source       string
	Size         int
	Scale        int
	CopyFrom     BytesGetter
	CurrentValue []byte
}

func (f *outgoingField) SetBool(value bool) {
	if value {
		f.CurrentValue[0] = 1
	} else {
		f.CurrentValue[0] = 0
	}
}

func (f *outgoingField) SetNullBool() {
	f.CurrentValue[0] = 2
}

func (f *outgoingField) GetCurrentBool() (bool, bool) {
	if f.CurrentValue[0] == 2 {
		return false, true
	}
	return f.CurrentValue[0] == 1, false
}

func (f *outgoingField) SetByte(value int) {
	f.CurrentValue[0] = byte(value)
	f.CurrentValue[1] = 0
}

func (f *outgoingField) SetNullByte() {
	f.CurrentValue[1] = 1
}

func (f *outgoingField) GetCurrentByte() (int, bool) {
	if f.CurrentValue[1] == 1 {
		return 0, true
	}
	return int(f.CurrentValue[0]), false
}

type OutgoingBoolField interface {
	SetBool(bool)
	SetNullBool()
	GetCurrentBool() (bool, bool)
}

type OutgoingByteField interface {
	SetByte(int)
	SetNullByte()
	GetCurrentByte() (int, bool)
}

func (i *OutgoingRecordInfo) GetBoolField(name string) (OutgoingBoolField, error) {
	for _, field := range i.outgoingFields {
		if field.Name == name {
			if field.Type != `Bool` {
				return nil, fmt.Errorf(`the '%v' field is not a bool field, it is '%v'`, name, field.Type)
			}
			return field, nil
		}
	}
	return nil, fmt.Errorf(`there is no '%v' field in the record`, name)
}

func (i *OutgoingRecordInfo) GetByteField(name string) (OutgoingByteField, error) {
	for _, field := range i.outgoingFields {
		if field.Name == name {
			if field.Type != `Byte` {
				return nil, fmt.Errorf(`the '%v' field is not a byte field, it is '%v'`, name, field.Type)
			}
			return field, nil
		}
	}
	return nil, fmt.Errorf(`there is no '%v' field in the record`, name)
}
