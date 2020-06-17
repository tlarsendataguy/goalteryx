package recordblob

import "unsafe"

// RecordBlob is a struct that contains the record blob.  Having the record blob embedded in a Go struct prevents
// the unsafe package from being required in client code.
type RecordBlob interface {
	Blob() unsafe.Pointer
}

type recordBlob struct {
	blob unsafe.Pointer
}

func (blob *recordBlob) Blob() unsafe.Pointer {
	return blob.blob
}

func NewRecordBlob(record unsafe.Pointer) RecordBlob {
	return &recordBlob{blob: record}
}
