package recordinfo

func GeneratorFromXml(recordInfoXml string) (Generator, error) {
	return recordInfoFromXml(recordInfoXml)
}

func NewGenerator() Generator {
	return &recordInfo{
		fieldNames: map[string]int{},
		blobLen:    0,
	}
}

func (info *recordInfo) GenerateRecordInfo() RecordInfo {
	return info
}
