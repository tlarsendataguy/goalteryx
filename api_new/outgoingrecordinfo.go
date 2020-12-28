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
	intSetter    func(int, *outgoingField)
	intGetter    func(*outgoingField) int
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

func getByte(f *outgoingField) int {
	return int(f.CurrentValue[0])
}

func setByte(value int, f *outgoingField) {
	f.CurrentValue[0] = byte(value)
}

func (f *outgoingField) SetInt(value int) {
	f.intSetter(value, f)
	f.CurrentValue[f.Size] = 0
}

func (f *outgoingField) SetNullInt() {
	f.CurrentValue[f.Size] = 1
}

func (f *outgoingField) GetCurrentInt() (int, bool) {
	if f.CurrentValue[f.Size] == 1 {
		return 0, true
	}
	return f.intGetter(f), false
}

type OutgoingBoolField interface {
	SetBool(bool)
	SetNullBool()
	GetCurrentBool() (bool, bool)
}

type OutgoingIntField interface {
	SetInt(int)
	SetNullInt()
	GetCurrentInt() (int, bool)
}

func (i *OutgoingRecordInfo) GetBoolField(name string) (OutgoingBoolField, error) {
	return i.getField(name, []string{`Bool`}, `Int`)
}

func (i *OutgoingRecordInfo) GetIntField(name string) (OutgoingIntField, error) {
	return i.getField(name, []string{`Byte`}, `Int`)
}

func (i *OutgoingRecordInfo) getField(name string, types []string, label string) (*outgoingField, error) {
	for _, field := range i.outgoingFields {
		if field.Name == name {
			for _, ofType := range types {
				if field.Type == ofType {
					return field, nil
				}
			}
			return nil, fmt.Errorf(`the '%v' field is not a %v field, it is '%v'`, name, label, field.Type)
		}
	}
	return nil, fmt.Errorf(`there is no '%v' field in the record`, name)
}
