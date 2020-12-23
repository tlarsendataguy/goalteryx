package api_new

type EditingRecordInfo struct {
	fields []IncomingField
}

func (i EditingRecordInfo) NumFields() int {
	return len(i.fields)
}
