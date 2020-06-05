package recordcopier

import (
	"fmt"
	"goalteryx/recordinfo"
	"unsafe"
)

type RecordCopier struct {
	destination recordinfo.RecordInfo
	source      recordinfo.RecordInfo
	indexMaps   []IndexMap
}

type IndexMap struct {
	DestinationIndex int
	SourceIndex      int
}

func New(destination recordinfo.RecordInfo, source recordinfo.RecordInfo, indexMaps []IndexMap) (*RecordCopier, error) {
	srcFields := source.NumFields()
	dstFields := destination.NumFields()
	for _, indexMap := range indexMaps {
		if index := indexMap.SourceIndex; index < 0 || index >= srcFields {
			return nil, fmt.Errorf(`error creating RecordCopier: source index %v is not between 0 and %v`, index, srcFields)
		}
		if index := indexMap.DestinationIndex; index < 0 || index >= dstFields {
			return nil, fmt.Errorf(`error creating RecordCopier: destination index %v is not between 0 and %v`, index, dstFields)
		}
	}
	return &RecordCopier{
		destination: destination,
		source:      source,
		indexMaps:   indexMaps,
	}, nil
}

func (copier *RecordCopier) Copy(record unsafe.Pointer) error {
	for _, indexMap := range copier.indexMaps {
		value, err := copier.source.GetRawBytesFromIndex(indexMap.SourceIndex, record)
		if err != nil {
			return err
		}
		err = copier.destination.SetIndexFromRawBytes(indexMap.DestinationIndex, value)
		if err != nil {
			return err
		}
	}
	return nil
}
