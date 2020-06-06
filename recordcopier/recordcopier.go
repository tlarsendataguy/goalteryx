// Package recordcopier provides a fast and easy way to copy data from an incoming record blob into a RecordInfo object.
package recordcopier

import (
	"fmt"
	"goalteryx/recordinfo"
	"unsafe"
)

// RecordCopier is used to easily and quickly copy data between RecordInfo objects.
type RecordCopier struct {
	destination recordinfo.RecordInfo
	source      recordinfo.RecordInfo
	indexMaps   []IndexMap
}

// IndexMap maps a source field index to a destination field index for Copy.
type IndexMap struct {
	DestinationIndex int
	SourceIndex      int
}

// New generates a new RecordCopier from the provided RecordInfo objects and IndexMap slice.  Only a basic validation
// is performed to check that the indices in the IndexMap slice are in range.
// TODO: We may want to perform FieldType and size validation here as well.
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

// Copy copies the contents of the record blob into the destination RecordInfo.
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
